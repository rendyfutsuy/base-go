# 📘 Panduan Implementasi GORM - Module Auth

## ✅ Yang Sudah Selesai:
1. ✅ Setup GORM connection di `database/database.go`
2. ✅ Update models dengan GORM tags (User, JWTToken, PasswordHistory, ResetPasswordToken)
3. ✅ Update `main.go` untuk initialize GORM DB

## 🔧 Yang Perlu Dilakukan:

### 1. Install Dependencies
```bash
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go mod tidy
```

### 2. Update Router (`router/router.go`)

**Ubah signature function:**
```go
// BEFORE:
func InitializedRouter(db *sql.DB, timeoutContext time.Duration, v *validator.Validate, nrApp *newrelic.Application) *echo.Echo {

// AFTER:
func InitializedRouter(db *sql.DB, gormDB *gorm.DB, timeoutContext time.Duration, v *validator.Validate, nrApp *newrelic.Application) *echo.Echo {
```

**Update auth repository initialization:**
```go
// BEFORE:
emailService := services.NewEmailService()
authRepo := _authRepo.NewAuthRepository(db, emailService)

// AFTER:
emailService := services.NewEmailService()
authRepo := _authRepo.NewAuthRepository(gormDB, emailService)  // Pass gormDB instead of db
```

### 3. Update Auth Repository Struct

**File: `modules/auth/repository/auth_repository.go`**

```go
// BEFORE:
import (
	"context"
	"database/sql"
	// ... other imports
)

type authRepository struct {
	Conn         *sql.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
}

func NewAuthRepository(Conn *sql.DB, EmailService *services.EmailService) auth.Repository {
	// ...
}

// AFTER:
import (
	"context"
	// ... other imports (remove database/sql)
	"gorm.io/gorm"
)

type authRepository struct {
	DB           *gorm.DB  // Changed from Conn *sql.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
}

func NewAuthRepository(DB *gorm.DB, EmailService *services.EmailService) auth.Repository {
	QueueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	})

	return &authRepository{
		DB,           // Changed
		EmailService,
		QueueClient,
	}
}
```

## 📝 Pattern Konversi Query SQL → GORM

### Pattern 1: SELECT Single Row dengan WHERE

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) FindByEmailOrUsername(ctx context.Context, login string) (user models.User, err error) {
	query := `SELECT id, email, password FROM users WHERE (email = $1 OR username = $1) AND deleted_at IS NULL`
	err = repo.Conn.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) FindByEmailOrUsername(ctx context.Context, login string) (user models.User, err error) {
	err = repo.DB.WithContext(ctx).
		Select("id, email, password").
		Where("(email = ? OR username = ?) AND deleted_at IS NULL", login, login).
		First(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}
```

### Pattern 2: SELECT dengan Custom Scan

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	var storedHash string
	query := `SELECT password FROM users WHERE id = $1 AND deleted_at IS NULL`
	err := repo.Conn.QueryRowContext(ctx, query, userId).Scan(&storedHash)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("user not found")
		}
		return false, err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password))
	if err != nil {
		// Update counter
		_, _ = repo.Conn.ExecContext(ctx, `UPDATE users SET counter = counter + 1, updated_at = NOW() WHERE id = $1`, userId)
		return false, errors.New("invalid password")
	}
	
	// Reset counter on success
	_, _ = repo.Conn.ExecContext(ctx, `UPDATE users SET counter = 0, updated_at = NOW() WHERE id = $1`, userId)
	return true, nil
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) AssertPasswordRight(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	var user models.User
	err := repo.DB.WithContext(ctx).
		Select("password").
		Where("id = ? AND deleted_at IS NULL", userId).
		First(&user).Error
	
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found")
		}
		return false, err
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Update counter (increment)
		repo.DB.WithContext(ctx).
			Model(&models.User{}).
			Where("id = ?", userId).
			Updates(map[string]interface{}{
				"counter":    gorm.Expr("counter + ?", 1),
				"updated_at": time.Now(),
			})
		return false, errors.New("invalid password")
	}
	
	// Reset counter on success
	repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", userId).
		Updates(map[string]interface{}{
			"counter":    0,
			"updated_at": time.Now(),
		})
	return true, nil
}
```

### Pattern 3: SELECT Multiple Rows dengan Query

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) AssertPasswordNeverUsesByUser(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	query := `SELECT hashed_password FROM password_histories WHERE user_id = $1 ORDER BY created_at DESC LIMIT 5`
	rows, err := repo.Conn.QueryContext(ctx, query, userId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var hashedPassword string
		if err := rows.Scan(&hashedPassword); err != nil {
			return false, err
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err == nil {
			return false, nil // Password has been used before
		}
	}
	
	return true, nil // Password never used
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) AssertPasswordNeverUsesByUser(ctx context.Context, password string, userId uuid.UUID) (bool, error) {
	var histories []models.PasswordHistory
	err := repo.DB.WithContext(ctx).
		Select("hashed_password").
		Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(5).
		Find(&histories).Error
		
	if err != nil {
		return false, err
	}

	for _, history := range histories {
		if err := bcrypt.CompareHashAndPassword([]byte(history.HashedPassword), []byte(password)); err == nil {
			return false, nil // Password has been used before
		}
	}
	
	return true, nil // Password never used
}
```

### Pattern 4: INSERT

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) AddUserAccessToken(ctx context.Context, token string, userId uuid.UUID) error {
	query := `INSERT INTO jwt_tokens (user_id, access_token, created_at, updated_at) 
	          VALUES ($1, $2, NOW(), NOW())`
	_, err := repo.Conn.ExecContext(ctx, query, userId, token)
	return err
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) AddUserAccessToken(ctx context.Context, token string, userId uuid.UUID) error {
	jwtToken := models.JWTToken{
		UserId:      userId,
		AccessToken: token,
		CreatedAt:   time.Now(),
		UpdatedAt:   &[]time.Time{time.Now()}[0],
	}
	
	return repo.DB.WithContext(ctx).Create(&jwtToken).Error
}
```

### Pattern 5: UPDATE

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) UpdatePasswordById(ctx context.Context, id uuid.UUID, password string) error {
	query := `UPDATE users SET password = $1, updated_at = NOW(), password_expired_at = NOW() + INTERVAL '3 months' 
	          WHERE id = $2 AND deleted_at IS NULL`
	_, err := repo.Conn.ExecContext(ctx, query, password, id)
	return err
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) UpdatePasswordById(ctx context.Context, id uuid.UUID, password string) error {
	updates := map[string]interface{}{
		"password":            password,
		"updated_at":          time.Now(),
		"password_expired_at": time.Now().AddDate(0, 3, 0), // +3 months
	}
	
	return repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}
