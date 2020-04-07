## dfuse Auth Library

This repository contains stuff shared between JWT issuer service and any other
consumer that might needs to handle back JWT authentication/authorization (mainly
`eosws` for now).

### Plugins

* `null://`
* `secret://this-is-the-secret-as-the-hostname`
* `cloud-gcp://null/projects/eoscanada-public/locations/global/keyRings/eosws-api-auth/cryptoKeys/default/cryptoKeyVersions/1?quotaEnforce=true&quotaRedisAddr=1.2.3.4&quotaBlacklistUpdateInterval=5s`

### Philosophy

This library is the common interface between JWT issue and consumer of those
emitted JWT tokens.

This package should contains everything related to common data structures
between the two legs (issuer/consumer) as well as any method that make sense
to re-use on the consumer side.

All generation of actual JWT tokens is left out of this library.

### Reference

#### Credentials

This is the core data structure of our JWT token(s). A `Credentials` is a wrapper
around a JWT token containing various extra information like:

- `tier`
- `scopes`
- `start_block`

And any other contextual information our JWT tokens are dealing with.

#### Context

We provide `GetCredentials`,
