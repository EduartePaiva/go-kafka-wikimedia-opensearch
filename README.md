# go-kafka-wikimedia-opensearch

A real-time data pipeline written in **Go** that streams live edit events from the [Wikimedia Recent Changes feed](https://stream.wikimedia.org/v2/stream/recentchange) into **OpenSearch** using **Apache Kafka** as the message broker.

---

## What this project does

Wikimedia publishes a continuous Server-Sent Events (SSE) stream of every edit made across its platforms (Wikipedia, Wikidata, etc.). This project taps into that stream and builds a full end-to-end pipeline:

```
Wikimedia SSE Stream
       │
       ▼
  [ Producer ]  ── publishes events ──▶  Kafka Topic (wikimedia.recentchange)
                                                │
                                                ▼
                                        [ Consumer ]  ── indexes documents ──▶  OpenSearch
                                                                                      │
                                                                                      ▼
                                                                          OpenSearch Dashboards
```

1. **Producer** — connects to the Wikimedia SSE feed and forwards each event as a message to a Kafka topic.
2. **Consumer** — reads from that Kafka topic and bulk-indexes each event into OpenSearch for storage and querying.
3. **OpenSearch Dashboards** — lets you visually explore and analyze the indexed data in real time.

---

## Tech stack

| Layer | Technology |
|---|---|
| Language | Go |
| Kafka client | [IBM/sarama](https://github.com/IBM/sarama) |
| SSE client | [r3labs/sse](https://github.com/r3labs/sse) |
| Search & storage | OpenSearch 3.5 |
| Schema management | Confluent Schema Registry |
| Kafka UI | [Kafbat UI](https://github.com/kafbat/kafka-ui) |
| Infrastructure | Docker Compose |
| Task runner | [Task](https://taskfile.dev) |

---

## Project structure

```
.
├── cmd/
│   ├── producer/       # SSE → Kafka producer
│   └── consumer/       # Kafka → OpenSearch consumer
├── internal/           # Shared internal packages
├── compose.yaml        # Docker Compose for Kafka, OpenSearch, Schema Registry & UIs
├── Taskfile.yml        # Dev task shortcuts (fmt, tidy, update)
├── go.mod
└── go.sum
```

---

## Getting started

### Prerequisites

- [Go 1.22+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/)
- [Task](https://taskfile.dev/installation/) (optional, for convenience commands)

### 1. Start the infrastructure

```bash
docker compose up -d
```

This spins up:
- **Kafka** (KRaft mode, no Zookeeper) on `localhost:9094`
- **Kafka UI** on [http://localhost:8070](http://localhost:8070)
- **Schema Registry** on [http://localhost:8081](http://localhost:8081)
- **OpenSearch** on [http://localhost:9200](http://localhost:9200)
- **OpenSearch Dashboards** on [http://localhost:5601](http://localhost:5601)

### 2. Run the producer

```bash
go run ./cmd/producer
```

The producer will start consuming the Wikimedia SSE stream and publishing events to the `wikimedia.recentchange` Kafka topic.

### 3. Run the consumer

```bash
go run ./cmd/consumer
```

The consumer will read from Kafka and bulk-index events into OpenSearch. You can then explore the data at [http://localhost:5601](http://localhost:5601).

---

## Useful Task commands

```bash
task tidy      # Format code and tidy module dependencies
task update    # Update all Go packages to their latest versions
```

---

## Key concepts practiced

- **Kafka producer** in Go using the Sarama client
- **Kafka consumer** with offset management and bulk processing
- **SSE (Server-Sent Events)** consumption from a live public stream
- **OpenSearch indexing** via the official Go client
- **KRaft mode** Kafka setup (no Zookeeper dependency)
- **Schema Registry** integration for schema management
- **Docker Compose** orchestration of a multi-service data stack