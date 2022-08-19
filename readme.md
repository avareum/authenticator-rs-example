# Avareum Hubble Signer

Avareum fund operation signing modules

See: [Notes](/note.md)

## Example: App signer sequence diagram

```mermaid
sequenceDiagram
    autonumber
    Core-->>MessageQueue: Message data receive
    MessageQueue->>+App: Foward message data
    App->>+Signer: Parse message data <br/>to SignerRequest
    Note right of App: chain,id,fund,payload,<br/>signature,caller
    Signer->>ACL: Verify payload signature
    Signer->>Secret: Request signing key
    Secret->>Signer: Return signing key
    Signer->>Signer: Decode payload & sign
    Signer->>Blockchain: Broadcast
    Signer->>-App: Return tx signatures
    App-->>-MessageQueue: WIP: Publish tx status
```
