# Schema & Contract Management – Runink

Runink enables **data contracts** as native Go structs — giving you strong typing, version tracking, schema validation, and backward compatibility across pipelines.

This guide shows how to define, version, test, and enforce schema contracts in your pipelines.

---

## 📦 What Is a Contract?
A contract in Runink is a schema definition used to:
- Validate incoming and outgoing data
- Detect schema drift
- Provide PII and RBAC tagging
- Drive pipeline generation and testing

Contracts are generated from Go structs annotated with tags.

---

## ✍️ Defining a Contract
```go
package contracts

type Order struct {
  OrderID    string  `json:"order_id"`
  CustomerID string  `json:"customer_id"`
  Amount     float64 `json:"amount"`
  Timestamp  string  `json:"timestamp"`
  Notes      string  `json:"notes" pii:"true" access:"support"`
}
```

Then run:
```bash
runi contract gen --struct contracts.Order --out contracts/order.json
```

---

## ✅ Enforcing a Contract
```gherkin
Given the contract: contracts/order.json
```
Or:
```bash
runi run --verify-contract
```
Runink ensures that all records match the expected schema.

---

## 🔍 Schema Drift Detection
Compare current vs expected schema:
```bash
runi contract diff --old v1.json --new v2.json
```
Output shows added, removed, or changed fields, types, tags, and ordering.

---

## 📊 Hashing and Snapshotting
Each contract has a hash for:
- Version tracking
- Lineage graph integrity
- Change detection

```bash
runi contract hash contracts/order.json
```

Snapshot for reproducibility:
```bash
runi snapshot --contract contracts/order.json --out snapshots/order_v1.json
```

---

## 🧬 Advanced Tags
- `pii:"true"` – marks field as sensitive
- `access:"finance"` – restricts field to roles
- `enum:"pending,approved,rejected"` – enum constraint (optional)
- `required:"true"` – fail if field is null or missing

---

## 🧪 Contract Testing
Use golden tests to assert schema correctness:
```bash
runi test --scenario features/orders.dsl
```
And diff output against expected:
```bash
runi diff --gold testdata/orders.golden.json --new out/orders.json
```

---

## 🗃️ Contract Catalog
Generate an index of all contracts in your repo:
```bash
runi contract catalog --out docs/contracts.json
```
This can be plugged into:
- Docs browser
- Contract registry
- CI schema check

---

## 🧾 Example Contract Output
```json
{
  "name": "Order",
  "fields": [
    { "name": "order_id", "type": "string" },
    { "name": "amount", "type": "float64" },
    { "name": "notes", "type": "string", "pii": true, "access": "support" }
  ],
  "hash": "a94f3bc..."
}
```

---

## Summary
Contracts in Runink power everything:
- Schema validation
- RBAC and compliance
- Pipeline generation
- Test automation
- Lineage and snapshots

Use contracts to make your data:
- Safe
- Trustworthy
- Documented
- Governed