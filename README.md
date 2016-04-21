# snap plugin collector - keystone

snap plugin for collecting metrics from OpenStack Keystone module. 

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [snap's Global Config](#snaps-global-config)
  * [Examples](#examples)
  * [Roadmap](#roadmap)
3. [Community Support](#community-support)
4. [Contributing](#contributing)
5. [License](#license)
6. [Acknowledgements](#acknowledgements)

## Getting Started

Plugin collects metrics by communicating with OpenStack by REST API.
It can run locally on the host, or in proxy mode (communicating with the host via HTTP(S)). 

### System Requirements
* OpenStack deployment available
* Supports Keystone V2 and V3 authorization APIs
 
### Operating systems
All OSs currently supported by snap:
* Linux/amd64

### Installation
#### Download keystone plugin binary:
You can get the pre-built binaries for your OS and architecture at snap's [Github Releases](https://github.com/intelsdi-x/snap/releases) page.

#### To build the plugin binary:
Fork https://github.com/intelsdi-x/snap-plugin-collector-keystone

Clone repo into `$GOPATH/src/github/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-keystone
```
Build the plugin by running make in repo:
```
$ make
```
This builds the plugin in `/build/rootfs`

### Configuration and Usage
* Set up the [snap framework](https://github.com/intelsdi-x/snap/blob/master/README.md#getting-started).
* Create Global Config, see description in [snap's Global Config] (https://github.com/intelsdi-x/snap-plugin-collector-keystone/blob/master/README.md#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](https://github.com/intelsdi-x/snap-plugin-collector-keystone/blob/master/README.md#examples).

#### Suggestions
* It is not recommended to set interval for task less than 20 seconds. This may lead to overloading Keystone API with requests.

## Documentation

### Collected Metrics
This plugin has the ability to gather the following metrics:

Namespace | Data Type | Description
----------|-----------|-----------------------
intel/openstack/keystone/\<tenant_name\>/users_count | int | Total number of users for given tenant
intel/openstack/keystone/total_tenants_count | int | Total number of tenants
intel/openstack/keystone/total_users_count | int | Total number of users 
intel/openstack/keystone/total_endpoints_count | int | Total number of endpoints
intel/openstack/keystone/total_services_count | int | Total number of services

### snap's Global Config
Global configuration files are described in [snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "keystone" in "collector" section and then specify following options:
- `"admin_endpoint"` - URL for OpenStack Identity admin endpoint (ex. `"http://keystone.public.org:35357"`)
- `"admin_user"` -  administrator user name
- `"admin_password"` - administrator password
- `"admin_tenant"` - administration tenant

Example global configuration file for snap-plugin-collector-keystone plugin (exemplary file in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-keystone/blob/master/examples/cfg/)):
```
{
  "control": {
    "cache_ttl": "5s"
  },
  "scheduler": {
    "default_deadline": "5s",
    "worker_pool_size": 5
  },
  "plugins": {
    "collector": {
      "keystone": {
        "all": {
          "admin_endpoint": "https://public.fuel.local:35357",
          "admin_user": "admin",
          "admin_password": "admin",
          "admin_tenant": "admin"
        }
      }
    },
    "publisher": {},
    "processor": {}
  }
}
```

### Examples
Example running snap-plugin-collector-keystone plugin and writing data to a file.

Make sure that your `$SNAP_PATH` is set, if not:
```
$ export SNAP_PATH=<snapDirectoryPath>/build
```
Other paths to files should be set according to your configuration, using a file you should indicate where it is located.

Create Global Config, see example in [examples/cfg/] (https://github.com/intelsdi-x/snap-plugin-collector-keystone/blob/master/examples/cfg/).

In one terminal window, open the snap daemon (in this case with logging set to 1,  trust disabled and global configuration saved in cfg.json ):
```
$ $SNAP_PATH/bin/snapd -l 1 -t 0 --config cfg.json
```
In another terminal window:

Load snap-plugin-collector-keystone plugin
```
$ $SNAP_PATH/bin/snapctl plugin load snap-plugin-collector-keystone
```
Load file plugin for publishing:
```
$ $SNAP_PATH/bin/snapctl plugin load $SNAP_PATH/plugin/snap-publisher-file
```
See available metrics for your system

```
$ $SNAP_PATH/bin/snapctl metric list
```

Create a task manifest file to use snap-plugin-collector-keystone plugin (exemplary files in [examples/tasks/] (https://github.com/intelsdi-x/snap-plugin-collector-keystone/blob/master/examples/tasks/)):
```
{
    "version": 1,
    "schedule": {
        "type": "simple",
        "interval": "60s"
    },
    "workflow": {
        "collect": {
            "metrics": {
		        "/intel/openstack/keystone/total_tenants_count": {},
		        "/intel/openstack/keystone/total_users_count": {},
		        "/intel/openstack/keystone/total_endpoints_count": {},
		        "/intel/openstack/keystone/admin/users_count": {}
           },
            "config": {
            },
            "process": null,
             "publish": [
                {
                    "plugin_name": "file",
                    "config": {
                        "file": "/tmp/published_keystone"
                    }
                }
            ]
        }
    }
}
```
Create a task:
```
$ $SNAP_PATH/bin/snapctl task create -t examples/tasks/task.json
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support) or visit [snap Gitter channel](https://gitter.im/intelsdi-x/snap).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

And **thank you!** Your contribution, through code and participation, is incredibly important to us.

## License
[snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)