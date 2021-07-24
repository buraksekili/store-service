# store-service

> **CAVEAT:** This repo is still under development. Currently, I am working on [these tasks](https://github.com/buraksekili/store-service/issues/3).
> You can follow the development from [branches](https://github.com/buraksekili/store-service/branches).

`store-service` is a microservice-based web application. 


The purpose of the development of `store-service` is to gain hands-on experience in microservices architecture.


## Installation

```console
$ git clone https://github.com/buraksekili/store-service.git
$ cd store-service/
$ make
$ make containers
$ make run
```

## Endpoints
- **RabbitMQ Management**   : *localhost:15672/*
- **Prometheus**: *localhost:9090/*
- **Grafana**: *localhost:3000/*
- [users](https://github.com/buraksekili/store-service/tree/master/users): *localhost:8282*
- [products](https://github.com/buraksekili/store-service/tree/master/products): *localhost:8181*

The default username and password are `admin` for `Grafana`.

The default ports for `users` and `products` are defined in [.env file](https://github.com/buraksekili/store-service/blob/master/docker/.env). You can change it via [.env](https://github.com/buraksekili/store-service/blob/master/docker/.env), if you need.

## Local development

```
make                        : Compiles all services and saves them into `./build` folder.


make containers             : Creates Docker containers for all compiled services under `./build`.


make <service_name>         : Rebuilds <service_name>.

> For example;
    `make users`

make docker_<service_name>  : Creates a new Docker image for <service_name>.

make run                    : Runs all containers and other services defined in docker-compose.yaml file.
```

For all possible <service_name>, see [Makefile](https://github.com/buraksekili/store-service/blob/bd576d5ae6bd13867932bed7ec106a24b6f1625d/Makefile#L1).

## Architecture

| Service       |  Language      | |
| ------------- |:-------------:| -----|
| users  |   Go          | `users` provides REST endpoints to operate CRUD <br/> operations for users of the store. <br/> It emits a message for `RabbitMQ` after new user created. <br/> This emitted message is forward to email service which is <br/> responsible to sending emails to the newly added users.|
| products  |   Go          | `products` provides REST endpoints to operate <br/> CRUD operations for products and vendors of the store.   |
| emailer  |   Go          |  `emailer` sends an email to newly added users. <br /> See [`.env`](https://github.com/buraksekili/store-service/blob/master/docker/.env) for required configurations for SMTP. <br />Note: `to` field is not required, it is just for testing purposes.|

## Features

- [**AMQP - RabbitMQ**](https://github.com/buraksekili/store-service/tree/master/amqp)
- [**Prometheus & Grafana**](https://github.com/buraksekili/store-service/tree/master/metrics)
- [**MongoDB**](https://github.com/buraksekili/store-service/tree/master/db/mongo)

## HTTP API Examples

### [Users endpoints](https://github.com/buraksekili/store-service/tree/gokit/users)
### [Products endpoints](https://github.com/buraksekili/store-service/tree/gokit/products)

## Disclaimer

This project is built for *educational purposes*. It does not assert to be production ready. 

## Acknowledgments

Inspired by [GCP's microservices demo](https://github.com/GoogleCloudPlatform/microservices-demo)

Highly influenced by [Mainflux](https://github.com/mainflux/mainflux) 