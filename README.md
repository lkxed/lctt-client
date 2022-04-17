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

（进度：仅实现了自动选题。）

### 到哪里去

新的思考：

貌似不需要使用本地 Git，直接使用 GitHub 的 API 来交互即可。这样一来，可以减少 Clone 和 Pull/Push 操作，从而减少磁盘/网络 IO，继而降低出错概率。\
如此一来，网络 IO 是最频繁的操作，主要体现在上传/下载/更新文件上。但即使最频繁，均摊下来也是很少的，因为短暂的 IO 过后是相对漫长的翻译过程。\
对服务器要求降低（除了上网姿势/地区），服务端的中心化就有了可能，实现一个全平台的客户端就有了可能（不过移动端貌似没有必要）。

新的实现流程：

* 自动选题：获取 RSS Feed -> 选取最新的文章 -> 解析生成 Markdown 文件 -> Fork 仓库 -> 新建分支 -> 上传文件 -> 创建 PR -> 合并后删除分支。
* 申领原文：获取待翻译列表 -> 添加译者信息 -> 更新文件 -> 创建 PR -> 合并后删除分支。
* 提交译文：更新文件 -> 移动文件 -> 创建 PR -> 合并后删除分支。

（进度：未开启。）

## 快速上手

### 创建 GitHub Token
参考这篇 [官方指南][8]。

（待完善）

[1]: https://github.com/LCTT/lctt-scripts
[2]: https://github.com/lkxed
[3]: https://github.com/PuerkitoBio/goquery
[4]: https://github.com/go-git/go-git
[5]: https://github.com/google/go-github
[6]: https://github.com/urfave/cli
[7]: https://github.com/SlyMarbo/rss
[8]: https://docs.github.com/cn/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-token
