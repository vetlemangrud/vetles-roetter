BASE=http://localhost:8080

# list carrots
curl -s "$BASE/" | jq

# add a carrot
curl -s -X POST "$BASE/" | jq
