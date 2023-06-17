# AWS ENV

Lambda runtime wrapper for exporting SystemsManager ParameterStore params & SecretsManager JSON secrets to the ENV

## Purpose

Instead of storing sensitive ENV variables in your Lambda function configuration, you might use ParameterStore or SecretsManager to keep sensitive values. You can use this tool to load those resources into the ENV through a Lambda [runtime wrapper](https://docs.aws.amazon.com/lambda/latest/dg/runtimes-modify.html#runtime-wrapper) script.

## Usage

Download the latest version of the `aws-env` binary from the [releases](https://github.com/amancevice/aws-env/releases) page.

Or, build it yourself with `make build`.

Include the binary in your lambda package or create a layer from the binary.

Export the ENV variable `AWS_LAMBDA_EXEC_WRAPPER` with your desired invocation using an absolute path to the binary.

> Note that if you include the binary in a Lambda layer the path will be `/opt/aws-env`, otherwise it will be found under `/var/task` wherever in your package you have included it (eg, `bin/aws-env`).

## Configuration

You can use a the ENV variable `AWS_ENV_EXPORTS` and/or a YAML configuration file to export the desired resources to ENV.

### ENV variable

Set the variable `AWS_ENV_EXPORTS` as a comma-delimited list of resources to export.

A resource should formatted like a URI, using the scheme for the service where the resource lives.

Examples:

- `secretsmanager://my-secret/`
- `secretsmanager://my-other-secret/`
- `ssm://my/path/`
- `ssm://my/other/path/`

Example ENV var:

```bash
AWS_ENV_EXPORTS=secretsmanager://my-secret/,ssm://my/path/
```

> Note that `ssm://` resources _must_ end with a trailing `/`

### Config File

You can include a config file named `.aws` in your lambda package that contains the parameters/secrets you wish to export.

By default this file is expected to be found at `/var/task/.aws`, but this can be overridden using the ENV variable `AWS_ENV_CONFIG`, eg `AWS_ENV_CONFIG=/var/task/.config/aws`

Example Config:

```yaml
---
exports:
  - secretsmanager: my-secret
  - secretsmanager: my-other-secret
  - ssm: /my/path/
  - ssm: /my/other/path/
```
