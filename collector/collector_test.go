/*
http://www.apache.org/licenses/LICENSE-2.0.txt
Copyright 2016 Intel Corporation
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

package collector

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/rackspace/gophercloud/testhelper"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/intelsdi-x/snap/core/ctypes"

	"github.com/intelsdi-x/snap-plugin-utilities/str"
)

type CollectorSuite struct {
	suite.Suite
	Token string
}

func (s *CollectorSuite) SetupSuite() {
	th.SetupHTTP()
	registerRoot()
	registerAuthentication(s)
	registerTenants(s)
	registerUsers(s)
	registerServices(s)
	registerEndpoints(s)
	registerAdminUsers(s)
	registerDemoUsers(s)
}

func (suite *CollectorSuite) TearDownSuite() {
	th.TeardownHTTP()
}

func TestRunSuite(t *testing.T) {
	keystoneTestSuite := new(CollectorSuite)
	suite.Run(t, keystoneTestSuite)
}

func (s *CollectorSuite) TestGetMetricTypes() {
	Convey("Given config with enpoint, user and password defined", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "me", "secret", "admin")

		Convey("When GetMetricTypes() is called", func() {
			collector := New()
			mts, err := collector.GetMetricTypes(cfg)

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := []string{}
				for _, m := range mts {
					metricNames = append(metricNames, m.Namespace().String())
				}

				So(len(mts), ShouldEqual, 6)
				So(str.Contains(metricNames, "/intel/openstack/keystone/demo/users_count"), ShouldBeTrue)
				So(str.Contains(metricNames, "/intel/openstack/keystone/admin/users_count"), ShouldBeTrue)
				So(str.Contains(metricNames, "/intel/openstack/keystone/total_tenants_count"), ShouldBeTrue)
				So(str.Contains(metricNames, "/intel/openstack/keystone/total_users_count"), ShouldBeTrue)
				So(str.Contains(metricNames, "/intel/openstack/keystone/total_endpoints_count"), ShouldBeTrue)
				So(str.Contains(metricNames, "/intel/openstack/keystone/total_services_count"), ShouldBeTrue)
			})
		})
	})
}

func (s *CollectorSuite) TestCollectMetrics() {
	Convey("Given set of metric types", s.T(), func() {
		cfg := setupCfg(th.Endpoint(), "me", "secret", "admin")
		m1 := plugin.MetricType{
			Namespace_: core.NewNamespace("intel", "openstack", "keystone", "demo", "users_count"),
			Config_:    cfg.ConfigDataNode}
		m2 := plugin.MetricType{
			Namespace_: core.NewNamespace("intel", "openstack", "keystone", "total_services_count"),
			Config_:    cfg.ConfigDataNode}

		Convey("When ColelctMetrics() is called", func() {
			collector := New()

			mts, err := collector.CollectMetrics([]plugin.MetricType{m1, m2})

			Convey("Then no error should be reported", func() {
				So(err, ShouldBeNil)
			})

			Convey("and proper metric types are returned", func() {
				metricNames := map[string]interface{}{}
				for _, m := range mts {
					ns := m.Namespace().String()
					metricNames[ns] = m.Data()
				}
				fmt.Println(metricNames)
				So(len(mts), ShouldEqual, 2)

				val, ok := metricNames["/intel/openstack/keystone/demo/users_count"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 3)

				val, ok = metricNames["/intel/openstack/keystone/total_services_count"]
				So(ok, ShouldBeTrue)
				So(val, ShouldEqual, 4)
			})
		})
	})
}

func setupCfg(endpoint, user, password, tenant string) plugin.ConfigType {
	node := cdata.NewNode()
	node.AddItem("admin_endpoint", ctypes.ConfigValueStr{Value: endpoint})
	node.AddItem("admin_user", ctypes.ConfigValueStr{Value: user})
	node.AddItem("admin_password", ctypes.ConfigValueStr{Value: password})
	node.AddItem("admin_tenant", ctypes.ConfigValueStr{Value: tenant})
	return plugin.ConfigType{ConfigDataNode: node}
}

func registerRoot() {
	th.Mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"versions": {
						"values": [
							{
								"status": "experimental",
								"id": "v3.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							},
							{
								"status": "stable",
								"id": "v2.0",
								"links": [
									{ "href": "%s", "rel": "self" }
								]
							}
						]
					}
				}
				`, th.Endpoint()+"v3/", th.Endpoint()+"v2.0/")
	})
}

func registerAuthentication(s *CollectorSuite) {
	s.Token = "2ed210f132564f21b178afb197ee99e3"
	th.Mux.HandleFunc("/v2.0/tokens", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
				{
					"access": {
						"metadata": {
							"is_admin": 0,
							"roles": [
								"3083d61996d648ca88d6ff420542f324"
							]
						},
						"serviceCatalog": [],
						"token": {
							"expires": "2016-02-21T14:28:30Z",
							"id": "%s",
							"issued_at": "2016-02-21T13:28:30.656527",
							"tenant": {
								"description": null,
								"enabled": true,
								"id": "97ea299c37bb4e04b3779039ea8aba44",
								"name": "tenant"
							}
						}
					}
				}
			`, s.Token)
	})
}

func registerTenants(s *CollectorSuite) {
	th.Mux.HandleFunc("/v2.0/tenants", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"tenants": [
					{
						"description": "Test tenat",
						"enabled": true,
						"id": "11111",
						"name": "demo"
					},
					{
						"description": "admin tenant",
						"enabled": true,
						"id": "22222",
						"name": "admin"
					}
				],
				"tenants_links": []
			}
		`)
	})
}

func registerUsers(s *CollectorSuite) {
	th.Mux.HandleFunc("/v2.0/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"users": [
					{
						"email": "heat@localhost",
						"enabled": true,
						"id": "27b6b98022314a6b9c4524efaedf4694",
						"name": "heat",
						"username": "heat"
					},
					{
						"email": "heat-cfn@localhost",
						"enabled": true,
						"id": "60251a9059f84770acbd037468f2e414",
						"name": "heat-cfn",
						"username": "heat-cfn"
					},
					{
						"email": "cinder@localhost",
						"enabled": true,
						"id": "659a62b0da35495e85b08b11e5b6f092",
						"name": "cinder",
						"username": "cinder"
					}
				]
			}
	`)
	})
}

func registerServices(s *CollectorSuite) {
	th.Mux.HandleFunc("/v3/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			{
				"links": {
					"next": null,
					"previous": null,
					"self": "https://public.fuel.local:5000/v3/services"
				},
				"services": [
					{
						"description": "Openstack Orchestration Service",
						"enabled": true,
						"id": "13c6403db2cd4b029403740e30e43d88",
						"links": {
							"self": "https://public.fuel.local:5000/v3/services/13c6403db2cd4b029403740e30e43d88"
						},
						"name": "heat",
						"type": "orchestration"
					},
					{
						"description": "Openstack Cloudformation Service",
						"enabled": true,
						"id": "361c1c46035a414fb6024a7a5b8cabfb",
						"links": {
							"self": "https://public.fuel.local:5000/v3/services/361c1c46035a414fb6024a7a5b8cabfb"
						},
						"name": "heat-cfn",
						"type": "cloudformation"
					},
					{
						"description": "Openstack Metering Service",
						"enabled": true,
						"id": "615e06490105462cbbbab919bbe1c725",
						"links": {
							"self": "https://public.fuel.local:5000/v3/services/615e06490105462cbbbab919bbe1c725"
						},
						"name": "ceilometer",
						"type": "metering"
					},
					{
						"description": "Openstack Compute Service v3",
						"enabled": true,
						"id": "79b1d028220f47a5b0de7756f3a5b286",
						"links": {
							"self": "https://public.fuel.local:5000/v3/services/79b1d028220f47a5b0de7756f3a5b286"
						},
						"name": "novav3",
						"type": "computev3"
					}
				]
			}
		`)
	})
}

func registerEndpoints(s *CollectorSuite) {
	th.Mux.HandleFunc("/v3/endpoints", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
			{
				"endpoints": [
					{
						"enabled": true,
						"id": "035d8ea06d88496e929310a7cda173a6",
						"interface": "public",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/035d8ea06d88496e929310a7cda173a6"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "dc52eef88c2d470fb68912a2641eaab4",
						"url": "https://public.fuel.local:8773/services/Cloud"
					},
					{
						"enabled": true,
						"id": "0a587d5392a54d4e8eb8d7328a7acbf1",
						"interface": "public",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/0a587d5392a54d4e8eb8d7328a7acbf1"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "615e06490105462cbbbab919bbe1c725",
						"url": "https://public.fuel.local:8777"
					},
					{
						"enabled": true,
						"id": "0b74729c70814c5cb65be5e7b56d56ed",
						"interface": "admin",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/0b74729c70814c5cb65be5e7b56d56ed"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "efbf568dd1234f52a73869c8cab10d93",
						"url": "http://192.168.20.2:9696"
					},
					{
						"enabled": true,
						"id": "159572c2ce0d480db916fc4986812d0b",
						"interface": "internal",
						"links": {
							"self": "https://public.fuel.local:5000/v3/endpoints/159572c2ce0d480db916fc4986812d0b"
						},
						"region": "RegionOne",
						"region_id": "RegionOne",
						"service_id": "efbf568dd1234f52a73869c8cab10d93",
						"url": "http://192.168.20.2:9696"
					}
				]
			}
		`)
	})
}

func registerDemoUsers(s *CollectorSuite) {
	th.Mux.HandleFunc("/v2.0/tenants/11111/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"users": [
					{
						"email": "heat@localhost",
						"enabled": true,
						"id": "27b6b98022314a6b9c4524efaedf4694",
						"name": "heat",
						"username": "heat"
					},
					{
						"email": "heat-cfn@localhost",
						"enabled": true,
						"id": "60251a9059f84770acbd037468f2e414",
						"name": "heat-cfn",
						"username": "heat-cfn"
					},
					{
						"email": "cinder@localhost",
						"enabled": true,
						"id": "659a62b0da35495e85b08b11e5b6f092",
						"name": "cinder",
						"username": "cinder"
					}
				]
			}
	`)
	})
}

func registerAdminUsers(s *CollectorSuite) {
	th.Mux.HandleFunc("/v2.0/tenants/22222/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(s.T(), r, "GET")
		th.TestHeader(s.T(), r, "X-Auth-Token", s.Token)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
			{
				"users": [
					{
						"email": "heat@localhost",
						"enabled": true,
						"id": "27b6b98022314a6b9c4524efaedf4694",
						"name": "heat",
						"username": "heat"
					}
				]
			}
	`)
	})
}
