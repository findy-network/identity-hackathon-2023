# Go example

## Start dev container

Start dev container in VSCode.

Alternatively make sure you have Go and
[findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation) installed.

## Configure env variables

Configure environment for the first time by utilizing a script from the agency environment
that will set the needed environment variables:

```bash
source <(curl <agency_url>/set-env.sh)
```

For cloud installation, use the cloud URL e.g. <https://agency.example.com>
For local installation, use the web wallet URL: <http://localhost:3000>

This script will create `.envrc` that will contain needed variables.
If you are working in the dev container or have `direnv` installed,
you can type `direnv allow` which will auto-export the needed variables each time
the terminal is opened.
Otherwise you should manually enter `source .envrc` when opening a new terminal window.

## Start server

```bash
go mod tidy -e
go run .
```

## More examples

* [CLI tool](https://github.com/findy-network/findy-agent-cli)
