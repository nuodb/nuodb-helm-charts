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
from threading import Timer
from shutil import which

logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(message)s")
LOGGER = logging.getLogger(__name__)

ARCHIVE_DIR = os.environ.get("NUODB_ARCHIVE_DIR", "/mnt/archive")
JOURNAL_DIR = os.environ.get("NUODB_JOURNAL_DIR")
FREEZE_MODE = os.environ.get("FREEZE_MODE")
FREEZE_TIMEOUT = os.environ.get("FREEZE_TIMEOUT")

MODE_HOTSNAP = "hotsnap"
MODE_FSFREEZE = "fsfreeze"
MODE_SUSPEND = "suspend"


def from_dir(base_dir, *args):
    if base_dir is None:
        return None
    return os.path.join(base_dir, *args)


ARCHIVE_BACKUP_ID_FILE = from_dir(ARCHIVE_DIR, "backup.txt")
JOURNAL_BACKUP_ID_FILE = from_dir(JOURNAL_DIR, "backup.txt")
BACKUP_PAYLOAD_FILE = from_dir(ARCHIVE_DIR, "backup_payload.txt")
BACKUP_START_ID_FILE = from_dir(ARCHIVE_DIR, "backup_sid.txt")
RESTORED_FILE = from_dir(ARCHIVE_DIR, "restored.txt")


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

        # Flush buffers and invoke fsync() to make sure data is written to disk
        f.flush()
        os.fsync(f.fileno())

    # Make sure that file is accessible by nuodb user, which has uid 1000 by
    # default. Making the file group-writable ensures that it is accessible to
    # the nuodb user in OpenShift deployments where an arbitrary uid is used
    # with gid 0.
    if os.getuid() == 0:
        os.chown(path, 1000, 0)
        os.chmod(path, mode=0o660)


def get_nuodb_process_info():
    # Get the process info for the nuodb process, which should be sharing a
    # namespace with the backup-hooks sidecar
    processes = []
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
                    sid = None
                    for arg in args:
                        parts = arg.split(b":")
                        if len(parts) == 2 and parts[0] == b"sid":
                            sid = parts[1]
                    if not sid:
                        raise RuntimeError(
                            "Unable to find start ID for nuodb process PID " + pid
                        )
                    processes.append({"pid": pid, "sid": sid.decode("utf-8")})
        except FileNotFoundError:
            # Process may have exited
            pass

    if len(processes) > 1:
        LOGGER.warning("Multiple nuodb processes found: %s", processes)
    return processes


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


def freeze_archive(backup_id, processes, unfreeze=False, timeout=None):
    if FREEZE_MODE == MODE_SUSPEND:
        # Resume or suspend nuodb process. There should be a unique nuodb
        # process, but if there are multiple somehow, suspend all of them.
        for process in processes:
            pid = process["pid"]
            if unfreeze:
                LOGGER.info("Resuming nuodb process with pid %s", pid)
                os.kill(int(pid), signal.SIGCONT)
            else:
                LOGGER.info("Suspending nuodb process with pid %s", pid)
                os.kill(int(pid), signal.SIGSTOP)
                # Make sure all threads are suspended
                await_suspended(pid)
    elif FREEZE_MODE == MODE_HOTSNAP:
        # Freeze or unfreeze the archive using hotsnap
        sid = processes[0]["sid"]
        extra_args = []
        if unfreeze:
            LOGGER.info("Resuming archiving for nuodb process with start ID %s", sid)
            action = "resume"
        else:
            LOGGER.info("Pausing archiving for nuodb process with start ID %s", sid)
            action = "pause"
            if timeout and timeout > 0:
                extra_args += ["--timeout", "{}s".format(timeout)]
        try:
            subprocess.check_output(
                [
                    "nuocmd",
                    action,
                    "archiving",
                    "--start-id",
                    sid,
                    "--pause-id",
                    backup_id,
                ]
                + extra_args,
                stderr=subprocess.STDOUT,
            )
        except subprocess.CalledProcessError as e:
            # HotSnap supports timeout and will automatic resume. Check the
            # error message to find out if the automatic resume kicked in
            if action == "resume" and b"Archiving not paused" in e.output:
                raise ArchivingNotPausedError(e) from e
            raise
    elif FREEZE_MODE == MODE_FSFREEZE:
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
    else:
        raise RuntimeError("Unsupported freeze mode '{}'".format(FREEZE_MODE))


