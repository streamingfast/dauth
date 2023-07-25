# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 2020-07-25

* Added
- Added continuous authentication support, enable by setting `grpc://localhost:9018?continuous=true`

## 2020-03-21

### Changed

* *Breaking* `middleware/connect/NewAuthInterceptor` requires a `logger *zap.Logger` has its last argument now.
* *Breaking* `middleware/grpc/UnaryAuthChecker` requires a `logger *zap.Logger` has its last argument now.
* *Breaking* `middleware/grpc/StreamAuthChecker` requires a `logger *zap.Logger` has its last argument now.

* License changed to Apache 2.0