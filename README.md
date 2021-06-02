dp-areas-api
================
Digital Publishing areas API used to navigate profiles for different geographical areas across the UK

### Getting started

* Run `make debug`

### Dependencies

* No further dependencies other than those defined in `go.mod`

### Configuration

| Environment variable         | Default   | Description
| ---------------------------- | --------- | -----------
| BIND_ADDR                    | :25500    | The host and port to bind to
| GRACEFUL_SHUTDOWN_TIMEOUT    | 5s        | The graceful shutdown timeout in seconds (`time.Duration` format)
| HEALTHCHECK_INTERVAL         | 30s       | Time between self-healthchecks (`time.Duration` format)
| HEALTHCHECK_CRITICAL_TIMEOUT | 90s       | Time to wait until an unhealthy dependent propagates its state to make this app unhealthy (`time.Duration` format)
| MONGODB_AREAS_DATABASE       | areas     | The MongoDB areas database
| MONGODB_AREAS_COLLECTION     | areas     | The MongoDB areas collection
| MONGODB_USERNAME             | test      | The MongoDB Username
| MONGODB_PASSWORD             | test      | The MongoDB Password
| MONGODB_CA_FILE_PATH         | file-path | The MongoDB CA FilePath
| DEFAULT_LIMIT                | 20        | Default limit for pagination
| DEFAULT_OFFSET               | 0         | Default offset for pagination
| DEFAULT_MAXIMUM_LIMIT        | 1000      | Default maximum limit for pagination

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2021, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.

