# Data Lineage & Metadata Tracking – Runink

Runink pipelines are designed to be fully **traceable**, **auditable**, and **schema-aware**. With built-in lineage support, every pipeline can generate:

- Visual DAGs of data flow and dependencies
- Metadata snapshots with schema versions and field hashes
- Run-level logs for audit, debugging, and compliance

This guide walks through how Runink enables robust **data observability** and **governance by default**.

---

## 🔍 What Is Data Lineage?
Lineage describes where your data came from, what happened to it, and where it went.

In Runink, every pipeline run captures:
- **Sources**: file paths, streaming URIs, tags
- **Stages**: steps applied, transform versions
- **Contracts**: schema file, struct, and hash
- **Sinks**: output paths, filters, conditions
- **Run metadata**: timestamps, roles, record count

---

## 📈 Generate a Lineage Graph
```bash
runi lineage --scenario features/orders.feature --out lineage/orders.svg
```

The graph shows:
- Inputs and outputs
- All applied steps
- Contract versions and field diff hashes
- Optional labels (e.g., `role`, `source`, `drift`)

---

## 🧾 Per-Run Metadata Log
Every run emits a record like:
```json
{
  "run_id": "run-20240423-abc123",
  "stage": "JoinUsersAndOrders",
  "contract": "user_order_v2.json",
  "schema_hash": "b72cd1a",
  "records_processed": 9123,
  "timestamp": "2024-04-23T11:02:00Z",
  "role": "analytics",
  "drift_detected": false
}
```

---

## 🧪 Snapshotting & Version Tracking
You can snapshot inputs/outputs with:
```bash
runi snapshot --contract contracts/user.json --out snapshots/users_2024-04-23.json
```
And later compare against historical output.

---

## 🚨 Drift Detection
Runink detects when incoming data deviates from expected contract:
```bash
runi contract diff --old v1.json --new incoming.json
```
Or as part of a scenario run:
```bash
runi run --verify-contract
```
This flags:
- Missing/extra fields
- Type mismatches
- Tag mismatches (e.g., missing `pii`, `access`)

---

## 🔐 Metadata for Compliance
Attach metadata to every stage:
```go
type StageMetadata struct {
  RunID     string
  Role      string
  Contract  string
  Hash      string
  Source    string
  Timestamp string
}
```
Send this to a:
- Document DB (e.g. Mongo)
- Data lake (e.g. MinIO, S3)
- Audit stream (e.g. Kafka topic)

---

## 📡 Monitoring & Observability
Runink supports Prometheus metrics per stage:
- `runi_records_processed_total`
- `runi_stage_duration_seconds`
- `runi_schema_drift_detected_total`
- `runi_invalid_records_total`

---

## 🧠 Example Use Cases
| Role         | How Lineage Helps                  |
|--------------|------------------------------------|
| Data Engineer| Debug broken joins, drift, formats |
| Analyst      | Understand where numbers came from |
| Governance   | Prove schema conformance           |
| ML Engineer  | Snapshot training input lineage    |

---

## Summary
Runink provides end-to-end data lineage as a **first-class feature**, not an afterthought:
- Built-in visual DAGs
- Contract + transform metadata
- Auditable, role-aware stage outputs
- Real-time observability with metrics

Lineage lets you move fast **without breaking trust**.
