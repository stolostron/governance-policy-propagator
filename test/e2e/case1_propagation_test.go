// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package e2e

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	appsv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/placementrule/v1"

	policiesv1 "open-cluster-management.io/governance-policy-propagator/api/v1"
	"open-cluster-management.io/governance-policy-propagator/controllers/common"
	"open-cluster-management.io/governance-policy-propagator/test/utils"
)

const (
	case1PolicyName string = "case1-test-policy"
	case1PolicyYaml string = "../resources/case1_propagation/case1-test-policy.yaml"
)

var _ = Describe("Test policy propagation", func() {
	Describe("Create policy/pb/plc in ns:"+testNamespace+" and then update pb/plc", func() {
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed2", func() {
			By("Patching test-policy-plr with decision of cluster managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed2", true, defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPolicy,
				testNamespace+"."+case1PolicyName,
				"managed1",
				true,
				defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			plc = utils.GetWithTimeout(
				clientHubDynamic,
				gvrPolicy,
				testNamespace+"."+case1PolicyName,
				"managed2",
				false,
				defaultTimeoutSeconds,
			)
			Expect(plc).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-plr with decision of both managed1 and managed2")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1", "managed2")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a non existing plr")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plr-nonexists",
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plr")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace, true,
				defaultTimeoutSeconds,
			)
			pb.Object["placementRef"] = &policiesv1.Subject{
				APIGroup: "apps.open-cluster-management.io",
				Kind:     "PlacementRule",
				Name:     case1PolicyName + "-plr",
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a plc with wrong apigroup")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy1.open-cluster-management.io",
					Kind:     "Policy",
					Name:     case1PolicyName,
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace, true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy.open-cluster-management.io",
					Kind:     "Policy",
					Name:     case1PolicyName,
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should prevent user from creating subject with unsupported kind", func() {
			By("Patching test-policy-pb with a plc with wrong kind")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy.open-cluster-management.io",
					Kind:     "Policy1",
					Name:     case1PolicyName,
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy.open-cluster-management.io",
					Kind:     "Policy",
					Name:     case1PolicyName,
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-pb with a plc with wrong name")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy.open-cluster-management.io",
					Kind:     "Policy",
					Name:     case1PolicyName + "1",
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should propagate to cluster ns managed1 and managed2", func() {
			By("Patching test-policy-pb with correct plc")
			pb := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPlacementBinding,
				case1PolicyName+"-pb",
				testNamespace,
				true,
				defaultTimeoutSeconds,
			)
			pb.Object["subjects"] = []policiesv1.Subject{
				{
					APIGroup: "policy.open-cluster-management.io",
					Kind:     "Policy",
					Name:     case1PolicyName,
				},
			}
			_, err := clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Update(
				context.TODO(), pb, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 2, true, defaultTimeoutSeconds)
		})
		It("should remove policy from ns managed1 and managed2", func() {
			By("Patching test-policy-plr with no decision")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = &appsv1.PlacementRuleStatus{}
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should clean up", func() {
			utils.Kubectl("delete",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			opt := metav1.ListOptions{}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 10)
		})
	})

	Describe("Create policy/pb/plc in ns:"+testNamespace+" and then update policy", func() {
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should update replicated policy in ns managed1", func() {
			By("Patching test-policy with spec.remediationAction = enforce")
			rootPlc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(rootPlc).NotTo(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["remediationAction"]).To(Equal("inform"))
			rootPlc.Object["spec"].(map[string]interface{})["remediationAction"] = "enforce"
			rootPlc, err := clientHubDynamic.Resource(gvrPolicy).Namespace(testNamespace).Update(
				context.TODO(), rootPlc, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			Eventually(func() interface{} {
				replicatedPlc := utils.GetWithTimeout(
					clientHubDynamic,
					gvrPolicy,
					testNamespace+"."+case1PolicyName,
					"managed1",
					true,
					defaultTimeoutSeconds,
				)

				return replicatedPlc.Object["spec"]
			}, defaultTimeoutSeconds, 1).Should(utils.SemanticEqual(rootPlc.Object["spec"]))
		})
		It("should remove replicated policy in ns managed1", func() {
			By("Patching test-policy with spec.disabled = true")
			rootPlc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(rootPlc).NotTo(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["disabled"]).To(Equal(false))
			rootPlc.Object["spec"].(map[string]interface{})["disabled"] = true
			rootPlc, err := clientHubDynamic.Resource(gvrPolicy).Namespace(testNamespace).Update(
				context.TODO(), rootPlc, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			Expect(rootPlc.Object["spec"].(map[string]interface{})["disabled"]).To(Equal(true))
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, true, defaultTimeoutSeconds)
		})
		It("should be created in user ns", func() {
			By("Creating " + case1PolicyYaml)
			utils.Kubectl("apply",
				"-f", case1PolicyYaml,
				"-n", testNamespace)
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(plc).NotTo(BeNil())
		})
		It("should propagate to cluster ns managed1", func() {
			By("Patching test-policy-plr with decision of cluster managed1")
			plr := utils.GetWithTimeout(
				clientHubDynamic, gvrPlacementRule, case1PolicyName+"-plr", testNamespace, true, defaultTimeoutSeconds,
			)
			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())
			plc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, testNamespace+"."+case1PolicyName, "managed1", true, defaultTimeoutSeconds,
			)
			Expect(plc).ToNot(BeNil())
			opt := metav1.ListOptions{
				LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case1PolicyName,
			}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 1, true, defaultTimeoutSeconds)
		})
		It("should update test-policy to a different policy template", func() {
			By("Creating ../resources/case1_propagation/case1-test-policy2.yaml")
			utils.Kubectl("apply",
				"-f", "../resources/case1_propagation/case1-test-policy2.yaml",
				"-n", testNamespace)
			rootPlc := utils.GetWithTimeout(
				clientHubDynamic, gvrPolicy, case1PolicyName, testNamespace, true, defaultTimeoutSeconds,
			)
			Expect(rootPlc).NotTo(BeNil())
			yamlPlc := utils.ParseYaml("../resources/case1_propagation/case1-test-policy2.yaml")
			Eventually(func() interface{} {
				replicatedPlc := utils.GetWithTimeout(
					clientHubDynamic,
					gvrPolicy,
					testNamespace+"."+case1PolicyName,
					"managed1",
					true,
					defaultTimeoutSeconds,
				)

				return replicatedPlc.Object["spec"]
			}, defaultTimeoutSeconds, 1).Should(utils.SemanticEqual(yamlPlc.Object["spec"]))
		})
		It("should clean up", func() {
			utils.Kubectl("delete",
				"-f", "../resources/case1_propagation/case1-test-policy2.yaml",
				"-n", testNamespace)
			opt := metav1.ListOptions{}
			utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, 10)
		})
	})

	Describe("Test spec.copyPolicyMetadata", Ordered, func() {
		const policyName = "case1-metadata"
		policyClient := func() dynamic.ResourceInterface {
			return clientHubDynamic.Resource(gvrPolicy).Namespace(testNamespace)
		}

		getPolicy := func() *policiesv1.Policy {
			return &policiesv1.Policy{
				TypeMeta: metav1.TypeMeta{
					Kind:       policiesv1.Kind,
					APIVersion: policiesv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      policyName,
					Namespace: testNamespace,
				},
				Spec: policiesv1.PolicySpec{
					Disabled:        false,
					PolicyTemplates: []*policiesv1.PolicyTemplate{},
				},
			}
		}

		BeforeAll(func() {
			By("Creating a PlacementRule")
			plr := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps.open-cluster-management.io/v1",
					"kind":       "PlacementRule",
					"metadata": map[string]interface{}{
						"name":      policyName,
						"namespace": testNamespace,
					},
					"spec": map[string]interface{}{},
				},
			}
			var err error
			plr, err = clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).Create(
				context.TODO(), plr, metav1.CreateOptions{},
			)
			Expect(err).To(BeNil())

			plr.Object["status"] = utils.GeneratePlrStatus("managed1")
			_, err = clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
				context.TODO(), plr, metav1.UpdateOptions{},
			)
			Expect(err).To(BeNil())

			By("Creating a PlacementBinding")
			plb := &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": policiesv1.GroupVersion.String(),
					"kind":       "PlacementBinding",
					"metadata": map[string]interface{}{
						"name":      policyName,
						"namespace": testNamespace,
					},
					"placementRef": map[string]interface{}{
						"apiGroup": gvrPlacementRule.Group,
						"kind":     "PlacementRule",
						"name":     policyName,
					},
					"subjects": []interface{}{
						map[string]interface{}{
							"apiGroup": policiesv1.GroupVersion.Group,
							"kind":     policiesv1.Kind,
							"name":     policyName,
						},
					},
				},
			}

			_, err = clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Create(
				context.TODO(), plb, metav1.CreateOptions{},
			)
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			By("Removing the policy")
			err := policyClient().Delete(context.TODO(), policyName, metav1.DeleteOptions{})
			if !errors.IsNotFound(err) {
				Expect(err).To(BeNil())
			}

			replicatedPolicy := utils.GetWithTimeout(
				clientHubDynamic,
				gvrPolicy,
				testNamespace+"."+policyName,
				"managed1",
				false,
				defaultTimeoutSeconds,
			)
			Expect(replicatedPolicy).To(BeNil())
		})

		AfterAll(func() {
			By("Removing the PlacementRule")
			err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).Delete(
				context.TODO(), policyName, metav1.DeleteOptions{},
			)
			if !errors.IsNotFound(err) {
				Expect(err).To(BeNil())
			}

			By("Removing the PlacementBinding")
			err = clientHubDynamic.Resource(gvrPlacementBinding).Namespace(testNamespace).Delete(
				context.TODO(), policyName, metav1.DeleteOptions{},
			)
			if !errors.IsNotFound(err) {
				Expect(err).To(BeNil())
			}
		})

		It("Verifies that spec.copyPolicyMetadata defaults to unset", func() {
			By("Creating an empty policy")
			policy := getPolicy()
			policyMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
			Expect(err).To(BeNil())

			policyRV, err := policyClient().Create(
				context.TODO(), &unstructured.Unstructured{Object: policyMap}, metav1.CreateOptions{},
			)
			Expect(err).To(BeNil())

			_, found, _ := unstructured.NestedBool(policyRV.Object, "spec", "copyPolicyMetadata")
			Expect(found).To(BeFalse())
		})

		It("verifies that the labels and annotations are copied with spec.copyPolicyMetadata=true", func() {
			By("Creating a policy with labels and annotations")
			policy := getPolicy()
			copyPolicyMetadata := true
			policy.Spec.CopyPolicyMetadata = &copyPolicyMetadata
			policy.SetAnnotations(map[string]string{"do": "copy", "please": "do-copy"})
			policy.SetLabels(map[string]string{"do": "copy", "please": "do-copy"})

			policyMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
			Expect(err).To(BeNil())

			_, err = policyClient().Create(
				context.TODO(), &unstructured.Unstructured{Object: policyMap}, metav1.CreateOptions{},
			)
			Expect(err).To(BeNil())

			Eventually(func(g Gomega) {
				replicatedPlc := utils.GetWithTimeout(
					clientHubDynamic,
					gvrPolicy,
					testNamespace+"."+policyName,
					"managed1",
					true,
					defaultTimeoutSeconds,
				)

				annotations := replicatedPlc.GetAnnotations()
				g.Expect(annotations["do"]).To(Equal("copy"))
				g.Expect(annotations["please"]).To(Equal("do-copy"))
				// This annotation is always set.
				g.Expect(annotations["argocd.argoproj.io/compare-options"]).To(Equal("IgnoreExtraneous"))

				labels := replicatedPlc.GetLabels()
				g.Expect(labels["do"]).To(Equal("copy"))
				g.Expect(labels["please"]).To(Equal("do-copy"))
			}, defaultTimeoutSeconds, 1).Should(Succeed())
		})

		It("verifies that the labels and annotations are not copied with spec.copyPolicyMetadata=false", func() {
			By("Creating a policy with labels and annotations")
			policy := getPolicy()
			copyPolicyMetadata := false
			policy.Spec.CopyPolicyMetadata = &copyPolicyMetadata
			policy.SetAnnotations(map[string]string{"do-not": "copy", "please": "do-not-copy"})
			policy.SetLabels(map[string]string{"do-not": "copy", "please": "do-not-copy"})

			policyMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(policy)
			Expect(err).To(BeNil())

			_, err = policyClient().Create(
				context.TODO(), &unstructured.Unstructured{Object: policyMap}, metav1.CreateOptions{},
			)
			Expect(err).To(BeNil())

			Eventually(func(g Gomega) {
				replicatedPlc := utils.GetWithTimeout(
					clientHubDynamic,
					gvrPolicy,
					testNamespace+"."+policyName,
					"managed1",
					true,
					defaultTimeoutSeconds,
				)

				annotations := replicatedPlc.GetAnnotations()
				g.Expect(annotations["do-not"]).To(Equal(""))
				g.Expect(annotations["please"]).To(Equal(""))
				// This annotation is always set.
				g.Expect(annotations["argocd.argoproj.io/compare-options"]).To(Equal("IgnoreExtraneous"))

				labels := replicatedPlc.GetLabels()
				g.Expect(labels["do-not"]).To(Equal(""))
				g.Expect(labels["please"]).To(Equal(""))
			}, defaultTimeoutSeconds, 1).Should(Succeed())
		})
	})
})
