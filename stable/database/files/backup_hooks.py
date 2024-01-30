import argparse
import http
import inspect
import json
import logging
import os
import signal
import subprocess
import sys
import time
import urllib
import wsgiref
from wsgiref import simple_server

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
LOGGER = logging.getLogger(__name__)

ARCHIVE_DIR = os.environ.get("NUODB_ARCHIVE_DIR", "/mnt/archive")
JOURNAL_DIR = os.environ.get("NUODB_JOURNAL_DIR")
USE_SUSPEND = os.environ.get("USE_SUSPEND") == "true"


def from_dir(base_dir, *args):
    if base_dir is None:
        return None
    return os.path.join(base_dir, *args)


ARCHIVE_BACKUP_ID_FILE = from_dir(ARCHIVE_DIR, "backup.txt")
JOURNAL_BACKUP_ID_FILE = from_dir(JOURNAL_DIR, "backup.txt")
BACKUP_PAYLOAD_FILE = from_dir(ARCHIVE_DIR, "backup_payload.txt")


def write_file(path, content):
    # Write verbatim content if string or bytes, otherwise serialize as JSON
    mode = "w"
    if isinstance(content, (bytes, bytearray)):
        mode += "b"
    with open(path, mode) as f:
        if isinstance(content, (str, bytes, bytearray)):
            f.write(content)
        else:
            json.dump(content, f)

    # Make sure that file is accessible by nuodb user, which has uid 1000 by
    # default. Making the file group-writable ensures that it is accessible to
    # the nuodb user in OpenShift deployments where an arbitrary uid is used
    # with gid 0.
    if os.getuid() == 0:
        os.chown(path, 1000, 0)
        os.chmod(path, mode=0o660)


def get_nuodb_pids():
    # Get the pid for the nuodb process, which should be sharing a namespace
    # with the backup-hooks sidecar
    pids = []
    for pid in os.listdir("/proc"):
        # Skip any non-pid directories
        if not pid.isdigit():
            continue
        try:
            # Read file containing command-line arguments
            with open(f"/proc/{pid}/cmdline", "rb") as f:
                # Arguments are delimited by byte 0, and nuodb rewrites
                # command-line arguments to include state information on
                # argv[0] with entries delimited by spaces
                args = f.read().replace(b"\x00", b" ").split()
                if args and args[0] in [b"nuodb", b"/opt/nuodb/bin/nuodb"]:
                    pids.append(pid)
        except FileNotFoundError:
            # Process may have exited
            pass

    if len(pids) > 1:
        LOGGER.warning("Multiple nuodb processes found: %s", pids)
    return pids


def check_thread_suspended(pid, tid):
    try:
        # Read status file
        with open(f"/proc/{pid}/task/{tid}/status", "r") as f:
            # Look for line with prefix 'Status:'
            prefix = "State:"
            for line in f:
                if line.startswith(prefix):
                    state = line[len(prefix) :].strip()
                    if state != "T (stopped)":
                        LOGGER.info("Thread %s not stopped: %s", tid, state)
                        return False
                    return True
    except FileNotFoundError as e:
        LOGGER.warning("Unable to get state of thread %s: %s", tid, e)


def check_suspended(pid):
    # Check if any threads are not in suspended (stopped) state
    for tid in os.listdir(f"/proc/{pid}/task"):
        if tid.isdigit() and not check_thread_suspended(pid, tid):
            return False
    return True


def await_suspended(pid, interval=0.25, retries=16):
    for i in range(0, retries):
        if i == 0:
            LOGGER.info("Waiting for all threads of process %s to be suspended", pid)
        else:
            LOGGER.info(
                "Retrying in %s second(s) (iteration %s/%s)", interval, i + 1, retries
            )
            time.sleep(interval)

        # If all threads are suspended, return immediately
        if check_suspended(pid):
            return

    raise RuntimeError(
        "Process {} not suspended after {} seconds".format(pid, retries * interval)
    )


