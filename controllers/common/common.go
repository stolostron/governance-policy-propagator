// Copyright (c) 2020 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package common

import (
	"errors"
	"fmt"
	"strings"

	policiesv1 "github.com/stolostron/governance-policy-propagator/api/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	clusterv1 "open-cluster-management.io/api/cluster/v1"
)

const (
	APIGroup              string = "policy.open-cluster-management.io"
	ClusterNameLabel      string = APIGroup + "/cluster-name"
	ClusterNamespaceLabel string = APIGroup + "/cluster-namespace"
	RootPolicyLabel       string = APIGroup + "/root-policy"
)

var ErrInvalidLabelValue = errors.New("unexpected format of label value")

// IsInClusterNamespace check if policy is in cluster namespace
func IsInClusterNamespace(ns string, allClusters []clusterv1.ManagedCluster) bool {
	for _, cluster := range allClusters {
		if ns == cluster.GetName() {
			return true
		}
	}
	return false
}

// LabelsForRootPolicy returns the labels for given policy
func LabelsForRootPolicy(plc *policiesv1.Policy) map[string]string {
	return map[string]string{RootPolicyLabel: FullNameForPolicy(plc)}
}

// FullNameForPolicy returns the fully qualified name for given policy
// full qualified name: ${namespace}.${name}
func FullNameForPolicy(plc *policiesv1.Policy) string {
	return plc.GetNamespace() + "." + plc.GetName()
}

// CompareSpecAndAnnotation compares annotation and spec for given policies
// true if matches, false if doesn't match
func CompareSpecAndAnnotation(plc1 *policiesv1.Policy, plc2 *policiesv1.Policy) bool {
	annotationMatch := equality.Semantic.DeepEqual(plc1.GetAnnotations(), plc2.GetAnnotations())
	specMatch := equality.Semantic.DeepEqual(plc1.Spec, plc2.Spec)
	return annotationMatch && specMatch
}

// IsPbForPoicy compares group and kind with policy group and kind for given pb
func IsPbForPoicy(pb *policiesv1.PlacementBinding) bool {
	subjects := pb.Subjects
	found := false
	for _, subject := range subjects {
		if subject.Kind == policiesv1.Kind && subject.APIGroup == policiesv1.SchemeGroupVersion.Group {
			found = true
			break
		}
	}
	return found
}

// FindNonCompliantClustersForPolicy returns cluster in noncompliant status with given policy
func FindNonCompliantClustersForPolicy(plc *policiesv1.Policy) []string {
	clusterList := []string{}
	for _, clusterStatus := range plc.Status.Status {
		if clusterStatus.ComplianceState == policiesv1.NonCompliant {
			clusterList = append(clusterList, clusterStatus.ClusterName)
		}
	}
	return clusterList
}

func ParseRootPolicyLabel(rootPlc string) (name, namespace string, err error) {
	// namespaces can't have a `.` (but names can) so this always correctly pulls the namespace out
	namespace, name, found := strings.Cut(rootPlc, ".")
	if !found {
		err = fmt.Errorf("required at least one `.` in value of label `%v`: %w",
			RootPolicyLabel, ErrInvalidLabelValue)

		return "", "", err
	}

	return name, namespace, nil
}
