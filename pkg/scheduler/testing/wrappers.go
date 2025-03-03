/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testing

import (
	"fmt"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	imageutils "k8s.io/kubernetes/test/utils/image"
	"k8s.io/utils/pointer"
)

var zero int64

// NodeSelectorWrapper wraps a NodeSelector inside.
type NodeSelectorWrapper struct{ v1.NodeSelector }

// MakeNodeSelector creates a NodeSelector wrapper.
func MakeNodeSelector() *NodeSelectorWrapper {
	return &NodeSelectorWrapper{v1.NodeSelector{}}
}

// In injects a matchExpression (with an operator IN) as a selectorTerm
// to the inner nodeSelector.
// NOTE: appended selecterTerms are ORed.
func (s *NodeSelectorWrapper) In(key string, vals []string) *NodeSelectorWrapper {
	expression := v1.NodeSelectorRequirement{
		Key:      key,
		Operator: v1.NodeSelectorOpIn,
		Values:   vals,
	}
	selectorTerm := v1.NodeSelectorTerm{}
	selectorTerm.MatchExpressions = append(selectorTerm.MatchExpressions, expression)
	s.NodeSelectorTerms = append(s.NodeSelectorTerms, selectorTerm)
	return s
}

// NotIn injects a matchExpression (with an operator NotIn) as a selectorTerm
// to the inner nodeSelector.
func (s *NodeSelectorWrapper) NotIn(key string, vals []string) *NodeSelectorWrapper {
	expression := v1.NodeSelectorRequirement{
		Key:      key,
		Operator: v1.NodeSelectorOpNotIn,
		Values:   vals,
	}
	selectorTerm := v1.NodeSelectorTerm{}
	selectorTerm.MatchExpressions = append(selectorTerm.MatchExpressions, expression)
	s.NodeSelectorTerms = append(s.NodeSelectorTerms, selectorTerm)
	return s
}

// Obj returns the inner NodeSelector.
func (s *NodeSelectorWrapper) Obj() *v1.NodeSelector {
	return &s.NodeSelector
}

// LabelSelectorWrapper wraps a LabelSelector inside.
type LabelSelectorWrapper struct{ metav1.LabelSelector }

// MakeLabelSelector creates a LabelSelector wrapper.
func MakeLabelSelector() *LabelSelectorWrapper {
	return &LabelSelectorWrapper{metav1.LabelSelector{}}
}

// Label applies a {k,v} pair to the inner LabelSelector.
func (s *LabelSelectorWrapper) Label(k, v string) *LabelSelectorWrapper {
	if s.MatchLabels == nil {
		s.MatchLabels = make(map[string]string)
	}
	s.MatchLabels[k] = v
	return s
}

// In injects a matchExpression (with an operator In) to the inner labelSelector.
func (s *LabelSelectorWrapper) In(key string, vals []string) *LabelSelectorWrapper {
	expression := metav1.LabelSelectorRequirement{
		Key:      key,
		Operator: metav1.LabelSelectorOpIn,
		Values:   vals,
	}
	s.MatchExpressions = append(s.MatchExpressions, expression)
	return s
}

// NotIn injects a matchExpression (with an operator NotIn) to the inner labelSelector.
func (s *LabelSelectorWrapper) NotIn(key string, vals []string) *LabelSelectorWrapper {
	expression := metav1.LabelSelectorRequirement{
		Key:      key,
		Operator: metav1.LabelSelectorOpNotIn,
		Values:   vals,
	}
	s.MatchExpressions = append(s.MatchExpressions, expression)
	return s
}

// Exists injects a matchExpression (with an operator Exists) to the inner labelSelector.
func (s *LabelSelectorWrapper) Exists(k string) *LabelSelectorWrapper {
	expression := metav1.LabelSelectorRequirement{
		Key:      k,
		Operator: metav1.LabelSelectorOpExists,
	}
	s.MatchExpressions = append(s.MatchExpressions, expression)
	return s
}

// NotExist injects a matchExpression (with an operator NotExist) to the inner labelSelector.
func (s *LabelSelectorWrapper) NotExist(k string) *LabelSelectorWrapper {
	expression := metav1.LabelSelectorRequirement{
		Key:      k,
		Operator: metav1.LabelSelectorOpDoesNotExist,
	}
	s.MatchExpressions = append(s.MatchExpressions, expression)
	return s
}

