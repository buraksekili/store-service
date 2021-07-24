# Products Service

Products service presents HTTP endpoints for the product related operations.

The default service port for the Products service is `:8181`. You can modify the default value of it through the [docker/.env](https://github.com/buraksekili/store-service/blob/master/docker/.env) file (from the project's root directory).

```
S_PRODUCTS_PORT=8181
```

> The default value of the following `<service_port>` is `8181`.

### Create A Product

```bash
curl -s -S -X POST -H "Content-Type: application/json" http://localhost:<service_port>/products -d '{
    "name": <product_name>,
    "category": <product_category>,
    "description": <product_description>,
    "price": <product_price>,
    "imageUrl": <product_image_url>,
    "stock": <product_stock>,
    "vendor": {
        "id": <vendor_id>
    }
}'
```

### Get All Products

```bash
curl -s -S -i -X GET http://localhost:<service_port>/products
```

### Get A Product

```bash
curl -s -S -i -X GET http://localhost:<service_port>/products/<product_id>
```

### Get Products of a Vendor

```bash
curl -s -S -i -X GET http://localhost:<service_port>/products/venders/<vendor_id>
```


### Get Comments of a Product

```bash
curl -s -S -i -X GET http://localhost:<service_port>/products/comments/{product_id}
```

### Get All Comments 

```bash
curl -s -S -i -X GET http://localhost:<service_port>/products/comments
```

