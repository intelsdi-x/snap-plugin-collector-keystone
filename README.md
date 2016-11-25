# Snap plugin collector - keystone

Snap plugin for collecting metrics from OpenStack Keystone module.

1. [Getting Started](#getting-started)
  * [System Requirements](#system-requirements)
  * [Operating systems](#operating-systems)
  * [Installation](#installation)
  * [Configuration and Usage](configuration-and-usage)
2. [Documentation](#documentation)
  * [Collected Metrics](#collected-metrics)
  * [Snap's Global Config](#snaps-global-config)
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
All OSs currently supported by Snap:
* Linux/amd64

### Installation
#### Download the plugin binary:

You can get the pre-built binaries for your OS and architecture from the plugin's [GitHub Releases](https://github.com/intelsdi-x/snap-plugin-collector-keystone/releases) page. Download the plugin from the latest release and load it into `snapteld` (`/opt/snap/plugins` is the default location for Snap packages).

#### To build the plugin binary:

Fork https://github.com/intelsdi-x/snap-plugin-collector-keystone
Clone repo into `$GOPATH/src/github.com/intelsdi-x/`:

```
$ git clone https://github.com/<yourGithubID>/snap-plugin-collector-keystone.git
```

Build the Snap keystone plugin by running make within the cloned repo:
```
$ make
```
This builds the plugin in `./build/`

### Configuration and Usage
* Set up the [Snap framework](https://github.com/intelsdi-x/snap#getting-started).
* Create Global Config, see description in [Snap's Global Config] (#snaps-global-config).
* Load the plugin and create a task, see example in [Examples](#examples).

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

### Snap's Global Config
Global configuration files are described in [Snap's documentation](https://github.com/intelsdi-x/snap/blob/master/docs/SNAPD_CONFIGURATION.md). You have to add section "keystone" in "collector" section and then specify following options:
- `"admin_endpoint"` - URL for OpenStack Identity admin endpoint (ex. `"http://keystone.public.org:35357"`)
- `"admin_user"` -  administrator user name
- `"admin_password"` - administrator password
- `"admin_tenant"` - administration tenant

Example global configuration file for snap-plugin-collector-keystone plugin (exemplary file in [examples/cfg] (examples/cfg/cfg.json)):

### Examples
Example of running Snap keystone collector and writing data to file.

Download an [example Snap global config](examples/cfg/cfg.json) file.
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-keystone/master/examples/cfg/cfg.json
```
Ensure to provide your Keystone instance address and credentials.

Ensure [Snap daemon is running](https://github.com/intelsdi-x/snap#running-snap) with provided configuration file:
* command line: `snapteld -l 1 -t 0 --config cfg.json&`

Download and load Snap plugins:
```
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-collector-keystone/latest/linux/x86_64/snap-plugin-collector-keystone
$ wget http://snap.ci.snap-telemetry.io/plugins/snap-plugin-publisher-file/latest/linux/x86_64/snap-plugin-publisher-file
$ chmod 755 snap-plugin-*
$ snaptel plugin load snap-plugin-collector-keystone
$ snaptel plugin load snap-plugin-publisher-file
```

See all available metrics:

```
$ snaptel metric list
```

Download an [example task file](examples/tasks/task.json) and load it:
```
$ curl -sfLO https://raw.githubusercontent.com/intelsdi-x/snap-plugin-collector-keystone/master/examples/tasks/task.json
$ snaptel task create -t task.json
Using task manifest to create task
Task created
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
Name: Task-02dd7ff4-8106-47e9-8b86-70067cd0a850
State: Running
```

See realtime output from `snaptel task watch <task_id>` (CTRL+C to exit)
```
$ snaptel task watch 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

This data is published to a file `/tmp/published_keystone` per task specification

Stop task:
```
$ snaptel task stop 02dd7ff4-8106-47e9-8b86-70067cd0a850
Task stopped:
ID: 02dd7ff4-8106-47e9-8b86-70067cd0a850
```

### Roadmap
There isn't a current roadmap for this plugin, but it is in active development. As we launch this plugin, we do not have any outstanding requirements for the next release.

## Community Support
This repository is one of **many** plugins in **Snap**, a powerful telemetry framework. See the full project at http://github.com/intelsdi-x/snap.

To reach out to other users, head to the [main framework](https://github.com/intelsdi-x/snap#community-support).

## Contributing
We love contributions!

There's more than one way to give back, from examples to blogs to code updates. See our recommended process in [CONTRIBUTING.md](CONTRIBUTING.md).

## License
[Snap](http://github.com/intelsdi-x/snap), along with this plugin, is an Open Source software released under the Apache 2.0 [License](LICENSE).

## Acknowledgements
* Author: [Marcin Krolik](https://github.com/marcin-krolik)