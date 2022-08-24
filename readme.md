# Avareum Hubble Signer

Avareum fund operation signing modules

See: [Notes](/note.md)

## Run

1. Prepare `GOOGLE_APPLICATION_CREDENTIALS` in `.env`
2. Run using these methods

   2.1. Run using VSCode by `Run and Debug` -> `SignerApp`

   2.2. Run using docker container

```
$ ./make.bash build
$ ./make.bash run
```

3. Access swagger API docs on `http://localhost:8080/swagger/index.html`

## Example: App signer sequence diagram

```mermaid
sequenceDiagram
    autonumber
    Core-->>Endpoints: Message data receive
    Note left of Endpoints: support MQ/REST
    Endpoints->>+App: Foward message data
    App->>+Signer: Parse message data <br/>to SignerRequest
    Note right of App: chain,id,fund,payload,<br/>signature,caller
    Signer->>ACL: Verify payload signature
    Signer->>Secret: Request signing key
    Secret->>Signer: Return signing key
    Signer->>Signer: Decode payload & sign
    Signer->>Blockchain: Broadcast
    Signer->>-App: Return tx signatures
    App-->>-Endpoints: WIP: Publish tx status
```
