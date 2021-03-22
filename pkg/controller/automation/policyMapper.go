// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package automation

import (
	"context"

	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policies/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type policyMapper struct {
	client.Client
}

func (mapper *policyMapper) Map(obj handler.MapObject) []reconcile.Request {
	policy := obj.Object.(*policiesv1.Policy)
	var result []reconcile.Request
	cfgMapList := &corev1.ConfigMapList{}
	err := mapper.Client.List(context.TODO(), cfgMapList, &client.ListOptions{Namespace: policy.GetNamespace()})
	if err != nil {
		return nil
	}
	foundCfgMap := false
	cfgMap := corev1.ConfigMap{}
	for _, cfgMapTemp := range cfgMapList.Items {
		if cfgMapTemp.Data["policyRef"] == policy.GetName() {
			foundCfgMap = true
			cfgMap = cfgMapTemp
			break

		}
	}
	if foundCfgMap {
		if cfgMap.Data["mode"] == "scan" {
			// scan mode, do not queue
		} else if cfgMap.Data["mode"] == "once" {
			request := reconcile.Request{NamespacedName: types.NamespacedName{
				Name:      cfgMap.GetName(),
				Namespace: cfgMap.GetNamespace(),
			}}
			result = append(result, request)
		}
	}
	return result
}