def pre_backup(backup_id, payload):
    # Check that backup ID was specified
    if not backup_id:
        raise UserError("Backup ID not specified")

    # Get nuodb processes
    processes = get_nuodb_process_info()
    if not processes:
        raise RuntimeError("No nuodb process found")

    # Make sure we do not attempt to execute pre-backup hook if a backup is
    # already in progress, which may block indefinitely if writes to the
    # archive and journal directories are frozen
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        # Check if the SM process was restarted which will invalidate the
        # previous pre-backup operation using HotSnap.
        stored_start_id = get_start_id()
        if (
            FREEZE_MODE != MODE_FSFREEZE
            and stored_start_id
            and processes[0]["sid"] != stored_start_id
        ):
            # Remote the backup files from the previous pre-backup operation
            LOGGER.warning(
                "Unexpected start ID: current=%s, stored=%s. SM process restarted while executing backup ID %s.",
                processes[0]["sid"],
                stored_start_id,
                get_backup_id(),
            )
            remove_backup_files()
        else:
            raise UserError(
                "Backup ID file {} exists. Execute post-backup hook to complete current backup.".format(
                    ARCHIVE_BACKUP_ID_FILE
                )
            )
    # Get timeout for the operation if specified
    timeout = get_backup_timeout(payload)

    # Delete file that is used by restored database to signal that archive
    # preparation is complete. This may be present if this database was
    # restored from a backup and this is the first time that a backup has been
    # taken on it.
    if os.path.exists(RESTORED_FILE):
        os.remove(RESTORED_FILE)

    # Write the start_id of the SM to a file
    write_file(BACKUP_START_ID_FILE, processes[0]["sid"])

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
        freeze_archive(backup_id, processes, timeout=timeout)

    # Cancel backup operation after specified timeout
    schedule_cancellation(backup_id, timeout)


def schedule_cancellation(backup_id, timeout):
    if timeout and timeout > 0:

        def cancel_backup():
            try:
                current_backup_id = get_backup_id()
                if current_backup_id and current_backup_id == backup_id:
                    LOGGER.warning(
                        "Canceling backup with ID %s. Timeout after %ds.",
                        backup_id,
                        timeout,
                    )
                    post_backup(backup_id, query={})
            except ArchivingNotPausedError:
                # Suppress this error during backup cancellation
                pass

        Timer(timeout, cancel_backup).start()


def post_backup(backup_id, query):
    # Check that backup ID was specified
    if not backup_id:
        raise UserError("Backup ID not specified")

    # Check backup ID to make sure it matches current backup
    force = query is not None and query.get("force") in [True, "true"]
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

    # Get nuodb processes
    processes = get_nuodb_process_info()
    if not processes:
        raise RuntimeError("No nuodb process found")

    # If the archive and journal are on separate volumes, we must unblock
    # writes to the archive
    if JOURNAL_DIR is not None:
        try:
            freeze_archive(backup_id, processes, unfreeze=True)
        except ArchivingNotPausedError:
            # Delete backup metadata files to allow new backups
            remove_backup_files()
            raise

    # Delete backup metadata files
    remove_backup_files()


def remove_backup_files():
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        os.remove(ARCHIVE_BACKUP_ID_FILE)
    if JOURNAL_BACKUP_ID_FILE is not None and os.path.exists(JOURNAL_BACKUP_ID_FILE):
        os.remove(JOURNAL_BACKUP_ID_FILE)
    if os.path.exists(BACKUP_PAYLOAD_FILE):
        os.remove(BACKUP_PAYLOAD_FILE)
    if os.path.exists(BACKUP_START_ID_FILE):
        os.remove(BACKUP_START_ID_FILE)


def get_backup_id():
    if os.path.exists(ARCHIVE_BACKUP_ID_FILE):
        with open(ARCHIVE_BACKUP_ID_FILE, "r") as f:
            return f.read().strip()


def get_start_id():
    if os.path.exists(BACKUP_START_ID_FILE):
        with open(BACKUP_START_ID_FILE, "r") as f:
            return f.read().strip()


def get_backup_timeout(payload):
    if payload is not None and payload.get("timeout") is not None:
        try:
            return int(payload.get("timeout"))
        except ValueError as e:
            raise UserError("Invalid request timeout") from e
    if FREEZE_TIMEOUT:
        return int(FREEZE_TIMEOUT)


class ArchivingNotPausedError(subprocess.CalledProcessError):
    def __init__(self, e):
        super().__init__(e.returncode, e.cmd, e.output, e.stderr)

    def __str__(self):
        return self.output.decode("utf-8").strip()


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

    # If content length is non-0, read payload
    if content_length:
        return environ.get("wsgi.input").read(content_length).decode("utf-8")


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
        as_map[key] = values[0]
    return as_map


def normalize_path(path):
    return path[1:] if path.startswith("/") else path


