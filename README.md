# LCTT 客户端

> 「LCTT 客户端」是专门为「Linux 中国」翻译组开发的全流程客户端，是 LCTT 新手译者的最佳拍档。

目前状态：开发中，仍在自用阶段。

## 简介

受 [LCTT/lctt-scripts][1] 启发。

### 为了谁

希望通过翻译贡献开源社区，但是对 Git & GitHub 不大熟悉，在实际操作中会遇到各种问题的朋友。

### 依靠谁

目前只有我 @[lkxed][2]。

### 从哪里来

目前的实现方案：

* [PuerkitoBio/goquery][3]：类似 JQuery 的解析库，用于 HTML 解析。
* [go-git/go-git][4]：纯 Go 实现的 Git，用于本地操作。
* [google/go-github][5]：GitHub API 的 Go 语言封装库，用于创建 PR。
* [urfave/cli][6]：一个 CLI 库，用于实现命令行客户端。
* [SlyMarbo/rss][7]：用于解析 RSS 订阅。

（进度：实现了手动选题 & 申领原文 & 提交译文，满足了我自己的基本需求。）

### 到哪里去

新的思考：

貌似不需要使用本地 Git，直接使用 GitHub 的 API 来交互即可。这样一来，可以减少 Clone 和 Pull/Push 操作，从而减少磁盘/网络 IO，继而降低出错概率。

如此一来，网络 IO 是最频繁的操作，主要体现在上传/下载/更新文件上。但即使最频繁，均摊下来也是很少的，因为短暂的 IO 过后是相对漫长的翻译过程。

对服务器要求降低（除了上网姿势/地区），服务端的中心化就有了可能，实现一个全平台的客户端就有了可能（不过移动端貌似没有必要）。

新的实现流程：

* 自动选题：定时获取 RSS Feed -> 选取最新的文章 -> 解析生成 Markdown 文件 -> Fork 仓库 -> 新建分支 -> 上传文件 -> 创建 PR -> 合并后删除分支。
* 申领原文：获取待翻译列表 -> 添加译者信息 -> 更新文件 -> 创建 PR -> 合并后删除分支。
* 提交译文：更新文件 -> 移动文件 -> 创建 PR -> 合并后删除分支。

（进度：未开启。）

## 快速上手

### 创建 GitHub Token

参考这篇 [官方指南][8]。

### 构建并生成可执行文件

目前，客户端仍在开发阶段，只有我一个人使用，所以没有 release 包。

你可以 Clone 源码 / 下载源码压缩包，运行以下命令构建并生成可执行文件：

```shell
go build -o bin/lctt cmd/lctt/*
```

### 客户端配置

配置文件：`configs/settings.yml`。以下是你必须修改的配置项：

```yaml
git:
  # 你在 git commit 中显示的身份，可以任意指定，非空即可
  user:
    name: lkxed
    email: lkxed@outlook.com
  local:
    # 你的本地仓库存放路径 (不存在会自动创建).
    repository: /Users/lkxed/Documents/GitHub/TranslateProject
  hub:
    # 你的 GitHub 用户名
    username: lkxed
    # 你刚刚创建的 GitHub Token
    access-token: ghp_kUa1diAnFoLLoWGr1atLkXEd0NG1tHu68Ha0D

# 你使用的编辑器的打开命令（如果你需要预览选题文章的话）
editor: "code -n -w"
# 你使用的浏览器的打开命令
browser: "open -a safari"
```

### 基本使用

目前客户端支持以下命令：

```
COMMANDS:
   init      （初始化客户端）Initializes the client.
   feed      （获取网站的最新文章列表）Feeds you a list of articles published recently.
   collect   （手动选题）Collects an article with its <CATEGORY> and <LINK>.
   request   （申领原文）Requests to translate an article with its <CATEGORY> and <FILENAME>.
   complete  （提交译文）Completes the translating process of an article with its <CATEGORY> and <FILENAME>.
   help, h   （显示帮助）Shows a list of commands or help for one command.
```

#### 初始化客户端

```
USAGE:
   lctt init
```

这个命令很简单，不需要指定选项/参数。在执行其他操作前，最好 `init` 一下确保当前仓库是最新的。

#### 获取网站最新列表

```
USAGE:
   lctt feed [command options] [arguments...]

OPTIONS:
   --since <DATE>, -s <DATE>
   --verbose, -v
   --open, -o
   --help, -h 
```

具体来说，有以下几个使用场景：

显示网站今天刚发布的文章列表：
```shell
bin/lctt feed
```

显示网站今天刚发布的文章列表，并在浏览器中查看它们的原文：

```shell
bin/lctt feed --open/-o
```

显示网站今天刚发布的文章列表，显示它的标签和摘要：

```shell
bin/lctt feed --verbose/-v
```

显示网站 2006 年 1 月 2 日之后发布的所有文章：

```shell
bin/lctt feed --since/-s 2006-01-02
```

以上选项可结合使用，如：

```shell
bin/lctt feed -vos 2006-01-02
```

这将显示获取网站 2006 年 1 月 2 日之后发布的所有文章，显示它们的标签和摘要，并在浏览器中查看它们的原文。

但是，请不要同时指定 `-o` 和 `-s <某个久远的日期>`，这将在浏览器中打开数百个窗口/标签，你的浏览器和操作系统可能会不堪重负。

#### 手动选题

当你看中了某一篇文章，你可以尝试选题。你需要指定文章的类别和链接，就像下面这样：

```shell
bin/lctt collect -c tech https://opensource.com/article/xxx/
```

这将生成符合 LCTT 规范的原文 Markdown 文件，存放在 `previews` 目录中。

如果你想要同时在编辑器中预览这篇文章，你需要指定 `--preview/-p` 选项。

如果你已经下定决心，想要一气呵成地完成选题，你需要指定 `--upload/-u` 选项。

一个典型的命令是：

```shell
bin/lctt collect -puc tech https://opensource.com/article/xxx/
```

这将生成符合 LCTT 规范的原文 Markdown 文件，存放在 `previews` 目录中，同时在你的编辑器中打开它。

当你确认格式无误后，返回客户端，并根据提示按下回车，上传原文，完成选题。

#### 申领原文

请查阅 `bin/lctt request --help/-h`。

（待补充）

#### 提交译文

请查阅 `bin/lctt complete --help/-h`

（待补充）

## 致谢

首先，感谢本项目中使用到的所有 Go 项目的贡献者，是他们让这个项目成为可能。

其次，感谢 @[lujun9972][9]，他维护的 lctt-scripts 给了我很大的启发。

最后，感谢 LCTT 组长 @[wxy][10] 对我和这个项目的大力支持。

[1]: https://github.com/LCTT/lctt-scripts
[2]: https://github.com/lkxed
[3]: https://github.com/PuerkitoBio/goquery
[4]: https://github.com/go-git/go-git
[5]: https://github.com/google/go-github
[6]: https://github.com/urfave/cli
[7]: https://github.com/SlyMarbo/rss
[8]: https://docs.github.com/cn/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-token
[9]: https://github.com/lujun9972
[10]: https://github.com/wxy
