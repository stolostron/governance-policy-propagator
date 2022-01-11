// Copyright (c) 2020 Red Hat, Inc.
package controller

import (
	"github.com/stolostron/governance-policy-propagator/pkg/controller/propagator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, propagator.Add)
}