class ScriptHandler(object):

    DEFAULT_STATUS_MAPPING = {
        "0": 200,
        "*": 400,
    }

    def __init__(self, **kwargs):
        self.method = kwargs["method"]
        self.path = kwargs["path"]
        self.script = kwargs["script"]
        self.status_mappings = dict(self.DEFAULT_STATUS_MAPPING)
        if "statusMappings" in kwargs:
            self.status_mappings.update(kwargs["statusMappings"])

    def matches(self, method, path_components):
        # Check that method matches
        if method != self.method:
            return False, None

        # Check that path matches and extract path parameters
        path_params = {}
        handler_components = normalize_path(self.path).split("/")
        for i in range(0, len(path_components)):
            # If we have more path components in request than in handler path,
            # then request does not match
            if i >= len(handler_components):
                return False, None
            # If handler path component is a parameter, extract value from
            # request path
            handler_component = handler_components[i]
            if handler_component.startswith("{") and handler_component.endswith("}"):
                path_params[handler_component[1:-1]] = path_components[i]
                continue
            # If handler path component is literal, check that it matches
            if path_components[i] != handler_component:
                return False, None

        # Check that there are no remaining handler path components
        return i + 1 == len(handler_components), path_params

    def execute(self, path_params, query_params, payload):
        # Build environment variables to supply for script invocation
        env = dict(os.environ)
        if query_params:
            env.update(query_params)
        if payload:
            env.update(payload=payload)
        env.update(path_params)
        # Run script and capture output. Combine both stdout and stderr into
        # one stream, which is supplied as the HTTP response.
        proc = subprocess.run(
            self.script,
            shell=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.STDOUT,
            env=env,
        )
        # Map exit code to HTTP status
        return self.map_status(proc), proc.stdout

    def map_status(self, proc):
        status = self.status_mappings.get(str(proc.returncode))
        if not status:
            status = self.status_mappings["*"]
        for http_status in http.HTTPStatus:
            if str(http_status.value) == str(status):
                return http_status
        raise RuntimeError(
            "No HTTP status for mapping from {} to {}".format(proc.returncode, status)
        )


def read_handler_config(handler_config):
    if not handler_config:
        return []
    if not os.path.exists(handler_config):
        LOGGER.info("Ignoring non-existent handler config %s", handler_config)
        return []
    try:
        handlers = []
        with open(handler_config) as f:
            config_data = json.loads(f.read())
            # Account for `handlers` property being absent or explicitly
            # `null`, since this is how Helm renders the config if no custom
            # handlers are specified
            handlers_list = config_data.get("handlers")
            if handlers_list:
                for handler in handlers_list:
                    handlers.append(ScriptHandler(**handler))
        return handlers
    except Exception:
        LOGGER.exception("Unable to process handler config %s", handler_config)
        return []


class RequestInfo(object):
    def __init__(self, environ):
        self.method = environ.get("REQUEST_METHOD")
        # Parse URL and make sure it contains at least one path component
        uri = wsgiref.util.request_uri(environ)
        self.parsed = urllib.parse.urlparse(uri)
        self.components = normalize_path(self.parsed.path).split("/")
        if not self.components:
            raise UserError("No path provided")
        # Decode query parameters and payload
        self.query_params = get_query_params(self.parsed.query)
        self.payload = get_payload(environ)


def handle_method(req):
    # Find a handler that matches the request
    for method, path_prefix, handler in REGISTERED_HANDLERS:
        if req.method == method and req.components[0] == path_prefix:
            # Make sure the correct number of parameters were supplied by
            # inspecting method signature
            args = req.components[1:]
            pos_args = inspect.getfullargspec(handler).args
            # If `payload` is in the method signature, then read the request
            # payload as JSON and pass it at the corresponding index
            if "payload" in pos_args:
                try:
                    decoded_payload = json.loads(req.payload) if req.payload else None
                    args.insert(pos_args.index("payload"), decoded_payload)
                except json.JSONDecodeError as e:
                    raise UserError("Unable to decode request payload: " + str(e))
            # If `query` is in the method signature, then parse the query
            # parameters as a dictionary and pass them at the corresponding
            # index
            if "query" in pos_args:
                args.insert(pos_args.index("query"), req.query_params)
            if len(args) != len(pos_args):
                msg = "{} parameter(s) expected but {} supplied in request {}".format(
                    len(pos_args), len(args), req.parsed.path
                )
                if "payload" in pos_args or "query" in pos_args:
                    msg += ", including payload and query parameters"
                raise UserError(msg)
            # Request is valid. Send it to handler.
            handler(*args)
            return

    # Handler was not found
    raise UserError("No handler found for path " + req.parsed.path)


