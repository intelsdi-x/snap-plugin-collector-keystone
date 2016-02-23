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

package openstack

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
	"github.com/rackspace/gophercloud/openstack/identity/v2/users"
	"github.com/rackspace/gophercloud/openstack/identity/v3/endpoints"
	"github.com/rackspace/gophercloud/openstack/identity/v3/services"

	"github.com/intelsdi-x/snap-plugin-collector-keystone/openstack/tenantusers"
	"github.com/intelsdi-x/snap-plugin-collector-keystone/types"
)

// GetTenants is used to retrieve list of available tenant for authenticated user
func GetAllTenants(provider *gophercloud.ProviderClient) ([]types.Tenant, error) {
	tnts := []types.Tenant{}

	client := openstack.NewIdentityV2(provider)

	opts := tenants.ListOpts{}
	pager := tenants.List(client, &opts)

	page, err := pager.AllPages()
	if err != nil {
		return tnts, err
	}

	tenantList, err := tenants.ExtractTenants(page)
	if err != nil {
		return tnts, err
	}

	for _, t := range tenantList {
		tnts = append(tnts, types.Tenant{Name: t.Name, ID: t.ID})
	}

	return tnts, nil
}

// GetAllUsers is used to retrieve list of available users
func GetAllUsers(provider *gophercloud.ProviderClient) ([]types.User, error) {
	userList := []types.User{}

	client := openstack.NewIdentityV3(provider)

	pager := users.List(client)
	page, err := pager.AllPages()
	if err != nil {
		return userList, err
	}

	usrs, err := users.ExtractUsers(page)
	if err != nil {
		return userList, err
	}

	for _, u := range usrs {
		userList = append(userList, types.User{
			ID:       u.ID,
			Name:     u.Name,
			Username: u.Username,
		})
	}

	return userList, nil
}

// GetAllServices is used to retrieve list of available services for authenticated admin
func GetAllServices(provider *gophercloud.ProviderClient) ([]types.Service, error) {
	serviceList := []types.Service{}

	client := openstack.NewIdentityV3(provider)

	opts := services.ListOpts{}
	pager := services.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return serviceList, err
	}

	srvs, err := services.ExtractServices(page)
	if err != nil {
		return serviceList, err
	}

	for _, s := range srvs {
		serviceList = append(serviceList, types.Service{
			ID:          s.ID,
			Name:        s.Name,
			Type:        s.Type,
			//Description: s.Description,
		})
	}

	return serviceList, nil
}

// GetAllServices is used to retrieve list of available services for authenticated admin
func GetAllEndpoints(provider *gophercloud.ProviderClient) ([]types.Endpoint, error) {
	endpointList := []types.Endpoint{}

	client := openstack.NewIdentityV3(provider)

	opts := endpoints.ListOpts{}
	pager := endpoints.List(client, opts)
	page, err := pager.AllPages()
	if err != nil {
		return endpointList, err
	}

	endpts, err := endpoints.ExtractEndpoints(page)
	if err != nil {
		return endpointList, err
	}

	for _, endpt := range endpts {
		endpointList = append(endpointList, types.Endpoint{
			ID:           endpt.ID,
			ServiceID:    endpt.ServiceID,
			URL:          endpt.URL,
			Region:       endpt.Region,
			//Availability: endpt.Availability,
			Name:         endpt.Name,
		})
	}

	return endpointList, nil
}

func GetUsersPerTenant(provider *gophercloud.ProviderClient, tenantList []types.Tenant) (map[string]int, error) {
	tenantUsersCount := map[string]int{}

	client := openstack.NewIdentityV2(provider)

	for _, tnt := range tenantList {
		usrs, err := tenantusers.Get(client, tnt.ID).Extract()
		if err != nil {
			return tenantUsersCount, err
		}

		tenantUsersCount[tnt.Name] = len(usrs)
	}

	return tenantUsersCount, nil
}
