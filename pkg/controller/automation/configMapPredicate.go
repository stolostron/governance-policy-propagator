// Copyright Contributors to the Open Cluster Management project

package automation

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// we only want to watch for pb contains policy as subjects
var configMapPredicateFuncs = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		cfgObjNew := e.ObjectNew.(*corev1.ConfigMap)
		cfgObjOld := e.ObjectOld.(*corev1.ConfigMap)
		if cfgObjNew.Data["policyRef"] == "" {
			return false
		}
		if cfgObjNew.ObjectMeta.Annotations["policy.open-cluster-management.io/rerun"] == "true" {
			return true
		}
		return !equality.Semantic.DeepEqual(cfgObjNew.Data, cfgObjOld.Data)
	},
	CreateFunc: func(e event.CreateEvent) bool {
		cfgObjNew := e.Object.(*corev1.ConfigMap)
		return cfgObjNew.Data["policyRef"] != ""
	},
	DeleteFunc: func(e event.DeleteEvent) bool {
		return false
	},
}