// Obj returns the inner LabelSelector.
func (s *LabelSelectorWrapper) Obj() *metav1.LabelSelector {
	return &s.LabelSelector
}

// PodWrapper wraps a Pod inside.
type PodWrapper struct{ v1.Pod }

// MakePod creates a Pod wrapper.
func MakePod() *PodWrapper {
	return &PodWrapper{v1.Pod{}}
}

// Obj returns the inner Pod.
func (p *PodWrapper) Obj() *v1.Pod {
	return &p.Pod
}

// Name sets `s` as the name of the inner pod.
func (p *PodWrapper) Name(s string) *PodWrapper {
	p.SetName(s)
	return p
}

// UID sets `s` as the UID of the inner pod.
func (p *PodWrapper) UID(s string) *PodWrapper {
	p.SetUID(types.UID(s))
	return p
}

// SchedulerName sets `s` as the scheduler name of the inner pod.
func (p *PodWrapper) SchedulerName(s string) *PodWrapper {
	p.Spec.SchedulerName = s
	return p
}

// Namespace sets `s` as the namespace of the inner pod.
func (p *PodWrapper) Namespace(s string) *PodWrapper {
	p.SetNamespace(s)
	return p
}

// OwnerReference updates the owning controller of the pod.
func (p *PodWrapper) OwnerReference(name string, gvk schema.GroupVersionKind) *PodWrapper {
	p.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion: gvk.GroupVersion().String(),
			Kind:       gvk.Kind,
			Name:       name,
			Controller: pointer.BoolPtr(true),
		},
	}
	return p
}

// Container appends a container into PodSpec of the inner pod.
func (p *PodWrapper) Container(s string) *PodWrapper {
	p.Spec.Containers = append(p.Spec.Containers, v1.Container{
		Name:  fmt.Sprintf("con%d", len(p.Spec.Containers)),
		Image: s,
	})
	return p
}

// Priority sets a priority value into PodSpec of the inner pod.
func (p *PodWrapper) Priority(val int32) *PodWrapper {
	p.Spec.Priority = &val
	return p
}

// Terminating sets the inner pod's deletionTimestamp to current timestamp.
func (p *PodWrapper) Terminating() *PodWrapper {
	now := metav1.Now()
	p.DeletionTimestamp = &now
	return p
}

// ZeroTerminationGracePeriod sets the TerminationGracePeriodSeconds of the inner pod to zero.
func (p *PodWrapper) ZeroTerminationGracePeriod() *PodWrapper {
	p.Spec.TerminationGracePeriodSeconds = &zero
	return p
}

// Node sets `s` as the nodeName of the inner pod.
func (p *PodWrapper) Node(s string) *PodWrapper {
	p.Spec.NodeName = s
	return p
}

// NodeSelector sets `m` as the nodeSelector of the inner pod.
func (p *PodWrapper) NodeSelector(m map[string]string) *PodWrapper {
	p.Spec.NodeSelector = m
	return p
}

// NodeAffinityIn creates a HARD node affinity (with the operator In)
// and injects into the inner pod.
func (p *PodWrapper) NodeAffinityIn(key string, vals []string) *PodWrapper {
	if p.Spec.Affinity == nil {
		p.Spec.Affinity = &v1.Affinity{}
	}
	if p.Spec.Affinity.NodeAffinity == nil {
		p.Spec.Affinity.NodeAffinity = &v1.NodeAffinity{}
	}
	nodeSelector := MakeNodeSelector().In(key, vals).Obj()
	p.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = nodeSelector
	return p
}

// NodeAffinityNotIn creates a HARD node affinity (with the operator NotIn)
// and injects into the inner pod.
func (p *PodWrapper) NodeAffinityNotIn(key string, vals []string) *PodWrapper {
	if p.Spec.Affinity == nil {
		p.Spec.Affinity = &v1.Affinity{}
	}
	if p.Spec.Affinity.NodeAffinity == nil {
		p.Spec.Affinity.NodeAffinity = &v1.NodeAffinity{}
	}
	nodeSelector := MakeNodeSelector().NotIn(key, vals).Obj()
	p.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution = nodeSelector
	return p
}

// StartTime sets `t` as .status.startTime for the inner pod.
func (p *PodWrapper) StartTime(t metav1.Time) *PodWrapper {
	p.Status.StartTime = &t
	return p
}

