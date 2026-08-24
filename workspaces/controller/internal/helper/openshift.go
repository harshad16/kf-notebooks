/*
Copyright 2024.

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

package helper

import (
	corev1 "k8s.io/api/core/v1"
)

// AdjustSecurityContextForOpenShift updates pod and container security contexts so
// workspace pods comply with OpenShift restricted SCCs (restricted-v2 / restricted-v3).
//
// WorkspaceKind manifests often pin low UIDs (e.g. fsGroup: 100) that are incompatible
// with namespace-assigned UID ranges on OpenShift. Omitting those fields lets the SCC
// admission controller assign values from the namespace range.
func AdjustSecurityContextForOpenShift(podSecurityContext **corev1.PodSecurityContext, containerSecurityContext **corev1.SecurityContext) {
	podSC := &corev1.PodSecurityContext{}
	if *podSecurityContext != nil {
		podSC = (*podSecurityContext).DeepCopy()
	}

	podSC.FSGroup = nil
	podSC.RunAsUser = nil
	podSC.RunAsGroup = nil
	podSC.SupplementalGroups = nil
	podSC.SupplementalGroupsPolicy = nil
	if podSC.SeccompProfile == nil {
		podSC.SeccompProfile = &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		}
	}

	*podSecurityContext = podSC

	containerSC := &corev1.SecurityContext{}
	if *containerSecurityContext != nil {
		containerSC = (*containerSecurityContext).DeepCopy()
	}

	containerSC.RunAsUser = nil
	containerSC.RunAsGroup = nil

	*containerSecurityContext = containerSC
}
