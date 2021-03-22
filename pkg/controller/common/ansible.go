// Copyright Contributors to the Open Cluster Management project

package common

import (
	"context"

	"github.com/ghodss/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// CreateAnsibleJob creates ansiblejob with given config map
func CreateAnsibleJob(cfgMap *corev1.ConfigMap, dyamicClient dynamic.Interface) error {
	ansibleJob := &unstructured.Unstructured{}
	err := yaml.Unmarshal([]byte(cfgMap.Data["ansibleJob.yaml"]), ansibleJob)
	if err != nil {
		return err
	}
	ansibleJobRes := schema.GroupVersionResource{Group: "tower.ansible.com", Version: "v1alpha1", 
		Resource: "ansiblejobs"}
	ansibleJob.SetGenerateName(cfgMap.GetName() + "-")
	ansibleJob.SetOwnerReferences([]metav1.OwnerReference{
		*metav1.NewControllerRef(cfgMap, cfgMap.GroupVersionKind()),
	})
	_, err = dyamicClient.Resource(ansibleJobRes).Namespace(cfgMap.GetNamespace()).
		Create(context.TODO(), ansibleJob, v1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}
