# Kotlin example

## Setup authentication to GitHub registry

Create personal access token with packages:read-permission.

Declare env variables

```bash
export USERNAME=<your_gh_username>
export TOKEN=<personal_access_token>
```

## Define env variables

[Described here](../README.md#setup-env-variables-for-the-agency-connection)

Note that API server cert path can be empty if using trusted issuer.

```bash
# agency API server cert path (leave empty if trusted issuer)
export AGENCY_API_SERVER_CERT_PATH='/path/to/self-issued-cert'
```

## Install CLI

## Run server

```bash
gradle bootRun
```
