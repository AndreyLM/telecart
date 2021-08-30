# Тестовое задание Go + MQTT + SQLite

## Runnig service
1. Run **Eclipse mosquitto**. Can install and run locally or using docker-composer in folder **./build**. Ideally there must be Dockerfile for our service. If running locally **allow_anonymous true** 
2. Update config file **./bin/telecart.yalm**
3. Update all dependencies:
> go mod tidy
3. Run service:
> make run

## Tests
> cd ./internal && go test 

## Should do
1. Random clientID for mqtt
2. When service is starting, check first message for avoiding duplication or could be better to use **ID** from **mqtt** and store it in db (now **id** is auto)

