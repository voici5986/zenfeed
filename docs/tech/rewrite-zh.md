> 适用版本：v0.2.2

`rewrite` 组件是 zenfeed 中负责对信息流内容进行动态处理和转换的核心模块。它允许用户通过声明式的规则配置，利用大型语言模型 (LLM) 等工具，对内容的元数据（标签）进行修改、丰富、过滤，甚至决定是否丢弃某条信息。

## 1. 设计理念与哲学

*   **Prometheus 的 `relabel_config`**: 借鉴其强大的标签重写能力。在 Prometheus 中，`relabel_config` 允许用户在采集指标前后动态地修改标签集，实现服务发现、指标过滤和路由等高级功能。`rewrite` 组件将此思想应用于信息流处理，将每一条信息（如一篇文章、一个帖子）视为一个标签集，通过规则来操作这些标签。
*   **管道 (Pipeline) 处理模式**: 信息的处理过程被设计成一个可配置的管道。每个规则是管道中的一个处理阶段，信息流经这些规则，逐步被转换和打标。这种模式使得复杂的处理逻辑可以被分解为一系列简单、独立的步骤，易于理解和维护。
*   **AI 能力的模块化与按需应用**: 大型语言模型 (LLM) 被视为一种强大的"转换函数"。用户可以根据需求，在规则中指定使用哪个 LLM、配合什么样的提示词 (Prompt) 来处理特定的文本内容（例如，从文章正文生成摘要、分类、评分等）。这种设计使得 AI 能力可以灵活地嵌入到信息处理的任意环节。
*   **内容即标签 (Content as Labels)**: 这是 zenfeed 的一个核心抽象。原始信息（如标题、正文、链接、发布时间）和经过 AI 或规则处理后产生的衍生信息（如类别、标签、评分、摘要）都被统一表示为键值对形式的"标签"。这种统一表示简化了后续的查询、过滤、路由和展示逻辑。
*   **声明式配置优于命令式代码**: 用户通过 YAML 配置文件定义重写规则，而不是编写代码来实现处理逻辑。这降低了使用门槛，使得非程序员也能方便地定制自己的信息处理流程，同时也使得配置更易于管理和版本控制。

## 2. 业务流程

内容重写组件的核心工作流程是接收一个代表信息单元的标签集 (`model.Labels`)，然后按顺序应用预定义的重写规则 (`Rule`)，最终输出一个经过修改的标签集，或者指示该信息单元应被丢弃。

其处理流程可以概括为：

1.  **接收标签集**: 组件的入口是一个 `model.Labels` 对象，代表待处理的信息单元。
2.  **顺序应用规则**: 系统会遍历用户配置的每一条 `Rule`。
3.  **规则评估与执行**: 对于每一条规则，系统会：
    *   **定位源文本**: 根据规则指定的 `SourceLabel` (默认为 `content`)，找到相应的文本内容。
    *   **条件检查**: 检查源文本是否满足规则中声明的 `SkipTooShortThreshold`（最小长度，默认为300字符）。若不满足，则跳过当前规则。
    *   **文本转换 (可选)**: 若规则声明了 `Transform` (如通过 `ToText` 使用 LLM 和特定 `Prompt` 进行处理)，则源文本会被转换为新文本。此转换结果将用于后续的匹配。
    *   **模式匹配**: 使用规则中声明的 `Match` 正则表达式 (默认为 `.*`) 来匹配（可能已被转换过的）文本。若不匹配，则跳过当前规则。
    *   **执行动作**: 若文本匹配成功，则执行规则声明的 `Action`：
        *   `ActionDropFeed`: 指示应丢弃当前信息单元，处理流程终止。
        *   `ActionCreateOrUpdateLabel`: 使用（可能已被转换过的）匹配文本，为规则中指定的 `Label` 创建或更新标签值。
