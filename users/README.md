# Users service

Users service presents HTTP endpoints for the user related operations.

The default service port for the Users service is `:8282`. You can modify the default value of it through the [docker/.env](https://github.com/buraksekili/store-service/blob/master/docker/.env) file (from the project's root directory).

 ```
S_USERS_PORT=8282
```

> The default value of the following `<service_port>` is `8282`.


### Login

```bash
curl -s -S -i -X POST -H "Content-Type: application/json" http://localhost/users/login -d '{"email":"<user_email>", "password":"<user_password>"}'
```

### Create User

```bash
curl -s -S -X POST -H "Content-Type: application/json" http://localhost:<service_port>/users/signup -d '{"username": "<user_name>", "email":"<user_email>", "password":"<user_password>"}'
```

### Get User

```bash
curl -s -S -X GET http://localhost:<service_port>/users/<user_id>
```

### Get All Users

```bash
curl -s -S -X GET http://localhost:<service_port>/users
```

### Create Vendor

```bash
curl -s -S -X POST -H "Content-Type: application/json" http://localhost:<service_port>/vendors -d '{"name": "<vendor_name">, "description":"<vendor_description>"}'
```

### Get All Vendors

```bash
curl -s -S -X GET http://localhost:<service_port>/vendors
```