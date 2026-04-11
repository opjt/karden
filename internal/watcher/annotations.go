package watcher

import "karden/internal/domain/workload"

const (
	AnnotationInject       = "karden.io/inject"
	AnnotationType         = "karden.io/type"
	AnnotationSecretName   = "karden.io/secret-name"
	AnnotationDBType       = "karden.io/db-type"
	AnnotationDBHost       = "karden.io/db-host"
	AnnotationDBPort       = "karden.io/db-port"
	AnnotationRotationDays = "karden.io/rotation-days"
)

// defaultDBPort returns the well-known port for each DB type.
func defaultDBPort(dbType workload.DBType) int {
	switch dbType {
	case workload.DBTypePostgres:
		return 5432
	case workload.DBTypeMySQL:
		return 3306
	case workload.DBTypeMongoDB:
		return 27017
	case workload.DBTypeRedis:
		return 6379
	default:
		return 0
	}
}
