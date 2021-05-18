// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package propagator

import (
	"fmt"
	"strings"

	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type policyMapper struct {
	client.Client
	bUpdater *batchUpdater
}

func (mapper *policyMapper) Map(obj handler.MapObject) []reconcile.Request {
	return getOwnerReconcileRequest(mapper.bUpdater, obj.Meta)
}

// getOwnerReconcileRequest looks at object and returns a slice of reconcile.Request to reconcile
// owners of object from label: policy.open-cluster-management.io/root-policy
func getOwnerReconcileRequest(bu *batchUpdater, object metav1.Object) []reconcile.Request {
#	log.Info("izhang enter getOwnerReconcileRequest()")

	var req reconcile.Request

#	defer func() {
#		log.Info(fmt.Sprintf("izhang exit getOwnerReconcileRequest(), request: %s", req))
#	}()

	var result []reconcile.Request
	rootPlcName := object.GetLabels()[common.RootPolicyLabel]
	var name string
	var namespace string

	if rootPlcName != "" {
		// policy.open-cluster-management.io/root-policy exists, should be a replicated policy
		log.Info("Found reconciliation request from replicated policy...", "Namespace", object.GetNamespace(),
			"Name", object.GetName())
		name = strings.Split(rootPlcName, ".")[1]
		namespace = strings.Split(rootPlcName, ".")[0]

		req = reconcile.Request{NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		}}

		// bu.add(req)

		return []reconcile.Request{}
	} else {
		// policy.open-cluster-management.io/root-policy doesn't exist, should be a root policy
		log.Info("Found reconciliation request from root policy...", "Namespace", object.GetNamespace(),
			"Name", object.GetName())
		name = object.GetName()
		namespace = object.GetNamespace()
	}
	req = reconcile.Request{NamespacedName: types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}}

	result = append(result, req)
	return result
}
