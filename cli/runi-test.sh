#!/bin/bash
# runi-test.sh — run golden validation tests for fdc3events scenario

SCENARIO_FILE="features/cdm_trade/fdc3events.dsl"
CONTRACT_FILE="contracts/cdm_trade/fdc3events.go"
ENV_FILE="contracts/cdm_trade/fdc3events.env"
GOLDEN_DIR="golden/cdm_trade/"
HERD="finance"

echo "[runi test] 🔐 Loading environment from $ENV_FILE"
export $(grep -v '^#' "$ENV_FILE" | xargs)

echo "[runi test] 🧪 Running golden test for: $SCENARIO_FILE"
runi test \
  --scenario "$SCENARIO_FILE" \
  --contract "$CONTRACT_FILE" \
  --golden "$GOLDEN_DIR" \
  --herd "$HERD"

echo "[runi test] ✅ Completed golden test validation"
