# AWS SecretsManager ENV

Lambda runtime wrapper for exporting a SecretsManager JSON secret to the ENV

## Purpose

Instead of storing sensitive ENV variables in your Lambda function configuration, you might store a JSON document containing sensitive variables and their values and then load that secret in a Lambda [runtime wrapper](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-modify.html#runtime-wrapper) script

Example secret JSON:

```json

{
  "SOME_API_KEY": "F!ZZ",
  "SOME_SECRET": "B@ZZ"
}
```

Including `aws-secretsmanager-env` in your functions allows you to export the contents of a secret JSON document into the ENV during the runtime init phase of the function lifecycle.

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
