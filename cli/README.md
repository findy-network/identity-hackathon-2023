# CLI example

## Flow

### Initialization

```mermaid
sequenceDiagram
  autonumber
    participant CLI
    participant Agency
    participant Web Wallet

    CLI->>Agency: Register
    CLI->>Agency: Login
    Agency-->>CLI: JWT token
    CLI->>Agency: Create schema
    Agency-->>CLI: Schema id
    CLI->>Agency: Create cred def
    Agency-->>CLI: Cred def id
    Web Wallet->>Agency: Register
    Web Wallet->>Agency: Login
```

### Issue credential

```mermaid
sequenceDiagram
  autonumber
    participant CLI
    participant Issue Bot
    participant Agency
    participant Web Wallet

    CLI->>Agency: Create invitation
    Agency-->>CLI: Invitation URL
    CLI->>Issue Bot: Start
    CLI-->>Web Wallet: <<show QR code>
    Web Wallet->>Agency: Read QR code
    Agency-->>Issue Bot: Connection ready!
    Issue Bot->>Agency: Issue credential
    Agency-->>Web Wallet: Accept offer?
    Web Wallet->>Agency: Accept
    Agency-->>Issue Bot: Issue ready!
    Issue Bot->>Issue Bot: Terminate
```

### Verify proof

```mermaid
sequenceDiagram
  autonumber
    participant CLI
    participant Verify Bot
    participant Agency
    participant Web Wallet

    CLI->>Agency: Create invitation
    Agency-->>CLI: Invitation URL
    CLI->>Verify Bot: Start
    CLI-->>Web Wallet: <<show QR code>
    Web Wallet->>Agency: Read QR code
    Agency-->>Verify Bot: Connection ready!
    Verify Bot->>Agency: Proof request
    Agency-->>Web Wallet: Accept request?
    Web Wallet->>Agency: Accept
    Agency-->>Verify Bot: Proof paused
    Verify Bot->>Agency: Resume proof
    Agency-->>Verify Bot: Proof ready!
    Verify Bot->>Verify Bot: Terminate
```


## Requirements

* [findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation)
* `qrencode`

    Mac:

    ```bash
    brew install qrencode
    ```

## Define env variables

[Described here](../README.md#setup-env-variables-for-the-agency-connection)

## Run the script

```bash
./run.sh
```

You can ignore the rpc errors that CLI prints out after canceling issuer/verifier bot.

## Testing

1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the credential sent from the CLI. Hit C-c once the credential is saved.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the proof request sent from this server. Hit C-c once the proof request is accepted.

## More examples

* [Findy Agency demo](https://github.com/findy-network/findy-agency-demo)
* [CLI usage examples](https://github.com/findy-network/findy-agent-cli#cli-usage-examples)
