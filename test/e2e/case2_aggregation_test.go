// Copyright (c) 2020 Red Hat, Inc.
// +build integration

package e2e_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/open-cluster-management/governance-policy-propagator/test/e2e"
)

var _ = Describe("Test policy status aggregation", func() {
	It("should be created in user ns", func() {
		By("Creating ../resources/test-policy.yaml")
		Kubectl("apply",
			"-f", "../resources/test-policy.yaml",
			"-n", testNamespace)
		plc := GetWithTimeout(clientHubDynamic, gvrPolicy, "test-policy", testNamespace, true, 15)
		Expect(plc).NotTo(BeNil())
	})
})
