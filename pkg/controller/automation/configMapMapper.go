// Copyright Contributors to the Open Cluster Management project

package automation

import (
	"context"

	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policy/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type configMapMapper struct {
	client.Client
}

func (mapper *configMapMapper) Map(obj handler.MapObject) []reconcile.Request {
	cfgMap := obj.Object.(*corev1.ConfigMap)
	var result []reconcile.Request
	policyRef := cfgMap.Data["policyRef"]
	if policyRef != "" {
		policy := &policiesv1.Policy{}
		err := mapper.Client.Get(context.TODO(), types.NamespacedName{
			Name:      policyRef,
			Namespace: cfgMap.GetNamespace(),
		}, policy)
		if err == nil {
			log.Info("Found reconciliation request from config map ...",
				"Namespace", cfgMap.GetNamespace(), "Name", cfgMap.GetName(), "Policy-Name", policyRef)
			request := reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      cfgMap.GetName(),
				Namespace: cfgMap.GetNamespace(),
			}}
			result = append(result, request)
		} else {
			log.Info("Failed to retrieve policyRef from config map...ignoring it...",
				"Namespace", cfgMap.GetNamespace(), "Name", cfgMap.GetName(), "Policy-Name", policyRef)
		}
	}
	return result
}
