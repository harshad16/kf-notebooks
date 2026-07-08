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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/kubeflow/notebooks/workspaces/backend/api/constants"
)

const nonAdminUser = "non-admin-user"

var _ = Describe("User Handler", func() {

	BeforeEach(func() {
		a.Config.UserIdHeader = userIdHeader
		a.Config.DisableAuth = false
	})

	It("should report a cluster-admin user as clusterAdmin", func() {
		req, err := http.NewRequest(http.MethodGet, constants.UserPath, http.NoBody)
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set(userIdHeader, adminUser)

		rr := httptest.NewRecorder()
		a.GetUserHandler(rr, req, httprouter.Params{})
		rs := rr.Result()
		defer rs.Body.Close()

		Expect(rs.StatusCode).To(Equal(http.StatusOK), descUnexpectedHTTPStatus, rr.Body.String())

		body, err := io.ReadAll(rs.Body)
		Expect(err).NotTo(HaveOccurred())

		var response UserEnvelope
		Expect(json.Unmarshal(body, &response)).To(Succeed())
		Expect(response.Data.UserId).To(Equal(adminUser))
		Expect(response.Data.ClusterAdmin).To(BeTrue())
	})

	It("should report a non-admin user as not clusterAdmin", func() {
		req, err := http.NewRequest(http.MethodGet, constants.UserPath, http.NoBody)
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set(userIdHeader, nonAdminUser)

		rr := httptest.NewRecorder()
		a.GetUserHandler(rr, req, httprouter.Params{})
		rs := rr.Result()
		defer rs.Body.Close()

		Expect(rs.StatusCode).To(Equal(http.StatusOK), descUnexpectedHTTPStatus, rr.Body.String())

		body, err := io.ReadAll(rs.Body)
		Expect(err).NotTo(HaveOccurred())

		var response UserEnvelope
		Expect(json.Unmarshal(body, &response)).To(Succeed())
		Expect(response.Data.UserId).To(Equal(nonAdminUser))
		Expect(response.Data.ClusterAdmin).To(BeFalse())
	})

	It("should fail closed when auth is disabled", func() {
		a.Config.DisableAuth = true

		req, err := http.NewRequest(http.MethodGet, constants.UserPath, http.NoBody)
		Expect(err).NotTo(HaveOccurred())
		req.Header.Set(userIdHeader, adminUser)

		rr := httptest.NewRecorder()
		a.GetUserHandler(rr, req, httprouter.Params{})
		rs := rr.Result()
		defer rs.Body.Close()

		Expect(rs.StatusCode).To(Equal(http.StatusOK), descUnexpectedHTTPStatus, rr.Body.String())

		body, err := io.ReadAll(rs.Body)
		Expect(err).NotTo(HaveOccurred())

		var response UserEnvelope
		Expect(json.Unmarshal(body, &response)).To(Succeed())
		Expect(response.Data.UserId).To(Equal(adminUser))
		Expect(response.Data.ClusterAdmin).To(BeFalse())
	})
})
