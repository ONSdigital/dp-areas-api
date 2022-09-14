Update topic release date
================

### Dependencies

* No further dependencies other than configuration

### Configuration

| Environment variable         | Default                                           | Description
| ---------------------------- | ------------------------------------------------- | -----------
| MONGODB_BIND_ADDR            | localhost:27017                                   | The MongoDB bind address
| MONGODB_USERNAME             |                                                   | MongoDB Username
| MONGODB_PASSWORD             |                                                   | MongoDB Password
| MONGODB_DATABASE             | topics                                            | The MongoDB topics database
| MONGODB_COLLECTIONS          | TopicsCollection:topics,ContentCollection:content | MongoDB collections
| MONGODB_ENABLE_READ_CONCERN  | false                                             | Switch to use (or not) majority read concern
| MONGODB_ENABLE_WRITE_CONCERN | true                                              | Switch to use (or not) majority write concern
| MONGODB_CONNECT_TIMEOUT      | 5s                                                | The timeout when connecting to MongoDB (`time.Duration` format)
| MONGODB_QUERY_TIMEOUT        | 15s                                               | The timeout for querying MongoDB (`time.Duration` format)
| MONGODB_IS_SSL               | false                                             | Switch to use (or not) TLS when connecting to mongodb
