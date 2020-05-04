package propagator

import (
	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policies/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type placementBindingMapper struct {
	client.Client
}

func (mapper *placementBindingMapper) Map(obj handler.MapObject) []reconcile.Request {
	object := obj.Object.(*policiesv1.PlacementBinding)
	log.Info("Found reconciliation request from placmenet binding...", "Namespace", object.GetNamespace(), "Name", object.GetName())

	var result []reconcile.Request
	request := reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      object.Spec.Subject.Name,
		Namespace: object.GetNamespace(),
	}}
	result = append(result, request)
	return result
}
