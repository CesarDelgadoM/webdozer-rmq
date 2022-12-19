# Webdozer-RMQ

## RabbitMQ application

Application to extract urls using different data types of redis, consist of a producer and a consumer with objective of filter the searchs words and write a file with the urls found.

## Deploy

In the deployments folder is the docker compose file for the image RabbitMQ and Redis database

```bash
docker-compose up
```

## app-producer

Extracts urls from redis database that the webdozer application generates

```bash
go run main.go
```

## app-consumer

Read the queues that the app-producer create, filter by search words the urls that is in the redis database that the webdozer application generates and write a file with the urls found.

```bash
go run main.go
```

## Note

The application webdozer is in this github profile
