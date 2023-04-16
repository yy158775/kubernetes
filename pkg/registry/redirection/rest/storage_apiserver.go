package rest

import (
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/extendedapis/redirection"
	redirectionv1 "k8s.io/kubernetes/pkg/extendedapis/redirection/v1"
	redirectioncheckconfigurationstorage "k8s.io/kubernetes/pkg/registry/redirection/redirectioncheckconfiguration/storage"
)

type RESTStorageProvider struct{}

func (p RESTStorageProvider) GroupName() string {
	return redirection.GroupName
}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, error) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(redirection.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if storageMap, err := p.v1(apiResourceConfigSource, restOptionsGetter); err != nil {
		return genericapiserver.APIGroupInfo{}, err
	} else if len(storageMap) > 0 {
		apiGroupInfo.VersionedResourcesStorageMap[redirection.SchemeGroupVersion.Version] = storageMap
	}

	return apiGroupInfo, nil
}

func (p RESTStorageProvider) v1(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (map[string]rest.Storage, error) {
	storage := map[string]rest.Storage{}

	// redirectioncheckconfigurations
	if resource := "redirectioncheckconfigurations"; apiResourceConfigSource.ResourceEnabled(redirectionv1.SchemeGroupVersion.WithResource(resource)) {
		redirectStorage, err := redirectioncheckconfigurationstorage.NewREST(restOptionsGetter)
		if err != nil {
			return storage, err
		}
		storage[resource] = redirectStorage
	}

	return storage, nil
}
