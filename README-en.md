zenfeed: Empower RSS with AI, automatically filter, summarize, and push important information for you, say goodbye to information overload, and regain control of reading.

## Preface

RSS (Really Simple Syndication) was born in the Web 1.0 era to solve the problem of information fragmentation, allowing users to aggregate and track updates from multiple websites in one place without frequent visits. It pushes website updates in summary form to subscribers for quick information access.

However, with the rise of Web 2.0, social media, and algorithmic recommendations, RSS didn't become mainstream. The shutdown of Google Reader in 2013 was a landmark event. As Zhang Yiming pointed out at the time, RSS demands a lot from users: strong information filtering skills and self-discipline to manage feeds, otherwise it's easy to get overwhelmed by information noise. He believed that for most users, the easier "personalized recommendation" was a better solution, which later led to Toutiao and TikTok.

Algorithmic recommendations indeed lowered the bar for information acquisition, but their excessive catering to human weaknesses often leads to filter bubbles and addiction to entertainment. If you want to get truly valuable content from the information stream, you actually need stronger self-control to resist the algorithm's "feeding".

So, is pure RSS subscription the answer? Not necessarily. Information overload and filtering difficulties (information noise) remain pain points for RSS users.

Confucius advocated the doctrine of the mean in all things. Can we find a middle ground that combines the sense of control and high-quality sources from active RSS subscription with technological means to overcome its information overload drawbacks?

Try zenfeed! **AI + RSS** might be a better way to acquire information in this era. zenfeed aims to leverage AI capabilities to help you automatically filter and summarize the information you care about, allowing you to maintain Zen (calmness) amidst the Feed (information flood).

