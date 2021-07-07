# go-redis
A sample API using Redis SortedSet for caching.

## Setup
- `docker-compose up`
- `make run`

## Datastores
- DynamoDB: http://localhost:28000/
- Redis: `redis-cli -h 127.0.0.1 -p 26379`

## Endpoints
```bash
curl -XPOST -H "Content-Type: application/json" \
-d '{"title": "Threaddy"}' \
"localhost:8888/threads" | jq
```
```bash
curl "localhost:8888/messages?thread_id=01FA0VNJWRCA1TRWGB4BSJEE57&min=1625674158866341&max=1625674165835997&limit=10" \
| jq
```
```bash
curl -XPOST -H "Content-Type: application/json" \
-d '{
  "thread_id": 
  "01FA0VNJWRCA1TRWGB4BSJEE57", 
  "content": "Hello 1", 
  "user_id": 1
}' "localhost:8888/messages" | jq
```
```bash
curl -XPOST -H "Content-Type: application/json" \
-d '{
  "thread_id": "01FA0VNJWRCA1TRWGB4BSJEE57", 
  "sent_at": 1625674158866341, 
  "kind": "+1", 
  "user_id": 1
}' "localhost:8888/reactions" | jq
```
