# TODO:

- [x] create tx mockup which signed by `secrets/signer.{fund-pub}/version/latest`
- [x] create payload & signature mockup which signed by mock `service-caller`
- [x] initiate message queue connection
- [x] receive & parse request from message queue to SignerRequest
- [x] publish on-chain signature
- [x] SignerApp unit tests
- [x] ACL unit tests
- [x] ~~reload ACL list before verifying payload signature~~
- [x] GCP cloud log (hook from `SignerRequestedResponseHandler`)
- [x] signer app support partial signing
- [x] implement fund wallet creation APIs (feature in `/http`)
- [x] support both REST/PubSub endpoints
- [x] solana signer unit tests
- [x] implement Ethereum signer
- [x] ethereum signer unit tests
- [x] deploy scripts
- [ ] CI
- [ ] parse SignerRequest in `/internal/server/api/rest.go#42`
- [ ] setting IAM
  - [ ] read in `Keystores` section
  - [ ] 0.0.0.0:8080 access availabled for core service

## Keystores

- **Service keys**:
  - name: `SERVICE_{service_name}`
  - permissions:
    - read: service(s), app_signer
    - write: admin
- **Fund wallet keys**:
  - name: `WALLET_{wallet}`
  - permissions:
    - read: app_signer
    - write: app_signer
