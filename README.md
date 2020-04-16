dfuse Auth Library
------------------

This library is the common interface between JWT issue and consumer of those
emitted JWT tokens. It is part of [dfuse](https://github.com/dfuse-io/dfuse).

It should contain everything related to common data structures
between the two legs (issuer/consumer) as well as any method that make sense
to re-use on the consumer side.

All generation of actual JWT tokens is left out of this library.


## Plugins

* `null://`
* `secret://this-is-the-secret-as-the-hostname`
* `cloud-gcp://null/projects/eoscanada-public/locations/global/keyRings/eosws-api-auth/cryptoKeys/default/cryptoKeyVersions/1?quotaEnforce=true&quotaRedisAddr=1.2.3.4&quotaBlacklistUpdateInterval=5s`


## Reference

### Credentials

This is the core data structure of our JWT token(s). A `Credentials` is a wrapper
around a JWT token containing various extra information like:

- `tier`
- `scopes`
- `start_block`

And any other contextual information our JWT tokens are dealing with.

### Context

We provide `GetCredentials`,


## Contributing

**Issues and PR in this repo related strictly to the dauth library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/dfuse-io/dfuse#protocols)

**Please first refer to the general
[dfuse contribution guide](https://github.com/dfuse-io/dfuse#contributing)**,
if you wish to contribute to this code base.


## License

[Apache 2.0](LICENSE)
