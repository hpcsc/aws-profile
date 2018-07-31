# AWS Profile Utils

### Build Status
[![Build Status](https://travis-ci.org/hpcsc/aws-profile-utils.png)](https://travis-ci.org/hpcsc/aws-profile-utils)

Commands:

### Get current AWS profile

```
aws-profile-utils get
```

Parameters

- `credentials-path` (optional, default to ~/.aws/credentials): path to AWS CLI credentials file
- `config-path` (optional, default to ~/.aws/config): path to AWS CLI config file

### Set default AWS profile

**Note**: this command assumes `fzf` is setup and `fzf-tmux` executable is available in your shell

```
aws-profile-utils set
```

Parameters

- `credentials-path` (optional, default to ~/.aws/credentials): path to AWS CLI credentials file
- `config-path` (optional, default to ~/.aws/config): path to AWS CLI config file
- `pattern` (optional, default to empty string): start the fzf finder with the given pattern

### Get version of this util

```
aws-profile-utils version
```
