# store-service


> *CAVEAT*: This repo is still under development. Currently, I am working on [these tasks](https://github.com/buraksekili/store-service/issues/3).
> You can follow the development from [branches](https://github.com/buraksekili/store-service/branches).

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
$ make
$ make images
$ make run
```

## Endpoints
- **RabbitMQ Management**   : *localhost:15672/*
- **Prometheus**: *localhost:9090/*
- **Grafana**: *localhost:3000/*
- [userservice](https://github.com/buraksekili/store-service/blob/master/src/userservice/rest/rest.go): *localhost:8282/users/{relative_path}*
- [productservice](https://github.com/buraksekili/store-service/blob/master/src/productservice/main.go): <span>localhost:8282/{products | comments | vendors}/{relative_path}

The default username and password are `admin` for `Grafana`.

### Local development


`make`: Compiles all services.

`make images`: Creates Docker images for all compiled services.

`make <service_name>`: Rebuilds `<service_name>`.

`make docker_<service_name>`: Creates new Docker image for `<service_name>`.

`make run`: Runs all compiled images and other services.

## Architecture

| Service       |  Language      | |
| ------------- |:-------------:| -----|
| userservice   |   Go          | `userservice` provides REST endpoints to operate CRUD <br/> operations for users of the store. <br/> It emits a message for `RabbitMQ` after new user created. <br/> This emitted message is forward to email service which is <br/> responsible to sending emails to the newly added users.|
| productservice|   Go          | `productservice` provides REST endpoints to operate <br/> CRUD operations for products and vendors of the store.   |
| emailservice  |   Go          |  `emailservice` sends an email to newly added users. <br /> See [`.env`](https://github.com/buraksekili/store-service/blob/master/docker/.env) for required configurations for SMTP. <br />Note: `to` field is not required, it is just for testing purposes.|

## Features

- [**AMQP - RabbitMQ**](https://github.com/buraksekili/store-service/tree/master/amqp)
- [**Prometheus & Grafana**](https://github.com/buraksekili/store-service/tree/master/metrics)
- [**MongoDB**](https://github.com/buraksekili/store-service/tree/master/db/mongo)

### Examples
Create new user:
```
curl --location --request POST 'http://localhost:8282/users/signup' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "NewUser",
    "email": "NewUserMail@newUser.com",
    "password": "NewUserPassword"
}'
```

Add a vendor:
```
curl --location --request POST 'http://localhost:8181/vendors' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Apple Inc.",
    "description": "Apple Inc. is an American multinational technology company that specializes in consumer electronics, computer software, and online services."
}'
```

Add a product:
```
curl --location --request POST 'http://localhost:8181/products' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "iPhone 12 Pro",
    "category": "Electronics",
    "description": "Meet the ultimate iPhone. With the fastest smartphone chip. 5G speed.",
    "price": 1799,
    "imageUrl": "https://store.storeimages.cdn-apple.com/4668/as-images.apple.com/is/iphone-12-pro-family-hero?wid=470&hei=556&fmt=jpeg&qlt=95&.v=1604021663000",
    "stock": 50,
    "vendor": {
        "name": "Apple Inc."
    }
}'
```

Login: 
```
curl --location --request POST 'http://localhost:8282/users/login' \
--header 'Content-Type: application/json' \
--data-raw '{
    "username": "<Username>",
    "password": "<Password>"
}'
```

Add a comment:
```
curl --location --request POST 'http://localhost:8181/comments/60c9e0510fdfd3501763408c' \
--header 'Content-Type: application/json' \
--data-raw '{
    "Owner": "<Owner Name>",
    "Content": "<Comment>"
}'
```