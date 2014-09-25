# cfn-clone

This tool allows you to clone a CloudFormation stack.

## Installation

Just grab the [latest binary](https://github.com/thbishop/mmfa/releases/latest) and you should be good to go.

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

You have the ability to override the template for the new stack.
```sh
cfn-clone -s source-stack-name -n new-stack-name -t ./new_template.json
```

### Config

You can use the normal aws cli environment variables for controlling credentials, etc. When cfn-clone invokes the aws cli, these will be made available.

Examples:

* `AWS_DEFAULT_PROFILE`
* `AWS_DEFAULT_REGION`
* `AWS_ACCESS_KEY_ID`
* `AWS_SECRET_ACCESS_KEY`
* `AWS_SECURITY_TOKEN`

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

See `LICENSE` for details
