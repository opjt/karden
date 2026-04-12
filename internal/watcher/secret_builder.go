package watcher

import "karden/internal/domain/workload"

// buildSecretData generates initial secret key-value pairs based on the workload type.
func buildSecretData(wl *workload.ManagedWorkload) map[string]string {
	switch wl.Type {
	case workload.TypeDatabase:
		return buildDBSecretData(wl)
	default:
		return map[string]string{}
	}
}

func buildDBSecretData(wl *workload.ManagedWorkload) map[string]string {
	username := buildUsername(wl.SecretName)
	pw := generatePassword()

	switch wl.DBType {
	case workload.DBTypePostgres:
		return map[string]string{
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": pw,
			"POSTGRES_DB":       "app",
		}
	case workload.DBTypeMySQL:
		return map[string]string{
			"MYSQL_USER":          username,
			"MYSQL_PASSWORD":      pw,
			"MYSQL_ROOT_PASSWORD": generatePassword(),
		}
	default:
		return map[string]string{
			"USERNAME": username,
			"PASSWORD": pw,
		}
	}
}
