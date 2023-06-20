# StreamingFast Auth Library

[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/dauth)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

## Standard headers

These headers can be passed to the middleware and exposed to the app as **TrustedHeaders**

* `x-sf-user-id`
* `x-sf-api-key-id`
* `x-real-ip`

More headers can be passed.

## Plugins

The following plugins are provided by this package:

* `trust://` This will trust the incoming HTTP headers as-is
* `grpc://hostname:port` This will send the incoming HTTP headers to a grpc service which returns trusted headers (see [proto definitions](/proto/sf/authentication/v1/authentication.proto)

## Contributing

**Issues and PR in this repo related strictly to the dauth library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[dfuse contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

## License

[Apache 2.0](LICENSE)
