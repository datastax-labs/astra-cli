# astra-cli

Apache 2.0 licensed Astra Cloud Management CLI 

[![.github/workflows/go.yaml](https://github.com/rsds143/astra-cli/actions/workflows/go.yaml/badge.svg)](https://github.com/rsds143/astra-cli/actions/workflows/go.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/rsds143/astra-cli)](https://goreportcard.com/report/github.com/rsds143/astra-cli)
[![Latest Version](https://img.shields.io/github/v/tag/rsds143/astra-cli?label=version)]

## status

- Alpha

## How to install

* download a [release](https://github.com/rsds143/astra-cli/releases)
* tar zxvf <download>
* cd <extracted folder>
* ./astra-cli

## How to build

* Install [Go 1.16](https://golang.org/dl/)
* run `git clone git@github.com:rsds143/astra-cli.git`
* run `./scripts/build` or `go build -o ./bin/astra-cli .`

## How to use

* login
* execute commands on your database

### login

After creating a service account on the Astra page 

```
./bin/astra-cli login -id "changed" -name "changed" -secret "changed"
Login information saved
```

## login with json

```
./bin/astra-cli login -json '{"clientId":"changed","clientName":"change@me.com","clientSecret":"changed"}'
Login information saved

```

### creating database

```
./bin/astra-cli db create -user dbuser -password test234  -keyspace myks -name mydb 
2021/02/24 18:23:24 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 19 more times
2021/02/24 18:23:29 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 18 more times
2021/02/24 18:23:35 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 17 more times
2021/02/24 18:23:40 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 16 more times
2021/02/24 18:23:45 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 15 more times
2021/02/24 18:23:50 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 14 more times
2021/02/24 18:23:55 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 13 more times
2021/02/24 18:24:00 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PENDING but expected ACTIVE trying again 12 more times
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b created
```

### listing databases

```
./bin/astra-cli db list                                                             
name id                                   status
mydb 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b ACTIVE
```

### listing databases in json

```
./bin/astra-cli db list -format json
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
./bin/astra-cli db get 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
name id                                   status
mydb 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b ACTIVE
```

### getting database by id in json

```
./bin/astra-cli db get 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b -format json 
json
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

```
./bin/astra-cli db park 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b         
starting to park database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
2021/02/24 18:31:26 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PARKING but expected PARKED trying again 29 more times
2021/02/24 18:31:56 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PARKING but expected PARKED trying again 28 more times
2021/02/24 18:32:26 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PARKING but expected PARKED trying again 27 more times
2021/02/24 18:32:57 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PARKING but expected PARKED trying again 26 more times
2021/02/24 18:33:27 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state PARKING but expected PARKED trying again 25 more times
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b parked
```

### unparking database

```
./bin/astra-cli db unpark 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
starting to unpark database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
2021/02/25 08:41:02 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 59 more times
2021/02/25 08:41:32 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 58 more times
2021/02/25 08:42:02 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 57 more times
2021/02/25 08:42:32 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 56 more times
2021/02/25 08:43:02 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 55 more times
2021/02/25 08:43:32 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 54 more times
2021/02/25 08:44:02 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 53 more times
2021/02/25 08:44:32 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 52 more times
2021/02/25 08:45:03 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 51 more times
2021/02/25 08:45:33 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 50 more times
2021/02/25 08:46:03 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 49 more times
2021/02/25 08:46:33 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 48 more times
2021/02/25 08:47:03 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 47 more times
2021/02/25 08:47:33 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 46 more times
2021/02/25 08:48:03 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 45 more times
2021/02/25 08:48:33 db 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b in state UNPARKING but expected ACTIVE trying again 44 more times
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b unparked
```

### deleteting database

```
./bin/astra-cli db delete 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
starting to delete database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b
database 2c3bc0d6-5e3e-4d77-81c8-d95a35bdc58b deleted
```

### resizing

I did not have a paid account to verify this works, but you can see it succesfully starts the process

```
./bin/astra-cli db resize 72c4d35b-1875-495a-b5f1-97329d90b6c5 2                    
unable to unpark '72c4d35b-1875-495a-b5f1-97329d90b6c5' with error expected status code 2xx but had: 400 error was [map[ID:2.000009e+06 message:resizing is not supported for this database tier]]
