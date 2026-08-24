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
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestAdjustSecurityContextForOpenShift(t *testing.T) {
	podSC := &corev1.PodSecurityContext{
		FSGroup: ptr.To(int64(100)),
	}
	containerSC := &corev1.SecurityContext{
		AllowPrivilegeEscalation: ptr.To(false),
		RunAsNonRoot:             ptr.To(true),
		RunAsUser:                ptr.To(int64(100)),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
	}

	AdjustSecurityContextForOpenShift(&podSC, &containerSC)

	if podSC.FSGroup != nil {
		t.Fatalf("expected fsGroup to be cleared, got %v", *podSC.FSGroup)
	}
	if podSC.RunAsUser != nil {
		t.Fatalf("expected pod runAsUser to be cleared, got %v", *podSC.RunAsUser)
	}
	if podSC.SeccompProfile == nil || podSC.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Fatalf("expected seccompProfile RuntimeDefault, got %#v", podSC.SeccompProfile)
	}
	if containerSC.RunAsUser != nil {
		t.Fatalf("expected container runAsUser to be cleared, got %v", *containerSC.RunAsUser)
	}
	if containerSC.RunAsNonRoot == nil || !*containerSC.RunAsNonRoot {
		t.Fatalf("expected runAsNonRoot to be preserved")
	}
}

func TestAdjustSecurityContextForOpenShiftNilInputs(t *testing.T) {
	var podSC *corev1.PodSecurityContext
	var containerSC *corev1.SecurityContext

	AdjustSecurityContextForOpenShift(&podSC, &containerSC)

	if podSC == nil {
		t.Fatal("expected pod security context to be created")
	}
	if podSC.SeccompProfile == nil || podSC.SeccompProfile.Type != corev1.SeccompProfileTypeRuntimeDefault {
		t.Fatalf("expected seccompProfile RuntimeDefault, got %#v", podSC.SeccompProfile)
	}
	if containerSC == nil {
		t.Fatal("expected container security context to be created")
	}
}
