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
)

// Authenticate is used to authenticate user for given tenant. Request is send to provided Keystone endpoint
// Returns authenticated provider client, which is used as a base for service clients.
func Authenticate(endpoint, user, password, tenant, domain_name, domain_id string) (*gophercloud.ProviderClient, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: endpoint,
		Username:         user,
		Password:         password,
		TenantName:       tenant,
		AllowReauth:      true,
	}
	if domain_name != "" && domain_id == "" {
		authOpts.DomainName = domain_name
	}
	if domain_id != "" && domain_name == "" {
		authOpts.DomainID = domain_id
	}

	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	return provider, nil
}
