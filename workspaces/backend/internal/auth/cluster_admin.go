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

package auth

import (
	"context"
	"fmt"

	"k8s.io/apiserver/pkg/authentication/user"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

// IsClusterAdmin checks whether the user has cluster-admin-equivalent permissions
// by performing a SubjectAccessReview with wildcard verb and resource, matching
// the SelfSubjectAccessReview pattern used by other ODH mod-arch BFFs.
func IsClusterAdmin(ctx context.Context, authZ authorizer.Authorizer, u user.Info) (bool, error) {
	if authZ == nil {
		return false, fmt.Errorf("authorizer is not configured")
	}
	if u == nil || u.GetName() == "" {
		return false, nil
	}

	attributes := authorizer.AttributesRecord{
		User:            u,
		Verb:            "*",
		Resource:        "*",
		ResourceRequest: true,
	}

	decision, _, err := authZ.Authorize(ctx, attributes)
	if err != nil {
		return false, fmt.Errorf("failed to verify cluster-admin permissions: %w", err)
	}

	return decision == authorizer.DecisionAllow, nil
}
