# Go API Gateway

Proyek ini adalah contoh implementasi API Gateway menggunakan bahasa pemrograman Go (Golang) dengan framework Gin Gonic. API Gateway ini berfungsi sebagai pintu masuk tunggal untuk berbagai layanan backend, menyediakan fitur-fitur penting seperti routing, otentikasi JWT, rate limiting, dan logging.

## Fitur Utama

* **Routing Dinamis:** Meneruskan request ke layanan backend yang sesuai berdasarkan path URL.
* **Reverse Proxy:** Menggunakan `net/http/httputil` untuk meneruskan request.
* **Otentikasi JWT:** Mengamankan endpoint menggunakan JSON Web Tokens. Termasuk endpoint `/auth/login` untuk menghasilkan token.
* **Rate Limiting:** Pembatasan jumlah request per IP untuk mencegah penyalahgunaan (menggunakan `golang.org/x/time/rate`).
* **Middleware:**
    * Logging request HTTP.
    * Validasi token JWT.
    * Penanganan CORS.
    * Rate Limiting.
* **Manajemen Konfigurasi:** Menggunakan Viper untuk memuat konfigurasi dari file (`config.yaml`) dan variabel environment.
* **Struktur Proyek Modular:** Kode diorganisir ke dalam package-package untuk kemudahan pengelolaan dan skalabilitas.

## Struktur Proyek
api-gateway-go/
├── cmd/
│   └── api-gateway/
│       └── main.go                # Titik masuk aplikasi
├── pkg/
│   ├── config/
│   │   ├── config.go              # Logika untuk memuat konfigurasi
│   │   └── config.yaml            # File konfigurasi default
│   ├── handlers/
│   │   ├── auth_handler.go        # Handler untuk otentikasi (login, generate token)
│   │   ├── proxy_handler.go       # Handler untuk meneruskan request (reverse proxy)
│   │   └── health_handler.go      # Handler untuk health check
│   ├── middleware/
│   │   ├── auth_middleware.go     # Middleware untuk validasi token JWT
│   │   ├── logging_middleware.go  # Middleware untuk logging request
│   │   └── ratelimit_middleware.go# Middleware untuk rate limiting
│   ├── routes/
│   │   └── router.go              # Definisi semua rute API Gateway
│   └── services/                  # (Opsional) Logika untuk interaksi dengan service backend
├── dummy-services/                # Contoh layanan backend sederhana untuk pengujian
│   ├── product-service/
│   │   └── main.go
│   └── user-service/
│       └── main.go
├── go.mod                         # File dependensi Go Modules
├── go.sum
└── README.md                      # File ini

## Prasyarat

