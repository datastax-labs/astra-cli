# {{classname}}

All URIs are relative to *https://api.astra.datastax.com*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddKeyspace**](OperationsApi.md#AddKeyspace) | **Post** /v2/databases/{databaseID}/keyspaces/{keyspaceName} | Adds keyspace into database
[**CreateDatabase**](OperationsApi.md#CreateDatabase) | **Post** /v2/databases | Create a new database
[**GenerateSecureBundleURL**](OperationsApi.md#GenerateSecureBundleURL) | **Post** /v2/databases/{databaseID}/secureBundleURL | Obtain zip for connecting to the database
[**GetDatabase**](OperationsApi.md#GetDatabase) | **Get** /v2/databases/{databaseID} | Finds database by ID
[**ListAvailableRegions**](OperationsApi.md#ListAvailableRegions) | **Get** /v2/availableRegions | Returns supported regions and availability for a given user and organization
[**ListDatabases**](OperationsApi.md#ListDatabases) | **Get** /v2/databases | Returns a list of databases
[**ParkDatabase**](OperationsApi.md#ParkDatabase) | **Post** /v2/databases/{databaseID}/park | Parks a database
[**ResetPassword**](OperationsApi.md#ResetPassword) | **Post** /v2/databases/{databaseID}/resetPassword | Resets Password
[**ResizeDatabase**](OperationsApi.md#ResizeDatabase) | **Post** /v2/databases/{databaseID}/resize | Resizes a database
[**TerminateDatabase**](OperationsApi.md#TerminateDatabase) | **Post** /v2/databases/{databaseID}/terminate | Terminates a database
[**UnparkDatabase**](OperationsApi.md#UnparkDatabase) | **Post** /v2/databases/{databaseID}/unpark | Unparks a database

# **AddKeyspace**
> AddKeyspace(ctx, databaseID, keyspaceName)
Adds keyspace into database

Adds the specified keyspace to the database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 
  **keyspaceName** | **string**| Name of database keyspace | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateDatabase**
> CreateDatabase(ctx, body)
Create a new database

Takes a user provided databaseInfo and returns the uuid for a new database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**DatabaseInfoCreate**](DatabaseInfoCreate.md)| Definition of new database | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GenerateSecureBundleURL**
> CredsUrl GenerateSecureBundleURL(ctx, databaseID)
Obtain zip for connecting to the database

Returns a temporary URL to download a zip file with certificates for connecting to the database. The URL expires after five minutes.<p>There are two types of the secure bundle URL: <ul><li><b>Internal</b> - Use with VPC peering connections to use private networking and avoid public internet for communication.</li> <li><b>External</b> - Use with any connection where the public internet is sufficient for communication between the application and the Astra database with MTLS.</li></ul> Both types support MTLS for communication via the driver.</p>

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 

### Return type

[**CredsUrl**](CredsURL.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetDatabase**
> Database GetDatabase(ctx, databaseID)
Finds database by ID

Returns specified database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 

### Return type

[**Database**](Database.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListAvailableRegions**
> []AvailableRegionCombination ListAvailableRegions(ctx, )
Returns supported regions and availability for a given user and organization

Returns all supported tier, cloud, region, count, and capacitity combinations

### Required Parameters
This endpoint does not need any parameter.

### Return type

[**[]AvailableRegionCombination**](AvailableRegionCombination.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListDatabases**
> []Database ListDatabases(ctx, optional)
Returns a list of databases

Get a list of databases visible to the user

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***OperationsApiListDatabasesOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OperationsApiListDatabasesOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **include** | **optional.String**| Allows filtering so that databases in listed states are returned | [default to nonterminated]
 **provider** | **optional.String**| Allows filtering so that databases from a given provider are returned | [default to ALL]
 **startingAfter** | **optional.String**| Optional parameter for pagination purposes. Used as this value for starting retrieving a specific page of results | 
 **limit** | **optional.Int32**| Optional parameter for pagination purposes. Specify the number of items for one page of data | [default to 25]

### Return type

[**[]Database**](Database.md)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ParkDatabase**
> ParkDatabase(ctx, databaseID)
Parks a database

Parks a database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResetPassword**
> ResetPassword(ctx, body, databaseID)
Resets Password

Sets a database password to the one specified in POST body

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**UserPassword**](UserPassword.md)| Map containing username and password. The specified password will be updated for the specified database user | 
  **databaseID** | **string**| String representation of the database ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ResizeDatabase**
> ResizeDatabase(ctx, body, databaseID)
Resizes a database

Resizes a database. Total number of capacity units desired should be specified. Reducing a size of a database is not supported at this time.

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **body** | [**CapacityUnits**](CapacityUnits.md)| Map containing capacityUnits key with a value greater than the current number of capacity units (max increment of 3 additional capacity units) | 
  **databaseID** | **string**| String representation of the database ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **TerminateDatabase**
> TerminateDatabase(ctx, databaseID, optional)
Terminates a database

Terminates a database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 
 **optional** | ***OperationsApiTerminateDatabaseOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a OperationsApiTerminateDatabaseOpts struct
Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **preparedStateOnly** | **optional.Bool**| For internal use only.  Used to safely terminate prepared databases. | [default to false]

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnparkDatabase**
> UnparkDatabase(ctx, databaseID)
Unparks a database

Unparks a database

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **databaseID** | **string**| String representation of the database ID | 

### Return type

 (empty response body)

### Authorization

[BearerAuth](../README.md#BearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

