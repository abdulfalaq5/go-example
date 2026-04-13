# Go REST API (Clean Architecture)

REST API yang dibangun menggunakan spesifikasi industri standar, mengusung arsitektur *Clean Architecture* tanpa ikatan (*decoupled*) yang kuat dengan framework sehingga sangat rapi, andal, dan mudah dalam penambahan fitur (*scalable*).

---

## 🏗️ Struktur Proyek (Clean Architecture)

Prosedur arsitektur kita memecah logika fungsional dari infrastruktur ke berbagai lapisan:

```text
/
├── cmd/
│   └── api/
│       └── main.go       # Entry point utama aplikasi (koneksi DB, wiring router)
├── docs/                 # Auto-generated dokumen API OpenAPI/Swagger
├── internal/
│   ├── config/           # Membawa / memparsing file .env ke struct global 
│   ├── handler/          # HTTP Layer (Menerima request, validasi input body)
│   ├── middleware/       # Slog logger, JWT Auth Keycloak, Panic recovery
│   ├── model/            # Definisi Entitas / Struct & skema request
│   ├── repository/       # Data Layer (Ekseskusi query raw SQL native - pgxpool)
│   ├── service/          # Business Logic (Logika if-else aplikasi terpusat di sini)
│   └── storage/          # Bootstrap database pool dan file storage (minio/local)
├── migrations/           # Direktori raw query untuk migrasi database PostgreSQL
├── pkg/
│   └── response/         # Format standar JSON pembungkus standar sukses / error
├── .env.example          # Template enviroment variabel konfigurasi
├── docker-compose.yml    # Manajemen service Docker terisolasi
├── Dockerfile            # Alpine based Multi-stage Build blueprint
└── Makefile              # Shortcut alias command untuk run, build, swagger, dsb.
```

---

## 🛠️ Cara Menambahkan Fitur Barumu

Dalam Clean Architecture, alur kontrol berjalan terbalik dari luar ke dalam: `Handler` -> `Service` -> `Repository`. Jika kamu ingin menambahkan fitur/domain baru, misalnya "Product", jalurnya seperti ini:

1. **Buat Model Dulu:** (`internal/model/product.go`)
   Buat definisi struct database dan struct JSON Request inputannya.
   
2. **Buat Repository:** (`internal/repository/product_repository.go`)
   Tulis kueri PostgreSQL langsung ke database. Fungsi me-return hasil struct model.
   
3. **Buat Service:** (`internal/service/product_service.go`)
   Taruh semua business logic (misal: Diskon, pengecekan ketersediaan). Service akan memanggil fungsi di repository.
   
4. **Buat Handler Layer:** (`internal/handler/product_handler.go`)
   Ini adalah jembatan HTTP Gin. Ambil body, validasi, dan lempar ke `Service`. Kembalikan hasil `Service` via `pkg/response`.
   
5. **Daftarkan Router:** (`internal/handler/product_handler.go`)
   Buat router khusus: `func (h *ProductHandler) RegisterRoutes(rg *gin.RouterGroup) { ... }`
   
6. **Hubungkan di Entry Point:** (`cmd/api/main.go`)
   Di file utama, lakukan '*wiring/inject*':
   `repo := repository.New...` 👉 `svc := service.New...(repo)` 👉 `h := handler.New...(svc)` 👉 `h.RegisterRoutes(v1)`

---

## 🗄️ Panduan Migrasi Database

Aplikasi ini menggunakan skema SQL murni untuk struktur tabel tanpa bantuan auto-migration ORM demi performa terbaik. Sangat disarankan menginstall alat standar seperti **[golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)**.

- **Membuat file Migrasi baru:**
  ```bash
  migrate create -ext sql -dir migrations -seq create_products_table
  ```
  *(Command ini akan menciptakan file .up.sql dan .down.sql otomatis di folder `/migrations`)*

- **Menjalankan/Menyuntik (UP) Schema ke PostgreSQL:**
  ```bash
  migrate -database "postgres://postgres:secret@localhost:5432/main_db?sslmode=disable" -path migrations up
  ```

- **Rollback (DOWN) Skema 1 langkah mundur:**
  ```bash
  migrate -database "postgres://postgres:secret@localhost:5432/main_db?sslmode=disable" -path migrations down 1
  ```

---

## 🚀 Setup & Cara Menjalankan

### Cara 1: Menjalankan Secara Asli (*Local Tanpa Docker*)

**Syarat Tambahan:** Golang 1.22+, dan memiliki PostgreSQL lokal di Local machine (Port standar 5432 dsb), dan opsional Minio/Keycloak.

1. **Clone dan Setup `.env`**
   ```bash
   cp .env.example .env
   ```
   **PENTING**: Buka file `.env`. Di sini isi host-host `DB_MAIN_DSN`, dan lainnya dengan `localhost`. *(Gunakan `STORAGE_TYPE=local` untuk mengunggah berkas ke dalam direktori `./uploads` tanpa pelengkap MinIO server).*
   
2. **Setup Dependencies**
   ```bash
   go mod tidy
   ```
   
3. **Migrasikan Tabel (`users` dsb)**
   Gunakan instruksi CLI `migrate` seperti panduan migrasi di atas ke database kamu.

4. **Kompilasi / Run Server**
   Anda bisa menggunakan `make` yang ada:
   ```bash
   make run    # (Ini otomatis menjalankan go run ./cmd/api/main.go)
   ```

📝 *Catatan:* Perbarui Swagger *Documentation* sebelum rilis apabila terdapat *Handler* baru:
```bash
make swag
```

---

### Cara 2: Menjalankan Lewat Container (*Dengan Docker*)

Docker sangat berguna saat peluncuran Produksi (VPS/Server) maupun testing bersih tanpa harus mencemari Environment komputer host anda dengan software Go compiler. 

**Perhatian Khusus (`.env` networking Docker):**
Karena API ini dijalankan *di dalam kurungan container*, `localhost` bagi si *container* bermaksud *container itu sendiri*. Oleh karena itu, jika Aplikasi ini perlu menghubungi Postgres Database atau Storage Minio yang hidup sebagai sistem program komputer mu di luar docker, perbarui host koneksi kalian di `.env`:
`localhost` ➡️ `host.docker.internal`

*Contoh di `.env`:* `DB_MAIN_DSN=postgres://postgres:secret@host.docker.internal:5432/main_db...`

1. **Build image & Jalankan (detached)**
   ```bash
   docker-compose up -d --build
   ```
2. **Cek Log Container Cepat**
   ```bash
   docker logs -f go_example_api
   ```
3. **Menghentikan Server**
   ```bash
   docker-compose down
   ```

---

## 📖 Swagger API Documentation

Server siap digunakan saat konsol menampilkan `🚀 Server running on port...`. 
Buka *browser* anda pada tautan berikut untuk membaca interaktif *interface* API UI:

- **Alamat Swagger:** `http://localhost:8080/swagger/index.html` *(Atau port lain yang anda pakai)*

**Fitur Tersedia di Dokumen:**
- Endpoint `GET /health` (Status monitoring aplikasi)
- Endpoint `GET, POST, PUT, DELETE` dari ranah `/api/v1/users` (Memerintahkan proteksi *Auth Keycloak JWT*)
- Endpoint `POST /api/v1/upload` (Form-data berkas Multipart)
