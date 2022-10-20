
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

* [bind9](https://www.isc.org/bind/). 

The component **Client** is an empty docker containers used for completeness and submitting the DNS queries respectively.

The IP addresses are assigned as follows:

|Component	|IP			 |	
|------------	|----------|
|root NS		|172.20.0.2|
|com			|172.20.0.3|
|net			|172.20.0.4|
|target		|172.20.0.7|
|inter			|172.20.0.8|
|resolver		|172.20.0.9|
|client		|172.20.0.10|


## Installation
1. Install [Docker](https://docs.docker.com/get-docker/)
2. Start the docker containers with the command



## Run
The docker containers are build and run with docker compose.

* Start the testbed: 

```bash
compose -f testbed/docker/buildContext/docker-compose.yml up -d
```
* Create a new testbed.

```go
testbed := testbed.NewTestbed()
```

* Start a dns implementation

```go
err = testbed.Start(component.Bind9)
if err != nil {
	log.Fatal(err)
}
```

* Create zone files by running an attack

```go
_, err := attack.NewTemplateAttack().WriteZoneFilesAndReturnEntryZone(0, testbed.Nameservers["sld"])
if err != nil {
	log.Fatal(err)
}
```

* Query zone

```go
queryResult, err := testbed.Query("target.com	")
if err != nil {
	log.Fatal(err)
}
fmt.Print(queryResult)
```