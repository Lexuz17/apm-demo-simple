# apm-demo

Companion repo for the article **"Your App Works — But Is It Observable?"**

→ [Read the article on Medium](https://medium.com/@YOUR_USERNAME/your-article-slug)

---

## Stack

| Service | Port |
|---|---|
| demo-api (Go) | 8080 |
| trigger-web | 80 |
| prometheus | 9090 |
| loki | 3100 |
| alloy | 12345 |
| grafana | 3000 |

## Quick Start

```bash
git clone https://github.com/YOUR_USERNAME/apm-demo
cd apm-demo

cd demo-api && go mod tidy && cd ..
docker compose up -d --build
```

Open:
- Playground: http://localhost
- Grafana: http://localhost:3000 (admin / admin)

## Endpoints

| Endpoint | Behavior |
|---|---|
| `/normal` | 200 OK ~10ms |
| `/slow` | 200 OK, 2–5s delay |
| `/error` | 500 Error |
| `/burst` | 10 parallel requests |

## Stop

```bash
docker compose down        # stop
docker compose down -v     # stop + delete data
```