
# DNS Testbed
DNS testbed is an environment to test `dns` zone configurations and queries. We provide a complete `dns` environment from `root` nameserver to `resolver`. We also provide tools to automatically send queries for variable zone configurations and measure the resulting volume or delay at a specified target component. 

The current environment contains the following components:

* **root NS**: Root server
* **com**: Authoritative name server for top-level-domain `.com`
* **net**: Authoritative name server for top-level-domain `.net`
* **target-com**: Authoritative name server for second-level-domain `target.com`
* **inter-net**: Authoritative name server for second-level-domain `inter.net`
* **resolver**: Recursive resolver
* **client**: Client machine

The **client** can query the **resolver** for zones provided by target.com or inter.net. The **resolver** then resolves that query recursively according to the DNS specifications.

## Setup
The individual components are simulated by running a separate docker container for each component. The nameservers (**com**, **net**, **target.com** and **inter.net**) run bind9 (version [bind-9.18.4](https://bind9.readthedocs.io/en/v9_18_4/notes.html)). The resolver is duplicated to support multiple resolver implementations. We provide one container for each implementation we support. Currently, the following implementations are supported:

* [bind-9.18.4](https://bind9.readthedocs.io/en/v9_18_4/notes.html).
* [unbound-1.16.0](https://www.nlnetlabs.nl/news/2022/Jun/02/unbound-1.16.0-released/)
* [unbound-1.10.0](https://www.nlnetlabs.nl/news/2020/Feb/20/unbound-1.10.0-released/)
* [powerDNS-4.7.3](https://docs.powerdns.com/recursor/changelog/4.7.html#change-4.7.3)

The **client** is a container running `dig`  to submit the DNS queries.

The file `docker/buildContext/docker-compose.yml` specifies a network, the properties of each component, including build context, IP addresses, and mounted volumes.

The IP addresses are assigned as follows:

| Component	                  | IP			       |	
|-----------------------------|-------------|
| root NS	                    | 172.20.0.2  |
| com			                 | 172.20.0.3  |
| net			                      | 172.20.0.4  |
| target		                    | 172.20.0.5  |
| inter			                    | 172.20.0.6  |
| client		                    | 172.20.0.9  |
| resolver-unbound-1.16.0		   | 172.20.0.12 |
| resolver-unbound-1.10.0		   | 172.20.0.13 |
| resolver-powerDNS-4.7.3		   | 172.20.0.14 |
| resolver-bind-9.18.4		      | 172.20.0.15 |


## Installation
1. Install [Docker](https://docs.docker.com/get-docker/)
2. Install [docker-compose](https://docs.docker.com/compose/install/linux/)

## Run
### Manually run queries
The docker containers are build and run with docker compose.

To start the testbed use the following command: 

```bash
docker-compose -f docker/buildContext/docker-compose.yml up -d
```

This builds and runs each component and starts the `dns` service. The testbed is now ready to process `dns` queries. E.g. one can attach to the client container and query a resolver for `target.com`  by specifying the resolvers `ip` address.

```bash
$ docker exec -it client bash

root@64a7f6a32a3a:/# dig @172.20.0.10 target.com A
```

The zone configurations can be updated by changing the `active.zone` file in the build context of the respective component, e.g. `docker/buildContext/nameserver/target-com/zones/active.zone`. The `dns` service of that component then needs to be restarted. To restart a service connect to that component and run the implementation specific restart command:

```bash
$docker exec -it target-com bash

# bind
root@64a7f6a32a3a:/# service named restart

# unbound
root@64a7f6a32a3a:/# unbound-control reload

# powerDNS
root@64a7f6a32a3a:/# /etc/init.d/pdns-recursor restart
```

The configuration of a resolver can also be changed by changing the configuration file in the respective build context, e.g. `docker/buildContext/resolver/resolver-bind-9.18.4/named.conf.options`. Again the `dns` service needs to be restarted.

### Automatically run experiments
We provide the possibility to automate experiments involving `dns` queries for different zone configurations. We define an experiment to be a sequence of measurements. Each measurement involves sending a `dns` query to a specific resolver implementation and measuring and effect. We distinguish between two types of measurements based on the type of effect we want to observe:

* **Volume based measurement**: We vary the zone configurations for the nameservers involved and measure the volume of queries received by a target nameserver. Currently, it is possible to simulate a scenario involving one or two nameservers.
* **Time based measurement**: We delay the responses of a nameserver, vary that delay and measure the time a resolver takes to process a query to that nameserver. To make things easier we measure the processing time as the time between the first and the last query the delayed nameserver receives. Currently, we only support scenarios with one nameserver involved.

We provide a selection of experiments that are already implemented and ready to run. Specifically we provide variations for two kind of attacks, **Subquery Unchained** and **Slow DNS**. 

**Subquery Unchained** is an instance of a volume based experiment involving two nameservers which we call target and inter. The zone configurations contain NS delegations to CNAME chains alternating between the target and inter nameserver. For each query for the entry zone we vary the number of NS delegations and measure the volume of queries received by the target resolver. 

**SlowDNS** is an instance of a time based experiment involving one nameserver which we call target. The zone configuration contains a CNAME chain with each element of the chain residing in the target zone. For each query we vary the delay the target nameserver responds to queries and measure the processing time of the resolver resolving that query.

Both experiments can be combined with different variations:

* `scrubbing` only involves a target nameserver and  each element of the CNAME chain resides in the target zone
* `DNAME` instead of CNAME chains the zone configuration contains DNAME chains

The concrete implementations of all these variations can be accessed as variables of the `lab` package e.g.`lab.Subquery_CNAME_Scrubbing)`

### Creating new experiments
We provide the possibility to create new experiments based on templates we provide for the two types of measurements. For both of these two measurements we provide a `go` struct that implements the necessary functionallity to measure the effect.

* `lab/volumetric_experiment.go`
* `lab/timing_experiment.go`

To create a new experiment run the following steps:

* Create a new instance of the respective experiment by specifying the name of the experiment, the entry zone triggering the query resolution, the base directory of the zone configuration files and the prefix of the zone name (for volumetric experiments only). The configuration files should be collected in a subdirectory of the base directory with the same name as the experiment. If more than one nameserver is involved, the zone configuration files for each nameserver should be contained in another subdirectory with the name of the nameserver.

```
volumetricExperiment := lab.newVolumetricExperiment("subquery+CNAME", "del.inter.net", "zones", "ns-del")

/*
directory structure:
/zones
	/subquery+CNAME
		/target-com
			ns-del-1.zone
			...
			ns-del-n.zone
		/inter-net
			ns-del-1.zone
			...
			ns-del-n.zone
*/

newTimingExperiment("slowDNS+CNAME", "a1.target.com", "zones")

/*
directory structure:
/zones
	/slowDNS+CNAME
		target.zone
*/
```
* Create a new instance of lab specifying the directory where the results should be stored in:
 
 ```
dnsTestLab := lab.New("results")
```
* Conduct the experiment by specifying the resolver implementations that should be tested and a data range. The data range corresponds to the delay (in ms) at the nameserver for the timing experiment and the zone configuration (more specifically the `n` in the name of the zone configuration) of the nameservers for the volumetric experiment.

```
dnsTestLab.Conduct(
		volumetricExperiment,
		lab.NewDataIterator(
			[]string{
				"bind-9.18.4",
				"powerDNS-4.7.3",
			},
			lab.MakeRange(1, 10, 1)
		),
)
dnsTestLab.SaveResults()
```

* The results will be stored in the specified directory as a .csv.