package watcher

import "karden/internal/domain/workload"

// kardenSecretSpec mirrors the CRD spec fields.
type kardenSecretSpec struct {
	Type         workload.Type   `json:"type"`
	DBType       workload.DBType `json:"dbType,omitempty"`
	DBHost       string          `json:"dbHost,omitempty"`
	DBPort       int             `json:"dbPort,omitempty"`
	RotationDays int             `json:"rotationDays,omitempty"`
}

// kardenSecret is the watcher-internal representation of the CRD resource.
type kardenSecret struct {
	Name      string
	Namespace string
	Spec      kardenSecretSpec
}

// toWorkload converts a kardenSecret into a ManagedWorkload for SQLite persistence.
func (ks *kardenSecret) toWorkload() *workload.ManagedWorkload {
	port := ks.Spec.DBPort
	if port == 0 {
		port = defaultDBPort(ks.Spec.DBType)
	}
	days := ks.Spec.RotationDays
	if days == 0 {
		days = 30
	}
	return &workload.ManagedWorkload{
		PodName:      ks.Name,
		Namespace:    ks.Namespace,
		SecretName:   ks.Name,
		Type:         ks.Spec.Type,
		DBType:       ks.Spec.DBType,
		DBHost:       ks.Spec.DBHost,
		DBPort:       port,
		RotationDays: days,
		Status:       workload.StatusActive,
	}
}

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
