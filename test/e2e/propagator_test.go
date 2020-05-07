package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/apps/v1"
	. "github.com/open-cluster-management/governance-policy-propagator/test/e2e"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Propagator", func() {
	Describe("Creating a policy in user ns", func() {
		It("should be created in user ns", func() {
			By("Creating ../resources/test-policy.yaml")
			Kubectl("apply",
				"-f", "../resources/test-policy.yaml",
				"-n", testNamespace)
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, "test-policy", testNamespace, true, 15)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns", func() {
			By("Patch test-policy-plr with a decision")
			plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, "test-policy-plr", testNamespace, true, 15)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{
				Decisions: []appsv1.PlacementDecision{
					{
						ClusterName:      "managed1",
						ClusterNamespace: "managed1",
					},
				},
			}
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, 30)
		})
		It("should remove replicated policy from cluster ns", func() {
			By("Deleting placement binding")
			Kubectl("delete",
				"placementbindings.policies.open-cluster-management.io", "test-policy-pb",
				"-n", testNamespace)
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, 30)
		})
	})
})
