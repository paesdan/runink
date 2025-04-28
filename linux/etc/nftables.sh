#!/usr/bin/env bash

set -euo pipefail

sysctl -w net.ipv6.conf.all.disable_ipv6=1
sysctl -w net.ipv6.conf.default.disable_ipv6=1

# Flush old rules
nft flush ruleset

# Create table
nft add table inet filter

# Define chains
nft add chain inet filter input { type filter hook input priority 0 \; policy drop \; }
nft add chain inet filter forward { type filter hook forward priority 0 \; policy drop \; }
nft add chain inet filter output { type filter hook output priority 0 \; policy accept \; }

# Allow loopback
nft add rule inet filter input iif lo accept

# Allow already established connections
nft add rule inet filter input ct state established,related accept

# Accept SSH (optional if debugging, otherwise remove)
nft add rule inet filter input tcp dport 22 ct state new accept

# === Allow Essential Control Plane Ports ===
# API Server
nft add rule inet filter input tcp dport 8443 accept

# Herd Directory
nft add rule inet filter input tcp dport 8500 accept

# === Herd Partition Ports ===
# Barn (KV Store)
nft add rule inet filter input tcp dport 8700 accept
nft add rule inet filter input tcp dport 8701 accept

# Raft Communication (mTLS)
nft add rule inet filter input tcp dport 8800 accept
nft add rule inet filter input tcp dport 8801 accept

# Herd Scheduler
nft add rule inet filter input tcp dport 8900 accept
nft add rule inet filter input tcp dport 8901 accept

# === Worker Node Ports ===
# Runi Agent
nft add rule inet filter input tcp dport 9000 accept

# Prometheus Exporter (Runi)
nft add rule inet filter input tcp dport 9100 accept

# Optional Metrics (Barn, Herd)
nft add rule inet filter input tcp dport 9200 accept
nft add rule inet filter input tcp dport 9300 accept

# === Ephemeral Slice Ports ===
# Only allow from internal (localhost or VPC, optional strictness)
nft add rule inet filter input tcp dport 30000-32767 ip saddr 10.0.0.0/8 accept

# Apply conntrack limits to ephemeral slice services
nft add rule inet filter input tcp dport 30000-32767 ct state new ct count over 50 drop comment "Limit 50 concurrent connections per slice"

# === ICMP (Ping) (optional) ===
nft add rule inet filter input icmp type echo-request accept

# === Log + Drop Anything Else ===
nft add rule inet filter input counter drop
