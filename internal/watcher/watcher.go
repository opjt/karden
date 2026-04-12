package watcher

import (
	"context"
	"encoding/json"
	"log/slog"

	"karden/internal/domain/audit"
	"karden/internal/domain/workload"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

var kardenSecretGVR = schema.GroupVersionResource{
	Group:    "karden.io",
	Version:  "v1",
	Resource: "kardensecrets",
}

type Watcher struct {
	dynClient dynamic.Interface
	store     workload.SecretStore
	repo      workload.Repository
	auditRepo audit.Repository
	stopCh    chan struct{}
}

func New(dynClient dynamic.Interface, store workload.SecretStore, repo workload.Repository, auditRepo audit.Repository) *Watcher {
	return &Watcher{
		dynClient: dynClient,
		store:     store,
		repo:      repo,
		auditRepo: auditRepo,
		stopCh:    make(chan struct{}),
	}
}

func (w *Watcher) Start() {
	factory := dynamicinformer.NewDynamicSharedInformerFactory(w.dynClient, 0)
	informer := factory.ForResource(kardenSecretGVR).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj any) {
			ks := toKardenSecret(obj)
			if ks != nil {
				w.handleAdd(ks)
			}
		},
		UpdateFunc: func(_, newObj any) {
			ks := toKardenSecret(newObj)
			if ks != nil {
				w.handleAdd(ks)
			}
		},
		DeleteFunc: func(obj any) {
			ks := toKardenSecret(obj)
			if ks == nil {
				// handle tombstone
				tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
				if !ok {
					return
				}
				ks = toKardenSecret(tombstone.Obj)
			}
			if ks != nil {
				w.handleDelete(ks)
			}
		},
	})

	factory.Start(w.stopCh)
	factory.WaitForCacheSync(w.stopCh)

	slog.Info("watcher started")
	<-w.stopCh
	slog.Info("watcher stopped")
}

func (w *Watcher) Stop() {
	close(w.stopCh)
}

func (w *Watcher) handleAdd(ks *kardenSecret) {
	slog.Info("detected KardenSecret",
		"name", ks.Name,
		"namespace", ks.Namespace,
		"type", ks.Spec.Type,
	)

	ctx := context.Background()
	wl := ks.toWorkload()

	id := w.upsertWorkload(ctx, wl)
	if id > 0 {
		w.ensureSecret(ctx, wl, id)
	}
}

func (w *Watcher) handleDelete(ks *kardenSecret) {
	slog.Info("KardenSecret deleted",
		"name", ks.Name,
		"namespace", ks.Namespace,
	)

	ctx := context.Background()
	if err := w.repo.SetInactive(ctx, ks.Name, ks.Namespace); err != nil {
		slog.Error("failed to mark workload inactive",
			"name", ks.Name,
			"namespace", ks.Namespace,
			"err", err,
		)
	}
}

// upsertWorkload persists the workload and returns its DB id (0 on error).
func (w *Watcher) upsertWorkload(ctx context.Context, wl *workload.ManagedWorkload) int64 {
	id, err := w.repo.Upsert(ctx, wl)
	if err != nil {
		slog.Error("failed to upsert workload",
			"name", wl.PodName,
			"namespace", wl.Namespace,
			"err", err,
		)
		return 0
	}
	slog.Info("workload upserted", "name", wl.PodName, "namespace", wl.Namespace)
	return id
}

// ensureSecret creates the K8s Secret if it doesn't exist yet, then writes an audit log.
func (w *Watcher) ensureSecret(ctx context.Context, wl *workload.ManagedWorkload, workloadID int64) {
	existing, err := w.store.GetData(ctx, wl.Namespace, wl.SecretName)
	if err == nil && len(existing) > 0 {
		slog.Info("secret already exists, skipping", "secret", wl.SecretName)
		return
	}

	data := buildSecretData(wl)
	if err := w.store.Set(ctx, wl.Namespace, wl.SecretName, data); err != nil {
		slog.Error("failed to create secret", "secret", wl.SecretName, "err", err)
		_ = w.auditRepo.Save(ctx, &audit.AuditLog{
			TargetID: int(workloadID),
			Action:   audit.ActionCreate,
			Actor:    "karden",
			Result:   audit.ResultFailure,
			Reason:   err.Error(),
		})
		return
	}

	slog.Info("secret created", "secret", wl.SecretName, "namespace", wl.Namespace)
	_ = w.auditRepo.Save(ctx, &audit.AuditLog{
		TargetID: int(workloadID),
		Action:   audit.ActionCreate,
		Actor:    "karden",
		Result:   audit.ResultSuccess,
	})
}

// toKardenSecret converts an unstructured K8s object to a KardenSecret.
func toKardenSecret(obj any) *kardenSecret {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return nil
	}

	specRaw, _, _ := unstructured.NestedMap(u.Object, "spec")
	specBytes, err := json.Marshal(specRaw)
	if err != nil {
		return nil
	}

	var spec kardenSecretSpec
	if err := json.Unmarshal(specBytes, &spec); err != nil {
		return nil
	}

	return &kardenSecret{
		Name:      u.GetName(),
		Namespace: u.GetNamespace(),
		Spec:      spec,
	}
}
