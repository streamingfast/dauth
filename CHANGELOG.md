# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Unreleased

### Added

* added `x-trace-id` header in grpc request to dauth

* Added back `secret` plugin support with the form `secret://this-is-the-secret-and-fits-in-the-host-field?[user_id=<value>]&[api_key_id=<value>]`.

* Added continuous authentication support, enable by setting `grpc://localhost:9018?continuous=true`

* Added `x-sf-meta` header to pass arbitrary metadata as trusted headers.

## 2020-03-21

### Changed

* *Breaking* `middleware/connect/NewAuthInterceptor` requires a `logger *zap.Logger` has its last argument now.
* *Breaking* `middleware/grpc/UnaryAuthChecker` requires a `logger *zap.Logger` has its last argument now.
* *Breaking* `middleware/grpc/StreamAuthChecker` requires a `logger *zap.Logger` has its last argument now.

* License changed to Apache 2.0
