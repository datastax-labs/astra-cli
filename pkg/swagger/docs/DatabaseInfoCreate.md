# DatabaseInfoCreate

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** | Name of the database--user friendly identifier | [default to null]
**Keyspace** | **string** | Keyspace name in database | [default to null]
**CloudProvider** | **string** | This is the cloud provider where the database lives. | [default to null]
**Tier** | **string** | With the exception of classic databases, all databases are serverless. Classic databases can no longer be created with the DevOps API. | [default to null]
**CapacityUnits** | **int32** | Capacity units were used for classic databases, but are not used for serverless databases. Enter 1 CU for serverless databases. Classic databases can no longer be created with the DevOps API. | [default to null]
**Region** | **string** | Region refers to the cloud region. | [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

