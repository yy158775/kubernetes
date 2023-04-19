package install

import (
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/extendedapis/redirection"
	redirectionv1 "k8s.io/kubernetes/pkg/extendedapis/redirection/v1"
)

func init() {
	Install(legacyscheme.Scheme)
}

func Install(scheme *runtime.Scheme) {
	utilruntime.Must(redirection.AddToScheme(scheme))
	utilruntime.Must(redirectionv1.AddToScheme(scheme))
	utilruntime.Must(scheme.SetVersionPriority(redirectionv1.SchemeGroupVersion))
}