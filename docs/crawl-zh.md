# 使用 Zenfeed 爬虫功能

Zenfeed 提供了将网页内容抓取并转换为 Markdown 格式的功能。这主要通过重写规则 (`rewrites` rule) 中的 `transform.to_text.type` 配置项实现。

## 如何使用

在你的配置文件中，找到 `storage.feed.rewrites` 部分。当你定义一条重写规则时，可以通过 `transform` 字段来启用爬虫功能。

具体配置如下：

```yaml
storage:
  feed:
    rewrites:
      - if: ["source=xxx", ...]
        source_label: "link" # 指定包含 URL 的标签，例如 feed 中的 'link' 标签
        transform:
          to_text:
            type: "crawl" # 或 "crawl_by_jina"
            # llm: "your-llm-name" # crawl 类型不需要 llm
            # prompt: "your-prompt" # crawl 类型不需要 prompt
        # match: ".*" # 可选：对抓取到的 Markdown 内容进行匹配
        action: "create_or_update_label" # 对抓取到的内容执行的动作
        label: "crawled_content" # 将抓取到的 Markdown 存储到这个新标签
    # ... 其他配置 ...
jina: # 如果使用 crawl_by_jina，并且需要更高的速率限制（匿名ip: 20 RPM），请配置 Jina API Token
  token: "YOUR_JINA_AI_TOKEN" # 从 https://jina.ai/api-dashboard/ 获取
```

### 转换类型 (`transform.to_text.type`)

你有以下几种选择：

1.  **`crawl`**:
    *   Zenfeed 将使用内置的本地爬虫尝试抓取 `source_label` 中指定的 URL。
    *   它会尝试遵循目标网站的 `robots.txt` 协议。
    *   适用于静态网页或结构相对简单的网站。

2.  **`crawl_by_jina`**:
    *   Zenfeed 将通过 [Jina AI Reader API](https://jina.ai/reader/) 来抓取和处理 `source_label` 中指定的 URL。
    *   Jina AI 可能能更好地处理动态内容和复杂网站结构。
    *   同样遵循目标网站的 `robots.txt` 协议。
    *   **依赖 Jina AI 服务**：
        *   建议在配置文件的顶层添加 `jina.token` (如上示例) 来提供你的 Jina AI API Token，以获得更高的服务速率限制。
        *   如果未提供 Token，将以匿名用户身份请求，速率限制较低。
        *   请查阅 Jina AI 的服务条款和隐私政策。

### 关键配置说明

*   `source_label`: 此标签的值**必须是一个有效的 URL**。例如，如果你的 RSS Feed 中的 `link` 标签指向的是一篇包含完整文章的网页，你可以将 `source_label` 设置为 `link`。
*   `action`: 通常设置为 `create_or_update_label`，将抓取并转换后的 Markdown 内容存入一个新的标签中（由 `label` 字段指定）。
*   `label`: 指定存储抓取到的 Markdown 内容的新标签名称。

## 使用场景

**全文内容提取**:
很多 RSS 源只提供文章摘要和原文链接。使用爬虫功能可以将原文完整内容抓取下来，转换为 Markdown 格式，方便后续的 AI 处理（如总结、打标签、分类等）或直接阅读。

## 免责声明

**在使用 Zenfeed 的爬虫功能（包括 `crawl` 和 `crawl_by_jina` 类型）前，请仔细阅读并理解以下声明。您的使用行为即表示您已接受本声明的所有条款。**

1.  **用户责任与授权**:
    *   您将对使用爬虫功能的所有行为承担全部责任。
    *   您必须确保拥有访问、抓取和处理所提供 URL 内容的合法权利。
    *   请严格遵守目标网站的 `robots.txt` 协议、服务条款 (ToS)、版权政策以及所有相关的法律法规。
    *   不得使用本功能处理、存储或分发任何非法、侵权、诽谤、淫秽或其他令人反感的内容。

2.  **内容准确性与完整性**:
    *   网页抓取和 Markdown 转换过程的结果可能不准确、不完整或存在偏差。这可能受到目标网站结构、反爬虫机制、动态内容渲染、网络问题等多种因素的影响。
    *   Zenfeed 项目作者和贡献者不对抓取内容的准确性、完整性、及时性或质量作任何保证。

3.  **第三方服务依赖 (`crawl_by_jina`)**:
    *   `crawl_by_jina` 功能依赖于 Jina AI 提供的第三方服务。
    *   Jina AI 服务的可用性、性能、数据处理政策、服务条款以及可能的费用（超出免费额度后）均由 Jina AI 自行决定。
    *   项目作者和贡献者不对 Jina AI 服务的任何方面负责。请在使用前查阅 [Jina AI 的相关条款](https://jina.ai/terms/) 和 [隐私政策](https://jina.ai/privacy/)。

4.  **无间接或后果性损害赔偿**:
    *   在任何情况下，无论基于何种法律理论，项目作者和贡献者均不对因使用或无法使用爬虫功能而导致的任何直接、间接、偶然、特殊、惩戒性或后果性损害负责，包括但不限于利润损失、数据丢失、商誉损失或业务中断。

5.  **法律与合规风险**:
    *   未经授权抓取、复制、存储、处理或传播受版权保护的内容，或违反网站服务条款的行为，可能违反相关法律法规，并可能导致法律纠纷或处罚。
    *   用户需自行承担因使用爬虫功能而产生的所有法律风险和责任。

6.  **"按原样"提供**:
    *   爬虫功能按"现状"和"可用"的基础提供，不附带任何形式的明示或默示担保。

**强烈建议您在启用和配置爬虫功能前，仔细评估相关风险，并确保您的使用行为完全合法合规。对于任何因用户滥用或不当使用本软件（包括爬虫功能）而引起的法律纠纷、损失或损害，Zenfeed 项目作者和贡献者不承担任何责任。**
