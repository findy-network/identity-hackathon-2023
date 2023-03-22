# Kotlin example

## Setup authentication to GitHub registry

[Create personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic) with read:packages-permission.

Declare env variables

```bash
export USERNAME=<your_gh_username>
export TOKEN=<personal_access_token>
```

## Start dev container

Start dev container in VSCode.

Alternatively make sure you have Java, Gradle and
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

## Run server

```bash
gradle bootRun
```

## More examples

* [Kotlin sample](https://github.com/findy-network/findy-common-kt)
