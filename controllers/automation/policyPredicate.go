// Copyright Contributors to the Open Cluster Management project

package automation

import (
	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policy/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// we only want to watch for pb contains policy as subjects
var policyPredicateFuncs = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		plcObjNew := e.ObjectNew.(*policiesv1.Policy)
		if _, ok := plcObjNew.Labels["policy.open-cluster-management.io/root-policy"]; ok {
			return false
		}
		plcObjOld := e.ObjectOld.(*policiesv1.Policy)
		return plcObjNew.Status.ComplianceState != plcObjOld.Status.ComplianceState
	},
	CreateFunc: func(e event.CreateEvent) bool {
		return false
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return false
	},
}
