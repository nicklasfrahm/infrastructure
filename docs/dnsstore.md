# DNS store

An experimental approach to use DNS as a distributed database.

`dnsstore` uses TXT records to store data. All it requires is a domain name and a DNS server. All data records are stored under the `_dnsstore` zone, so if your domain is `example.com`, all records will be records within the zone `_dnsstore.example.com`. You should be able to use `dnsstore` similar to a document store database.

The typical use case is to use `dnsstore` as a distributed database for a cluster of servers. It should only be used to store small amounts of data, where concurrent access is not an issue. It is not recommended to use `dnsstore` as a general purpose database.

## Metadata

The root record, `_dnsstore.example.com`, is used to store metadata about the `dnsstore`. The metadata is stored as a JSON string. The following metadata is stored:

```go
// Metadata decribes the metadata stored in the dnsstore.
type Metadata struct {
  // Version is the version of the schema.
  Version int `json:"version"`
  // Chunks contains metadata about available chunks.
  Chunks *ChunkMetadata `json:"chunks"`
}

// ChunkMetadata describes the metadata for a chunk in the dnsstore.
type ChunkMetadata struct {
  // Total is the total number of chunks.
  Total int `json:"total"`
}
```

Each child chunk of the `dnsstore` contains metadata about its databases. The metadata is stored as a JSON string. The following metadata is stored:

```go
// DatabaseChunk describes the metadata for a chunk in the dnsstore.
type DatabaseChunk struct {
  // ID is a positive integer that uniquely identifies the chunk.
  ID int `json:"id"`
  // Databases contains metadata about available databases.
  Databases map[string]*DatabaseMetadata `json:"database"`
}

// DatabaseMetadata describes the metadata for a database in the dnsstore.
type DatabaseMetadata struct {
  // Databases contains metadata about available databases. We use a map to speed up discovery.
  Databases map[string]*DatabaseMetadata `json:"database"`
}
```

The overall format may roughly be described as:

`<chunk>.<record>.<chunk>.<collection>.<chunk>.<database>.<chunk>._dnsstore.example.com`

## Databases

A database is a subdomain of the `_dnsstore` zone. For example, if you want to create a database for a service called `warehouse`, you would create a database under the zone `warehouse.0._dnsstore.example.com`. As previously mentioned, the root record `_dnsstore.example.com` is used to store metadata about the database. The database contains metadata about its collections.

### Collections

`dnsstore` uses the concept of collections to group data. A collection is a subdomain of the `_dnsstore` zone. For example, if you want to store data for a collection called `users`, you would store the data under the zone `users._dnsstore.example.com`. The `TXT` record `_dnsstore.example.com` is used to store the list of collections for automatic discovery of the available collections. The data is stored as a JSON string.
