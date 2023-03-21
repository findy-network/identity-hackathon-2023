# Typescript example

## Define env variables

[Described here](../README.md#setup-env-variables-for-the-agency-connection)

Note that API server cert path can be empty if using trusted issuer.

```bash
# agency API server cert path (leave empty if trusted issuer)
export FCLI_TLS_PATH='/path/to/self-issued-cert'
```

## Optional: Start dev container

## Configure env variables

Running for the first time:

```bash
source <(curl localhost:3000/set-env.sh)
```

This script will create `.envrc` that will export needed variables.
You can recreate the environment with `direnv allow` or `source .envrc`

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
