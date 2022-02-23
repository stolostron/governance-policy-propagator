// Copyright (c) 2022 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package e2e

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	policiesv1 "github.com/stolostron/governance-policy-propagator/api/v1"
	"github.com/stolostron/governance-policy-propagator/controllers/common"
	"github.com/stolostron/governance-policy-propagator/test/utils"
)

const (
	case13PolicyName    string = "case13-test-policy"
	case13PolicySetName string = "case13-test-policyset"
	case13PolicySetYaml string = "../resources/case13_scalability/case13-test-policyset.yaml"
	CCLHiddenMsg        string = "Compliant cluster list is not displayed due to its large size"
	NCCLHiddenMsg       string = "Compliant and NonCompliant cluster " +
		"lists are not displayed due to its large size"
)

// batchProcess process for loop in parallel
func batchProcess(total int, maxBatchSize int, f func(int)) {
	if maxBatchSize == 0 {
		// default to 25
		maxBatchSize = 25
	}

	items := make([]int, total)
	for i := 0; i < total; i++ {
		items[i] = i
	}

	skip := 0
	filesAmount := total
	batchAmount := int(float64(filesAmount / maxBatchSize))

	for i := 0; i <= batchAmount; i++ {
		lowerBound := skip
		upperBound := skip + maxBatchSize

		if upperBound > filesAmount {
			upperBound = filesAmount
		}

		batchItems := items[lowerBound:upperBound]

		skip += maxBatchSize

		var itemProcessingGroup sync.WaitGroup

		itemProcessingGroup.Add(len(batchItems))

		for idx := range batchItems {
			go func(currentItem int) {
				defer itemProcessingGroup.Done()
				f(currentItem)
			}(batchItems[idx])
		}

		itemProcessingGroup.Wait()
	}
}

// testPlcsetStatusLimit tests the policyset status limit by incrementally adding the policy to policyset
func testPlcsetStatusLimit(startNumber int, cCLThreshold int, nCCLThreshold int) {
	steps := 10

	for i := startNumber - 1; i < nCCLThreshold+1; i++ {
		replicatedPlcClient := clientHubDynamic.Resource(gvrPolicySet).Namespace(testNamespace)

		if i == startNumber-1 {
			By("Adding policy #" + strconv.Itoa(i+1))

			value := []string{}
			for j := 0; j < startNumber; j++ {
				value = append(value, case13PolicyName)
			}

			patch := []interface{}{
				map[string]interface{}{
					"op":    "replace",
					"path":  "/spec/policies",
					"value": value,
				},
			}

			payload, err := json.Marshal(patch)
			Expect(err).To(BeNil())
			_, err = replicatedPlcClient.Patch(context.TODO(), case13PolicySetName,
				types.JSONPatchType, payload, metav1.PatchOptions{})
			Expect(err).To(BeNil())
		} else {
			patches := []interface{}{}
			singlePatch := map[string]interface{}{
				"op":    "add",
				"path":  "/spec/policies/-",
				"value": case13PolicyName,
			}
			if i < cCLThreshold && (cCLThreshold-i) > steps || i > cCLThreshold && (nCCLThreshold-i) > steps {
				for j := 0; j < steps; j++ {
					patches = append(patches, singlePatch)
				}
				i = i + steps - 1
			} else {
				patches = append(patches, singlePatch)
			}
			By("Adding policy #" + strconv.Itoa(i+1))

			payload, err := json.Marshal(patches)
			Expect(err).To(BeNil())
			_, err = replicatedPlcClient.Patch(context.TODO(), case13PolicySetName,
				types.JSONPatchType, payload, metav1.PatchOptions{})
			Expect(err).To(BeNil())
		}

		Eventually(
			func() interface{} {
				plcset := utils.GetWithTimeout(
					clientHubDynamic, gvrPolicySet, case13PolicySetName, testNamespace, true,
					defaultTimeoutSeconds,
				)
				results, ok, err := unstructured.NestedSlice(plcset.Object, "status", "results")
				Expect(ok).To(BeTrue())
				Expect(err).To(BeNil())

				return len(results)
			},
			defaultTimeoutSeconds,
			1,
		).Should(Equal(i + 1))

		plcset := utils.GetWithTimeout(
			clientHubDynamic, gvrPolicySet, case13PolicySetName, testNamespace, true, defaultTimeoutSeconds*2,
		)
		lastResult := plcset.Object["status"].(map[string]interface{})["results"].([]interface{})[i-1]

		if i < cCLThreshold {
			By("Checking if compliant clusters and nonCompliant clusters list exist")
			Expect(lastResult.(map[string]interface{})["message"]).To(BeNil())
			Expect(lastResult.(map[string]interface{})["compliantClusters"]).NotTo(BeNil())
			Expect(lastResult.(map[string]interface{})["nonCompliantClusters"]).NotTo(BeNil())
		} else if i < nCCLThreshold {
			By("Checking if compliant clusters list is hidden")
			Expect(lastResult.(map[string]interface{})["message"]).To(Equal(CCLHiddenMsg))
			Expect(lastResult.(map[string]interface{})["compliantClusters"]).To(BeNil())
			Expect(lastResult.(map[string]interface{})["nonCompliantClusters"]).NotTo(BeNil())
		} else {
			By("Checking if both lists are hidden")
			Expect(lastResult.(map[string]interface{})["message"]).To(Equal(NCCLHiddenMsg))
			Expect(lastResult.(map[string]interface{})["compliantClusters"]).To(BeNil())
			Expect(lastResult.(map[string]interface{})["nonCompliantClusters"]).To(BeNil())
		}
	}
}

