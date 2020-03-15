# AWS Profile

### Build Status
![Build Status](https://github.com/hpcsc/aws-profile/workflows/Pipeline/badge.svg)

[![Demo](https://github.com/hpcsc/aws-profile/raw/master/aws-profile.gif)](https://github.com/hpcsc/aws-profile/raw/master/aws-profile.gif)

### Installation

- Latest build from master branch: [Bintray](https://dl.bintray.com/hpcsc/aws-profile)

- Release build [Github Releases](https://github.com/hpcsc/aws-profile/releases)

After downloading binary file, rename it to `aws-profile` (or `aws-profile.exe` on Windows), `chmod +x` and move the executable to a location in your `PATH` (.e.g. `/usr/local/bin` for Linux/MacOS):

```
chmod +x aws-profile && mv ./aws-profile /usr/local/bin
```

### How it works

- The `get` and `set` commands work primarily on AWS `config` and `credentials` files.
- `set` command sets `default` profile in either `config` or `credentials` file with values (e.g. `aws_access_key_id` and `aws_secret_access_key` or `role_arn` and `source_profile`) from selected profile.
- `get` command first checks whether AWS credentials environment variables (e.g. `AWS_ACCESS_KEY_ID`, `AWS_SESSION_TOKEN`) are set. If yes, it will do a call to STS to get caller identity and cache the result locally. If those environment variables are not set, it compares values of `default` profile with other profiles in `config` and `credentials` files and returns the matched profile
- `export` command prints out suitable command for your OS (`export` in Linux/MacOS or `$env:VAR` setting in Windows Powershell). These printed commands can be copied and executed directly in your terminal to set suitable AWS environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`, `AWS_SESSION`). The purpose of this command is to support some of the tools that don't work well with AWS `config` and `credentials` files, e.g.
    - Terraform aws provider with role that requires MFA [https://github.com/terraform-providers/terraform-provider-aws/issues/2420](https://github.com/terraform-providers/terraform-provider-aws/issues/2420)
    - Or when you want to execute AWS CLI commands inside a container and it's not convenient to mount host machine `~/.aws` folder

### Usage

```
usage: aws-profile [<flags>] <command> [<args> ...]

simple tool to help switching among AWS profiles more easily

Flags:
  -h, --help  Show context-sensitive help (also try --help-long and --help-man).
      --credentials-path="~/.aws/credentials"
              Path to AWS Credentials file
      --config-path="~/.aws/config"
              Path to AWS Config file

Commands:
  help [<command>...]
    Show help.

  get
    get current AWS profile

  set [<pattern>]
    set default profile with credentials of selected profile

  export [<flags>] [<pattern>]
    print commands to set environment variables for assuming a AWS role

    To execute the command without printing it to console:

    - For Linux/MacOS, execute: "eval $(aws-profile export)"

    - For Windows, execute: "Invoke-Expression (path\to\aws-profile.exe export)"

  unset
    print commands to unset AWS credentials environment variables

    To execute the command without printing it to console:

    - For Linux/MacOS, execute: "eval $(aws-profile unset)"

    - For Windows, execute: "Invoke-Expression (path\to\aws-profile.exe unset)"

  version
    show aws-profile version
```

### Sample AWS config/credentials files

`~/.aws/credentials`

```
[default]

[some-profile]
aws_access_key_id     = xxx
aws_secret_access_key = yyy
```

`~/.aws/config`

```
[default]

[profile role-with-mfa]
role_arn       = arn:aws:iam::xxxxxxxxxxxx:role/role-with-mfa-enabled
source_profile = hpcsc
mfa_serial     = arn:aws:iam::xxxxxxxxxxxx:mfa/my-mfa-device

[profile role-without-mfa]
role_arn       = arn:aws:iam::xxxxxxxxxxxx:role/role-without-mfa
source_profile = hpcsc
```