## Project Introduction

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/1b51f1087558402d85496fbe7bddde89)](https://app.codacy.com/gh/glidea/zenfeed/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=glidea_zenfeed&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=glidea_zenfeed)
[![Go Report Card](https://goreportcard.com/badge/github.com/glidea/zenfeed)](https://goreportcard.com/report/github.com/glidea/zenfeed)

zenfeed is your intelligent information assistant. It automatically collects, filters, and summarizes news or topics you follow, then sends them to you. But we're not just building another "Toutiao"... 🤔

![Zenfeed](docs/images/arch.png)

**For [RSS](https://en.wikipedia.org/wiki/RSS) Veterans** 🚗
* zenfeed can be your AI-powered RSS reader (works with [zenfeed-web](https://github.com/glidea/zenfeed-web))
* An [MCP](https://mcp.so/) Server for [RSSHub](https://github.com/DIYgod/RSSHub)
* A customizable, trusted RSS data source and an incredibly fast AI search engine
* Similar to [Feedly AI](https://feedly.com/ai)
<details>
  <summary>Preview</summary>
  <img src="docs/images/feed-list-with-web.png" alt="Feed list with web UI" width="600">
  <img src="docs/images/chat-with-feeds.png" alt="Chat with feeds" width="500">
</details>

**For Seekers of [WWZZ](https://www.wwzzai.com/) Alternatives** 🔍
* zenfeed also offers [information tracking capabilities](https://github.com/glidea/zenfeed/blob/main/docs/config.md#schedule-configuration-schedules), emphasizing high-quality, customizable data sources.
* Think of it as an RSS-based, flexible, more PaaS-like version of [AI Chief Information Officer](https://github.com/TeamWiseFlow/wiseflow?tab=readme-ov-file).
<details>
  <summary>Preview</summary>
  <img src="docs/images/monitoring.png" alt="Monitoring preview" width="500">
  <img src="docs/images/notification-with-web.png" alt="Notification with web UI" width="500">
</details>

**For Information Anxiety Sufferers (like me)** 😌
* "zenfeed" combines "zen" and "feed," signifying maintaining calm (zen) amidst the information flood (feed).
* If you feel anxious and tired from constantly checking information streams, it's because context switching costs more than you think and hinders entering a flow state. Try the briefing feature: receive a summary email at a fixed time each day covering the relevant period. This allows for a one-time, quick, comprehensive overview. Ah, a bit of a renaissance feel, isn't it? ✨
<details>
  <summary>Preview</summary>
  <img src="docs/images/daily-brief.png" alt="Daily brief preview" width="500">
</details>

**For Explorers of AI Content Processing** 🔬
* zenfeed features a custom mechanism for pipelining content processing, similar to Prometheus [Relabeling](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config).
* Each piece of content is abstracted as a set of labels (e.g., title, source, body... are labels). At each node in the pipeline, you can process specific label values based on custom prompts (e.g., scoring, classifying, summarizing, filtering, adding new labels...). Subsequently, you can filter based on label queries, [route](https://github.com/glidea/zenfeed/blob/main/docs/config.md#notification-route-configuration-notifyroute-and-notifyroutesub_routes), and [display](https://github.com/glidea/zenfeed/blob/main/docs/config.md#notification-channel-email-configuration-notifychannelsemail)... See [Rewrite Rules](https://github.com/glidea/zenfeed/blob/main/docs/config.md#rewrite-rule-configuration-storagefeedrewrites).
* Crucially, you can flexibly orchestrate all this, giving zenfeed a strong tooling and personalization flavor. Welcome to integrate private data via the Push API and explore more possibilities.
<details>
  <summary>Preview</summary>
  <img src="docs/images/update-config-with-web.png" alt="Update config with web UI" width="500">
</details>

**For Onlookers** 🍉

Just for the exquisite email styles, install and use it now!

<img src="docs/images/monitoring.png" alt="Monitoring email style" width="400">

[More Previews](docs/preview.md)

## Installation and Usage

### 1. Installation

By default, uses SiliconFlow's Qwen/Qwen2.5-7B-Instruct (free) and Pro/BAAI/bge-m3. If you don't have a SiliconFlow account yet, use this [invitation link](https://cloud.siliconflow.cn/i/U2VS0Q5A) to get a ¥14 credit.

Support for other vendors or models is available; follow the instructions below.

#### Mac/Linux

```bash
curl -L -O https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml

# If you need to customize more configuration parameters, directly edit docker-compose.yml#configs.zenfeed_config.content BEFORE running the command below.
# Configuration Docs: https://github.com/glidea/zenfeed/blob/main/docs/config.md
API_KEY=your_apikey TZ=your_local_IANA LANG=English  docker-compose -p zenfeed up -d
```

#### Windows
> Use PowerShell to execute
```powershell
Invoke-WebRequest -Uri "https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml" -OutFile ([System.IO.Path]::GetFileName("https://raw.githubusercontent.com/glidea/zenfeed/main/docker-compose.yml"))

# If you need to customize more configuration parameters, directly edit docker-compose.yml#configs.zenfeed_config.content BEFORE running the command below.
# Configuration Docs: https://github.com/glidea/zenfeed/blob/main/docs/config.md
$env:API_KEY = "your_apikey"; $env:TZ = "your_local_IANA"; $env:LANG = "English"; docker-compose -p zenfeed up -d
```

### 2. Using the Web UI

Access https://zenfeed-web.pages.dev
> If deployed in an environment like a VPS, access https://vps_public_ip:1400 (remember to open the security group port). Do not use the public frontend above.
> ⚠️ zenfeed currently lacks authentication. Exposing it to the public internet might leak your API Key. Please configure your security groups carefully. If you have security concerns, please open an Issue.

#### Add RSS Feeds

<img src="docs/images/web-add-source.png" alt="Add source via web UI" width="400">

> To migrate from Follow, refer to [migrate-from-follow.md](docs/migrate-from-follow.md)
> Requires access to the respective source sites; ensure network connectivity.
> Wait a few minutes after adding, especially if the model has strict rate limits.

#### Configure Daily Briefings, Monitoring, etc.

<img src="docs/images/notification-with-web.png" alt="Configure notifications via web UI" width="400">

### 3. Configure MCP (Optional)
Using Cherry Studio as an example, configure MCP and connect to Zenfeed, see [Cherry Studio MCP](docs/cherry-studio-mcp.md)
> Default address http://localhost:1301/sse

## Roadmap
* P0 (Very Likely)
  * Support generating podcasts, male/female dialogues, similar to NotebookLM
  * More data sources
    * Email
    * Web clipping Chrome extension
* P1 (Possible)
  * Keyword search
  * Support search engines as data sources
  * App?
  * The following are temporarily not prioritized due to copyright risks:
    * Webhook notifications
    * Web scraping

## Notice
* Compatibility is not guaranteed before version 1.0.
* The project uses the AGPLv3 license; any forks must also be open source.
* For commercial use, please contact for registration; reasonable support can be provided. Note: Legal commercial use only, gray area activities are not welcome.
* Data is not stored permanently; default retention is 8 days.

## Acknowledgments
* Thanks to [eryajf](https://github.com/eryajf) for providing the [Compose Inline Config](https://github.com/glidea/zenfeed/issues/1) idea, making deployment easier to understand.

## 👏🏻 Contributions Welcome
* No formal guidelines yet, just one requirement: "Code Consistency" – it's very important.

## Disclaimer

**Before using the `zenfeed` software (hereinafter referred to as "the Software"), please read and understand this disclaimer carefully. Your download, installation, or use of the Software or any related services signifies that you have read, understood, and agreed to be bound by all terms of this disclaimer. If you do not agree with any part of this disclaimer, please cease using the Software immediately.**

1.  **Provided "AS IS":** The Software is provided on an "AS IS" and "AS AVAILABLE" basis, without any warranties of any kind, either express or implied. The project authors and contributors make no warranties or representations regarding the Software's merchantability, fitness for a particular purpose, non-infringement, accuracy, completeness, reliability, security, timeliness, or performance.

2.  **User Responsibility:** You are solely responsible for all actions taken using the Software. This includes, but is not limited to:
    *   **Data Source Selection:** You are responsible for selecting and configuring the data sources (e.g., RSS feeds, potential future Email sources) you connect to the Software. You must ensure you have the right to access and process the content from these sources and comply with their respective terms of service, copyright policies, and applicable laws and regulations.
    *   **Content Compliance:** You must not use the Software to process, store, or distribute any content that is unlawful, infringing, defamatory, obscene, or otherwise objectionable.
    *   **API Key and Credential Security:** You are responsible for safeguarding the security of any API keys, passwords, or other credentials you configure within the Software. The authors and contributors are not liable for any loss or damage arising from your failure to maintain proper security.
    *   **Configuration and Use:** You are responsible for correctly configuring and using the Software's features, including content processing pipelines, filtering rules, notification settings, etc.

3.  **Third-Party Content and Services:** The Software may integrate with or rely on third-party data sources and services (e.g., RSSHub, LLM providers, SMTP service providers). The project authors and contributors are not responsible for the availability, accuracy, legality, security, or terms of service of such third-party content or services. Your interactions with these third parties are governed by their respective terms and policies. Copyright for third-party content accessed or processed via the Software (including original articles, summaries, classifications, scores, etc.) belongs to the original rights holders, and you assume all legal liability arising from your use of such content.

4.  **No Warranty on Content Processing:** The Software utilizes technologies like Large Language Models (LLMs) to process content (e.g., summarization, classification, scoring, filtering). These processed results may be inaccurate, incomplete, or biased. The project authors and contributors are not responsible for any decisions made or actions taken based on these processed results. The accuracy of semantic search results is also affected by various factors and is not guaranteed.

5.  **No Liability for Indirect or Consequential Damages:** In no event shall the project authors or contributors be liable under any legal theory (whether contract, tort, or otherwise) for any direct, indirect, incidental, special, exemplary, or consequential damages arising out of the use or inability to use the Software. This includes, but is not limited to, loss of profits, loss of data, loss of goodwill, business interruption, or other commercial damages or losses, even if advised of the possibility of such damages.

6.  **Open Source Software:** The Software is licensed under the AGPLv3 License. You are responsible for understanding and complying with the terms of this license.

7.  **Not Legal Advice:** This disclaimer does not constitute legal advice. If you have any questions regarding the legal implications of using the Software, you should consult a qualified legal professional.

8.  **Modification and Acceptance:** The project authors reserve the right to modify this disclaimer at any time. Continued use of the Software following any modifications will be deemed acceptance of the revised terms.

**Please be aware: Using the Software to fetch, process, and distribute copyrighted content may carry legal risks. Users are responsible for ensuring their usage complies with all applicable laws, regulations, and third-party terms of service. The project authors and contributors assume no liability for any legal disputes or losses arising from user misuse or improper use of the Software.**