var _ = Describe("Test policy scalability", func() {
	matrix := [][]int{
		{3000, 38, 75, 38},
		{2000, 58, 114, 55},
	}

	for i := 0; i < len(matrix); i++ {
		totalManagedClusters := matrix[i][0]
		// CompliantClusterList threshold
		CCListThreshold := matrix[i][1]
		// NontCompliantClusterList threshold
		NCCListThreshold := matrix[i][2]
		// startNumber time saver
		startNumber := matrix[i][3]
		Describe("Test "+strconv.Itoa(totalManagedClusters)+" managed clusters ", func() {
			It("should create ns and manages clusters", func() {
				By("Creating " + strconv.Itoa(totalManagedClusters) + " ns and managed cluster ")
				batchProcess(totalManagedClusters, 100, func(currentItem int) {
					managedCluster := "cluster" + strconv.Itoa(currentItem)

					_, err := clientHub.CoreV1().Namespaces().Create(context.TODO(), &corev1.Namespace{
						ObjectMeta: metav1.ObjectMeta{
							Name: managedCluster,
							Labels: map[string]string{
								"grc-test": "true",
							},
						},
					}, metav1.CreateOptions{})

					if err != nil {
						Expect(k8serrors.IsAlreadyExists(err)).To(BeTrue())
					} else {
						Expect(err).To(BeNil())
					}

					cluster := utils.GenerateManagedCluters(managedCluster, managedCluster)
					unstructuredManagedCluster, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cluster)
					Expect(err).To(BeNil())
					_, err = clientHubDynamic.Resource(gvrManagedCluster).Create(context.TODO(),
						&unstructured.Unstructured{Object: unstructuredManagedCluster}, metav1.CreateOptions{})

					if err != nil {
						Expect(k8serrors.IsAlreadyExists(err)).To(BeTrue())
					} else {
						Expect(err).To(BeNil())
					}
				})
			})
			It("should be created in user ns", func() {
				By("Creating " + case13PolicyName)
				utils.Kubectl("create",
					"-f", case13PolicySetYaml,
					"-n", testNamespace)
				plc := utils.GetWithTimeout(
					clientHubDynamic, gvrPolicy, case13PolicyName, testNamespace, true, defaultTimeoutSeconds,
				)
				Expect(plc).NotTo(BeNil())
				plcset := utils.GetWithTimeout(
					clientHubDynamic, gvrPolicySet, case13PolicySetName, testNamespace, true, defaultTimeoutSeconds,
				)
				Expect(plcset).NotTo(BeNil())
			})
			It("should propagate to "+strconv.Itoa(totalManagedClusters)+" cluster ns", func() {
				clusters := []string{}
				for i := 0; i < totalManagedClusters; i++ {
					managedCluster := "cluster" + strconv.Itoa(i)
					clusters = append(clusters, managedCluster)
				}
				By("Patching " + case13PolicySetName + "-plr with decision of " +
					strconv.Itoa(totalManagedClusters) + " clusters")
				plr := utils.GetWithTimeout(
					clientHubDynamic, gvrPlacementRule, case13PolicySetName+"-plr", testNamespace, true,
					defaultTimeoutSeconds,
				)
				plr.Object["status"] = utils.GeneratePlrStatus(clusters...)
				_, err := clientHubDynamic.Resource(gvrPlacementRule).Namespace(testNamespace).UpdateStatus(
					context.TODO(), plr, metav1.UpdateOptions{},
				)
				Expect(err).To(BeNil())
				By("Checking if policy " + case13PolicySetName + " has been replicated to " +
					strconv.Itoa(totalManagedClusters) + " cluster ns")
				opt := metav1.ListOptions{
					LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case13PolicyName,
				}
				utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, totalManagedClusters,
					true, defaultTimeoutSeconds*2)
			})
			It("should generate replicated status in "+strconv.Itoa(totalManagedClusters)+" cluster ns", func() {
				By("Patching " + strconv.Itoa(totalManagedClusters) +
					" replicated policies status to Compliant and NonCompliant")
				batchProcess(totalManagedClusters, 100, func(i int) {
					replicatedPlc := utils.GetWithTimeout(clientHubDynamic, gvrPolicy,
						testNamespace+"."+case13PolicyName,
						"cluster"+strconv.Itoa(i), true, defaultTimeoutSeconds)
					compliantState := policiesv1.NonCompliant
					if i%2 == 0 {
						compliantState = policiesv1.Compliant
					}
					replicatedPlc.Object["status"] = &policiesv1.PolicyStatus{
						ComplianceState: compliantState,
					}

					_, err := clientHubDynamic.Resource(gvrPolicy).Namespace(replicatedPlc.GetNamespace()).UpdateStatus(
						context.TODO(), replicatedPlc, metav1.UpdateOptions{},
					)
					Expect(err).To(BeNil())
				})
			})
			It("Add policies according to threshold", func() {
				testPlcsetStatusLimit(startNumber, CCListThreshold, NCCListThreshold)
			})
		})

		Describe("Clean up", func() {
			It("Should delete all policy manifests", func() {
				By("Deleting policyset/policy/pb/plr")
				_, err := utils.KubectlWithOutput("delete",
					"-f", case13PolicySetYaml,
					"-n", testNamespace)
				Expect(err).To(BeNil())
				opt := metav1.ListOptions{
					LabelSelector: common.RootPolicyLabel + "=" + testNamespace + "." + case13PolicyName,
				}
				By("Checking if replicated policies have been deleted")
				utils.ListWithTimeout(clientHubDynamic, gvrPolicy, opt, 0, false, defaultTimeoutSeconds*6)
			})
			It("Should delete all managed clusters", func() {
				By("Deleting managed cluster")
				batchProcess(totalManagedClusters, 100, func(i int) {
					managedCluster := "cluster" + strconv.Itoa(i)
					err := clientHubDynamic.Resource(gvrManagedCluster).Delete(
						context.TODO(), managedCluster, metav1.DeleteOptions{})
					Expect(err).To(BeNil())
				})
			})
		})
	}
})
