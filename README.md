# ⚡ AI Cost-Control & Caching Gateway

**Slash your OpenAI & LLM API bills by 30%–50% in 60 seconds with a single zero-dependency Go binary.**

Modern AI applications waste massive amounts of money sending identical prompts and system instructions to expensive LLM endpoints repeatedly. This ultra-fast reverse proxy sits seamlessly between your application and upstream AI providers, caching identical requests in memory and tracking exact dollar savings in real time.

---

## 🚀 Quick Start (Open Core Edition - Apache 2.0)

### 1. Start the Gateway
```bash
go run main.go
```
*The proxy will immediately start listening on `:8888`.*

### 2. Point Your App to the Proxy
Simply change your OpenAI base URL from `https://openai.com` to `http://localhost:8888`. Zero code changes required.

```bash
curl -X POST http://localhost:8888/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_OPENAI_KEY" \
  -d '{"model": "gpt-4o", "messages": [{"role": "user", "content": "Summarize our Q3 policy."}]}'
```

### 3. Check Real-Time ROI & Savings
Query the built-in telemetry endpoint at any time to see exact dollar savings:
```bash
curl http://localhost:8888/metrics
```

---

## 📊 Community vs. Enterprise Edition

| Feature | Open Core (Apache 2.0) | Enterprise Edition |
| :--- | :---: | :---: |
| **Zero-Dependency Go Proxy** | ✅ | ✅ |
| **In-Memory SHA-256 Prompt Caching** | ✅ (1-Hour TTL) | ✅ (Custom / 24-Hour+ TTL) |
| **Real-Time Savings `/metrics` API** | ✅ | ✅ |
| **Distributed Redis Cluster Sync** | ❌ | ✅ |
| **Automated Semantic Caching (Fuzzy Match)** | ❌ | ✅ |
| **Smart Model Fallback Routing (Cost Optimization)** | ❌ | ✅ |
| **SOC2 Compliance & PII Redaction Logs** | ❌ | ✅ |
| **Dedicated 24/7 SLA & Architecture Support** | ❌ | ✅ |

---

## 💎 Unlock Enterprise Edition

Running AI in production at scale? Upgrade to the **Enterprise Edition** to unlock 24-hour+ persistent caching, multi-node cluster synchronization, automated PII scrubbing, and smart model routing.

**To request an Enterprise License Key (`ENT-...`) and schedule an architecture review, contact:**
* **Email:** `brentangel@gmail.com` *(Replace with your actual email before publishing)*

---

## 📄 License
The core engine in this repository is licensed under the **Apache License 2.0**. See the `LICENSE` file for full details.
