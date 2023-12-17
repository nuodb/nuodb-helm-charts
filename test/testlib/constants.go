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
const TEARDOWN_RESTORE = "restore"
const TEARDOWN_SECRETS = "secrets"
const TEARDOWN_YCSB = "ycsb"
const TEARDOWN_COLLECTOR = "nuocollector"
const TEARDOWN_VAULT = "vault"
const TEARDOWN_MULTICLUSTER = "multicluster"
const TEARDOWN_NGINX = "nginx"
const TEARDOWN_HAPROXY = "haproxy"
const TEARDOWN_CSIDRIVER_FS = "csidriver-fs"

const ADMIN_HELM_CHART_PATH = "../../stable/admin"
const DATABASE_HELM_CHART_PATH = "../../stable/database"
const RESTORE_HELM_CHART_PATH = "../../stable/restore"
const THP_HELM_CHART_PATH = "../../stable/transparent-hugepage"
const YCSB_HELM_CHART_PATH = "../../incubator/demo-ycsb"

const RESULT_DIR = "../../results"
const INJECT_FILE = "../../versionInject.yaml"
const UPGRADE_INJECT_FILE = "../../upgradeInject.yaml"
const INJECT_VALUES_FILE = "../../valuesInject.yaml"
const INJECT_CLUSTERS_FILE = "../../clustersInject.yaml"

const IMPORT_ARCHIVE_FILE = "../files/restore.bak.tz"
const IMPORT_SG0_BACKUP_FILE = "../files/backup-sg0.bak.tz"
const IMPORT_SG1_BACKUP_FILE = "../files/backup-sg1.bak.tz"

// suffix "m" for spec.containers[].resources.requests.cpu denotes "millicores",
// and 1 CPU is equivalent to 1000m
const MINIMAL_VIABLE_ENGINE_CPU = "500m"
const MINIMAL_VIABLE_ENGINE_MEMORY = "500Mi"

const SNAPSHOTABLE_STORAGE_CLASS = "csi-hostpath-sc"
const VOLUME_SNAPSHOT_CLASS = "csi-hostpath-snapclass"

const K8S_EVENT_LOG_FILE = "kubernetes_event.log"

const ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS = "ALWAYS_RUN_DIAGNOSTIC_TEARDOWNS"

const NGINX_DEPLOYMENT = "nginx"

// FQDNs used for Ingress testing that point to minikube IP address; they must
// be aligned with the /etc/hosts entries in the `install_deps.sh` script
const ADMIN_API_INGRESS_HOSTNAME = "api.nuodb.local"
const ADMIN_SQL_INGRESS_HOSTNAME = "sql.nuodb.local"
const DATABASE_TE_INGRESS_HOSTNAME = "demo.nuodb.local"

type LicenseType string

const (
	UNLICENSED LicenseType = "UNLICENSED"
	LIMITED    LicenseType = "LIMITED"
	ENTERPRISE LicenseType = "ENTERPRISE"
)
