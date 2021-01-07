// Copyright (c) 2021 Red Hat, Inc.

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/open-cluster-management/governance-policy-propagator/test/utils"
)

const case4PolicyName string = "case4-test-policy"
const case4PolicyYaml string = "../resources/case4_unexpected_policy/case4-test-policy.yaml"

var _ = Describe("Test unexpect policy handling", func() {
	It("Unexpected policy in cluster namespace should be deleted", func() {
		By("Creating " + case4PolicyYaml + "in cluster namespace: managed1")
		out, _ := utils.KubectlWithOutput("apply",
			"-f", case4PolicyYaml,
			"-n", "managed1")
		Expect(out).Should(ContainSubstring(case4PolicyName + " created"))
		Eventually(func() interface{} {
			return utils.GetWithTimeout(clientHubDynamic, gvrPolicy, case4PolicyName, "managed1", false, defaultTimeoutSeconds)
		}, defaultTimeoutSeconds, 1).Should(BeNil())
	})
})
