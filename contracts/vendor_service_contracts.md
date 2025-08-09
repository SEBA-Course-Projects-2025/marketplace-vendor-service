# Vendor Service API Contracts

---

__Swagger API Documentation: https://marketplace-vendor-service.onrender.com/swagger/index.html__

## 0. Vendor Login

__POST ```/auth/login```__

__Body:__

```json
{
  "email": "string",
  "password": "string"
}
```

__Response:__

```json
{
  "token": "jwt-token"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```401 Unauthorized```
- ```500 Internal Server Error```

## 1. Get Vendor Account Information

__GET ```api/account```__

__Response:__

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "description": "string",
  "logo": "string",
  "address": "string",
  "website": "string",
  "catalogId": "uuid"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---

## 2. Update Vendor Account Information

__PUT ```api/account```__

__Body:__

```json
{
  "email": "string",
  "name": "string",
  "description": "string",
  "logo": "string",
  "address": "string",
  "website": "string"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---

## 3. Modify Vendor Account Information

__PATCH ```api/account```__

__Body:__

```json
{
  "email": "string", (optional)
  "name": "string", (optional)
  "description": "string", (optional)
  "logo": "string", (optional)
  "address": "string", (optional)
  "website": "string" (optional)
}
```

__Response:__

```json
{
  "id": "uuid",
  "email": "string",
  "name": "string",
  "description": "string",
  "logo": "string",
  "address": "string",
  "website": "string"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---

## 4. Get Reviews On Products

__GET ```api/reviews```__

__Response:__

```json
[
  {
    "review_id": "uuid",
    "vendor_id": "uuid",
    "product_id": "uuid",
    "reviewer_id": "uuid",
    "reviewer_name": "string",
    "rating": "float64",
    "comment": "string",
    "date": "date"
  }
]
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```500 Internal Server Error```

---

## 5. Get Review On Product

__GET ```api/reviews/:reviewId```__

__Response:__

```json
{
  "review_id": "uuid",
  "vendor_id": "uuid",
  "product_id": "uuid",
  "reviewer_id": "uuid",
  "reviewer_name": "string",
  "rating": "float64",
  "comment": "string",
  "date": "date",
  "replies": [
    {
      "reply_id": "uuid",
      "replier_id": "uuid",
      "name": "string",
      "comment": "string",
      "date": "date"
    }
  ]
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---

## 6. Add a Reply On a Review

__POST ```api/reviews/:reviewId/replies```__

__Body:__

```json
{
  "comment": "string"
}
```

__Response:__

```json
{
  "review_id": "uuid",
  "reply_id": "uuid",
  "replier_id": "uuid",
  "name": "string",
  "comment": "string",
  "date": "date"
}
```

__Status codes:__
- ```201 Created (success)```
- ```400 Bad Request```
- ```500 Internal Server Error```

---

## 7. Updating Review Reply Comment

__PATCH ```api/reviews/:reviewId/replies/:replyId```__

__Body:__

```json
{
  "comment": "string"
}
```

__Response:__

```json
{
  "review_id": "uuid",
  "reply_id": "uuid",
  "replier_id": "uuid",
  "name": "string",
  "comment": "string",
  "date": "date"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---


## 8. Check Orders

__GET ```api/orders```__

__Response:__

```json
{
  "order_id": "uuid",
  "customer_id": "uuid",
  "vendor_id": "uuid",
  "items": ["string"],
  "total_price": "float64",
  "status": "string",
  "date": "date"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```500 Internal Server Error```

---

## 9. Get One Order

__GET ```api/orders/:orderId```__

__Response:__

```json
{
  "orderId": "uuid",
  "customerId": "uuid",
  "items": [
    {
      "productId": "uuid",
      "product_name": "string",
      "quantity": "int",
      "image_url": "string",
      "unit_price": "float64"
    }
  ],
  "totalPrice": "float64",
  "status": "string",
  "date": "date"
}
```

__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---

## 10. Approve Orders

__PATCH ```api/orders/:orderId```__

__Body:__

```json
{
  "status": "string"
}
```

__Response:__

```json
{
  "orderId": "uuid",
  "customerId": "uuid",
  "items": [
    {
      "productId": "uuid",
      "product_name": "string",
      "quantity": "int",
      "image_url": "string",
      "unit_price": "float64"
    }
  ],
  "totalPrice": "float64",
  "status": "string",
  "date": "date"
}
```



__Status codes:__
- ```200 OK (success)```
- ```400 Bad Request```
- ```404 Not Found```
- ```500 Internal Server Error```

---