def freeze_archive(unfreeze=False):
    if USE_SUSPEND:
        # Resume or suspend nuodb process. There should be a unique nuodb
        # process, but if there are multiple somehow, suspend all of them.
        pids = get_nuodb_pids()
        if not pids:
            raise RuntimeError("No nuodb process found")
        for pid in pids:
            if unfreeze:
                LOGGER.info("Resuming nuodb process with pid %s", pid)
                os.kill(int(pid), signal.SIGCONT)
            else:
                LOGGER.info("Suspending nuodb process with pid %s", pid)
                os.kill(int(pid), signal.SIGSTOP)
                # Make sure all threads are suspended
                await_suspended(pid)
    else:
        # Freeze or unfreeze the archive filesystem using fsfreeze
        if unfreeze:
            LOGGER.info("Unfreezing writes to archive volume")
            action = "--unfreeze"
        else:
            LOGGER.info("Freezing writes to archive volume")
            action = "--freeze"
        subprocess.check_output(
            ["fsfreeze", action, ARCHIVE_DIR], stderr=subprocess.STDOUT
        )


def pre_backup(backup_id, payload):
    # Check that backup ID was specified
    if not backup_id:
        raise UserError("Backup ID not specified")

    # Make sure we do not attempt to execute pre-backup hook if a backup is
    # already in progress, which may block indefinitely if writes to the
    # archive and journal directories are frozen
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        raise UserError(
            "Backup ID file {} exists. Execute post-backup hook to complete current backup.".format(
                ARCHIVE_BACKUP_ID_FILE
            )
        )

    # Write backup ID to archive directory
    write_file(ARCHIVE_BACKUP_ID_FILE, backup_id)

    # Write backup ID to journal directory if it is separate from archive
    if JOURNAL_DIR is not None:
        write_file(JOURNAL_BACKUP_ID_FILE, backup_id)

    # Write opaque data in request payload that should be backed up
    # (snapshotted) with archive
    if payload is not None and payload.get("opaque"):
        write_file(BACKUP_PAYLOAD_FILE, payload["opaque"])
    elif os.path.exists(BACKUP_PAYLOAD_FILE):
        # Remove opaque data from last backup
        os.remove(BACKUP_PAYLOAD_FILE)

    # If the archive and journal are on separate volumes, block writes to the
    # archive to allow consistent snapshots to be obtained of both
    if JOURNAL_DIR is not None:
        freeze_archive()


def post_backup(backup_id, query):
    # Check that backup ID was specified
    if not backup_id:
        raise UserError("Backup ID not specified")

    # Check backup ID to make sure it matches current backup
    force = query is not None and query.get("force") == True
    current_backup_id = get_backup_id()
    if current_backup_id != backup_id:
        msg = "Unexpected backup ID: current={}, supplied={}".format(
            current_backup_id, backup_id
        )
        # If force=true, then do not fail on mismatching backup ID
        if force:
            LOGGER.warning(msg)
        else:
            raise UserError(msg)

    # If the archive and journal are on separate volumes, we must unblock
    # writes to the archive
    if JOURNAL_DIR is not None:
        freeze_archive(unfreeze=True)

    # Delete backup ID files and payload file
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        os.remove(ARCHIVE_BACKUP_ID_FILE)
    if JOURNAL_BACKUP_ID_FILE is not None and os.path.exists(JOURNAL_BACKUP_ID_FILE):
        os.remove(JOURNAL_BACKUP_ID_FILE)
    if os.path.exists(BACKUP_PAYLOAD_FILE):
        os.remove(BACKUP_PAYLOAD_FILE)


def get_backup_id():
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        with open(ARCHIVE_BACKUP_ID_FILE, "r") as f:
            return f.read().strip()


class HttpError(RuntimeError):
    def __init__(self, status, message):
        super().__init__(message)
        self.status = status
        self.message = message


class UserError(HttpError):
    def __init__(self, message):
        super().__init__(http.HTTPStatus.BAD_REQUEST, message)


def log_error(exc):
    if isinstance(exc, subprocess.CalledProcessError):
        stdout = exc.stdout.decode("utf-8")
        LOGGER.exception(
            "Failed command while handling request\n---\n%s\n---", stdout.strip()
        )
    else:
        LOGGER.exception("Unexpected error while handling request")


def create_error(exc):
    # HttpError specifies status and message explicitly
    if isinstance(exc, HttpError):
        return exc.status, dict(success=False, message=exc.message)

    # Some unexpected error occured. Respond with '500 Internal Server Error'.
    log_error(exc)
    return http.HTTPStatus.INTERNAL_SERVER_ERROR, dict(success=False, message=str(exc))


REGISTERED_HANDLERS = [
    ("POST", "pre-backup", pre_backup),
    ("POST", "post-backup", post_backup),
]


