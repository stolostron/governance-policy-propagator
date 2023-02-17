package propagator

import (
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	appsv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"

	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
	"open-cluster-management.io/governance-policy-propagator/controllers/common"
)

const argoCDCompareOptionsAnnotation = "argocd.argoproj.io/compare-options"

// equivalentReplicatedPolicies compares replicated policies. Returns true if they match.
func equivalentReplicatedPolicies(plc1 *policiesv1.Policy, plc2 *policiesv1.Policy) bool {
	// Compare annotations
	if !equality.Semantic.DeepEqual(plc1.GetAnnotations(), plc2.GetAnnotations()) {
		return false
	}

	// Compare labels
	if !equality.Semantic.DeepEqual(plc1.GetLabels(), plc2.GetLabels()) {
		return false
	}

	// Compare the specs
	return equality.Semantic.DeepEqual(plc1.Spec, plc2.Spec)
}

// buildReplicatedPolicy constructs a replicated policy based on a root policy and a placementDecision.
// In particular, it adds labels that the policy framework uses.
func (r *PolicyReconciler) buildReplicatedPolicy(
	root *policiesv1.Policy, decision appsv1.PlacementDecision,
) *policiesv1.Policy {
	replicatedName := common.FullNameForPolicy(root)

	replicated := root.DeepCopy()
	replicated.SetName(replicatedName)
	replicated.SetNamespace(decision.ClusterNamespace)
	replicated.SetResourceVersion("")
	replicated.SetFinalizers(nil)
	replicated.SetOwnerReferences(nil)

	labels := root.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	if root.Spec.CopyPolicyMetadata != nil && !*root.Spec.CopyPolicyMetadata {
		originalLabels := replicated.GetLabels()

		for label := range originalLabels {
			if !strings.HasPrefix(label, policiesv1.GroupVersion.Group+"/") {
				delete(labels, label)
			}
		}
	}

	// Extra labels on replicated policies
	labels[common.ClusterNameLabel] = decision.ClusterName
	labels[common.ClusterNamespaceLabel] = decision.ClusterNamespace
	labels[common.RootPolicyLabel] = replicatedName

	replicated.SetLabels(labels)

	annotations := replicated.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	if root.Spec.CopyPolicyMetadata != nil && !*root.Spec.CopyPolicyMetadata {
		originalAnnotations := replicated.GetAnnotations()

		for annotation := range originalAnnotations {
			if !strings.HasPrefix(annotation, policiesv1.GroupVersion.Group+"/") {
				delete(annotations, annotation)
			}
		}
	}

	// Always set IgnoreExtraneous to avoid ArgoCD managing the replicated policy.
	annotations[argoCDCompareOptionsAnnotation] = "IgnoreExtraneous"

	replicated.SetAnnotations(annotations)

	return replicated
}
