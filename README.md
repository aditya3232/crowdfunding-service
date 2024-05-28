## Architecture 

Ini adalah Layered Architecture kalau dilihat secara struktur partisinya

![Layered Architecture](architecture.png)

## Berikut adalah penjelasan dari setiap layer yang ada di dalam Layered Architecture ini:

1. Sistem eksternal melakukan permintaan (HTTP, gRPC, Messaging, dll) ke delivery, contoh permintaan dari sistem eksternal adalah request HTTP get dari client, request gRPC dari client, baca pesan dari message broker, atau baca cache dari redis
2. Delivery akan mengakses model request untuk setiap request yang masuk, dan memanggil model response, untuk mengembalikan response data dari use case ke sistem eksternal
3. Delivery akan memanggil use case, didalam use case berisi bisnis logic (transaksi, validasi, dll) 
4. Use case membuat instance baru dari entity dan mengisi data entity dari request model
5. Use case memanggil repository, untuk menyimpan data instance entity yg telah dibuat ke database
6. Repository menggunakan data entity untuk melakukan operasi database
7. Repository melakukan operasi basis data ke database
8. Use case memberikan response yang sesuai berdasarkan model response
9. Lalu jika ada data yang perlu dikirim ke sistem eksternal, maka use case akan memanggil gateway, dan gateway akan mengirimkan data ke sistem eksternal
10. Gateway juga akan mengaksesl model untuk data yang dikirim ke sistem eksternal
11. Gateway akan melakukan permintaan mengirimkan data ke sistem eksternal

## Tech Stack

- Golang : https://github.com/golang/go
- MySQL : https://github.com/mysql/mysql-server
- Minio : https://github.com/minio/minio

## Framework & Library

- GoFiber (HTTP Framework) : https://github.com/gofiber/fiber
- GoFiber CORS (CORS Middleware) : https://docs.gofiber.io/api/middleware/cors
- GORM (ORM) : https://github.com/go-gorm/gorm
- Viper (Configuration) : https://github.com/spf13/viper
- Go Migrate (Database Migration) : https://github.com/golang-migrate/migrate
- Go Playground Validator (Validation) : https://github.com/go-playground/validator
- Logrus (Logger) : https://github.com/sirupsen/logrus
- Go Midtrans : https://github.com/veritrans/go-midtrans
- Go Oauth2 : https://pkg.go.dev/golang.org/x/oauth2#section-readme
- Go Minio : https://pkg.go.dev/github.com/minio/minio-go/v7#section-readme

# Erd

![Erd](erd.png)

## Configuration

Semua konfigurasi add di file `config.json`.

## API Spec

Semua api spec ada di folder `api-spec`.

## Database Migration

Semua database migration ada di folder `db/migrations`.

### Install Migrate Mysql
```shell
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Create Migration

```shell
migrate create -ext sql -dir db/migrations create_table_xxx
```

### Run Migration

```shell
migrate -database "mysql://root:root_password@tcp(localhost:3306)/crowdfunding_db?charset=utf8mb4&parseTime=True&loc=Local" -path db/migrations up
```

## Run Application

### Install Dependency

```bash
go mod tidy
```

### Run unit test

```bash
go test -v ./test/
```

### Run web server

```bash
go run cmd/web/main.go
```
