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
	"k8s.io/apimachinery/pkg/types"
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

	// Fetch the ConfigMap instance
	cfgMap := &corev1.ConfigMap{}
	err := r.client.Get(context.TODO(), request.NamespacedName, cfgMap)
	if err != nil {
		if errors.IsNotFound(err) {
			reqLogger.Info("Automation was deleted, doing nothing...")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	reqLogger = log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "policyRef", cfgMap.Data["policyRef"])
	// cfgMapList := &corev1.ConfigMapList{}
	// err = r.client.List(context.TODO(), cfgMapList, &client.ListOptions{Namespace: request.Namespace})
	// if err != nil {
	// 	// failed to query configmap
	// 	return reconcile.Result{}, err
	// }
	// if len(cfgMapList.Items) == 0 {
	// 	// configmap not found could have been deleted
	// 	// do nothing
	// 	return reconcile.Result{}, nil
	// }
	// foundCfgMap := false
	// cfgMap := corev1.ConfigMap{}
	// for _, cfgMapTemp := range cfgMapList.Items {
	// 	if cfgMapTemp.Data["policyRef"] == request.Name {
	// 		foundCfgMap = true
	// 		cfgMap = cfgMapTemp
	// 		break

	// 	}
	// }
	// if !foundCfgMap {
	// 	// configmap not found could have been deleted
	// 	// do nothing
	// 	return reconcile.Result{}, nil
	// }

	reqLogger.Info("Handling automation...")
	if cfgMap.Annotations["policy.open-cluster-management.io/rerun"] == "true" {
		reqLogger.Info("Triggering manual run...")
		err = common.CreateAnsibleJob(cfgMap, r.dyamicClient)
		if err != nil {
			return reconcile.Result{}, err
		}
		// manual run suceeded, remove annotation
		delete(cfgMap.Annotations, "policy.open-cluster-management.io/rerun")
		r.client.Update(context.TODO(), cfgMap, &client.UpdateOptions{})
		reqLogger.Info("Manual run complete...")
		return reconcile.Result{}, nil
	} else if cfgMap.Data["mode"] == "disabled" {
		reqLogger.Info("Automation is disabled, doing nothing...")
		return reconcile.Result{}, nil
	} else {
		policy := &policiesv1.Policy{}
		err := r.client.Get(context.TODO(), types.NamespacedName{
			Name:      cfgMap.Data["policyRef"],
			Namespace: cfgMap.GetNamespace(),
		}, policy)
		if err != nil {
			if errors.IsNotFound(err) {
				//policy is gone, need to delete automation
				return reconcile.Result{}, nil
			}
			// Error reading the object - requeue the request.
			return reconcile.Result{}, err
		}
		if policy.Spec.Disabled {
			reqLogger.Info("Policy is disabled, doing nothing...")
			return reconcile.Result{}, nil
		}
		if cfgMap.Data["mode"] == "scan" {
			reqLogger.Info("Triggering scan mode...")
			requeueAfter, err := time.ParseDuration(cfgMap.Data["rescanAfter"])
			if err != nil {
				requeueAfter = 10 * time.Minute
			}

			targetList := common.FindNonCompliantClustersForPolicy(policy)
			if len(targetList) > 0 {
				reqLogger.Info("Creating ansible job with targetList", "targetList", targetList)
				err = common.CreateAnsibleJob(cfgMap, r.dyamicClient)
				if err != nil {
					return reconcile.Result{RequeueAfter: requeueAfter}, err
				}

			} else {
				reqLogger.Info("No cluster is in noncompliant status, doing nothing...")
			}

			// no violations found, doing nothing
			r.counter++
			reqLogger.Info("RequeueAfter.", "RequeueAfter", requeueAfter.String(), "counter", fmt.Sprintf("%d", r.counter))
			return reconcile.Result{RequeueAfter: requeueAfter}, nil
		} else if cfgMap.Data["mode"] == "instant" {
			reqLogger.Info("Triggering instant mode...")
			targetList := common.FindNonCompliantClustersForPolicy(policy)
			if len(targetList) > 0 {
				reqLogger.Info("Creating ansible job with targetList", "targetList", targetList)
				err = common.CreateAnsibleJob(cfgMap, r.dyamicClient)
				if err != nil {
					return reconcile.Result{}, err
				}
			} else {
				reqLogger.Info("No cluster is in noncompliant status, doing nothing...")
			}
		}
	}
	return reconcile.Result{}, nil
}
