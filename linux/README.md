# 📋 Runink Cluster Service Port Table

| Service | Port | Protocol | Purpose |
|:--------|:-----|:---------|:--------|
| **Global API Server (`apiserver`)** | `8443` | `gRPC/TLS` | Entry point for CLI/API requests |
| **Herd Directory (`herd-directory`)** | `8500` | `gRPC/TLS` | Resolve Herd-to-Raft mappings |
| **Finance Barn (`barn`)** | `8700` | `gRPC` | Raft-backed KV operations (contracts, metadata, secrets) |
| **Analytics Barn (`barn`)** | `8701` | `gRPC` | Separate instance for Analytics Herd |
| **Finance Raft Cluster Communication** | `8800` | `gRPC/mTLS` | Raft peer-to-peer log replication (5 nodes, gossip+raft) |
| **Analytics Raft Cluster Communication** | `8801` | `gRPC/mTLS` | Raft peer-to-peer log replication (Analytics partition) |
| **Finance Scheduler (`herd`)** | `8900` | `gRPC` | DAG planning and slice placement |
| **Analytics Scheduler (`herd`)** | `8901` | `gRPC` | Separate scheduler for Analytics Herd |
| **Runi Agent (`runi`)** | `9000` | `HTTP` | Slice registration and management |
| **Runi Prometheus Exporter** | `9100` | `HTTP` | Expose metrics endpoint (`/metrics`) for Prometheus scraping |
| **Barn Metrics Exporter (optional)** | `9200` | `HTTP` | Barn (Badger/Raft) internal metrics |
| **Herd Metrics Exporter (optional)** | `9300` | `HTTP` | Herd scheduler slice/dag planning stats |
| **Slice Ephemeral Services (per slice)** | `30000-32767` | `TCP` | Dynamic ephemeral ports for slice sidecars / mini services (like K8s NodePort range) |

---

# 📦 Quick Visual Port Grouping

| Range | Usage |
|:------|:------|
| `8400-8500` | Global Control Plane Ports (API, Directory) |
| `8700-8900` | Herd-Specific Metadata and Raft |
| `9000-9400` | Worker Metrics (runi, barn, herd exporters) |
| `30000-32767` | Dynamic ephemeral Slice mini-services |

---

# 🔥 Bonus Recommendations

If you want **ultimate production-grade resilience**,  
later you can:

| Recommendation | Benefit |
|:---------------|:--------|
| Enable strict firewall rules (iptables/nftables) | Only expose required ports (e.g., block dynamic slice ports externally) |
| Run node-local Prometheus sidecar | Lower network usage for scraping |
| TLS rotate all mTLS certificates via short-lived tokens | Safer Raft mesh management |
| Rate-limit slice mini-service ports | Prevent slice-level DDoS scenarios |