// NominatedNodeName sets `n` as the .Status.NominatedNodeName of the inner pod.
func (p *PodWrapper) NominatedNodeName(n string) *PodWrapper {
	p.Status.NominatedNodeName = n
	return p
}

// Toleration creates a toleration (with the operator Exists)
// and injects into the inner pod.
func (p *PodWrapper) Toleration(key string) *PodWrapper {
	p.Spec.Tolerations = append(p.Spec.Tolerations, v1.Toleration{
		Key:      key,
		Operator: v1.TolerationOpExists,
	})
	return p
}

// HostPort creates a container with a hostPort valued `hostPort`,
// and injects into the inner pod.
func (p *PodWrapper) HostPort(port int32) *PodWrapper {
	p.Spec.Containers = append(p.Spec.Containers, v1.Container{
		Ports: []v1.ContainerPort{{HostPort: port}},
	})
	return p
}

// PVC creates a Volume with a PVC and injects into the inner pod.
func (p *PodWrapper) PVC(name string) *PodWrapper {
	p.Spec.Volumes = append(p.Spec.Volumes, v1.Volume{
		Name: name,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{ClaimName: name},
		},
	})
	return p
}

// PodAffinityKind represents different kinds of PodAffinity.
type PodAffinityKind int

const (
	// NilPodAffinity is a no-op which doesn't apply any PodAffinity.
	NilPodAffinity PodAffinityKind = iota
	// PodAffinityWithRequiredReq applies a HARD requirement to pod.spec.affinity.PodAffinity.
	PodAffinityWithRequiredReq
	// PodAffinityWithPreferredReq applies a SOFT requirement to pod.spec.affinity.PodAffinity.
	PodAffinityWithPreferredReq
	// PodAffinityWithRequiredPreferredReq applies HARD and SOFT requirements to pod.spec.affinity.PodAffinity.
	PodAffinityWithRequiredPreferredReq
	// PodAntiAffinityWithRequiredReq applies a HARD requirement to pod.spec.affinity.PodAntiAffinity.
	PodAntiAffinityWithRequiredReq
	// PodAntiAffinityWithPreferredReq applies a SOFT requirement to pod.spec.affinity.PodAntiAffinity.
	PodAntiAffinityWithPreferredReq
	// PodAntiAffinityWithRequiredPreferredReq applies HARD and SOFT requirements to pod.spec.affinity.PodAntiAffinity.
	PodAntiAffinityWithRequiredPreferredReq
)

// PodAffinityExists creates an PodAffinity with the operator "Exists"
// and injects into the inner pod.
func (p *PodWrapper) PodAffinityExists(labelKey, topologyKey string, kind PodAffinityKind) *PodWrapper {
	if kind == NilPodAffinity {
		return p
	}

	if p.Spec.Affinity == nil {
		p.Spec.Affinity = &v1.Affinity{}
	}
	if p.Spec.Affinity.PodAffinity == nil {
		p.Spec.Affinity.PodAffinity = &v1.PodAffinity{}
	}
	labelSelector := MakeLabelSelector().Exists(labelKey).Obj()
	term := v1.PodAffinityTerm{LabelSelector: labelSelector, TopologyKey: topologyKey}
	switch kind {
	case PodAffinityWithRequiredReq:
		p.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
			term,
		)
	case PodAffinityWithPreferredReq:
		p.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			v1.WeightedPodAffinityTerm{Weight: 1, PodAffinityTerm: term},
		)
	case PodAffinityWithRequiredPreferredReq:
		p.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
			term,
		)
		p.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			v1.WeightedPodAffinityTerm{Weight: 1, PodAffinityTerm: term},
		)
	}
	return p
}

// PodAntiAffinityExists creates an PodAntiAffinity with the operator "Exists"
// and injects into the inner pod.
func (p *PodWrapper) PodAntiAffinityExists(labelKey, topologyKey string, kind PodAffinityKind) *PodWrapper {
	if kind == NilPodAffinity {
		return p
	}

	if p.Spec.Affinity == nil {
		p.Spec.Affinity = &v1.Affinity{}
	}
	if p.Spec.Affinity.PodAntiAffinity == nil {
		p.Spec.Affinity.PodAntiAffinity = &v1.PodAntiAffinity{}
	}
	labelSelector := MakeLabelSelector().Exists(labelKey).Obj()
	term := v1.PodAffinityTerm{LabelSelector: labelSelector, TopologyKey: topologyKey}
	switch kind {
	case PodAntiAffinityWithRequiredReq:
		p.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
			term,
		)
	case PodAntiAffinityWithPreferredReq:
		p.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			v1.WeightedPodAffinityTerm{Weight: 1, PodAffinityTerm: term},
		)
	case PodAntiAffinityWithRequiredPreferredReq:
		p.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution,
			term,
		)
		p.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution = append(
			p.Spec.Affinity.PodAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution,
			v1.WeightedPodAffinityTerm{Weight: 1, PodAffinityTerm: term},
		)
	}
	return p
}

