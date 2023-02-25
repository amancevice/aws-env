# AWS SecretsManager ENV

Lambda runtime wrapper for exporting a SecretsManager JSON secret to the ENV

## Usage

Download the latest version of the `aws-secretsmanager-env` binary from the releases page.

Or, build it yourself with `make build`.

Include the binary in your lambda package or create a layer from the binary.

Set the following environmental variables in your Lambda:

| ENV                       | Example                                             |
|:--------------------------|:----------------------------------------------------|
| `AWS_SECRET`              | _[your secret name]_                                |
| `AWS_LAMBDA_EXEC_WRAPPER` | `/opt/aws-secretsmanager-env` (Lambda layer)        |
| `AWS_LAMBDA_EXEC_WRAPPER` | `/var/task/aws-secretsmanager-env` (Lambda package) |

With these variables set, your lambda will export the given SecretsManager secret JSON to the Lambda runtime ENV.
