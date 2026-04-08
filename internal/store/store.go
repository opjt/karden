package store

import "context"

// SecretStore is the abstraction for storing and retrieving secret values.
// The default implementation uses K8s Secrets, but can be swapped for Vault etc.
type SecretStore interface {
	Get(ctx context.Context, namespace, name, key string) (string, error)
	Set(ctx context.Context, namespace, name string, data map[string]string) error
	Delete(ctx context.Context, namespace, name string) error
}
