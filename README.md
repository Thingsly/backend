# backend

The plug-in IoT platform developed by Go language has high performance, low entry and easy expansion. Support MQTT, multi type device access and visualization, automation, alarm, rule engine and other functions.

```cmd

docker-compose -f docker-compose-dev.yml up 
docker-compose -f docker-compose-dev.yml down

System Administrator: demo-super@thingsly.vn / 123456
Tenant Administrator: demo-tenant@thingsly.vn / 123456

psql -U postgres -d Thingsly
docker exec -it backend-postgres-1 psql -U postgres -d Thingsly

docker exec -it backend-redis-1 redis-cli
AUTH redis
SELECT 1

CONFIG SET notify-keyspace-events Ex

# Xem số lượng keys
DBSIZE

# Xem tất cả keys
KEYS *

# Xem type của key
TYPE key_name

# Xem value (cho String)
GET key_name

# Xem tất cả fields (cho Hash)
HGETALL key_name

# Xem TTL
TTL key_name

http://localhost:9999/swagger/index.html
http://localhost:9999/metrics-viewer

```
