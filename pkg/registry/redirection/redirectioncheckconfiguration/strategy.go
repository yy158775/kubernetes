package redirectioncheckconfiguration

import (
	"context"
	genericvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/extendedapis/redirection"
	"reflect"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// RedirectionStrategy implements verification logic for redirectionCheckConfiguration..
type RedirectionStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating redirectionCheckConfiguration
// objects via the REST API.
var Strategy = RedirectionStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

func (s RedirectionStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	return map[fieldpath.APIVersion]*fieldpath.Set{
		"redirection.k8s.io/v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}
}

func (s RedirectionStrategy) NamespaceScoped() bool {
	return false
}

// PrepareForCreate clears the status of an validatingWebhookConfiguration before creation.
func (s RedirectionStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	vc := obj.(*redirection.RedirectionCheckConfiguration)
	vc.Generation = 1
}

func (s RedirectionStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	vc := obj.(*redirection.RedirectionCheckConfiguration)
	return genericvalidation.ValidateObjectMeta(&vc.ObjectMeta, false, genericvalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
}

func (s RedirectionStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (s RedirectionStrategy) Canonicalize(obj runtime.Object) {
}

func (s RedirectionStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (s RedirectionStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newVC := obj.(*redirection.RedirectionCheckConfiguration)
	oldVC := old.(*redirection.RedirectionCheckConfiguration)

	if !reflect.DeepEqual(oldVC.Spec.AllowedRedirectionHosts, newVC.Spec.AllowedRedirectionHosts) {
		newVC.Generation = oldVC.Generation + 1
	}
}

func (s RedirectionStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	vc := obj.(*redirection.RedirectionCheckConfiguration)
	return genericvalidation.ValidateObjectMeta(&vc.ObjectMeta, false, genericvalidation.NameIsDNSSubdomain, field.NewPath("metadata"))
}

func (s RedirectionStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (s RedirectionStrategy) AllowUnconditionalUpdate() bool {
	return false
}
