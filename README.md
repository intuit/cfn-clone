# cfn-clone

This tool allows you to clone a CloudFormation stack.

## Requirements

* [aws cli](http://aws.amazon.com/cli/) installed and in your PATH

## Usage

### Basic Clone

By default, cfn-clone will use the same template and parameters as the existing stack
```sh
cfn-clone -s source-stack-name -n new-stack-name
```

### Override Parameters

You have the ability to override parameters for the new stack.
```sh
cfn-clone -s source-stack-name -n new-stack-name -a FOO=BAR
```

### Override Template

You have the ability to overrid the template for the new stack.
```sh
cfn-clone -s source-stack-name -n new-stack-name -t ./new_template.json
```

### Config

cfn-clone will pass through the relevant AWS related environment variables to the aws cli.

These are:

* `AWS_DEFAULT_PROFILE`
* `AWS_DEFAULT_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_SECURITY_TOKEN`

