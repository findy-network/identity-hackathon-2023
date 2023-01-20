# Typescript example

## Setup env variables for the agency connection

Use values for your cloud agency installation. Default values point to localhost installation.

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

# agency API server cert path (leave empty if trusted issuer)
export AGENCY_API_SERVER_CERT_PATH='/path/to/self-issued-cert'
```

## Start server

```bash
nvm use         # or use whichever compatible node version
npm install
npm run build
npm run dev     # start server in watch mode
```

Note that server start may take a while at first run, because new credential definition
is registered on the ledger on server start.

## Testing

1. Open URL <http://localhost:3001/issue> with browser.
1. Read QR code with your wallet application.
1. Accept the credential sent from this server.
1. Open URL <http://localhost:3001/verify> with browser.
1. Read QR code with your wallet application.
1. Accept the proof request sent from this server.