4.  **输出结果**:
    *   若所有规则处理完毕且未触发 `ActionDropFeed`，则返回最终修改并排序后的 `model.Labels`。
    *   若任一规则触发 `ActionDropFeed`，则返回 `nil`，表示丢弃。
    *   处理过程中若发生错误（如 LLM 调用失败），则会中止并返回错误。


## 3. 使用示例

以下是一些如何使用 `rewrite` 规则的场景示例：

### 示例 1: 内容分类打标

*   **目标**: 根据文章内容，自动为其添加一个 `category` 标签，如 "Technology", "Finance" 等。
*   **规则配置 (概念性)**:
    ```yaml
    - sourceLabel: "content" # 使用文章正文作为分析源
      transform:
        toText:
          llm: "qwen-default" # 使用名为 "qwen-default" 的 LLM 配置
          prompt: "category"  # 使用预设的 "category" prompt 模板
      match: ".+"             # 匹配 LLM 返回的任何非空分类结果
      action: "create_or_update_label"
      label: "category"       # 新标签的键为 "category"
    ```
*   **效果**: 如果一篇文章内容是关于人工智能的，LLM 可能会返回 "Technology"。经过此规则处理后，文章的标签集会增加或更新一个标签，例如 `{"category", "Technology"}`。

### 示例 2: 基于 LLM 评分过滤低质量内容

*   **目标**: 让 LLM 对文章内容进行评分 (0-10)，如果评分低于 4，则丢弃该文章。
*   **规则配置 (包含两条规则)**:

    *   **规则 2.1: 内容评分**
        ```yaml
        - source_label: "content"
          transform:
            to_text:
              llm: "qwen-default"
              prompt: "score" # 使用预设的 "score" prompt 模板
          match: "^([0-9]|10)$" # 确保 LLM 返回的是 0-10 的数字
          action: "create_or_update_label"
          label: "ai_score"  # 将评分结果存入 "ai_score" 标签
        ```
    *   **规则 2.2: 根据评分过滤**
        ```yaml
        - source_label: "ai_score" # 使用上一条规则生成的评分作为判断依据
          # 无需 Transform
          match: "^[0-3]$"       # 匹配 0, 1, 2, 3 分
          action: "drop_feed"     # 丢弃这些低分文章
        ```
*   **效果**: 文章首先会被 LLM 评分并打上 `ai_score` 标签。如果该评分值在 0 到 3 之间，第二条规则会将其丢弃。

### 示例 3: 基于特定标签值添加新标签

*   **目标**: 如果文章的 `source` 标签值是 "Hacker News"，则添加一个新标签 `source_type: "community"`。
    *   **注意**: 当前 `ActionCreateOrUpdateLabel` 会将匹配成功的 `text` （即 `sourceLabel` 的值或其转换结果）作为新标签的值。若要实现固定值标签，需要通过 LLM 转换。
*   **规则配置 (通过 LLM 实现映射)**:
    ```yaml
    - source_label: "source" # 源标签是 "source"
      transform:
        to_text:
          llm: "qwen-mini"
          # Prompt 需要精心设计，告诉 LLM 如何根据输入映射到输出
          # 例如，Prompt 可以包含类似 "If input is 'Hacker News', output 'community'. If input is 'GitHub Trending', output 'code'." 的逻辑
          prompt: |
            Analyze the input, which is a news source name.
            If the source is "Hacker News", output "community".
            If the source is "GitHub Trending", output "code".
            If the source is "V2EX", output "community".
            Otherwise, output "unknown".
            Return ONLY the type, no other text.
      match: "^(community|code|unknown)$" # 确保 LLM 输出的是预期的类型
      action: "create_or_update_label"
      label: "source_type" # 新标签的键
    ```
*   **效果**: 如果某文章的 `source` 标签值为 "Hacker News"，经过 LLM 处理后（理想情况下）会输出 "community"，然后 `source_type` 标签会被设置为 `{"source_type", "community"}`。

这些示例展示了 `rewrite` 组件的灵活性和强大功能，通过组合不同的源标签、转换、匹配条件和动作，可以实现复杂的内容处理和信息增强逻辑。
