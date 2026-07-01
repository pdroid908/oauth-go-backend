# 🚀 OAuth Golang - Progress

## ✅ Progress

### Database
- [x] PostgreSQL (Supabase)
- [x] Koneksi Database
- [x] Tabel `users`
- [x] Tabel `email_verifications`
- [x] Tabel `password_resets`
- [x] Tabel `refresh_tokens`

### Backend
- [x] Gin Framework
- [x] Load `.env`
- [x] Package `database`
- [x] Package `security`
- [x] Package `mail` (SMTP Gmail)
- [x] Input Validation
- [x] Input Sanitization
- [x] Password Hash (`bcrypt`)
- [x] Register User
- [x] Generate Verification Token
- [x] Simpan Verification Token ke Database
- [x] Kirim Email Verifikasi (SMTP)
- [x] Verifikasi Email
- [x] Login Email & Password
- [x] JWT Authentication
- [x] HttpOnly Cookie
- [x] Google OAuth 2.0 Login
- [x] OAuth State (CSRF Protection)
- [x] Ambil Data User Google
- [x] Generate JWT Setelah Login Google

---

# 📦 Struktur Package

## `database`
Mengelola koneksi PostgreSQL (Supabase) menggunakan `pgxpool` dan menyediakan objek database global yang digunakan seluruh handler.

## `security`
Berisi fungsi keamanan yang dapat dipanggil di seluruh project, seperti sanitasi input, validasi, hashing password, dan helper keamanan lainnya.

## `mail`
Mengirim email menggunakan SMTP Gmail, digunakan untuk verifikasi email dan nantinya reset password.

## `handler`
Berisi seluruh endpoint Authentication seperti Register, Login, Verify Email, Google OAuth, dan endpoint auth lainnya.

---

# 📌 Register Flow

```text
POST /register
        │
        ▼
Validasi Input
        │
        ▼
Sanitasi Input
        │
        ▼
Hash Password (bcrypt)
        │
        ▼
Insert User
(email_verified = false)
        │
        ▼
Generate Verification Token
        │
        ▼
Simpan Token ke email_verifications
        │
        ▼
Kirim Email Verifikasi
        │
        ▼
User Klik Link
/verify-email?token=...
        │
        ▼
VerifyEmail
        │
        ▼
Update email_verified = true
        │
        ▼
Hapus Verification Token
        │
        ▼
Registrasi Selesai
```

---

# 📌 Login Email Flow

```text
POST /login
        │
        ▼
Validasi Input
        │
        ▼
Cari User Berdasarkan Email
        │
        ▼
Cek Email Sudah Verified
        │
        ▼
Bandingkan Password (bcrypt)
        │
        ▼
Generate JWT
        │
        ▼
Set HttpOnly Cookie
        │
        ▼
Return JWT
        │
        ▼
Login Berhasil
```

---

# 📌 Google OAuth Flow

```text
GET /auth/google
        │
        ▼
Generate OAuth State
        │
        ▼
Simpan State ke Cookie
        │
        ▼
Redirect ke Google Login
        │
        ▼
User Login Google
        │
        ▼
Google Redirect
/auth/google/callback
        │
        ▼
Validasi State (CSRF)
        │
        ▼
Exchange Authorization Code
        │
        ▼
Ambil Profile Google
        │
        ▼
Insert / Update User Database
        │
        ▼
Generate JWT
        │
        ▼
Set HttpOnly Cookie
        │
        ▼
Login Google Berhasil
```

---

# 🔐 Security

- Password menggunakan `bcrypt`.
- JWT menggunakan `HS256`.
- JWT memiliki waktu kedaluwarsa.
- Token JWT disimpan pada `HttpOnly Cookie`.
- Input divalidasi sebelum diproses.
- Input disanitasi untuk mengurangi karakter berbahaya.
- Email wajib diverifikasi sebelum login menggunakan email & password.
- Google OAuth menggunakan `state` untuk mencegah serangan CSRF.
- Secret dan konfigurasi disimpan pada file `.env`.

---

# 📂 Endpoint

| Method | Endpoint | Fungsi |
|---------|----------|--------|
| POST | `/register` | Registrasi akun |
| GET | `/verify-email` | Verifikasi email |
| POST | `/login` | Login email & password |
| GET | `/auth/google` | Login menggunakan Google |
| GET | `/auth/google/callback` | Callback Google OAuth |

---

# 📈 Status Project

## ✅ Sudah Selesai

- Sistem Registrasi
- Verifikasi Email
- Login Email
- Login Google OAuth
- JWT Authentication
- HttpOnly Cookie
- SMTP Gmail
- PostgreSQL (Supabase)
- Input Validation
- Input Sanitization
- Password Hashing
- CSRF Protection pada Google OAuth

---
struk login go
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
struk registrasi 
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
lupa passsword
func ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	var userID string
	err := database.DB.QueryRow(c, "SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)

	if err != nil {
		// Tetap berikan response sukses agar tidak membocorkan email yang terdaftar
		c.JSON(http.StatusOK, gin.H{"message": "If the email exists, a reset link has been sent."})
		return
	}

	// Buat Token
	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)
	expiresAt := time.Now().Add(time.Hour)

	// Hapus token lama jika ada agar tidak duplikat
	_, _ = database.DB.Exec(c, "DELETE FROM password_resets WHERE user_id = $1", userID)

	// Insert token baru (4 kolom: id, user_id, token, expires_at)
	_, err = database.DB.Exec(
		c,
		`INSERT INTO password_resets (id, user_id, token, expires_at) 
		 VALUES (gen_random_uuid(), $1, $2, $3)`,
		userID,
		token,
		expiresAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reset token"})
		return
	}

	// Kirim email
	err = mail.SendResetPasswordEmail(req.Email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reset link sent"})
}
 rest password
 func ResetPassword(c *gin.Context) {
    var req struct {
        Token       string `json:"token" binding:"required"`
        NewPassword string `json:"new_password" binding:"required,min=6"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // 1. Cari data di password_resets
    // PERBAIKAN: userID diubah dari int menjadi string (karena UUID di DB adalah text/uuid)
    var userID string 
    var expiresAt time.Time
    err := database.DB.QueryRow(c, 
        "SELECT user_id, expires_at FROM password_resets WHERE token = $1", req.Token).Scan(&userID, &expiresAt)
    
    if err != nil {
        // Log error untuk debug di terminal
        fmt.Println("Query token error:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired token"})
        return
    }

    // 2. Cek expired
    if time.Now().After(expiresAt) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Token expired"})
        return
    }

    // 3. Hash password baru
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Hashing failed"})
        return
    }

    // 4. Transaksi: Update password dan Hapus token
    tx, err := database.DB.Begin(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
        return
    }
    defer tx.Rollback(c)

    // Update password
    _, err = tx.Exec(c, "UPDATE users SET password_hash = $1 WHERE id = $2", string(hashedPassword), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Update password failed"})
        return
    }

    // Hapus token setelah digunakan
    _, err = tx.Exec(c, "DELETE FROM password_resets WHERE token = $1", req.Token)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cleanup failed"})
        return
    }

    // Komit transaksi
    err = tx.Commit(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Commit failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}