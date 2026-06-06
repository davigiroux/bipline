#!/usr/bin/env bash
set -euo pipefail
: "${BUFFER_API_KEY:?BUFFER_API_KEY must be set}"

QUERY='{"query":"query IntrospectionQuery { __schema { queryType { name } mutationType { name } subscriptionType { name } types { ...FullType } directives { name description locations args { ...InputValue } } } } fragment FullType on __Type { kind name description fields(includeDeprecated: true) { name description args { ...InputValue } type { ...TypeRef } isDeprecated deprecationReason } inputFields { ...InputValue } interfaces { ...TypeRef } enumValues(includeDeprecated: true) { name description isDeprecated deprecationReason } possibleTypes { ...TypeRef } } fragment InputValue on __InputValue { name description type { ...TypeRef } defaultValue } fragment TypeRef on __Type { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name ofType { kind name } } } } } } } }"}'

mkdir -p internal/buffer

curl -sf -X POST https://api.buffer.com \
  -H "Authorization: Bearer $BUFFER_API_KEY" \
  -H "Content-Type: application/json" \
  -d "$QUERY" \
  | jq '.data' \
  > internal/buffer/schema.graphql

echo "schema.graphql written ($(wc -l < internal/buffer/schema.graphql) lines)"
