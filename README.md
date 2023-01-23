# identity-hackathon-2023

Sample codes for Findy Agency service agents.

These sample servers demonstrate how
an agency client can issue and verify credentials using the Findy Agency API.
The sample server exposes two endpoints `/issue` and `/verify` that both
render an HTML page with QR code. The QR code can be read using [web wallet](https://github.com/findy-network/findy-wallet-pwa).
Once a pairwise connection is established between the server and the wallet user,
servers sends either a credential (`/issue`) or proof request (`/verify`) to the user.

## Run the sample

### Setup env variables for the agency connection

Use values for your cloud agency installation.
Default values point to localhost installation.

```bash
# agency authentication service URL
export AGENCY_AUTH_URL='http://localhost:8088'

# agency authentication origin
export AGENCY_AUTH_ORIGIN='http://localhost:3000'

# desired agent user name
export AGENCY_USER_NAME='ts-example'

# desired agent authentication key (create new key: 'findy-agent-cli new-key')
export AGENCY_KEY='15308490f1e4026284594dd08d31291bc8ef2aeac730d0daf6ff87bb92d4336c'

# agency API server
export AGENCY_API_SERVER_ADDRESS='localhost'

# agency API server port
export AGENCY_API_SERVER_PORT='50052'

# agency API server cert path
export AGENCY_API_SERVER_CERT_PATH='/path/to/self-issued-cert'
```

If you need to download the server cert from a cloud installation, you can use script:

```bash
./tools/dl-cert.sh <server_address>:<server_port>

# example:
./tools/dl-cert.sh agency-api.example.com:50051
```

### Run the sample server

* [Go](./go/README.md)
* [Typescript](./ts/README.md)

Note that server start may take a while at first run, because new credential definition
is registered on the ledger.

### Testing

1. Open URL <http://localhost:3001/issue> with browser.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the credential sent from this server.
1. Open URL <http://localhost:3001/verify> with browser.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the proof request sent from this server.

## Running agency on localhost

```bash
cd tools/local-env
docker-compose up
```