def create_response(start_response, status, data=None):
    status_str = "{0.value} {0.phrase}".format(status)
    headers = []
    if isinstance(data, bytes):
        # If data can be decoded as JSON, use application/json
        try:
            json.loads(data)
            headers.append(("Content-Type", "application/json"))
        except json.JSONDecodeError:
            headers.append(("Content-Type", "text/plain"))
    else:
        # Encode data as JSON
        data = json.dumps(data, indent=4).encode("utf-8") + b"\n"
        headers.append(("Content-Type", "application/json"))
    start_response(status_str, headers)
    return [data]


class HooksHandler(object):
    def __init__(self, handler_config=None):
        self.handlers = read_handler_config(handler_config)
        self.log_handlers()

    def log_handlers(self, indent=4):
        # Log all built-in handlers
        builtin_handlers = []
        for method, path_prefix, handler in REGISTERED_HANDLERS:
            path = path_prefix
            for arg in inspect.getfullargspec(handler).args:
                if arg not in ["payload", "query"]:
                    path += "/{" + arg + "}"
            builtin_handlers.append(
                "{}{} /{}".format(indent * " ", method, normalize_path(path))
            )
        LOGGER.info("Built-in handlers:\n%s", "\n".join(builtin_handlers))
        # Log all custom handlers, if there are any
        custom_handlers = []
        for handler in self.handlers:
            custom_handlers.append(
                "{}{} /{}".format(
                    indent * " ", handler.method, normalize_path(handler.path)
                )
            )
        if custom_handlers:
            LOGGER.info("Custom handlers:\n%s", "\n".join(custom_handlers))

    def handle(self, environ):
        try:
            # Dispatch to matching script handler
            req = RequestInfo(environ)
            for handler in self.handlers:
                matches, path_params = handler.matches(req.method, req.components)
                if matches:
                    return handler.execute(path_params, req.query_params, req.payload)
            # Dispatch to matching built-in handler
            handle_method(req)
            return http.HTTPStatus.OK, dict(success=True)
        except Exception as e:
            return create_error(e)

    def __call__(self, environ, start_response):
        status, json_response = self.handle(environ)
        return create_response(start_response, status, json_response)


def start_server(port, handler_config):
    hooks_handler = HooksHandler(handler_config)
    with simple_server.make_server("", port, hooks_handler) as httpd:
        LOGGER.info("Starting backup hooks server on port %s", port)
        httpd.serve_forever()


def verify_prerequisites():
    if FREEZE_MODE == MODE_HOTSNAP:
        if which("nuocmd") is None:
            raise RuntimeError("'nuocmd' command not found")
        try:
            subprocess.run(
                ["nuocmd", "pause", "archiving", "-h"],
                check=True,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
            )
        except subprocess.CalledProcessError as e:
            if e.returncode == 2:
                raise RuntimeError(
                    "'nuocmd pause archiving' command not supported"
                ) from e
            raise RuntimeError(
                "'nuocmd pause archiving' command failed: " + e.output("utf-8")
            ) from e
    elif FREEZE_MODE == MODE_FSFREEZE and which("fsfreeze") is None:
        raise RuntimeError("'fsfreeze' command not found")


if __name__ == "__main__":
    # Create CLI parser for direct invocation
    parser = argparse.ArgumentParser(sys.argv[0])
    subparsers = parser.add_subparsers(dest="subcommand")
    # Register server subcommand
    subparser = subparsers.add_parser("server")
    subparser.add_argument("--port", type=int, default=80)
    subparser.add_argument("--handler-config", default="/etc/nuodb/handlers.json")
    # Register pre-hook subcommand
    subparser = subparsers.add_parser("pre-hook")
    subparser.add_argument("--backup-id", required=True)
    subparser.add_argument("--opaque-file", type=argparse.FileType("r"))
    subparser.add_argument("--timeout", default=0)
    # Register post-hook subcommand
    subparser = subparsers.add_parser("post-hook")
    subparser.add_argument("--backup-id", required=True)
    subparser.add_argument("--force", action="store_true")

    # Parse arguments and invoke correct handler
    args = parser.parse_args(sys.argv[1:])
    if args.subcommand:
        verify_prerequisites()
    if args.subcommand == "server":
        start_server(args.port, args.handler_config)
    if args.subcommand == "pre-hook":
        # read opaque data and pass it to pre-hook
        opaque = None
        if args.opaque_file:
            opaque = args.opaque_file.read()
        pre_backup(args.backup_id, dict(opaque=opaque, timeout=args.timeout))
    elif args.subcommand == "post-hook":
        post_backup(args.backup_id, dict(force=args.force))
