# authenticator

## Development

```shell
cargo build
# check output in /target/debug/authenticator
```

## Usage

```shell
# view usage
authenticator

# create a new key, then import key to authenticator app
authenticator create <email>
authenticator create user@email.com

# test key
authenticator test <secret>
authenticator test S7JS64F5LOGMR6MWXJMSV3GEV2BC6BE7

```
