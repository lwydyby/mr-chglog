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

#### 参数说明

``` bash
USAGE:
  mr-chglog [options] <tag query>

    There are the following specification methods for <tag query>.

    1. <old>..<new> - MR contained in <old> tags from <new>.
    2. <name>..     - MR from the <name> to the latest tag.
    3. ..<name>     - MR from the oldest tag to <name>.
    4. <name>       - MR contained in <name>.
```

## TODO

- [ ] 支持更丰富的模板
- [x] 支持解析MR的描述,更完善的CHANGELOG
- [x] 直接将MR Diff发送给Chatgpt(或者其他AI) 生成更准确的CHANGELOG
- [x] 支持自动生成TAG 并上传CHANGELOG

## Thanks
[git-chlog](https://github.com/git-chglog/git-chglog)