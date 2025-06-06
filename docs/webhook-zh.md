# Zenfeed Webhook 通知对接指南

Zenfeed 支持通过 Webhook 将分组和总结后的 Feed 通知推送到您指定的 HTTP(S) 端点。这允许您将 Zenfeed 的通知集成到自定义的应用或工作流程中。

## 1. 配置方法

要在 Zenfeed 中配置 Webhook 通知，您需要在配置文件的 `notify.receivers` 部分定义一个或多个接收者，并为每个 Webhook 接收者指定其唯一的 `name` 和 `webhook` 配置块。

**示例配置 (`config.yaml`):**

```yaml
notify:
  # ... 其他通知配置 ...

  receivers:
    - name: my_awesome_webhook # 接收者的唯一名称，将在路由规则中引用
      webhook:
        url: "https://your-service.com/webhook-endpoint" # 您的 Webhook 接收端点 URL

  # 示例：路由规则中如何使用此接收者
  route: # or sub_routes..
    receivers:
      - my_awesome_webhook # 引用上面定义的接收者名称
    # ... 其他路由配置 ...
```

在上述示例中：
- 我们定义了一个名为 `my_awesome_webhook` 的接收者。
- `webhook.url` 字段指定了当有匹配此接收者的通知时，Zenfeed 将向哪个 URL 发送 POST 请求。

## 2. 数据格式详解

当 Zenfeed 向您的 Webhook 端点发送通知时，它会发送一个 `POST` 请求，请求体为 JSON 格式。

请求体结构如下：

```json
{
  "group": "string",
  "labels": {
    "label_key1": "label_value1",
    "label_key2": "label_value2"
  },
  "summary": "string",
  "feeds": [
    {
      "labels": {
        "title": "Feed Title 1",
        "link": "http://example.com/feed1",
        "content": "Feed content snippet 1...",
        "source": "example_source",
        "pub_time": "2024-07-30T10:00:00Z"
        // ... 其他自定义或标准标签
      },
      "time": "2024-07-30T10:00:00Z",
      "related": [
        // 可选：与此 Feed 相关的其他 Feed 对象，结构同父 Feed
      ]
    }
    // ...更多 Feed 对象
  ]
}
```

**字段说明:**

-   `group` (`string`):
    当前通知所属的组名。这个名称是根据通知路由配置中 `group_by` 定义的标签值组合而成的。例如，如果 `group_by: ["source", "category"]`，且一个 Feed 组的 `source` 是 `github_trending`，`category` 是 `golang`，那么 `group` 可能类似于 `"github_trending/golang"`。

-   `labels` (`object`):
    一个键值对对象，表示当前通知组的标签。这些标签是根据通知路由配置中 `group_by` 所指定的标签及其对应的值。
    例如，如果 `group_by: ["source"]` 且当前组的 `source` 标签值为 `rsshub`，则 `labels` 会是 `{"source": "rsshub"}`。

-   `summary` (`string`):
    由大语言模型 (LLM) 为当前这一组 Feed 生成的摘要文本。如果通知路由中没有配置 LLM 总结，此字段可能为空字符串或省略 (取决于具体的实现细节，但通常会是空字符串)。

-   `feeds` (`array` of `object`):
    一个数组，包含了属于当前通知组的所有 Feed 对象。每个 Feed 对象包含以下字段：
    *   `labels` (`object`): Feed 的元数据。这是一个键值对对象，包含了该 Feed 的所有标签，例如：
        *   `title` (`string`): Feed 的标题。
        *   `link` (`string`): Feed 的原始链接。
        *   `content` (`string`): Feed 的内容摘要或全文 (取决于抓取和重写规则)。
        *   `source` (`string`): Feed 的来源标识。
        *   `pub_time` (`string`): Feed 的发布时间 (RFC3339 格式的字符串，例如 `2025-01-01T00:00:00Z`)。
        *   ...以及其他在抓取或重写过程中添加的自定义标签。
    *   `time` (`string`): Feed 的时间戳，通常是其发布时间，采用 RFC3339 格式 (例如 `2025-01-01T00:00:00Z`)。此字段与 `labels.pub_time` 通常一致，但 `time` 是系统内部用于时间序列处理的主要时间字段。
    *   `related` (`array` of `object`, 可选):
        一个数组，包含了与当前 Feed 语义相关的其他 Feed 对象。这通常在通知路由中启用了 `compress_by_related_threshold` 选项时填充。每个相关的 Feed 对象结构与父 Feed 对象完全相同。如果未启用相关性压缩或没有相关的 Feed，此字段可能为空数组或不存在。

## 3. 请求示例

以下是一个发送到您的 Webhook 端点的 JSON 请求体示例：

```json
{
  "group": "my_favorite_blogs",
  "labels": {
    "category": "tech_updates",
  },
  "summary": "今天有多篇关于最新 AI 技术进展的文章，重点关注了大型语言模型在代码生成方面的应用，以及其对未来软件开发模式的潜在影响。",
  "feeds": [
    {
      "labels": {
        "content": "AlphaCode X 展示了惊人的代码理解和生成能力，在多个编程竞赛中超越了人类平均水平...",
        "link": "https://example.blog/alphacode-x-details",
        "pub_time": "2024-07-30T14:35:10Z",
        "source": "Example Tech Blog",
        "title": "AlphaCode X: 下一代 AI 编码助手",
        "type": "blog_post"
      },
      "time": "2024-07-30T14:35:10Z",
      "related": []
    },
    {
      "labels": {
        "content": "讨论了当前 LLM 在实际软件工程项目中落地所面临的挑战，包括成本、可控性和安全性问题。",
        "link": "https://another.blog/llm-in-swe-challenges",
        "pub_time": "2024-07-30T11:15:00Z",
        "source": "Another Tech Review",
        "title": "LLM 在软件工程中的应用：机遇与挑战",
        "type": "rss"
      },
      "time": "2024-07-30T11:15:00Z",
      "related": [
        {
          "labels": {
            "content": "一篇关于如何更经济有效地部署和微调大型语言模型的指南。",
            "link": "https://some.other.blog/cost-effective-llm",
            "pub_time": "2024-07-30T09:00:00Z",
            "source": "AI Infra Weekly",
            "title": "经济高效的 LLM 部署策略",
            "type": "rss"
          },
          "time": "2024-07-30T09:00:00Z",
          "related": []
        }
      ]
    }
  ]
}
```

## 4. 响应要求

Zenfeed 期望您的 Webhook 端点在成功接收并处理通知后，返回 HTTP `200 OK` 状态码。
如果 Zenfeed 收到任何非 `200` 的状态码，它会将该次通知尝试标记为失败，并可能根据重试策略进行重试 (具体重试行为取决于 Zenfeed 的内部实现)。

请确保您的端点能够及时响应，以避免超时。
