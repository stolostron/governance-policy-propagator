// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package e2e

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/stolostron/governance-policy-propagator/controllers/common"
	"github.com/stolostron/governance-policy-propagator/test/utils"
)

const (
	path                               string = "../resources/case10_policyset_propagation/"
	case10PolicyName                   string = "case10-test-policy"
	case10PolicySetName                string = "case10-test-policyset"
	case10PolicySetYaml                string = path + "case10-test-policyset.yaml"
	case10PolicySetPlacementYaml       string = path + "case10-test-policyset-placement.yaml"
	case10PolicySetPolicyYaml          string = path + "case10-test-policyset-policy.yaml"
	case10PolicySetPolicyPlacementYaml string = path + "case10-test-policyset-policy-placement.yaml"
)

var _ = Describe("Test policyset propagation", func() {
	Describe("Test policy propagation through policyset placementbinding with placementrule", func() {
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed2", func() {
			By("Patching test-policy-plr with decision of cluster managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPolicy,
				testNamespace+"."+case10PolicyName,
				"managed1",
				true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic,
				gvrPolicy,
				testNamespace+"."+case10PolicyName,
				"managed2",
				false,
				defaultTimeoutSeconds,
			)
			Expect(plc).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting policyset")
			_, err := utils.KubectlWithOutput("delete", "policyset", case10PolicySetName, "-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting placementbinding")
			_, err := utils.KubectlWithOutput("delete", "PlacementBinding", case10PolicySetName+"-pb", "-n",
				testNamespace)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting placementrule")
			_, err := utils.KubectlWithOutput("delete", "PlacementRule", case10PolicySetName+"-plr", "-n",
				testNamespace)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should clean up", func() {
			By("Deleting " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("delete",
				"-f", case10PolicySetYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
	})

	Describe("Test policy propagation through both policy and policyset placementbinding with placementrule", func() {
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPolicyYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching " + case10PolicySetName + "-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching " + case10PolicyName + "-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicyName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching " + case10PolicyName + "-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicyName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching " + case10PolicySetName + "-plr with decision of cluster managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case10PolicySetName+"-plr", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should clean up", func() {
			By("Deleting " + case10PolicySetYaml)
			_, err := utils.KubectlWithOutput("delete",
				"-f", case10PolicySetPolicyYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
	})

	Describe("Test policy propagation through policyset placementbinding with placement", func() {
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetPlacementYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plm with decision of cluster managed1")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed2", func() {
			By("Patching test-policy-plm with decision of cluster managed2")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to both cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plm with decision of cluster managed2")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting policyset")
			_, err := utils.KubectlWithOutput("delete", "policyset", case10PolicySetName, "-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetPlacementYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to both cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plm with decision of cluster managed2")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting placementbinding")
			_, err := utils.KubectlWithOutput("delete", "PlacementBinding", case10PolicySetName+"-pb", "-n",
				testNamespace)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetPlacementYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to both cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plm with decision of cluster managed2")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Deleting placementDecision")
			_, err := utils.KubectlWithOutput("delete", "PlacementDecision", case10PolicySetName+"-plm-decision", "-n",
				testNamespace)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetPlacementYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should cleanup", func() {
			By("Deleting " + case10PolicySetPlacementYaml)
			_, err := utils.KubectlWithOutput("delete",
				"-f", case10PolicySetPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
	})

	Describe("Test policy propagation through both policy and policyset placementbinding with placement", func() {
		It("should be created in user ns", func() {
			By("Creating " + case10PolicySetPolicyPlacementYaml)
			_, err := utils.KubectlWithOutput("apply",
				"-f", case10PolicySetPolicyPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plcSet).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plm with decision of cluster managed1")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicyName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to both cluster ns managed1 and managed2", func() {
			By("Patching test-policyset-plm with decision of cluster managed2")
			plm := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementDecision, case10PolicySetName+"-plm-decision", testNamespace, true,
				defaultTimeoutSeconds,
			)
			plm.Object["status"] = utils.GeneratePldStatus(plm.GetName(), plm.GetNamespace(), "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementDecision).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plm, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed1", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case10PolicyName, "managed2", true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should cleanup", func() {
			By("Deleting " + case10PolicySetPolicyPlacementYaml)
			_, err := utils.KubectlWithOutput("delete",
				"-f", case10PolicySetPolicyPlacementYaml,
				"-n", testNamespace)
			Expect(err).To(BeNil())
			plcSet := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicySet, case10PolicySetName, testNamespace, false, defaultTimeoutSeconds,
			)
			Expect(plcSet).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case10PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
	})
})
