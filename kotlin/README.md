# Kotlin example

## Setup authentication to GitHub registry

[Create personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic) with read:packages-permission.

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
export FCLI_TLS_PATH='/path/to/self-issued-cert'
```

## Install CLI

Currently, agency Kotlin wrapper uses Findy Agency CLI tool for authentication.
The environment should have [findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation) in `PATH`.

## Run server

```bash
gradle bootRun
```
