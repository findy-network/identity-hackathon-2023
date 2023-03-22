# Typescript example

## Start dev container

Start dev container in VSCode.

Alternatively make sure you have Node.js (nvm) and
[findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation) installed.

## Configure env variables

Configure environment for the first time by utilizing a script from the agency environment
that will set the needed environment variables:

```bash
source <(curl <agency_url>/set-env.sh)
```

For cloud installation, use the cloud URL e.g. https://agency.example.com
For local installation, use the web wallet URL: http://localhost:3000

This script will create `.envrc` that will contain needed variables.
If you are working in the dev container or have `direnv` installed,
you can type `direnv allow` which will auto-export the needed variables each time
the terminal is opened.
Otherwise you should manually enter `source .envrc` when opening a new terminal window.

## Start server

```bash
nvm use         # or use whichever compatible node version
npm install
npm run build
npm run dev     # start server in watch mode
```

## More TS+JS examples

* [Decentralized identity demo](https://github.com/findy-network/agency-demo)
* [Issuer tool](https://github.com/findy-network/findy-issuer-tool)
* [OIDC IdP](https://github.com/findy-network/findy-oidc-provider)
