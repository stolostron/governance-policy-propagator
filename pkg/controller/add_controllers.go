// Copyright Contributors to the Open Cluster Management project

package controller

import (
	"github.com/stolostron/governance-policy-propagator/pkg/controller/automation"
	"github.com/stolostron/governance-policy-propagator/pkg/controller/propagator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, propagator.Add, automation.Add)
}
