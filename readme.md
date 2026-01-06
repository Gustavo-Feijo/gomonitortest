# GoMonitor

___

This is a project with main purpose of improving some of my existing knowledge in Golang, Redis, Docker and other technologies.

I also aim to learn new technologies, specially in observability instrumentation, using Tempo, OTEL, Prometheus, Loki, Alloy and Grafana.
The application will be a simple API with user creation and basic authentication (Also learning how to implement from zero)

---

## Services

- #### APP
  - Gin API with user creation and authentication.
  - Observability implementation (Logging, Metrics and Tracing)
- #### Redis
  - Shared cache.
- #### PostgreSQL
  - Shared database.
- #### Prometheus
  - Application metrics scraping.
  - Redis and Postgres metrics.
- #### Postgres-exporter
  - Postgres metrics to prometheus.
- #### Redis-exporter
  - Redis metrics to prometheus.
- #### Tempo
  - Application traces (Gin, Gorm, Redis and any other dependency)
- #### Loki
  - Application logs storage.
- #### Alloy
  - Application logs scraping from docker services and transformation.
- #### Grafana
  - Centralized data visualization.
  - Postgres and redis dashboards.
  - Loki, Prometheus and Tempo datasources.
  - Loki and Tempo setted up to link to each other.

## Setup

Create a '.env' file in simmilar format to 'example.env', if needed, change the configurations (Should run fine without any changes).

Run `docker compose up --build`