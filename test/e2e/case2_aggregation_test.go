// Copyright (c) 2020 Red Hat, Inc.
// +build integration

package e2e

import (
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "github.com/open-cluster-management/governance-policy-propagator/pkg/apis/apps/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	// . "github.com/open-cluster-management/governance-policy-propagator/test/e2e"
)

const case2PolicyName string = "case2-test-policy"
const case2PolicyYaml string = "../resources/case2_aggregation/case2-test-policy.yaml"

var _ = Describe("Test policy status aggregation", func() {
	Describe("Create policy/pb/plc in ns:"+testNamespace+" and then update pb and", func() {
		It("should be created in user ns", func() {
			By("Creating " + case2PolicyYaml)
			Kubectl("apply",
				"-f", case2PolicyYaml,
				"-n", testNamespace)
			plc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
			Expect(plc).NotTo(BeNil())
		})
	})
	It("should contain status.placement with managed1", func() {
		By("Patch test-policy-plr with decision of cluster managed1")
		plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case2PolicyName+"-plr", testNamespace, true, 15)
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
		plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case2PolicyName, "managed1", true, 15)
		Expect(plc).ToNot(BeNil())
		opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + "." + case2PolicyName}
		ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, 30)
		By("Checking the status.placement of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed1-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should contain status.placement with both managed1 and managed2", func() {
		By("Patch test-policy-plr with decision of cluster managed1 and managed2")
		plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case2PolicyName+"-plr", testNamespace, true, 15)
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
		plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case2PolicyName, "managed2", true, 15)
		Expect(plc).ToNot(BeNil())
		opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + "." + case2PolicyName}
		ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, 30)
		By("Patch checking the status.placement of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should contain status.placement with managed2", func() {
		By("Patch test-policy-plr with decision of cluster managed2")
		plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case2PolicyName+"-plr", testNamespace, true, 15)
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
		plc := GetWithTimeout(clientHubDynamic, gvrPolicy, testNamespace+"."+case2PolicyName, "managed2", true, 15)
		Expect(plc).ToNot(BeNil())
		opt := metav1.ListOptions{LabelSelector: "root-policy=" + testNamespace + "." + case2PolicyName}
		ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, 30)
		By("Checking the status.placement of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed2-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should contain status.placement with two pb/plr", func() {
		By("Creating pb-plr-2 to binding second set of placement")
		Kubectl("apply",
			"-f", "../resources/case2_aggregation/pb-plr-2.yaml",
			"-n", testNamespace)
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-placement-single-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should contain status.placement with two pb/plr and both status", func() {
		By("Creating pb-plr-2 to binding second set of placement")
		plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case2PolicyName+"-plr2", testNamespace, true, 15)
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
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-placement-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should still contain status.placement with two pb/plr and both status", func() {
		By("Patch" + case2PolicyName + "-plr2 with both managed1 and managed2")
		plr := GetWithTimeout(clientHubDynamic, gvrPlacementRule, case2PolicyName+"-plr2", testNamespace, true, 15)
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
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-placement-status.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should still contain status.placement with two pb, one plr and both status", func() {
		By("Remove" + case2PolicyName + "-plr")
		Kubectl("delete",
			"placementrule", case2PolicyName+"-plr",
			"-n", testNamespace)
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-placement-status-missing-plr.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should clear out status.status", func() {
		By("Remove" + case2PolicyName + "-plr2")
		Kubectl("delete",
			"placementrule", case2PolicyName+"-plr2",
			"-n", testNamespace)
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		yamlFile, err := ioutil.ReadFile("../resources/case2_aggregation/managed-both-placementbinding.yaml")
		Expect(err).To(BeNil())
		yamlPlc := &unstructured.Unstructured{}
		err = yaml.Unmarshal(yamlFile, yamlPlc)
		Expect(err).To(BeNil())
		equal := equality.Semantic.DeepEqual(rootPlc.Object["status"], yamlPlc.Object["status"])
		Expect(equal).To(Equal(true))
	})
	It("should clear out status", func() {
		By("Remove" + case2PolicyName + "-pb and " + case2PolicyName + "-pb2")
		Kubectl("delete",
			"placementbinding", case2PolicyName+"-pb",
			"-n", testNamespace)
		Kubectl("delete",
			"placementbinding", case2PolicyName+"-pb2",
			"-n", testNamespace)
		By("Patch checking the status of root policy")
		time.Sleep(2 * time.Second)
		rootPlc := GetWithTimeout(clientHubDynamic, gvrPolicy, case2PolicyName, testNamespace, true, 15)
		emptyStatus := map[string]interface{}{}
		Expect(rootPlc.Object["status"]).To(Equal(emptyStatus))
	})
	It("should clean up", func() {
		Kubectl("delete",
			"-f", case2PolicyYaml,
			"-n", testNamespace)
		opt := metav1.ListOptions{}
		ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 10)
	})
})
