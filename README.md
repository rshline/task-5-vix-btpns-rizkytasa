# Task V - VIX BTPNS

## Deskripsi:

Tim developer bertanggung jawab untuk mengembangkan fitur ini untuk merancang API pada fitur upload, dan menghapus gambar. Beberapa ketentuannya antara lain :
- User dapat menambahkan foto profile
- Sistem dapat mengidentifikasi User ( log in / sign up)
- Hanya user yang telah login / sign up yang dapat melakukan delete / tambah foto profil
- User dapat menghapus gambar yang telah di post
- User yang berbeda tidak dapat menghapus / mengubah foto yang telah di buat oleh user lain

# Tools :
- [Gin Gonic Framework](https://github.com/gin-gonic/gin)
- [Gorm](https://gorm.io/index.html)
- [JWT Go](https://github.com/golang-jwt/jwt)
- [Go Validator](http://github.com/asaskevich/govalidator)

## Endpoint 

### User Endpoint:
1. POST: /users/register
    - ID (primary key, required)
    - Username (required)
    - Email (unique & required) 
    - Password (required & minlength 6)
    - Relasi dengan model Photo (Gunakan constraint cascade)
    - Created At (timestamp)
    - Updated At (timestamp)
2. POST: /users/login
    - Using email & password (required)
3. PUT: /users/:userId (Update User)
4. DELETE: /users/:userId (Delete User)

### Photos Endpoint
1. POST: /photos 
    - ID
    - Title
    - Caption
    - PhotoUrl
    - UserID
    - Relasi dengan model User
2. GET: /photos
3. PUT: /photoId
4. DELETE: /:photoId

# Struktur dokumen / environment:
- app: menampung pembuatan struct dalam kasus ini menggunakan struct User 
untuk keperluan data dan authentication
-  controllers: logic database yaitu models dan query
- database: konfigurasi database serta digunakan untuk menjalankan koneksi database 
dan migration
- helpers: fungsi-fungsi yang dapat digunakan di setiap tempat dalam hal ini jwt, 
bcrypt, headerValue
- middlewares: fungsi yang digunakan untuk proses otentikasi jwt yang digunakan untuk 
proteksi api
- models: models yang digunakan untuk relasi database 
- router: konfigurasi routing / endpoint yang akan digunakan untuk mengakses api
- go mod: manajemen package / dependency berupa library