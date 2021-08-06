// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"fmt"
	"os"
	"strings"

	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/automation"
	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/policymetrics"
	"github.com/open-cluster-management/governance-policy-propagator/pkg/controller/propagator"
)

// reportMetrics returns a bool on whether to report GRC metrics from the propagator
func reportMetrics() bool {
	metrics, found := os.LookupEnv("DISABLE_REPORT_METRICS")
	if found && strings.ToLower(metrics) == "true" {
		return false
	}
	return true
}

func init() {
	fmt.Println("-- in init --")
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, propagator.Add, automation.Add)
	if reportMetrics() {
		fmt.Println("--adding metrics --")
		AddToManagerFuncs = append(AddToManagerFuncs, policymetrics.Add)
	}
}
