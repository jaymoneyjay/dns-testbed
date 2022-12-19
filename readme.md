
# DNS Testbed
DNS testbed is an environment to test `dns` zone configurations and queries. A testbed consists of one or more **zones**, one or more **resolvers** and one **client**. The **client** can query the **resolver** for zones provided by target.com or inter.net. The **resolver** then resolves that query recursively according to the DNS specifications.
We provide a simple command line interface to:

* initialize a custom testbed
* change zone configurations
* delay zone responses
* flush the resolvers cache
* query a resolver for a dns record and measure the query duration or query volume at a specific component

We also provide functionallity to run a series of dns queries and measurements with different testbed settings such as different zone configurations and delays.

## Setup
The individual components are simulated by running a separate docker container for each component. The nameservers run bind9 (version [bind-9.18.4](https://bind9.readthedocs.io/en/v9_18_4/notes.html)). We provide different implementations for resolvers. Currently, the following implementations are supported:

* [bind-9.18.4](https://bind9.readthedocs.io/en/v9_18_4/notes.html).
* [unbound-1.16.0](https://www.nlnetlabs.nl/news/2022/Jun/02/unbound-1.16.0-released/)
* [unbound-1.10.0](https://www.nlnetlabs.nl/news/2020/Feb/20/unbound-1.10.0-released/)
* [powerDNS-4.7.3](https://docs.powerdns.com/recursor/changelog/4.7.html#change-4.7.3)

The **client** is a container running `dig`  to submit the DNS queries.

## Installation
1. Install [Docker](https://docs.docker.com/get-docker/)
2. Install [docker-compose](https://docs.docker.com/compose/install/linux/)
3. Inside the project root run


   ```
   $ go mod tidy
   $ go install
   ```

## Usage

### Initialization
A new testbed can be initialized by providing a `.yaml` file containing the configuration of the testbed. This configuration file is structured as follows:

```yaml
zones:
- qname:           string
  ip:              string
  defaultZoneFile: path
  [...]

resolvers:
- implementation:  string
  version:         string
  ip:              string
  [...]

client:
  id:              string
  ip:              string
  nameserver:	   string
```

We provide a configuration for a basic testbed `validation/testbed-basic.yaml`. The testbed contains the following components

* **zones**: root, com, net, target.com, inter.net, 
* four **resolvers** each with a different dns implementations: bind-9.18.4, unbound-1.16.0, unbound-1.10.0, powerDNS-4.7.3
* a **client**

To initialize and start a testbed run 

```
$ testbed init example/testbed.yaml
$ testbed start
```
Make sure that docker desktop is running.

### Manual Usage
For the full list of commands run `testbed -h`:

```
Usage:
  testbed [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delay       Delay the responses of a zone by the specified duration (in ms)
  flush       Flush the cache of all resolvers
  help        Help about any command
  init        Initialize a dns testbed
  query       Query a resolver for a specific qname and record
  run         Run an experiment according to the specified configuration
  start       Build and run the dns testbed
  stop        Stop the dns testbed
  zones       Set zone files

Flags:
  -h, --help   help for testbed

Use "testbed [command] --help" for more information about a command.

```

### DNS experiment
We define a dns experiment to be a series of dns measurements. To run a dns measurement a `.yaml` file with the configuration of the experiment has to be provided.

```yaml
name:            string, name of the experiment
resolverIDs:     []string, ids of resolvers to query
zonesDir:        path, path to zone configurations; 
	         the zone files should be named after the id of the zone
	         and be collected in a directory specifying the measurement; e.g.
	         /zonesDir
	              /1
	                   target-com.zone
	                   inter-net.zon
	              /n
                      ...
target:          string, specifying the component where to measure the query volume or duration
measure:         ["volume", "duration"], specifying the kind of measurement
query:
  recordType:    string (default="A")
  qname:         string
qmin:            bool (default=false), whether qmin is enabled
delay:           []int (default=[0]), delays to be set at the delayedZones before measurements
delayedZones:    []string
warmup:          []string, list of qnames to query before a measurement to warm up the resolver
dest:            path, location where to save the results
saveLogs:        bool, whether to save the logs
```

We provide a series of example configurations in `validation/experiment`.

Some commands require specifying the id of the component. For the zones the id corresponds to the qname with intermediary dots replaced by dashes. Zone files must always be named after the id of the respective zone.

|qname				|	id|
|-------------|----|
|.| root|
|.com| com|
|target.com| target-com|

The id of the resolvers corresponds to the concatenation of the identifier *resolver*, the implementation and the version.


|implementation|version	|	id|
|-------------|----|----|
|bind| 9.18.4|resolver-bind-9.18.4|
|unbound| 1.10.0|resolver-unbound-1.10.0|

