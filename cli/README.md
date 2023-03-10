# CLI example

Read [the introduction text](https://findy-network.github.io/blog/2023/01/30/getting-started-with-ssi-service-agent-development/).

## Requirements

* [findy-agent-cli](https://github.com/findy-network/findy-agent-cli#installation)
* `qrencode`

    Mac:

    ```bash
    brew install qrencode
    ```

## Define env variables

[Described here](../README.md#setup-env-variables-for-the-agency-connection)

## Run the script

```bash
./run.sh
```

## Testing

1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the credential sent from the CLI.
1. Read the QR code with your wallet application or
paste the invitation url to the "Add Connection" dialog input field.
1. Accept the proof request sent from this server.

## More examples

* [Findy Agency demo](https://github.com/findy-network/findy-agency-demo)
* [CLI usage examples](https://github.com/findy-network/findy-agent-cli#cli-usage-examples)

## Troubleshooting

If your CLI issuer fails to send the credential offer,
try to delete the generated `CRED_DEF_ID`-file and rerun the script.
Missing `CRED_DEF_ID`-file will force the script to recreate the credential definition.
You may encounter this problem if, for example,
you have initially created the credential definition with another user.
