# identity-hackathon-2023

Sample codes for [Findy Agency](https://findy-network.github.io) service agents.

These sample scripts and servers demonstrate how
an agency client can issue and verify credentials using the Findy Agency API.

## Run the sample

### Setup env variables for the agency connection

Use values for your cloud agency installation.
You can get the cloud agency configuration from the cloud agency maintainer.
Default values point to localhost installation.

```bash
# agency authentication service URL
export FCLI_URL='http://localhost:8088'

# agency authentication origin
export FCLI_ORIGIN='http://localhost:3000'

# desired agent user name
# note: this should be an unique string within agency context,
# use for example your email address
export FCLI_USER='my-very-own-issuer@example.com'

# desired agent authentication key (create new key: 'findy-agent-cli new-key')
# note: this key authenticates your client to agency, so keep it secret
export FCLI_KEY='15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c'

# agency API server
export AGENCY_API_SERVER='localhost'

# agency API server port
export AGENCY_API_SERVER_PORT='50052'

# API server address for CLI
export FCLI_SERVER="$AGENCY_API_SERVER:$AGENCY_API_SERVER_PORT"

# agency API server cert path
export FCLI_TLS_PATH='/path/to/self-issued-cert'
```

If you need to download the server cert from a cloud installation, you can use the script `dl-cert.sh`:

```bash
./tools/dl-cert.sh "$FCLI_SERVER"
```

### Run the CLI example

The sample script utilizes [findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation)
CLI tool for issuing and verifying credentials.

* [CLI](./cli/README.md)

### ...or run the sample server

The sample server exposes two endpoints `/issue` and `/verify` that both
render an HTML page with QR code. The QR code can be read using [web wallet](https://github.com/findy-network/findy-wallet-pwa).
Once a pairwise connection is established between the server and the wallet user,
servers sends either a credential (`/issue`) or proof request (`/verify`) to the user.

* [Go](./go/README.md)
* [Kotlin](./kotlin/README.md)
* [Typescript](./ts/README.md)

Note that server start may take a while at first run, because the new credential definition
is registered on the ledger.

#### Testing the server

![Server](https://user-images.githubusercontent.com/29113682/215501289-29fbf029-f796-487b-8370-6255d036e50d.gif)

1. Open URL <http://localhost:3001/issue> with browser.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the credential sent from this server.
1. Open URL <http://localhost:3001/verify> with browser.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the proof request sent from this server.

## Running agency on localhost

NOTE: setting up a localhost agency is not needed when using a cloud agency.

```bash
cd tools/local-env
docker-compose up
```
