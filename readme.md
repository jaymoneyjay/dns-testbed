
# DNS Testbed
DNS testbed is an environment to test `dns` zone configurations and queries. We provide a complete `dns` environment from `root` nameserver to `resolver`. 

The testbed can be configured dynamically. We provide a default configuration file in `config.yaml` resulting in the following environment:

| Component	                  | IP			       |	
|-----------------------------|-------------|
| root NS	                    | 172.20.0.2  |
| com			                 | 172.20.0.3  |
| net			                      | 172.20.0.4  |
| target		                    | 172.20.0.5  |
| inter			                    | 172.20.0.6  |
| resolver-bind-9.18.4		      | 172.20.0.51 |
| resolver-unbound-1.16.0		   | 172.20.0.52 |
| resolver-unbound-1.10.0		   | 172.20.0.53 |
| resolver-powerDNS-4.7.3		   | 172.20.0.54 |
| client		                    | 172.20.0.99  |

The **client** can query the **resolver** for zones provided by target.com or inter.net. The **resolver** then resolves that query recursively according to the DNS specifications.

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
To initialize and start the defautl testbed run 

```
$ testbed init config.yaml zones/default
$ testbed start
```

For the full list of commands run `testbed -h`:

```
Usage:
  testbed [command]

Available Commands:
  delay       Delay the responses of a zone by the specified duration (in ms)
  flush       Flush the cache of all resolvers
  help        Help about any command
  init        Initialize a dns testbed
  query       Query a resolver for a specific qname and record
  start       Build and run the dns testbed
  stop        Stop the dns testbed
  zones       Set zone files

Flags:
  -h, --help   help for testbed

```

### Customization
The testbed can be customized by providing a custom testbed configuration as `.yaml` in the `testbed init [config, default zones]` command. The configuration is structured as follows:

```yaml
zones:
- qname:
  ip: 
  [...]

resolvers:
- implementation:
  version:
  ip:
  [...]

client:
  id:
  ip:
  nameserver:
```

The file `config.yaml` can be used as a template.

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

