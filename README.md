# mr-chglog

## 简介

读取Gitlab的MR title生成CHANGELOG

## 安装

```bash
bash build.sh
```

## 使用说明

### 初始化

```bash
mr-chlog --init
```

### 生成CHANGELOG

```bash
mr-chlog
```

### 发送Release通知到飞书群

#### 命令行方式
```bash
mr-chlog --bot --app_id xxxx --app_secret xxxxxxx \
 --chat_id xxxxxxx --bot_title xxxxxxx --config /etc/config/mr_config.yml \
 --repository-url xxxx --token xxxx {tag}
```

#### 参数说明

``` bash
USAGE:
  mr-chglog [options] <tag query>

    There are the following specification methods for <tag query>.

    1. <old>..<new> - MR contained in <old> tags from <new>.
    2. <name>..     - MR from the <name> to the latest tag.
    3. ..<name>     - MR from the oldest tag to <name>.
    4. <name>       - MR contained in <name>.

OPTIONS:
  --init                         generate the mr-chglog configuration file in interactive (default: false)
  --app_id value                 feishu robot app secret
  --app_secret value             feishu robot app secret
  --chat_id value                feishu robot send group chat_id,Please use , to separate multiple
  --bot_title value              feishu robot send release title
  --path value [ --path value ]  Filter commits by path(s). Can use multiple times.
  --config value, -c value       specifies a different configuration file to pick up (default: ".chglog/mr_config.yml")
  --template value, -t value     specifies a template file to pick up. If not specified, use the one in config
  --repository-url value         specifies git repo URL. If not specified, use 'repository_url' in config
  --token value                  specifies git repo token. If not specified, use 'token' in config
  --output value, -o value       output path and filename for the changelogs. If not specified, output to stdout
  --next-tag value               treat unreleased commits as specified tags (EXPERIMENTAL)
  --create-tag value             create tag by CHANGELOG
  --ai                           use ai create CHANGELOG (default: false)
  --ai-type value                which ai API to use create CHANGELOG (default: poe)
  --bot                          push mr-chglog changelog to feishu group (default: false)
  --help, -h                     show help
  --version, -v                  print the version
  
EXAMPLE:

  $ mr-chglog

    If <tag query> is not specified, it corresponds to all tags.
    This is the simplest example.

  $ mr-chglog 1.0.0..2.0.0

    The above is a command to generate CHANGELOG including MR of 1.0.0 to 2.0.0.

  $ mr-chglog 1.0.0

    The above is a command to generate CHANGELOG including MR of only 1.0.0.

  $ mr-chglog $(git describe --tags $(git rev-list --tags --max-count=1))

    The above is a command to generate CHANGELOG with the MR included in the latest tag.

  $ mr-chglog --output CHANGELOG.md

    The above is a command to output to CHANGELOG.md instead of standard output.

  $ mr-chglog --config custom/dir/config.yml

    The above is a command that uses a configuration file placed other than ".chglog/config.yml".

  $ mr-chglog --path path/to/my/component --output CHANGELOG.component.md

    Filter commits by specific paths or files in git and output to a component specific changelog.
  $ mr-chglog --bot 
    Push mr-chglog Changelog to Feishu Group
  $ mr-chglog --ai 
    Use ai create CHANGELOG

```

## TODO

- [x] 支持更丰富的模板
- [x] 支持解析MR的描述,更完善的CHANGELOG
- [x] 直接将MR Diff发送给Chatgpt(或者其他AI) 生成更准确的CHANGELOG
- [x] 支持自动生成TAG 并上传CHANGELOG

## Thanks
[git-chlog](https://github.com/git-chglog/git-chglog)