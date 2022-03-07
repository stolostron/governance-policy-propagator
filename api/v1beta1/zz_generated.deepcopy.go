//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutomationDef) DeepCopyInto(out *AutomationDef) {
	*out = *in
	if in.ExtraVars != nil {
		in, out := &in.ExtraVars, &out.ExtraVars
		*out = new(runtime.RawExtension)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutomationDef.
func (in *AutomationDef) DeepCopy() *AutomationDef {
	if in == nil {
		return nil
	}
	out := new(AutomationDef)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyAutomation) DeepCopyInto(out *PolicyAutomation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyAutomation.
func (in *PolicyAutomation) DeepCopy() *PolicyAutomation {
	if in == nil {
		return nil
	}
	out := new(PolicyAutomation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicyAutomation) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyAutomationList) DeepCopyInto(out *PolicyAutomationList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PolicyAutomation, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyAutomationList.
func (in *PolicyAutomationList) DeepCopy() *PolicyAutomationList {
	if in == nil {
		return nil
	}
	out := new(PolicyAutomationList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicyAutomationList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyAutomationSpec) DeepCopyInto(out *PolicyAutomationSpec) {
	*out = *in
	in.Automation.DeepCopyInto(&out.Automation)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyAutomationSpec.
func (in *PolicyAutomationSpec) DeepCopy() *PolicyAutomationSpec {
	if in == nil {
		return nil
	}
	out := new(PolicyAutomationSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicyAutomationStatus) DeepCopyInto(out *PolicyAutomationStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicyAutomationStatus.
func (in *PolicyAutomationStatus) DeepCopy() *PolicyAutomationStatus {
	if in == nil {
		return nil
	}
	out := new(PolicyAutomationStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySet) DeepCopyInto(out *PolicySet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySet.
func (in *PolicySet) DeepCopy() *PolicySet {
	if in == nil {
		return nil
	}
	out := new(PolicySet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicySet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetList) DeepCopyInto(out *PolicySetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PolicySet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetList.
func (in *PolicySetList) DeepCopy() *PolicySetList {
	if in == nil {
		return nil
	}
	out := new(PolicySetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PolicySetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetResultCluster) DeepCopyInto(out *PolicySetResultCluster) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetResultCluster.
func (in *PolicySetResultCluster) DeepCopy() *PolicySetResultCluster {
	if in == nil {
		return nil
	}
	out := new(PolicySetResultCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetSpec) DeepCopyInto(out *PolicySetSpec) {
	*out = *in
	if in.Policies != nil {
		in, out := &in.Policies, &out.Policies
		*out = make([]NonEmptyString, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetSpec.
func (in *PolicySetSpec) DeepCopy() *PolicySetSpec {
	if in == nil {
		return nil
	}
	out := new(PolicySetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetStatus) DeepCopyInto(out *PolicySetStatus) {
	*out = *in
	if in.Placement != nil {
		in, out := &in.Placement, &out.Placement
		*out = make([]PolicySetStatusPlacement, len(*in))
		copy(*out, *in)
	}
	if in.Results != nil {
		in, out := &in.Results, &out.Results
		*out = make([]PolicySetStatusResult, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetStatus.
func (in *PolicySetStatus) DeepCopy() *PolicySetStatus {
	if in == nil {
		return nil
	}
	out := new(PolicySetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetStatusPlacement) DeepCopyInto(out *PolicySetStatusPlacement) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetStatusPlacement.
func (in *PolicySetStatusPlacement) DeepCopy() *PolicySetStatusPlacement {
	if in == nil {
		return nil
	}
	out := new(PolicySetStatusPlacement)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PolicySetStatusResult) DeepCopyInto(out *PolicySetStatusResult) {
	*out = *in
	if in.Clusters != nil {
		in, out := &in.Clusters, &out.Clusters
		*out = make([]PolicySetResultCluster, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PolicySetStatusResult.
func (in *PolicySetStatusResult) DeepCopy() *PolicySetStatusResult {
	if in == nil {
		return nil
	}
	out := new(PolicySetStatusResult)
	in.DeepCopyInto(out)
	return out
}
