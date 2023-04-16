package storage

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/kubernetes/pkg/extendedapis/redirection"
	"k8s.io/kubernetes/pkg/printers"
	printersinternal "k8s.io/kubernetes/pkg/printers/internalversion"
	printerstorage "k8s.io/kubernetes/pkg/printers/storage"
	"k8s.io/kubernetes/pkg/registry/redirection/redirectioncheckconfiguration"
)

// REST implements a RESTStorage for mutatingWebhookConfiguration against etcd
type REST struct {
	*genericregistry.Store
}

// NewREST returns a RESTStorage object that will work against mutatingWebhookConfiguration.
func NewREST(optsGetter generic.RESTOptionsGetter) (*REST, error) {
	store := &genericregistry.Store{
		NewFunc:     func() runtime.Object { return &redirection.RedirectionCheckConfiguration{} },
		NewListFunc: func() runtime.Object { return &redirection.RedirectionCheckConfigurationList{} },
		ObjectNameFunc: func(obj runtime.Object) (string, error) {
			return obj.(*redirection.RedirectionCheckConfiguration).Name, nil
		},
		DefaultQualifiedResource:  redirection.Resource("redirectioncheckconfigurations"),
		SingularQualifiedResource: redirection.Resource("redirectioncheckconfiguration"),

		CreateStrategy: redirectioncheckconfiguration.Strategy,
		UpdateStrategy: redirectioncheckconfiguration.Strategy,
		DeleteStrategy: redirectioncheckconfiguration.Strategy,

		TableConvertor: printerstorage.TableConvertor{TableGenerator: printers.NewTableGenerator().With(printersinternal.AddHandlers)},
	}

	options := &generic.StoreOptions{RESTOptions: optsGetter}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}

// Implement CategoriesProvider
var _ rest.CategoriesProvider = &REST{}

// Categories implements the CategoriesProvider interface. Returns a list of categories a resource is part of.
func (r *REST) Categories() []string {
	return []string{"all"}
}
