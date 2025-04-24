#!/bin/bash
# runi-run.sh — run FDC3 trade scenario with automatic DAG resolution

SCENARIO_FILE="features/fdc3events.dsl"
CONTRACT_FILE="contracts/cdm_trade/fdc3events.go"
ENV_FILE="contracts/cdm_trade/fdc3events.env"
HERD="finance"

echo "[runi run] 🔍 Compiling DAG for: $SCENARIO_FILE"
runi compile --scenario "$SCENARIO_FILE"

echo "[runi run] 🔐 Loading environment variables from $ENV_FILE"
export $(grep -v '^#' "$ENV_FILE" | xargs)

echo "[runi run] 🚀 Executing pipeline for herd: $HERD"
go run rendered/fdc3events_main.go

echo "[runi run] ✅ Execution completed"
