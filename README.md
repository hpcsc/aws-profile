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

### Usage

```
usage: aws-profile [<flags>] <command> [<args> ...]

simple tool to help switching among AWS profiles more easily

Flags:
  -h, --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  get [<flags>]
    get current AWS profile (that is set to default profile)

  set [<flags>] [<pattern>]
    set default profile with credentials of selected profile

  export [<flags>] [<pattern>]
    print commands to set environment variables for assuming a AWS role

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
