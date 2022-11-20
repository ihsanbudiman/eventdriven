# eventdriven


cara menjalankan code
```
docker compose up -d
```

mungkin akan menunggu sedikit lama karena download image, build, dan dependencies.
dapet mengakses dengan masuk ke url: localhost:8080

data yang tersedia
```
input: 1
{
  "user": {
    "id": 1,
    "name": "ihsan"
  },
  "message": "Success",
  "success": true
}

input: 2
{
  "user": {
    "id": 2,
    "name": "Tono"
  },
  "message": "Success",
  "success": true
}

input: 3
{
  "user": {
    "id": 3,
    "name": "Yadi"
  },
  "message": "Success",
  "success": true
}

jika tidak ada datanya
{
  "user": {
    "id": 0,
    "name": ""
  },
  "message": "Data Not Found",
  "success": false
}
```