def get_payload(environ):
    # Get content length
    content_length = environ.get("CONTENT_LENGTH")
    if content_length:
        content_length = int(content_length)

    # If content length is non-0, read payload and decode as JSON
    if content_length:
        content = environ.get("wsgi.input").read(content_length).decode("utf-8")
        try:
            return json.loads(content)
        except json.decoder.JSONDecodeError as e:
            raise UserError("Unable to decode request payload: " + str(e))


def get_query_params(query_str):
    as_multimap = urllib.parse.parse_qs(query_str)
    as_map = {}
    for key, values in as_multimap.items():
        if len(values) == 0:
            continue
        if len(values) > 1:
            raise UserError(
                "Multiple values specified for query parameter '{}'".format(key)
            )
        # Try to decode value as some token other than string, otherwise fall
        # back to string
        try:
            as_map[key] = json.loads(values[0])
        except json.decoder.JSONDecodeError:
            as_map[key] = values[0]
    return as_map


def dispatch_request(environ):
    # Parse URL and make sure it contains at least one path component
    uri = wsgiref.util.request_uri(environ)
    parsed = urllib.parse.urlparse(uri)
    path = parsed.path
    if path.startswith("/"):
        path = path[1:]
    components = path.split("/")
    if not components:
        raise UserError("No path provided")

    # Find a handler that matches the request
    req_method = environ.get("REQUEST_METHOD")
    found_handler = False
    for method, path_prefix, handler in REGISTERED_HANDLERS:
        if req_method == method and components[0] == path_prefix:
            # Make sure the correct number of parameters were supplied by
            # inspecting method signature
            args = components[1:]
            pos_args = inspect.getfullargspec(handler).args
            # If `payload` is in the method signature, then read the request
            # payload and pass it at the corresponding index
            if "payload" in pos_args:
                args.insert(pos_args.index("payload"), get_payload(environ))
            # If `query` is in the method signature, then parse the query
            # parameters as a dictionary and pass them at the corresponding
            # index
            if "query" in pos_args:
                args.insert(pos_args.index("query"), get_query_params(parsed.query))
            if len(args) != len(pos_args):
                msg = "{} parameter(s) expected but {} supplied in request {}".format(
                    len(pos_args), len(args), parsed.path
                )
                if "payload" in pos_args or "query" in pos_args:
                    msg += ", including payload and query parameters"
                raise UserError(msg)
            # Request is valid. Send it to handler.
            handler(*args)
            found_handler = True
            break

    # Make sure a handler was found
    if not found_handler:
        raise UserError("No handler found for path " + parsed.path)


def create_response(start_response, status, json_response):
    status_str = "{0.value} {0.phrase}".format(status)
    headers = [("Content-Type", "application/json")]
    start_response(status_str, headers)
    return [json.dumps(json_response, indent=4).encode("utf-8"), b"\n"]


def hooks_handler(environ, start_response):
    status = http.HTTPStatus.OK
    json_response = dict(success=True)
    try:
        dispatch_request(environ)
    except Exception as e:
        status, json_response = create_error(e)
    return create_response(start_response, status, json_response)


def start_server(port):
    with simple_server.make_server("", port, hooks_handler) as httpd:
        LOGGER.info("Starting backup hooks server on port %s", port)
        httpd.serve_forever()


if __name__ == "__main__":
    # Create CLI parser for direct invocation
    parser = argparse.ArgumentParser(sys.argv[0])
    subparsers = parser.add_subparsers(dest="subcommand")
    # Register server subcommand
    subparser = subparsers.add_parser("server")
    subparser.add_argument("--port", type=int, default=80)
    # Register pre-hook subcommand
    subparser = subparsers.add_parser("pre-hook")
    subparser.add_argument("--backup-id", required=True)
    subparser.add_argument("--opaque-file", type=argparse.FileType("r"))
    # Register post-hook subcommand
    subparser = subparsers.add_parser("post-hook")
    subparser.add_argument("--backup-id", required=True)
    subparser.add_argument("--force", action="store_true")

    # Parse arguments and invoke correct handler
    args = parser.parse_args(sys.argv[1:])
    if args.subcommand == "server":
        start_server(args.port)
    if args.subcommand == "pre-hook":
        # read opaque data and pass it to pre-hook
        opaque = None
        if args.opaque_file:
            opaque = args.opaque_file.read()
        pre_backup(args.backup_id, dict(opaque=opaque))
    elif args.subcommand == "post-hook":
        post_backup(args.backup_id, dict(force=args.force))
