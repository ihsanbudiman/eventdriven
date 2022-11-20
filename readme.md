# eventdriven

cara menjalankan code

```
docker compose up -d
```

mungkin akan menunggu sedikit lama karena download image, build, dan dependencies. dapat mengakses halaman website dengan masuk ke url: http://localhost:8080
sedangkan untuk url dari websocketnya adalah: ws://localhost:8080/ws

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

kalo ingin menambah data bisa mengakses redis container yang berjalan di docker compose key yang saya pakai untuk mendapatkan user adalah

```
user:<ID>
eg: user:4
```

dengan format json string

```
{
  "id": <ID>,
  "name": <Name>
}

eg:
{
  "id": 4,
  "name": "budiman"
}
```
