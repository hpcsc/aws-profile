# AWS Profile Utils

### Build Status
[![Build Status](https://travis-ci.org/hpcsc/aws-profile-utils.png)](https://travis-ci.org/hpcsc/aws-profile-utils)

### Installation

**Following steps assumes jq is available**

- Latest build from master branch


```
latest_build_number=$(curl https://api.travis-ci.org/repos/hpcsc/aws-profile-utils/branches/master | jq -r '.branch.number')
curl https://storage.googleapis.com/aws-profile-utils-master/aws-profile-utils-osx-${latest_build_number} -o aws-profile-utils
chmod +x aws-profile-utils && mv ./aws-profile-utils /usr/local/bin

```

- Release build

```
latest_release_tag=$(curl https://api.github.com/repos/hpcsc/aws-profile-utils/releases/latest | jq -r '.tag_name')
curl https://github.com/hpcsc/aws-profile-utils/releases/download/${latest_release_tag}/aws-profile-utils -o aws-profile-utils
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
