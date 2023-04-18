package redirection

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedirectionCheckConfigurationList is a list of RedirectionCheckConfiguration objects
type RedirectionCheckConfigurationList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []RedirectionCheckConfiguration
}

// RedirectionCheckSpec is the specification of a RedirectionCheckConfiguration
type RedirectionCheckSpec struct {
	AllowedRedirectionHosts []string
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RedirectionCheckConfiguration
type RedirectionCheckConfiguration struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// Spec represents the specification of the desired behavior of RedirectionCheckConfiguration.
	Spec RedirectionCheckSpec
}

// ScopeType specifies the type of scope being used
type ScopeType string

const (
	// ClusterScope means that scope is limited to cluster-scoped objects.
	// Namespace objects are cluster-scoped.
	ClusterScope ScopeType = "Cluster"
	// NamespacedScope means that scope is limited to namespaced objects.
	NamespacedScope ScopeType = "Namespaced"
	// AllScopes means that all scopes are included.
	AllScopes ScopeType = "*"
)
