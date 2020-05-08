// Copyright (c) 2020 Red Hat, Inc.

package e2e

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	policiesv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/policies/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// . "github.com/open-cluster-management/governance-policy-propagator/test/e2e"
)

const case3PolicyName string = "case3-test-policy"
const case3PolicyYaml string = "../resources/case3_mutation_recovery/case3-test-policy.yaml"

var _ = Describe("Test unexpected policy mutation", func() {
	BeforeEach(func() {
		It("should be created in user ns", func() {
			By("Creating " + case3PolicyYaml)
			Kubectl("apply",
				"-f", case3PolicyYaml,
				"-n", testNamespace)
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, case3PolicyName, testNamespace, true, defaultTimeoutSeconds)
			Expect(plc).NotTo(BeNil())
		})
		It("should contain status.placement with violation status from both managed1 and managed2", func() {
			By("Patch test-policy-plr with decision of cluster managed1 and managed2")
			plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case3PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds)
			plr.Object["status"] = GeneratePlrStatus("managed1", "managed2")
			plr, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(plr, metav1.UpdateOptions{})
			Expect(err).To(BeNil())
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case3PolicyName, "managed2", true, defaultTimeoutSeconds)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + "." + case3PolicyName}
			By("Patch both replicated policy status to compliant")
			replicatedPlcList := ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
			for _, replicatedPlc := range replicatedPlcList.Items {
				replicatedPlc.Object["status"] = &policiesv1.PolicyStatus{
					ComplianceState: policiesv1.Compliant,
				}
				_, err = clientHubDynamic.Resource(gvrPolicy).Namespace(replicatedPlc.GetNamespace()).UpdateStatus(&replicatedPlc, metav1.UpdateOptions{})
				Expect(err).To(BeNil())
			}
			By("Checking the status of root policy")
			time.Sleep(2 * time.Second)
			rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case3PolicyName, testNamespace, true, defaultTimeoutSeconds)
			yamlPlc := ParseYaml("../resources/case3_aggregation/managed-both-status-compliant.yaml")
			equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
			Expect(equal).To(Equal(true))
		})
	})
	AfterEach(func() {
		It("should clean up", func() {
			Kubectl("delete",
				"-f", case3PolicyYaml,
				"-n", testNamespace)
			opt := metav1.ListOptions{}
			ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 10)
		})
	})
})
