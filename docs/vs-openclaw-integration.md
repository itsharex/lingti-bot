# lingti-bot vs OpenClaw：简化 AI 集成的努力

> 本文对比 lingti-bot 和 OpenClaw 在简化 AI Bot 接入方面的设计理念和实现差异。

## 核心理念对比

| | OpenClaw | lingti-bot |
|---|---|---|
| **设计哲学** | 灵活优先，支持多种部署方式 | 简单优先，开箱即用 |
| **目标用户** | 开发者、技术爱好者 | 所有人（包括非技术用户） |
| **接入方式** | 主要依赖自建服务器 | 云中继 + 自建服务器双模式 |
| **入门门槛** | 需要服务器运维知识 | 3 条命令即可完成 |

## 企业微信接入对比

### OpenClaw 方式

OpenClaw 不原生支持企业微信，如需接入需要：

1. 准备公网服务器（或使用云服务器）
2. 配置域名和 DNS 解析
3. 申请 SSL 证书（企业微信可选，但推荐）
4. 部署 OpenClaw 或编写自定义回调服务
5. 配置防火墙规则
6. 编写消息转发逻辑到 OpenClaw

**预计耗时：数小时到数天**

### lingti-bot 方式

lingti-bot 提供两种接入方式：

#### 方式一：云中继模式（推荐，5分钟完成）

```bash
# 1. 安装
curl -fsSL https://cli.lingti.com/install.sh | bash -s -- --bot

# 2. 验证回调
lingti-bot verify --platform wecom --wecom-corp-id ... --wecom-token ...
# 去企业微信后台配置 URL: https://bot.lingti.com/wecom

# 3. 开始处理消息
lingti-bot relay --platform wecom --provider deepseek --api-key sk-xxx ...
```

**无需：**
- 公网服务器
- 域名/DNS
- SSL 证书
- 防火墙配置
- 任何运维知识

#### 方式二：自建服务器模式

```bash
lingti-bot router --wecom-corp-id ... --provider deepseek --api-key ...
```

与传统方式类似，但：
- 内置完整的企业微信回调处理
- 内置消息加解密
- 内置 Access Token 管理
- 无需编写任何代码

## 技术实现对比

### 消息流转

**传统方式（OpenClaw 等）：**
```
企业微信 → 公网服务器 → 消息处理 → AI API → 响应
            ↑
        需要用户准备
```

**lingti-bot 云中继：**
```
企业微信 → bot.lingti.com → WebSocket → 本地客户端 → AI API
                                            ↑
                                     用户本地运行
```

### 凭据安全

| | 传统方式 | lingti-bot 云中继 |
|---|---|---|
| AI API Key 存放 | 服务器 | 本地 |
| 企业微信凭据 | 服务器 | 动态传输，不持久化 |
| 消息内容 | 服务器处理 | 本地处理，云端仅转发 |

## 支持平台对比

| 平台 | OpenClaw | lingti-bot | 说明 |
|------|----------|------------|------|
| Slack | ✅ | ✅ | |
| Discord | ✅ | ✅ | |
| Telegram | ✅ | ✅ | |
| WhatsApp | ✅ | ❌ | lingti-bot 计划支持 |
| iMessage | ✅ | ❌ | |
| 飞书/Lark | ❌ | ✅ | **lingti-bot 独有** |
| 企业微信 | ❌ | ✅ | **lingti-bot 独有** |
| 微信公众号 | ❌ | ✅ | **lingti-bot 独有**（云中继）|
| 钉钉 | ❌ | 🚧 | 开发中 |

## 云中继技术详解

### 工作原理

1. **用户运行 `verify` 命令**
   - 通过 WebSocket 连接到 `wss://bot.lingti.com/ws`
   - 发送企业微信凭据（Token、AESKey、CorpID 等）

2. **云端接收验证请求**
   - 企业微信发送 GET 请求到 `https://bot.lingti.com/wecom`
   - 云端使用用户提供的凭据计算签名、解密 echostr
   - 返回明文完成验证

3. **用户运行 `relay` 命令**
   - 再次通过 WebSocket 连接
   - 发送凭据用于消息解密

4. **消息处理流程**
   - 企业微信发送加密消息到 `https://bot.lingti.com/wecom`
   - 云端使用对应的凭据解密
   - 通过 WebSocket 转发到用户本地客户端
   - 本地客户端调用 AI API 处理
   - 响应通过 Webhook 返回云端
   - 云端调用企业微信 API 发送消息

### 安全设计

1. **凭据动态传输**
   - 凭据通过加密的 WSS 连接传输
   - 仅在内存中保存，不持久化到数据库
   - 客户端断开后凭据自动清除

2. **AI 处理本地化**
   - AI API Key 始终在用户本地
   - 消息内容在本地处理
   - 云端仅负责消息中转

3. **单客户端限制**
   - 同一 user_id 只能有一个活跃连接
   - 防止凭据被滥用

## 适用场景

### 推荐使用 lingti-bot 云中继

- 个人用户快速体验 AI Bot
- 小团队内部使用
- 没有运维能力的用户
- 需要快速原型验证

### 推荐使用自建服务器

- 企业生产环境
- 对数据安全有严格要求
- 需要完全控制消息流转
- 高并发场景

### OpenClaw 适用场景

- 需要 WhatsApp/iMessage 等平台
- 已有成熟的服务器基础设施
- 需要深度定制消息处理逻辑
- 海外用户为主

## 总结

lingti-bot 在以下方面做出了显著改进：

1. **零门槛接入**：云中继模式让任何人都能在 5 分钟内接入企业微信等平台
2. **中国平台优先**：原生支持飞书、企业微信、微信公众号、钉钉等国内主流平台
3. **本地化优先**：AI 处理在本地完成，数据不上云
4. **单一二进制**：无需 Docker、数据库或其他依赖

这些努力的目标是：**让 AI Bot 接入像配置 Wi-Fi 一样简单**。

## 相关文档

- [企业微信集成指南](wecom-integration.md)
- [飞书集成指南](feishu-integration.md)
- [微信公众号接入指南](wechat-integration.md)
- [开发路线图](roadmap.md)
