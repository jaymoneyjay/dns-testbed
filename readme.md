
# DNS Testbed
DNS testbed is an environment to test DNS configurations and queries. The current environment contains the following components:

* **root NS**: Root server
* **com**: Authoritative name server for top-level-domain `.com`
* **net**: Authoritative name server for top-level-domain `.net`
* **target-com**: Authoritative name server for second-level-domain `target.com`
* **inter-net**: Authoritative name server for second-level-domain `inter.net`
* **resolver**: Recursive resolver
* **client**: Client machine

The goal is that a **client** can query the **resolver** about target.com and the **resover** responds with the address of **target-com** after recursively querying the higher level nameservers.

## Setup
The individual components are simulated by running a separate docker container for each component.
The testbed currently supports the following dns-intallations:

* [bind-9.11.3](https://www.isc.org/bind/).
* unbound-1.17.0
* unbound-1.16.0
* unbound-1.10.0
* powerDNS-4.7.3

The component **Client** is an empty docker containers used for completeness and submitting the DNS queries respectively.

The IP addresses are assigned as follows:

|Component	| IP			       |	
|------------	|-------------|
|root NS		| 172.20.0.2  |
|com			| 172.20.0.3  |
|net			| 172.20.0.4  |
|target		| 172.20.0.5  |
|inter			| 172.20.0.6  |
|resolver		| 172.20.0.9  |
|client		| 172.20.0.10 |


## Installation
1. Install [Docker](https://docs.docker.com/get-docker/)
2. Install [docker-compose](https://docs.docker.com/compose/install/linux/)
## Run
The docker containers are build and run with docker compose.

* Start the testbed: 

```bash
docker-compose -f testbed/docker/buildContext/docker-compose.yml up -d
```
* Create a new experiment

```go
subqueryCNAMEExperiment := experiment.NewSubqueryExperiment(experiment.SubqueryCNAME_QMIN)
```

* Run the experiment

```go
err := subqueryCNAMEExperiment.Run(utils.MakeRange(1, 10, 1))
if err != nil {
	log.Fatal(err)
}
```

* The results will be stored as a .csv in the results directory.