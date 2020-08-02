# airgap-ocp Tool
Tool to pull s subset of OCP Operators that work in an offline environment. 

## Requirements

1 - Docker needs to be installed
We currently use the Docker Go SDK to pull images to the local docker daemon.

Fedora:
```
$ dnf install docker
```

Docker should be installed, enabled and running before running the airgap-ocp CLI

2 - go1.14.4

3 - Make sure that the following entries are in the go.mod file if you are going to build the arigap-ocp tool.

	github.com/docker/docker master <= This entry will be replaced by github.com/docker/docker v17.12.0-ce-rc1.0.20200731192409-39691204f153+incompatible
	github.com/pkg/errors v0.9.1 // indirect

4 - The format for the .airgap-ocp.yaml file is as follows

rhnuser: "ACCESSUSER"
password: "REDACTED"
ocp-operators-dir: "/home/ocp-images/"
ocp-disconnected-operators:
  - registry.redhat.io/distributed-tracing/jaeger-ingester-rhel7
  - registry.redhat.io/rhpam-7/rhpam-rhel8-operator
  - registry.redhat.io/openshift4/ose-ansible-service-broker-operator
  - registry.redhat.io/openshift-service-mesh/istio-cni-rhel8
  - registry.redhat.io/openshift4/ose-ptp
  - registry.redhat.io/rhscl/postgresql-10-rhel7
  - registry.redhat.io/rhacm1-tech-preview/multicluster-operators-subscription-rhel8-operator
etc etc

This file can be anywhere.  Use the --config /PATH/airgap-ocp.yaml on the command line.


| Variable | Description | Sample Value |
| :------------- | :----------: | -----------: |
|  rhnuser  | User Name provided by Red Hat to access.redhat.com | rhsupport    |
|  password | User password | Self-explanatory |
| ocp-operators-dir | Directory where the tar balls for each image will be created | /ocp-images |
| ocp-disconnected-operators | YAML list of Ids for the images | registry.redhat.io/openshift4/ose-ptp |

## Building airgap-ocp

You will have to do the following:

1 - Clone the github repo

```
$ cd $GOPATH/go/src
$ git clone https://github.com/claudiol/airgap-ocp.git
```

2 - Build the binary

```
$ cd $GOPATH/go/src/airgap-ocp
$ go install airgap-ocp or go build airgap-ocp
```

## Run the airgap-ocp tool

```
[claudiol@fedora airgap-ocp]$ ./airgap-ocp 
A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  airgap-ocp [command]

Available Commands:
  help          Help about any command
  pullOperators A brief description of your command

Flags:
      --config string   config file (default is $HOME/.airgap-ocp.yaml)
  -h, --help            help for airgap-ocp
  -t, --toggle          Help message for toggle

Use "airgap-ocp [command] --help" for more information about a command.
[claudiol@fedora airgap-ocp]$
```

### Sample Run:

```
[claudiol@fedora airgap-ocp]$ sudo ./airgap-ocp pullOperators --config /home/claudiol/.airgap-ocp.yaml | tee airgap-ocp.log
[sudo] password for claudiol: 
Using config file: /home/claudiol/.airgap-ocp.yaml
rhn-gps-claudiol
/home/ocp-images/
{"status":"Pulling from distributed-tracing/jaeger-ingester-rhel7","id":"latest"}
{"status":"Already exists","progressDetail":{},"id":"a03401a44180"}
{"status":"Already exists","progressDetail":{},"id":"bb0da44cdbce"}
{"status":"Already exists","progressDetail":{},"id":"33a4b7b85f9b"}
{"status":"Digest: sha256:52fb63a6d89d0c9a44fe8560b44e429c8e40dafded207f02f8bde5bd6a108531"}
{"status":"Status: Downloaded newer image for registry.redhat.io/distributed-tracing/jaeger-ingester-rhel7:latest"}
Saving: sha256:fa735e06b2c28a09c6aa5db2eee362915c86fd6a1dddf1a3ec8778e8f2b6ef7a | [registry.redhat.io/distributed-tracing/jaeger-ingester-rhel7:latest] | 231775414
```

### Output for airgap-ocp tool

The output is a tar ball for each image downloaded.  The tarball can be added to an ISO or to a DVD and sent for scanning to an Air Gapped environment. Each tarball can be loaded to a registry by using *docker load -i tarball.tar* or *podman load -i tarball.tar*. 

```
[claudiol@fedora airgap-ocp]$ ls /home/ocp-images
amq-broker-latest.tar                           ocs-rhel8-operator-latest.tar                                           pluginbroker-artifacts-rhel8-latest.tar
amq-streams-bridge-rhel7-latest.tar             openshift-migration-controller-rhel8-latest.tar                         pluginbroker-metadata-rhel8-latest.tar
amq-streams-kafka-24-rhel7-latest.tar           openshift-migration-plugin-rhel8-latest.tar                             plugin-openshift-rhel8-latest.tar
[claudiol@fedora airgap-ocp]$ 
```

### Using *podman load* on one of the tar balls 
```
[claudiol@fedora ocp-images]$ podman load -i ./amq-broker-latest.tar
Getting image source signatures
Copying blob 563c2d198539 done  
Copying blob b9e30d7689db done  
Copying blob b1d42883ec50 done  
Copying blob abc8c32bad76 done  
Copying config 690a9b4d4e done  
Writing manifest to image destination
Storing signatures
Loaded image(s): registry.redhat.io/amq7/amq-broker:latest
[claudiol@fedora ocp-images]$ podman images
REPOSITORY                                                    TAG     IMAGE ID      CREATED      SIZE
registry.redhat.io/amq7/amq-broker                            latest  690a9b4d4e2f  2 weeks ago  653 MB
registry.redhat.io/distributed-tracing/jaeger-ingester-rhel7  latest  fa735e06b2c2  2 weeks ago  242 MB
[claudiol@fedora ocp-images]$ 
```

