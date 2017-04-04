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
	"sync"
	"time"

	"github.com/rackspace/gophercloud"

	"github.com/intelsdi-x/snap/control/plugin"
	"github.com/intelsdi-x/snap/control/plugin/cpolicy"

	"github.com/intelsdi-x/snap-plugin-utilities/config"
	"github.com/intelsdi-x/snap-plugin-utilities/str"

	openstackintel "github.com/intelsdi-x/snap-plugin-collector-keystone/openstack"
	"github.com/intelsdi-x/snap-plugin-collector-keystone/types"
	"github.com/intelsdi-x/snap/core"
)

const (
	name    = "keystone"
	version = 3
	plgtype = plugin.CollectorPluginType
	vendor  = "intel"
	fs      = "openstack"
)

var keystoneMetrics = []string{
	"total_tenants_count",
	"total_users_count",
	"total_services_count",
	"total_endpoints_count",
}

// New creates initialized instance of Glance collector
func New() *collector {
	return &collector{}
}

// GetMetricTypes returns list of available metric types
// It returns error in case retrieval was not successful
func (c *collector) GetMetricTypes(cfg plugin.ConfigType) ([]plugin.MetricType, error) {
	mts := []plugin.MetricType{}
	domain_name := ""
	domain_id := ""
	items, err := config.GetConfigItems(cfg, "admin_endpoint", "admin_user", "admin_password", "admin_tenant")
	if err != nil {
		return nil, err
	}

	endpoint := items["admin_endpoint"].(string)
	user := items["admin_user"].(string)
	password := items["admin_password"].(string)
	tenant := items["admin_tenant"].(string)
	dom_name, _ := config.GetConfigItem(cfg, "domain_name")
	dom_id, _ := config.GetConfigItem(cfg, "domain_id")
	if dom_name != nil {
		domain_name = dom_name.(string)
	}
	if dom_id != nil {
		domain_id = dom_id.(string)
	}

	if c.provider == nil {
		c.provider, err = openstackintel.Authenticate(endpoint, user, password, tenant, domain_name, domain_id)
		if err != nil {
			return nil, err
		}
	}

	// retrieve list of all available tenants for provided endpoint, user and password
	allTenants, err := openstackintel.GetAllTenants(c.provider)
	if err != nil {
		return nil, err
	}

	// Generate available namespace from tenants (user counts per tenant)
	for _, tenant := range allTenants {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace(vendor, fs, name, tenant.Name, "users_count"),
			Config_:    cfg.ConfigDataNode,
		})
	}

	// Generate available namespace from keystone metrics
	for _, keystoneMetric := range keystoneMetrics {
		mts = append(mts, plugin.MetricType{
			Namespace_: core.NewNamespace(vendor, fs, name, keystoneMetric),
			Config_:    cfg.ConfigDataNode,
		})
	}
	return mts, nil
}

// CollectMetrics returns list of requested metric values
// It returns error in case retrieval was not successful
func (c *collector) CollectMetrics(metricTypes []plugin.MetricType) ([]plugin.MetricType, error) {
	domain_name := ""
	domain_id := ""
	items, err := config.GetConfigItems(metricTypes[0], "admin_endpoint", "admin_user", "admin_password", "admin_tenant")
	if err != nil {
		return nil, err
	}

	endpoint := items["admin_endpoint"].(string)
	user := items["admin_user"].(string)
	password := items["admin_password"].(string)
	tenant := items["admin_tenant"].(string)
	dom_name, _ := config.GetConfigItem(metricTypes[0], "domain_name")
	dom_id, _ := config.GetConfigItem(metricTypes[0], "domain_id")
	if dom_name != nil {
		domain_name = dom_name.(string)
	}
	if dom_id != nil {
		domain_id = dom_id.(string)
	}

	if c.provider == nil {
		c.provider, err = openstackintel.Authenticate(endpoint, user, password, tenant, domain_name, domain_id)
		if err != nil {
			return nil, err
		}
	}

	var done sync.WaitGroup
	errCh := make(chan error, 4)

	// collect services and endpoint only once
	if c.endpoints == nil {
		done.Add(1)
		go func() {
			var err error
			if c.endpoints, err = openstackintel.GetAllEndpoints(c.provider); err != nil {
				errCh <- err
			}
			done.Done()
		}()
	}
	if c.services == nil {
		done.Add(1)
		go func() {
			var err error
			if c.services, err = openstackintel.GetAllServices(c.provider); err != nil {
				errCh <- err
			}
			done.Done()
		}()
	}

	done.Add(2)
	tenantList := []types.Tenant{}
	go func() {
		var err error
		if tenantList, err = openstackintel.GetAllTenants(c.provider); err != nil {
			errCh <- err
		}
		done.Done()
	}()

	userList := []types.User{}
	go func() {
		var err error
		if userList, err = openstackintel.GetAllUsers(c.provider); err != nil {
			errCh <- err
		}
		done.Done()
	}()

	done.Wait()
	close(errCh)

	if err = <-errCh; err != nil {
		return nil, err
	}

	tenantUsers, err := openstackintel.GetUsersPerTenant(c.provider, tenantList)
	if err != nil {
		return nil, err
	}

	metrics := []plugin.MetricType{}
	for _, metricType := range metricTypes {
		namespace := metricType.Namespace().Strings()
		metric := plugin.MetricType{
			Timestamp_: time.Now(),
			Namespace_: metricType.Namespace(),
		}

		if str.Contains(keystoneMetrics, namespace[3]) {
			switch namespace[3] {
			case "total_tenants_count":
				metric.Data_ = len(tenantList)
			case "total_users_count":
				metric.Data_ = len(userList)
			case "total_services_count":
				metric.Data_ = len(c.services)
			case "total_endpoints_count":
				metric.Data_ = len(c.endpoints)
			}
		} else {
			tenantName := namespace[3]
			val, ok := tenantUsers[tenantName]
			if ok {
				metric.Data_ = val
			}
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

// GetConfigPolicy returns config policy
// It returns error in case retrieval was not successful
func (c *collector) GetConfigPolicy() (*cpolicy.ConfigPolicy, error) {
	cp := cpolicy.New()
	return cp, nil
}

// Meta returns plugin meta data
func Meta() *plugin.PluginMeta {
	return plugin.NewPluginMeta(
		name,
		version,
		plgtype,
		[]string{plugin.SnapGOBContentType},
		[]string{plugin.SnapGOBContentType},
		plugin.RoutingStrategy(plugin.StickyRouting),
	)
}

type collector struct {
	provider  *gophercloud.ProviderClient
	endpoints []types.Endpoint
	services  []types.Service
}
