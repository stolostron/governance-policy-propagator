// Copyright (c) 2021 Red Hat, Inc.

package automation

import (
	"context"
	"fmt"
	"time"

	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policies/v1"
	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName string = "policy-automation"

var log = logf.Log.WithName(controllerName)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Policy Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	dyamicClient, err := dynamic.NewForConfig(mgr.GetConfig())
	if err != nil {
		panic(err)
	}
	return &ReconcilePolicy{client: mgr.GetClient(), scheme: mgr.GetScheme(), dyamicClient: dyamicClient,
		recorder: mgr.GetEventRecorderFor(controllerName)}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Policy
	err = c.Watch(&source.Kind{Type: &policiesv1.Policy{}},
		&common.EnqueueRequestsFromMapFunc{ToRequests: &policyMapper{mgr.GetClient()}},
		policyPredicateFuncs)
	if err != nil {
		return err
	}

	// Watch for changes to config map
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}},
		&common.EnqueueRequestsFromMapFunc{ToRequests: &configMapMapper{mgr.GetClient()}},
		configMapPredicateFuncs)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcilePolicy implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcilePolicy{}

// ReconcilePolicy reconciles a Policy object
type ReconcilePolicy struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client       client.Client
	dyamicClient dynamic.Interface
	scheme       *runtime.Scheme
	recorder     record.EventRecorder
	counter      int
}

// Reconcile reads that state of the cluster for a Policy object and makes changes based on the state read
// and what is in the Policy.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcilePolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)

	// Fetch the Policy instance
	instance := &policiesv1.Policy{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected.

			// reqLogger.Info("Policy clean up complete, reconciliation completed.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	cfgMapList := &corev1.ConfigMapList{}
	err = r.client.List(context.TODO(), cfgMapList, &client.ListOptions{Namespace: request.Namespace})
	if err != nil {
		// failed to query configmap
		return reconcile.Result{}, err
	}
	if len(cfgMapList.Items) == 0 {
		// configmap not found could have been deleted
		// do nothing
		return reconcile.Result{}, nil
	}
	foundCfgMap := false
	cfgMap := corev1.ConfigMap{}
	for _, cfgMapTemp := range cfgMapList.Items {
		if cfgMapTemp.Data["policyRef"] == request.Name {
			foundCfgMap = true
			cfgMap = cfgMapTemp
			break

		}
	}
	if !foundCfgMap {
		// configmap not found could have been deleted
		// do nothing
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Found automation from configmap ...",
		"Namespace", cfgMap.GetNamespace(), "Name", cfgMap.GetName(), "Policy-Name", request.Name)
	if cfgMap.Annotations["policy.open-cluster-management.io/run-immediately"] == "true" {
		reqLogger.Info("Triggering manual run mode ...",
			"Namespace", cfgMap.GetNamespace(), "Name", cfgMap.GetName(), "Policy-Name", request.Name)
		err = common.CreateAnsibleJob(cfgMap, r.dyamicClient)
		if err != nil {
			return reconcile.Result{}, err
		}
		// manual run suceeded, remove annotation
		delete(cfgMap.Annotations, "policy.open-cluster-management.io/run-immediately")
		r.client.Update(context.TODO(), &cfgMap, &client.UpdateOptions{})
		reqLogger.Info("Manual run complelte...",
			"Namespace", cfgMap.GetNamespace(), "Name", cfgMap.GetName(), "Policy-Name", request.Name)
		return reconcile.Result{}, nil

	} else if cfgMap.Data["rescanAfter"] != "" {
		requeueAfter, _ := time.ParseDuration(cfgMap.Data["rescanAfter"])
		if instance.Spec.Disabled {
			reqLogger.Info("Policy is disabled, doing nothing...",
				"Namespace", request.Namespace, "Policy-Name", request.Name)
		} else {
			targetList := common.FindNonCompliantClustersForPolicy(instance)
			if len(targetList) > 0 {
				reqLogger.Info("Creating ansible job with targetList", "targetList", targetList)
				err = common.CreateAnsibleJob(cfgMap, r.dyamicClient)
				if err != nil {
					return reconcile.Result{RequeueAfter: requeueAfter}, err
				}

			} else {
				reqLogger.Info("No cluster is in noncompliant status, doing nothing...")
			}
		}

		// no violations found, doing nothing
		r.counter++
		reqLogger.Info("RequeueAfter.", "RequeueAfter", fmt.Sprintf("%d", requeueAfter), "counter", fmt.Sprintf("%d", r.counter))
		return reconcile.Result{RequeueAfter: requeueAfter}, nil

		// run continously mode

		// if requeueAfter == 0 {
		// 	reqLogger.Info("RequeueAfter.", "RequeueAfter", fmt.Sprintf("%d", requeueAfter), "counter", fmt.Sprintf("%d", r.counter))
		// 	return reconcile.Result{}, nil
		// }

	}
	return reconcile.Result{}, nil
}
