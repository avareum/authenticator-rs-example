# TODO:

- [x] create tx mockup which signed by `secrets/signer.{fund-pub}/version/latest`
- [x] create payload & signature mockup which signed by mock `service-caller`
- [x] initiate message queue connection
- [x] receive & parse request from message queue to SignerRequest
- [x] publish on-chain signature
- [x] SignerApp unit tests
- [x] ACL unit tests
- [x] reload ACL list before verifying payload signature
- [x] GCP cloud log (hook from `SignerRequestedResponseHandler`)
- [x] Signer app support partial signing
- [ ] Solana signer unit tests
- [ ] integrate with GCP pubsub (wait: `core`)
- [ ] implement Fund wallet creation APIs (feature in `/http`)
- [ ] implement Ethereum signer
