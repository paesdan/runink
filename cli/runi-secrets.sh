#!/bin/bash
# runi-secrets.sh — inject secrets from .env into runink's herd secrets store

ENV_FILE="contracts/cdm_trade/fdc3events.env"
HERD="finance"

echo "[runi secrets] 🔐 Injecting secrets for herd: $HERD from $ENV_FILE"

# Export env file and loop over each line to inject
while IFS= read -r line || [ -n "$line" ]; do
  if [[ "$line" =~ ^[A-Za-z_]+= ]]; then
    KEY=$(echo "$line" | cut -d= -f1)
    VALUE=$(echo "$line" | cut -d= -f2-)
    echo "  → Injecting $KEY"
    runi secrets inject --herd "$HERD" --key "$KEY" --value "$VALUE"
  fi
done < "$ENV_FILE"

echo "[runi secrets] ✅ All secrets injected from $ENV_FILE"
