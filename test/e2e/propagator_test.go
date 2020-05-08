// Copyright (c) 2020 Red Hat, Inc.
// +build integration

package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/apps/v1"
	. "github.com/open-cluster-management/governance-policy-propagator/test/e2e"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Test policy propagation", func() {
	Describe("Creating policy/pb/plc in ns:"+testNamespace, func() {
		It("should be created in user ns", func() {
			By("Creating ../resources/test-policy.yaml")
			Kubectl("apply",
				"-f", "../resources/test-policy.yaml",
				"-n", testNamespace)
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, "test-policy", testNamespace, true, 15)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patch test-policy-plr with decision of cluster managed1")
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
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+".test-policy", "managed1", true, 15)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, 30)
		})
		It("should propagate to cluster ns managed2", func() {
			By("Patch test-policy-plr with decision of cluster managed2")
			plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, "test-policy-plr", testNamespace, true, 15)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{
				Decisions: []appsv1.PlacementDecision{
					{
						ClusterName:      "managed2",
						ClusterNamespace: "managed2",
					},
				},
			}
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+".test-policy", "managed2", true, 15)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, 30)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patch test-policy-plr with decision of both managed1 and managed2")
			plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, "test-policy-plr", testNamespace, true, 15)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{
				Decisions: []appsv1.PlacementDecision{
					{
						ClusterName:      "managed1",
						ClusterNamespace: "managed1",
					},
					{
						ClusterName:      "managed2",
						ClusterNamespace: "managed2",
					},
				},
			}
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, 30)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patch test-policy-plr with decision of both managed1 and managed2")
			plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, "test-policy-plr", testNamespace, true, 15)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{}
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + ".test-policy"}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, 30)
		})
	})
})
