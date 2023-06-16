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

You can use a command-line options and/or a YAML configuration file to export the desired resources to ENV.

Configure `aws-env` by exporting your desired invocation as the variable `AWS_LAMBDA_EXEC_WRAPPER` in your Lambda runtime.

### Command Line Options

You can use command line options to choose with params/secrets to export.

Export ParameterStore parameters from one or more parameter paths:

```bash
/opt/aws-env --path /my/params/ --path ...
```

Export one or more SecretsManager secrets:

```bash
/opt/aws-env --secret-id my-secret --secret-id ...
```

Export a combination of parameters and secrets:

```bash
/opt/aws-env --path /my/path/ --secret-id my-secret ...
```

### Config File

You can include a config file named `.aws` in your lambda package that contains the parameters/secrets you wish to export.

By default this file is expected to be found at `/var/task/.aws`, but this can be overridden on the command line using the `--config` or `-c` option

Example Config:

```yaml
---
paths:
  - /my/path/
  - /my/other/path/
secrets:
  - my-secret
  - my-other-secret
```