// SpreadConstraint constructs a TopologySpreadConstraint object and injects
// into the inner pod.
func (p *PodWrapper) SpreadConstraint(maxSkew int, tpKey string, mode v1.UnsatisfiableConstraintAction, selector *metav1.LabelSelector) *PodWrapper {
	c := v1.TopologySpreadConstraint{
		MaxSkew:           int32(maxSkew),
		TopologyKey:       tpKey,
		WhenUnsatisfiable: mode,
		LabelSelector:     selector,
	}
	p.Spec.TopologySpreadConstraints = append(p.Spec.TopologySpreadConstraints, c)
	return p
}

// Label sets a {k,v} pair to the inner pod.
func (p *PodWrapper) Label(k, v string) *PodWrapper {
	if p.Labels == nil {
		p.Labels = make(map[string]string)
	}
	p.Labels[k] = v
	return p
}

// Req adds a new container to the inner pod with given resource map.
func (p *PodWrapper) Req(resMap map[v1.ResourceName]string) *PodWrapper {
	if len(resMap) == 0 {
		return p
	}

	res := v1.ResourceList{}
	for k, v := range resMap {
		res[k] = resource.MustParse(v)
	}
	p.Spec.Containers = append(p.Spec.Containers, v1.Container{
		Name:  fmt.Sprintf("con%d", len(p.Spec.Containers)),
		Image: imageutils.GetPauseImageName(),
		Resources: v1.ResourceRequirements{
			Requests: res,
		},
	})
	return p
}

// PreemptionPolicy sets the give preemption policy to the inner pod.
func (p *PodWrapper) PreemptionPolicy(policy v1.PreemptionPolicy) *PodWrapper {
	p.Spec.PreemptionPolicy = &policy
	return p
}

// NodeWrapper wraps a Node inside.
type NodeWrapper struct{ v1.Node }

// MakeNode creates a Node wrapper.
func MakeNode() *NodeWrapper {
	w := &NodeWrapper{v1.Node{}}
	return w.Capacity(nil)
}

// Obj returns the inner Node.
func (n *NodeWrapper) Obj() *v1.Node {
	return &n.Node
}

// Name sets `s` as the name of the inner pod.
func (n *NodeWrapper) Name(s string) *NodeWrapper {
	n.SetName(s)
	return n
}

// UID sets `s` as the UID of the inner pod.
func (n *NodeWrapper) UID(s string) *NodeWrapper {
	n.SetUID(types.UID(s))
	return n
}

// Label applies a {k,v} label pair to the inner node.
func (n *NodeWrapper) Label(k, v string) *NodeWrapper {
	if n.Labels == nil {
		n.Labels = make(map[string]string)
	}
	n.Labels[k] = v
	return n
}

// Capacity sets the capacity and the allocatable resources of the inner node.
// Each entry in `resources` corresponds to a resource name and its quantity.
// By default, the capacity and allocatable number of pods are set to 32.
func (n *NodeWrapper) Capacity(resources map[v1.ResourceName]string) *NodeWrapper {
	res := v1.ResourceList{
		v1.ResourcePods: resource.MustParse("32"),
	}
	for name, value := range resources {
		res[name] = resource.MustParse(value)
	}
	n.Status.Capacity, n.Status.Allocatable = res, res
	return n
}

// Images sets the images of the inner node. Each entry in `images` corresponds
// to an image name and its size in bytes.
func (n *NodeWrapper) Images(images map[string]int64) *NodeWrapper {
	var containerImages []v1.ContainerImage
	for name, size := range images {
		containerImages = append(containerImages, v1.ContainerImage{Names: []string{name}, SizeBytes: size})
	}
	n.Status.Images = containerImages
	return n
}
