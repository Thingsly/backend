# backend

The plug-in IoT platform developed by Go language has high performance, low entry and easy expansion. Support MQTT, multi type device access and visualization, automation, alarm, rule engine and other functions.

```cmd

docker-compose -f docker-compose-dev.yml up 
docker-compose -f docker-compose-dev.yml down

System Administrator: demo-super@thingsly.vn / 123456
Tenant Administrator: demo-tenant@thingsly.vn / 123456

psql -U postgres -d Thingsly
docker exec -it backend-postgres-1 psql -U postgres -d Thingsly

http://localhost:9999/swagger/index.html
http://localhost:9999/metrics-viewer

```
