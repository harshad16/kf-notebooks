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

package config

import (
	"os"
	"strconv"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

// ResolveOpenShift returns whether the controller should apply OpenShift-specific
// workload adjustments. The OPENSHIFT env var overrides auto-detection when set.
func ResolveOpenShift(restConfig *rest.Config) bool {
	if value, exists := os.LookupEnv("OPENSHIFT"); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return detectOpenShift(restConfig)
}

func detectOpenShift(restConfig *rest.Config) bool {
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return false
	}
	apiGroups, err := discoveryClient.ServerGroups()
	if err != nil {
		return false
	}
	for _, group := range apiGroups.Groups {
		if group.Name == "route.openshift.io" {
			return true
		}
	}
	return false
}
