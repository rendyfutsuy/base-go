# 🔄 Status Implementasi GORM

## ✅ Completed:

### 1. Database Setup
- ✅ `database/database.go` - Added `ConnectToGORM()` function
- ✅ Connection pool configuration (max open conns, max idle conns, conn lifetime)

### 2. Models Updated dengan GORM Tags
- ✅ `models/user.go` - Added GORM tags + `TableName()` method
- ✅ `models/jwt_token.go` - Added GORM tags + `TableName()` method
- ✅ `models/password_history.go` - Added GORM tags + `TableName()` method  
- ✅ `models/reset_password_token.go` - Added GORM tags + `TableName()` method

### 3. Application Setup  
- ✅ `main.go` - Initialize both SQL DB and GORM DB
- ✅ `main.go` - Pass both DBs to router
- ✅ `router/router.go` - Updated signature to accept `gormDB *gorm.DB`
- ✅ `router/router.go` - Pass GORM DB to auth repository

## 📋 Next Steps (Manual):

### 1. Install Dependencies
```bash
cd /Users/rendyanggara/Projects/base-go
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
go mod tidy
```

### 2. Update Auth Repository Structure

**File: `modules/auth/repository/auth_repository.go`**

Lakukan perubahan berikut di bagian atas file:

```go
// REMOVE:
import (
	"context"
	"database/sql"  // <-- REMOVE THIS
	...
)

// ADD:
import (
	"context"
	"errors"  // Keep this
	...
	"gorm.io/gorm"  // <-- ADD THIS
)

// CHANGE STRUCT:
type authRepository struct {
	DB           *gorm.DB  // Changed from: Conn *sql.DB
	EmailService *services.EmailService
	QueueClient  *asynq.Client
}

// CHANGE CONSTRUCTOR:
func NewAuthRepository(DB *gorm.DB, EmailService *services.EmailService) auth.Repository {
	QueueClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr:     utils.ConfigVars.String("redis.address"),
		Password: utils.ConfigVars.String("redis.password"),
		DB:       utils.ConfigVars.Int("redis.db"),
	})

	return &authRepository{
		DB,  // Changed from Conn
		EmailService,
		QueueClient,
	}
}
```

### 3. Convert All Methods - Mapping

Berikut daftar lengkap method yang perlu diconvert:

#### `modules/auth/repository/auth_repository.go`:

| No | Method Name | Pattern | Status |
|----|-------------|---------|--------|
| 1  | `FindByEmailOrUsername` | SELECT Single + WHERE | ⏳ Pending |
| 2  | `AssertPasswordRight` | SELECT + UPDATE (conditional) | ⏳ Pending |
| 3  | `AssertPasswordNeverUsesByUser` | SELECT Multiple Rows + Loop | ⏳ Pending |
| 4  | `AssertPasswordExpiredIsPassed` | SELECT Single + Date Compare | ⏳ Pending |
| 5  | `AssertPasswordAttemptPassed` | SELECT Single + Counter Check | ⏳ Pending |
| 6  | `ResetPasswordAttempt` | UPDATE Counter | ⏳ Pending |
| 7  | `AddUserAccessToken` | INSERT | ⏳ Pending |
| 8  | `AddPasswordHistory` | INSERT | ⏳ Pending |
| 9  | `GetUserByAccessToken` | SELECT with JOIN | ⏳ Pending |
| 10 | `DestroyToken` | DELETE | ⏳ Pending |
| 11 | `FindByCurrentSession` | SELECT with WHERE | ⏳ Pending |
| 12 | `UpdateProfileById` | UPDATE Multiple Fields | ⏳ Pending |
| 13 | `UpdatePasswordById` | UPDATE Password + Expired Date | ⏳ Pending |
| 14 | `DestroyAllToken` | DELETE Multiple | ⏳ Pending |

#### `modules/auth/repository/reset_password_repository.go`:

| No | Method Name | Pattern | Status |
|----|-------------|---------|--------|
| 1  | `RequestResetPassword` | SELECT + INSERT (transaction) | ⏳ Pending |
| 2  | `GetUserByResetPasswordToken` | SELECT Single | ⏳ Pending |
| 3  | `AddResetPasswordToken` | INSERT | ⏳ Pending |
| 4  | `DestroyResetPasswordToken` | DELETE | ⏳ Pending |
| 5  | `DestroyAllResetPasswordToken` | DELETE Multiple | ⏳ Pending |
| 6  | `IncreasePasswordExpiredAt` | UPDATE Date | ⏳ Pending |

## 📝 Conversion Examples

Lihat file `GORM_IMPLEMENTATION_GUIDE.md` untuk contoh lengkap konversi patterns:
- Pattern 1: SELECT Single Row
- Pattern 2: SELECT dengan Custom Scan
- Pattern 3: SELECT Multiple Rows
- Pattern 4: INSERT
- Pattern 5: UPDATE
- Pattern 6: DELETE
- Pattern 7: Transaction

## ⚠️ Important Notes:

1. **Error Handling Changes:**
   - `sql.ErrNoRows` → `gorm.ErrRecordNotFound`
   - Check: `errors.Is(err, gorm.ErrRecordNotFound)`

2. **Context Handling:**
   - Always use: `repo.DB.WithContext(ctx)`

3. **Soft Delete:**
   - GORM automatically handles `deleted_at IS NULL`
   - Use `.Unscoped()` jika perlu include soft deleted

4. **Field Names:**
   - GORM menggunakan snake_case otomatis
   - Atau gunakan column name dari GORM tags

5. **Time Handling:**
   - `NOW()` → `time.Now()`
   - `INTERVAL '3 months'` → `time.Now().AddDate(0, 3, 0)`

## 🧪 Testing After Implementation:

```bash
# 1. Build
go build

# 2. Run
./base-go

# 3. Test API endpoints:
# - POST /api/auth/login
# - POST /api/auth/logout
# - GET /api/auth/profile
# - PUT /api/auth/profile
# - PUT /api/auth/password
# - POST /api/auth/reset-password
```

## 📊 Progress Tracker:

Total Methods: 20
- ✅ Completed: 0
- ⏳ Pending: 20
- 🔄 In Progress: 0

## 🚀 Quick Start Command:

```bash
# Install dependencies
go get -u gorm.io/gorm gorm.io/driver/postgres && go mod tidy

# Start implementing methods one by one using patterns from GORM_IMPLEMENTATION_GUIDE.md
```

## 💡 Tips:
- Convert & test 2-3 methods at a time
- Use `.Debug()` untuk see generated SQL: `repo.DB.Debug().Find(&users)`
- Check logs untuk verify queries
- Test each endpoint after conversion

## 🆘 Support:
Jika stuck, refer to:
1. `GORM_IMPLEMENTATION_GUIDE.md` - Detailed patterns
2. GORM docs: https://gorm.io/docs/
3. Check GORM SQL logs with `.Debug()`

