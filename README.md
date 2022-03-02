# astra-cli

Apache 2.0 licensed Astra Cloud Management CLI 

[![.github/workflows/go.yaml](https://github.com/rsds143/astra-cli/actions/workflows/go.yaml/badge.svg)](https://github.com/rsds143/astra-cli/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsds143/astra-cli)](https://goreportcard.com/report/github.com/rsds143/astra-cli)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/rsds143/astra-cli)](https://img.shields.io/github/go-mod/go-version/rsds143/astra-cli)
[![Latest Version](https://img.shields.io/github/v/release/rsds143/astra-cli?include_prereleases)](https://github.com/rsds143/astra-cli/releases)
[![Coverage Status](https://coveralls.io/repos/github/rsds143/astra-cli/badge.svg)](https://coveralls.io/github/rsds143/astra-cli)

## status

Ready for production

## How to install - install script

* /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/datastax-labs/astra-cli/main/script/install-astra.sh)"
* astra login

## How to install - docker

Instead of downloading the binary this trusts that you have docker installed

* make sure docker is installed
* /bin/bash -c "$(curl -fsSL  https://raw.githubusercontent.com/datastax-labs/astra-cli/main/script/install-astra-docker.sh)"
* astra.sh login

## How to install - Homebrew for Mac and Linux

* [install homebrew](https://brew.sh/) if you have not
* `brew tap rsds143/rsds && brew install astra-cli`

## How to install - Release Binaries

* download a [release](https://github.com/datastax-labs/astra-cli/releases)
* tar zxvf <download>
* cd <extracted folder>
* ./astra

## How to install - From Source

* Install [Go 1.17](https://golang.org/dl/)
* run `git clone git@github.com:datastax-labs/astra-cli.git`
* run `./scripts/build` or `go build -o ./bin/astra .`

## How to use

* login
* execute commands on your database

### login with token

After creating a token with rights to use the devops api 

```
astra login --token "changed"
Login information saved
```
### login service account

After creating a service account on the Astra page 

```
astra login --id "changed" --name "changed" --secret "changed"
Login information saved
```

## login service account with json

```
astra login --json '{"clientId":"changed","clientName":"change@me.com","clientSecret":"changed"}'
Login information saved

```

### creating database

```
astra db create -v --keyspace myks --name mydb 
............
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b created
```

### get secure connection bundle

```
astra db secBundle 3c577e51-4ff5-4551-86a4-41d475c61822 -d external -l external.zip            
file external.zip saved 12072 bytes written
astra db secBundle 3c577e51-4ff5-4551-86a4-41d475c61822 -d internal -l internal.zip            
file internal.zip saved 12066 bytes written
astra db secBundle 3c577e51-4ff5-4551-86a4-41d475c61822 -d proxy-internal -l proxy-internal.zip 
file proxy-internal.zip saved 348 bytes written
astra db secBundle 3c577e51-4ff5-4551-86a4-41d475c61822 -d proxy-external -l proxy-external.zip 
file proxy-external.zip saved 339 bytes written
```

### get secure connection bundle URLs

```
astra db secBundle 3c577e51-4ff5-4551-86a4-41d475c61822 -o list         
  external bundle: changed
  internal bundle: changed
  external proxy: changed
  internal proxy: changed
```

### listing databases

```
astra db list                                                             
name id                                   status
mydb 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b ACTIVE
```

### listing databases in json

```
astra db list -o json
[
  {
    "id": "2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b",
    "orgId": "changed",
    "ownerId": "changed",
    "info": {
      "name": "mydb",
      "keyspace": "myks",
      "cloudProvider": "GCP",
      "tier": "developer",
      "capacityUnits": 1,
      "region": "us-east1",
      "user": "dbuser",
      "password": "",
      "additionalKeyspaces": null,
      "cost": null
    },
    "creationTime": "2021-02-24T17:23:19Z",
    "terminationTime": "0001-01-01T00:00:00Z",
    "status": "ACTIVE",
    "storage": {
      "nodeCount": 1,
      "replicationFactor": 1,
      "totalStorage": 5,
      "usedStorage": 0
    },
    "availableActions": [
      "park",
      "getCreds",
      "resetPassword",
      "terminate",
      "addKeyspace",
      "removeKeyspace",
      "addTable"
    ],
    "message": "",
    "studioUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.studio.astra.datastax.com",
    "grafanaUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.dashboard.astra.datastax.com/d/cloud/dse-cluster-condensed?refresh=30s\u0026orgId=1\u0026kiosk=tv",
    "cqlshUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.apps.astra.datastax.com/cqlsh",
    "graphUrl": "",
    "dataEndpointUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.apps.astra.datastax.com/api/rest"
  }
]
```
### getting database by id

```
astra db get 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
name id                                   status
mydb 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b ACTIVE
```

### getting database by id in json

```
astra db get 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b -o json 
{
  "id": "2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b",
  "orgId": "changed",
  "ownerId": "changed",
  "info": {
    "name": "mydb",
    "keyspace": "myks",
    "cloudProvider": "GCP",
    "tier": "developer",
    "capacityUnits": 1,
    "region": "us-east1",
    "user": "dbuser",
    "password": "",
    "additionalKeyspaces": null,
    "cost": null
  },
  "creationTime": "2021-02-24T17:23:19Z",
  "terminationTime": "0001-01-01T00:00:00Z",
  "status": "ACTIVE",
  "storage": {
    "nodeCount": 1,
    "replicationFactor": 1,
    "totalStorage": 5,
    "usedStorage": 0
  },
  "availableActions": [
    "park",
    "getCreds",
    "resetPassword",
    "terminate",
    "addKeyspace",
    "removeKeyspace",
    "addTable"
  ],
  "message": "",
  "studioUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.studio.astra.datastax.com",
  "grafanaUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.dashboard.astra.datastax.com/d/cloud/dse-cluster-condensed?refresh=30s\u0026orgId=1\u0026kiosk=tv",
  "cqlshUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.apps.astra.datastax.com/cqlsh",
  "graphUrl": "",
  "dataEndpointUrl": "https://2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b-us-east1.apps.astra.datastax.com/api/rest"
}
```


### parking database

NOTE: Does not work on serverless

```
astra db park -v 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b         
starting to park database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
...........
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b parked
```

### unparking database

NOTE: Does not work on serverless

```
astra db unpark -v 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
starting to unpark database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
...........
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b unparked
```

### deleting database

```
astra db delete -v 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
starting to delete database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b deleted
```

### resizing

I did not have a paid account to verify this works, but you can see it succesfully starts the process

```
astra db resize -v 72c4d35b-1875-495a-b5f1-97329d90b6c5 2                    
unable to unpark '72c4d35b-1875-495a-b5f1-97329d90b6c5' with error expected status code 2xx but had: 400 error was [map[ID:2.000009e+06 message:resizing is not supported for this database tier]]
```


