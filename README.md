# store-service

`store-service` is a microservice-based web application. 
It does not provide fully featured services that you might expect to see in real web applications such as complex *DB queries*, *security*, or *advanced search*.
(For example, I did not even hash the user password). 

<br />

The purpose of the development of `store-service` is to gain hands-on experience in microservices.

`emailservice` consumes `add_user` events which can be published by `userservice` after new user created. Metrics are scraped by Prometheus with 15s interval. Also, you can check these metrics via Grafana.

## Installation

```console
$ git clone https://github.com/buraksekili/store-service.git
$ cd store-service/
$ docker-compose up -d
```

## Endpoints
- **RabbitMQ Management**   : *localhost:15672/*
- **Prometheus**: *localhost:9090/*
- **Grafana**: *localhost:3000/*
- [userservice](https://github.com/buraksekili/store-service/blob/master/src/userservice/rest/rest.go): *localhost:8282/users/{relative_path}*
- [productservice](https://github.com/buraksekili/store-service/blob/master/src/productservice/main.go): <span>localhost:8282/{products | comments | vendors}/{relative_path}

### Examples
To create new user;
```console
curl --location --request POST 'http://localhost:8282/users/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "NewUser",
    "email": "NewUserMail@newUser.com",
    "password": "NewUserPassword"
}'
```

## Architecture

| Service       | Language      | |
| ------------- |:-------------:| -----|
| userservice   |   Go          | `userservice` provides REST endpoints to operate CRUD <br/> operations for users of the store. <br/> It emits a message for `RabbitMQ` after new user created. <br/> This emitted message is forward to email service which is <br/> responsible to sending emails to the newly added users.|
| productservice|   Go          | `productservice` provides REST endpoints to operate <br/> CRUD operations for products and vendors of the store.   |
| emailservice  |   Go          |  `emailservice` sends an email to newly added users. <br /> See [`config.json`](https://github.com/buraksekili/store-service/blob/master/src/emailservice/config.json) for required configurations for SMTP. |

## Features

- [**AMQP - RabbitMQ**](https://github.com/buraksekili/store-service/tree/master/amqp)
- [**Prometheus & Grafana**](https://github.com/buraksekili/store-service/tree/master/metrics)
- [**MongoDB**](https://github.com/buraksekili/store-service/tree/master/db/mongo)
