# AWS Profile Utils

### Build Status
![Build Status](https://github.com/hpcsc/aws-profile-utils/workflows/Pipeline/badge.svg)

[![Demo](https://github.com/hpcsc/aws-profile-utils/raw/master/aws-profile-utils.gif)](https://github.com/hpcsc/aws-profile-utils/raw/master/aws-profile-utils.gif)

### Installation

- Latest build from master branch: [Bintray](https://dl.bintray.com/hpcsc/aws-profile-utils)

- Release build [Github Releases](https://github.com/hpcsc/aws-profile-utils/releases)

After downloading binary file, rename it to `aws-profile-utils`, `chmod +x` and move the executable to `/usr/local/bin`:

```
chmod +x aws-profile-utils && mv ./aws-profile-utils /usr/local/bin
```

### Usage

```
usage: aws-profile-utils [<flags>] <command> [<args> ...]

simple tool to help switching among AWS profiles more easily

Flags:
  -h, --help  Show context-sensitive help (also try --help-long and --help-man).

Commands:
  help [<command>...]
    Show help.

  get [<flags>]
    get current AWS profile (that is set to default profile)

  set [<flags>] [<pattern>]
    set default profile with credentials of selected profile (this command assumes fzf is already setup)

  version
    show aws-profile-utils version
```
