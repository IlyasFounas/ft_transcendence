#!/bin/sh
set -e
IP=$(ip route get 1.1.1.1 | awk '{for(i=1;i<=NF;i++) if($i=="src"){print $(i+1); exit}}')
[ -z "$IP" ] && { echo "IP introuvable" >&2; exit 1; }
ENV_FILE="${1:-./docker-config/.env}"
SECRET_SCALEWAY_API_KEY=$(grep '^SECRET_SCALEWAY_API_KEY=' "$ENV_FILE" | cut -d= -f2-)
[ -z "$SECRET_SCALEWAY_API_KEY" ] && { echo "clé API manquante" >&2; exit 1; }
curl -sS -f https://api.scaleway.com/domain/v2beta1/dns-zones/ocefunding.com/records \
  -X PATCH \
  -H "Content-Type: application/json" \
  -H "X-Auth-Token: $SECRET_SCALEWAY_API_KEY" \
  -d "{
    \"return_all_records\": false,
    \"changes\": [
      {
        \"set\": {
          \"id_fields\": { \"name\": \"campus\", \"type\": \"A\" },
          \"records\": [
            { \"name\": \"campus\", \"data\": \"$IP\", \"type\": \"A\", \"ttl\": 300 }
          ]
        }
      }
    ]
  }"