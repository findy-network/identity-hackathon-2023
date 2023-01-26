# Typescript example

## Define env variables

[Described here](../README.md#setup-env-variables-for-the-agency-connection)

Note that API server cert path can be empty if using trusted issuer.

```bash
# agency API server cert path (leave empty if trusted issuer)
export FCLI_TLS_PATH='/path/to/self-issued-cert'
```

## Start server

```bash
nvm use         # or use whichever compatible node version
npm install
npm run build
npm run dev     # start server in watch mode
```
