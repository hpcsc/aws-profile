# AWS Profile

### Build Status
![Build Status](https://github.com/hpcsc/aws-profile/workflows/Pipeline/badge.svg) [![codecov](https://codecov.io/gh/hpcsc/aws-profile/branch/master/graph/badge.svg?token=76OSPJNMON)](https://codecov.io/gh/hpcsc/aws-profile)

[![Demo](https://github.com/hpcsc/aws-profile/raw/master/aws-profile.gif)](https://github.com/hpcsc/aws-profile/raw/master/aws-profile.gif)

### Installation

#### MacOS/Linux users

- Release build
  ```
  curl -sL https://raw.githubusercontent.com/hpcsc/aws-profile/master/install | sh
  ```

- Latest build from master branch
  ```
  curl -sL https://raw.githubusercontent.com/hpcsc/aws-profile/master/install | PRERELEASE=true sh
  ```

This will download latest release from Github to `/usr/local/bin/aws-profile`

#### Manual Installation

- Release build [Github Releases](https://github.com/hpcsc/aws-profile/releases/latest)
- Latest build from master branch: [ ![Download](https://api.bintray.com/packages/hpcsc/aws-profile/master/images/download.svg) ](https://bintray.com/hpcsc/aws-profile/master/_latestVersion#files)

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

For more information, please refer to [aws-profile wiki](https://github.com/hpcsc/aws-profile/wiki)
