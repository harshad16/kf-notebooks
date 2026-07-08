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

package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/kubeflow/notebooks/workspaces/backend/internal/auth"
)

// UserResponse represents the user settings response
type UserResponse struct {
	UserId       string `json:"userId"`
	ClusterAdmin bool   `json:"clusterAdmin"`
}

// UserEnvelope wraps the UserResponse in the standard envelope format
type UserEnvelope = Envelope[*UserResponse]

// GetUserHandler returns the current user settings
//
//	@Summary		Get user settings
//	@Description	Returns the current user's settings including user ID and admin status
//	@Tags			user
//	@Produce		json
//	@Success		200	{object}	UserEnvelope
//	@Failure		500	{object}	ErrorEnvelope
//	@Router			/user [get]
func (a *App) GetUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	userId := a.resolveUserId(r)

	clusterAdmin, err := a.resolveClusterAdmin(r)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	response := UserEnvelope{
		Data: &UserResponse{
			UserId:       userId,
			ClusterAdmin: clusterAdmin,
		},
	}

	a.dataResponse(w, r, &response)
}

func (a *App) resolveUserId(r *http.Request) string {
	userId := r.Header.Get(a.Config.UserIdHeader)

	if a.Config.UserIdPrefix != "" {
		userId = strings.TrimPrefix(userId, a.Config.UserIdPrefix)
	}

	if userId == "" {
		userId = r.Header.Get("X-Auth-Request-User")
	}
	if userId == "" {
		userId = r.Header.Get("kubeflow-userid")
	}
	if userId == "" {
		userId = "anonymous"
	}

	return userId
}

func (a *App) resolveClusterAdmin(r *http.Request) (bool, error) {
	if a.Config.DisableAuth {
		return false, nil
	}

	res, ok, err := a.RequestAuthN.AuthenticateRequest(r)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return auth.IsClusterAdmin(r.Context(), a.RequestAuthZ, res.User)
}
