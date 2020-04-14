package testlib

const CA_CERT_FILE = "ca.cert"
const CA_CERT_FILE_NEW = "ca_new.cert"
const CA_CERT_SECRET = "nuodb-ca-cert"

const KEYSTORE_FILE = "nuoadmin.p12"
const KEYSTORE_SECRET = "nuodb-keystore"

const TRUSTSTORE_FILE = "nuoadmin-truststore.p12"
const TRUSTSTORE_SECRET = "nuodb-truststore"

const NUOCMD_FILE = "nuocmd.pem"
const NUOCMD_SECRET = "nuodb-client-pem"

const SECRET_PASSWORD = "changeIt"

const CERTIFICATES_GENERATION_PATH = "/tmp/keys"
const CERTIFICATES_BACKUP_PATH = "/tmp/keys_backup"

const TEARDOWN_ADMIN = "admin"
const TEARDOWN_BACKUP = "backup"
const TEARDOWN_DATABASE = "database"
const TEARDOWN_RESTORE = "database"
const TEARDOWN_SECRETS = "secrets"
const TEARDOWN_YCSB = "ycsb"

const ADMIN_HELM_CHART_PATH = "../../stable/admin"
const DATABASE_HELM_CHART_PATH = "../../stable/database"
const RESTORE_HELM_CHART_PATH = "../../stable/restore"
const THP_HELM_CHART_PATH = "../../stable/transparent-hugepage"
const YCSB_HELM_CHART_PATH = "../../incubator/demo-ycsb"

const RESULT_DIR = "../../results"
const INJECT_FILE = "../../versionInject.yaml"

const IMPORT_ARCHIVE_URL = "https://download.nuohub.org/ce_releases/restore.bak.tz"

// suffix "m" for spec.containers[].resources.requests.cpu denotes "millicores",
// and 1 CPU is equivalent to 1000m
const MINIMAL_VIABLE_ENGINE_CPU = "500m"
const MINIMAL_VIABLE_ENGINE_MEMORY = "500Mi"

const K8S_EVENT_LOG_FILE = "kubernetes_event.log"
