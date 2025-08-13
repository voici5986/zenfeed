# 使用 Zenfeed 将文章转换为播客

Zenfeed 的播客功能可以将任何文章源自动转换为一场引人入胜的多人对话式播客。该功能利用大语言模型（LLM）生成对话脚本和文本转语音（TTS），并将最终的音频文件托管在您自己的对象存储中。

## 工作原理

1.  **提取内容**: Zenfeed 首先通过重写规则提取文章的全文内容。
2.  **生成脚本**: 使用一个指定的 LLM（如 GPT-4o-mini）将文章内容改编成一个由多位虚拟主播对话的脚本。您可以定义每个主播的角色（人设）来控制对话风格。
3.  **语音合成**: 调用另一个支持 TTS 的 LLM（目前仅支持 Google Gemini）将脚本中的每一句对话转换为语音。
4.  **音频合并**: 将所有语音片段合成为一个完整的 WAV 音频文件。
5.  **上传存储**: 将生成的播客文件上传到您配置的 S3 兼容对象存储中。
6.  **保存链接**: 最后，将播客文件的公开访问 URL 保存为一个新的 Feed 标签，方便您在通知、API 或其他地方使用。

## 配置步骤

要启用播客功能，您需要完成以下三项配置：LLM、对象存储和重写规则。

### 1. 配置 LLM

您需要至少配置两个 LLM：一个用于生成对话脚本，另一个用于文本转语音（TTS）。

-   **脚本生成 LLM**: 可以是任何性能较好的聊天模型，例如 OpenAI 的 `gpt-4o-mini` 或 Google 的 `gemini-1.5-pro`。
-   **TTS LLM**: 用于将文本转换为语音。**注意：目前此功能仅支持 `provider` 为 `gemini` 的 LLM。**

**示例 `config.yaml`:**

```yaml
llms:
  # 用于生成播客脚本的 LLM
  - name: openai-chat
    provider: openai
    api_key: "sk-..."
    model: gpt-4o-mini
    default: true

  # 用于文本转语音 (TTS) 的 LLM
  - name: gemini-tts
    provider: gemini
    api_key: "..." # 你的 Google AI Studio API Key
    tts_model: "gemini-2.5-flash-preview-tts" # Gemini 的 TTS 模型
```

### 2. 配置对象存储

生成的播客音频文件需要一个地方存放。Zenfeed 支持任何 S3 兼容的对象存储服务。这里我们以 [Cloudflare R2](https://www.cloudflare.com/zh-cn/products/r2/) 为例。

首先，您需要在 Cloudflare R2 中创建一个存储桶（Bucket）。然后获取以下信息：

-   `endpoint`: 您的 R2 API 端点。通常格式为 `<account_id>.r2.cloudflarestorage.com`。您可以在 R2 存储桶的主页找到它。
-   `access_key_id` 和 `secret_access_key`: R2 API 令牌。您可以在 "R2" -> "管理 R2 API 令牌" 页面创建。
-   `bucket`: 您创建的存储桶的名称。
-   `bucket_url`: 存储桶的公开访问 URL。要获取此 URL，您需要将存储桶连接到一个自定义域，或者使用 R2 提供的 `r2.dev` 公开访问地址。

**示例 `config.yaml`:**

```yaml
storage:
  object:
    endpoint: "<your_account_id>.r2.cloudflarestorage.com"
    access_key_id: "..."
    secret_access_key: "..."
    bucket: "zenfeed-podcasts"
    bucket_url: "https://pub-xxxxxxxx.r2.dev"
```

### 3. 配置重写规则

最后一步是创建一个重写规则，告诉 Zenfeed 如何将文章转换为播客。这个规则定义了使用哪个标签作为源文本、由谁来对话、使用什么声音等。

**关键配置项:**

-   `source_label`: 包含文章全文的标签。
-   `label`: 用于存储最终播客 URL 的新标签名称。
-   `transform.to_podcast`: 播客转换的核心配置。
    -   `llm`: 用于生成脚本的 LLM 名称（来自 `llms` 配置）。
    -   `tts_llm`: 用于 TTS 的 LLM 名称（来自 `llms` 配置）。
    -   `speakers`: 定义播客的演讲者。
        -   `name`: 演讲者的名字。
        -   `role`: 演讲者的角色和人设，将影响脚本内容。
        -   `voice`: 演讲者的声音。请参考 [Gemini TTS 文档](https://ai.google.dev/gemini-api/docs/speech-generation#voices)。

**示例 `config.yaml`:**

```yaml
storage:
  feed:
    rewrites:
      - source_label: content # 基于原文
        transform:
          to_podcast:
            estimate_maximum_duration: 3m0s # 接近 3 分钟
            transcript_additional_prompt: 对话引人入胜，流畅自然，拒绝 AI 味，使用中文回复 # 脚本内容要求
            llm: xxxx # 负责生成脚本的 llm
            tts_llm: gemini-tts # 仅支持 gemini tts，推荐使用 https://github.com/glidea/one-balance 轮询
            speakers:
              - name: 小雅
                role: >-
                  一位经验丰富、声音甜美、风格活泼的科技播客主持人。前财经记者、媒体人出身，因为工作原因长期关注科技行业，后来凭着热爱和出色的口才转行做了全职内容创作者。擅长从普通用户视角出发，把复杂的技术概念讲得生动有趣，是她发掘了老王，并把他‘骗’来一起做播客的‘始作俑者’。
                voice: Autonoe
              - name: 老王
                role: >-
                  一位资深科技评论员，互联网老兵。亲身经历过中国互联网从草莽到巨头的全过程，当过程序员，做过产品经理，也创过业。因此他对行业的各种‘风口’和‘概念’有自己独到的、甚至有些刻薄的见解。观点犀利，一针见血，说话直接，热衷于给身边的一切产品挑刺。被‘忽悠’上了‘贼船’，表面上经常吐槽，但内心很享受这种分享观点的感觉。
                voice: Puck
        label: podcast_url
```

配置完成后，Zenfeed 将在每次抓取到新文章时，自动执行上述流程。可以在通知模版中使用 podcast_url label，或 Web 中直接收听（Web 固定读取 podcast_url label，若使用别的名称则无法读取）