* **Go:** Versi 1.18 atau lebih baru. [Instal Go](https://golang.org/doc/install)

## Setup dan Instalasi

1.  **Clone Proyek (atau buat struktur secara manual):**
    Jika ini adalah repositori, Anda akan meng-clone-nya. Jika Anda mengikuti panduan, Anda telah membuat struktur ini.

2.  **Inisialisasi Go Modules (jika belum):**
    ```bash
    go mod init api-gateway-go # Ganti dengan nama modul Anda jika berbeda
    ```

3.  **Instal Dependensi:**
    Navigate ke direktori root proyek dan jalankan:
    ```bash
    go get [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
    go get [github.com/spf13/viper](https://github.com/spf13/viper)
    go get [github.com/gin-contrib/cors](https://github.com/gin-contrib/cors)
    go get golang.org/x/time/rate
    go get [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt/v5)
    # Atau jalankan `go mod tidy` untuk mengambil semua dependensi yang terdaftar di kode
    # go mod tidy
    ```

4.  **Konfigurasi:**
    * Salin atau buat file `config.yaml` di direktori root proyek (`api-gateway-go/`). Contoh `config.yaml` disediakan di bawah.
    * Sesuaikan nilai-nilai dalam `config.yaml`, terutama `AUTH_SECRET` dan `SERVICE_ENDPOINTS`.
    * **PENTING:** `AUTH_SECRET` harus berupa string yang kuat dan rahasia. Jangan gunakan nilai default di lingkungan produksi.

    **Contoh `config.yaml`:**
    ```yaml
    SERVER_PORT: "8080"
    APP_ENV: "development" # "production" atau "development"
    AUTH_SECRET: "ganti-dengan-secret-key-anda-yang-sangat-aman-dan-panjang"

    SERVICE_ENDPOINTS:
      user_service: "http://localhost:8081/api/users" # Sesuaikan dengan URL backend service Anda
      product_service: "http://localhost:8082/api/products"
      # order_service: "http://localhost:8083/api/orders" # Contoh service lain

    RATE_LIMIT:
      ENABLED: true
      REQUESTS: 100    # Jumlah request
      WINDOW_SEC: 60   # Jendela waktu dalam detik (misal, 100 request per menit per IP)
    ```

## Menjalankan Aplikasi

1.  **Jalankan Layanan Backend (Dummy Services):**
    API Gateway ini dirancang untuk meneruskan request ke layanan backend. Untuk pengujian, Anda bisa menjalankan layanan dummy yang disediakan. Buka terminal terpisah untuk setiap layanan:

    * **User Service:**
        ```bash
        cd dummy-services/user-service
        go run main.go
        # Akan berjalan di http://localhost:8081 secara default
        ```
    * **Product Service:**
        ```bash
        cd dummy-services/product-service
        go run main.go
        # Akan berjalan di http://localhost:8082 secara default
        ```
    Pastikan URL di `config.yaml` (bagian `SERVICE_ENDPOINTS`) sesuai dengan alamat layanan backend Anda.

2.  **Jalankan API Gateway:**
    Buka terminal baru, navigasi ke direktori root proyek API Gateway, dan jalankan:
    ```bash
    go run cmd/api-gateway/main.go
    ```
    API Gateway akan berjalan di port yang ditentukan dalam `config.yaml` (default: `8080`).

## Endpoint API

Berikut adalah beberapa endpoint utama yang tersedia:

* **POST** `/auth/login`
    * Digunakan untuk otentikasi dan mendapatkan token JWT.
    * **Request Body:**
        ```json
        {
          "username": "user123",
          "password": "password123"
        }
        ```
        (Kredensial `user123`/`password123` atau `admin`/`adminpass` di-hardcode di `pkg/handlers/auth_handler.go` untuk contoh ini. Ganti dengan logika validasi database Anda.)
    * **Response Sukses (200 OK):**
        ```json
        {
          "token": "your.jwt.token",
          "expires_at": "...",
          "user_id": "...",
          "username": "..."
        }
        ```

* **GET** `/api/public/health`
    * Endpoint publik untuk memeriksa status kesehatan API Gateway.
    * **Response Sukses (200 OK):**
        ```json
        {
          "status": "UP",
          "message": "API Gateway is running smoothly!"
        }
        ```

* **ANY** `/api/v1/users/*proxyPath`
    * Meneruskan semua request (GET, POST, PUT, DELETE, dll.) ke `user_service` yang dikonfigurasi.
    * Membutuhkan token JWT di header `Authorization: Bearer <token>`.
    * Contoh: `GET /api/v1/users/profile` akan diteruskan ke `http://<user_service_url>/profile`.

* **GET** `/api/v1/products/*proxyPath`
    * Meneruskan request GET ke `product_service` yang dikonfigurasi.
    * Endpoint ini publik (tidak memerlukan otentikasi).
    * Contoh: `GET /api/v1/products/item/123` akan diteruskan ke `http://<product_service_url>/item/123`.

* **POST, PUT, DELETE** `/api/v1/products/*proxyPath`
    * Meneruskan request POST, PUT, DELETE ke `product_service` yang dikonfigurasi.
    * Membutuhkan token JWT di header `Authorization: Bearer <token>`.

## Contoh Pengujian dengan cURL

1.  **Login untuk Mendapatkan Token:**
    ```bash
    curl -X POST \
      http://localhost:8080/auth/login \
      -H 'Content-Type: application/json' \
      -d '{
        "username": "user123",
        "password": "password123"
      }'
    ```
    Simpan `token` yang diterima dari respons.

2.  **Mengakses Endpoint yang Dilindungi (User Service):**
    Ganti `YOUR_JWT_TOKEN` dengan token yang Anda dapatkan.
    ```bash
    curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/v1/users/profile
    ```

3.  **Mengakses Endpoint Publik (Product Service - GET):**
    ```bash
    curl http://localhost:8080/api/v1/products/
    ```

4.  **Mengakses Endpoint yang Dilindungi (Product Service - POST):**
    Ganti `YOUR_JWT_TOKEN` dengan token yang Anda dapatkan.
    ```bash
    curl -X POST \
      -H "Authorization: Bearer YOUR_JWT_TOKEN" \
      -H "Content-Type: application/json" \
      -d '{"name":"New Product", "price":99.99}' \
      http://localhost:8080/api/v1/products/
    ```

## Konfigurasi Detail

Konfigurasi utama diatur dalam file `config.yaml` atau melalui variabel environment yang sesuai (Viper mendukung ini). Variabel environment akan menimpa nilai dalam file config.

* `SERVER_PORT`: Port tempat API Gateway berjalan.
* `APP_ENV`: Lingkungan aplikasi (`development` atau `production`). Mempengaruhi mode Gin.
* `AUTH_SECRET`: Kunci rahasia untuk menandatangani dan memverifikasi token JWT. **SANGAT PENTING UNTUK DIJAGA KERAHASIAANNYA DAN DIGANTI DARI NILAI DEFAULT.**
* `SERVICE_ENDPOINTS`: Peta URL untuk layanan backend.
    * `user_service`: URL lengkap ke root endpoint layanan pengguna.
    * `product_service`: URL lengkap ke root endpoint layanan produk.
* `RATE_LIMIT`: Pengaturan untuk rate limiting.
    * `ENABLED`: `true` atau `false`.
    * `REQUESTS`: Jumlah maksimum request.
    * `WINDOW_SEC`: Jendela waktu (dalam detik) untuk batas request.

## Teknologi yang Digunakan

* **Go (Golang):** Bahasa pemrograman.
* **Gin Gonic:** Framework web HTTP berperforma tinggi.
* **Viper:** Library untuk manajemen konfigurasi.
* **JWT (github.com/golang-jwt/jwt/v5):** Implementasi JSON Web Token.
* **golang.org/x/time/rate:** Untuk implementasi rate limiting.

## Potensi Pengembangan Lebih Lanjut

* **Refresh Tokens:** Implementasi mekanisme refresh token untuk sesi pengguna yang lebih lama.
* **Validasi Kredensial Database:** Mengganti validasi login dummy dengan koneksi ke database pengguna.
* **Service Discovery:** Integrasi dengan alat service discovery seperti Consul atau etcd.
* **Caching:** Menambahkan lapisan caching untuk response yang sering diakses.
* **Request/Response Transformation:** Kemampuan untuk memodifikasi request atau response saat melewati gateway.
* **Observability:** Integrasi dengan Prometheus untuk metrics, Jaeger/OpenTelemetry untuk tracing.
* **Pengujian (Unit & Integrasi):** Menulis test suite yang komprehensif.
* **Granular Authorization (Roles/Permissions):** Memperluas klaim JWT untuk menyertakan peran atau izin, dan memvalidasinya di middleware.

---