```

### Pattern 6: DELETE

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) DestroyToken(ctx context.Context, accessToken string) error {
	query := `DELETE FROM jwt_tokens WHERE access_token = $1`
	_, err := repo.Conn.ExecContext(ctx, query, accessToken)
	return err
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) DestroyToken(ctx context.Context, accessToken string) error {
	return repo.DB.WithContext(ctx).
		Where("access_token = ?", accessToken).
		Delete(&models.JWTToken{}).Error
}
```

### Pattern 7: UPDATE dengan Profile (Multiple Fields)

**BEFORE (Raw SQL):**
```go
func (repo *authRepository) UpdateProfileById(ctx context.Context, id uuid.UUID, profile dto.UpdateProfileDto) error {
	query := `UPDATE users SET full_name = $1, email = $2, updated_at = NOW() 
	          WHERE id = $3 AND deleted_at IS NULL`
	_, err := repo.Conn.ExecContext(ctx, query, profile.FullName, profile.Email, id)
	return err
}
```

**AFTER (GORM):**
```go
func (repo *authRepository) UpdateProfileById(ctx context.Context, id uuid.UUID, profile dto.UpdateProfileDto) error {
	updates := map[string]interface{}{
		"full_name":  profile.FullName,
		"email":      profile.Email,
		"updated_at": time.Now(),
	}
	
	return repo.DB.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}
```

## ⚠️ Penting - Catatan Khusus:

### 1. **Context Handling**
Selalu gunakan `WithContext(ctx)` untuk mendukung timeout dan cancellation:
```go
repo.DB.WithContext(ctx).Find(&users)
```

### 2. **Error Handling**
- `sql.ErrNoRows` → `gorm.ErrRecordNotFound`
- Check dengan: `errors.Is(err, gorm.ErrRecordNotFound)`

### 3. **Soft Delete**
Model User sudah menggunakan `gorm.DeletedAt`, GORM otomatis handle soft delete:
```go
// Auto exclude deleted_at IS NOT NULL
repo.DB.Find(&users)

// Include soft deleted
repo.DB.Unscoped().Find(&users)
```

### 4. **Transaction**
Jika perlu transaction (contoh: reset password dengan multiple operations):
```go
err := repo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
	// Operation 1
	if err := tx.Model(&models.User{}).Where("id = ?", userId).Update("password", newPass).Error; err != nil {
		return err
	}
	
	// Operation 2
	if err := tx.Create(&passwordHistory).Error; err != nil {
		return err
	}
	
	return nil
})
```

## 📋 Checklist Files yang Perlu Diupdate:

### Auth Module:
- [ ] `modules/auth/repository/auth_repository.go` (~15 methods)
- [ ] `modules/auth/repository/reset_password_repository.go` (~6 methods)

### User Management Module (jika mau):
- [ ] `modules/user_management/repository/user_repository.go` (~20 methods)
- [ ] `modules/user_management/repository/user_password_repository.go` (~2 methods)

### Role Management Module (jika mau):
- [ ] `modules/role_management/repository/role_repository.go` (~15 methods)
- [ ] `modules/role_management/repository/permission_repository.go` (~6 methods)
- [ ] `modules/role_management/repository/permission_group_repository.go` (~6 methods)
- [ ] `modules/role_management/repository/role_assigment_repository.go` (~7 methods)
- [ ] `modules/role_management/repository/permission_group_assigment_repository.go` (~1 method)

## 🔍 Testing
Setelah implementasi, test dengan:
1. Run `go mod tidy`
2. Check linter: `golangci-lint run` atau built-in linter di IDE
3. Test API endpoints untuk memastikan semua query berfungsi
4. Check logs GORM untuk verify SQL queries yang dijalankan

## 💡 Tips:
1. Convert method per method, jangan sekaligus
2. Test setelah convert beberapa method
3. Gunakan GORM debug mode saat development: `repo.DB.Debug().Find(&users)`
4. Lihat generated SQL untuk debugging: GORM akan print SQL di logs

## 🆘 Jika Ada Error:
- `undefined: gorm` → Run `go get -u gorm.io/gorm`
- `cannot use gormDB (variable of type *gorm.DB)` → Update function signatures
- `field not found` → Check GORM tags di model
- SQL error → Use `.Debug()` untuk lihat generated SQL

