# TODO:

- [x] create tx mockup which signed by `secrets/signer.{fund-pub}/version/latest`
- [x] create payload & signature mockup which signed by mock `service-caller`
- [x] initiate message queue connection
- [x] receive & parse request from message queue to SignerRequest
- [x] publish on-chain signature
- [x] SignerApp unit tests
- [ ] ACL unit tests
- [ ] Solana signer unit tests
- [ ] integrate with GCP pubsub
- [ ] GCP cloud log (hook from `SignerRequestedResponseHandler`)
- [ ] implement Fund wallet creation APIs (mq features)
- [ ] implement Ethereum signer
