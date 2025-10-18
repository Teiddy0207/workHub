# WorkHub API - Login Endpoint

## API Login đã được implement theo đúng pattern của dự án

### 📋 Cấu trúc API Login

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response Success (200):**
```json
{
  "status": "success",
  "message": "login success",
  "data": {
    "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2024-01-01T12:00:00Z",
    "user": {
      "id": "uuid-string",
      "email": "user@example.com",
      "username": "username"
    }
  }
}
```

**Response Error (400):**
```json
{
  "status": "error",
  "code": 400,
  "message": "login failed"
}
```

### 🔧 Cấu trúc Implementation

#### 1. **DTO Layer** (`internal/dto/auth.go`)
- `LoginRequest`: Input validation cho email và password
- `LoginResponse`: Response với access token, refresh token và user info
- `UserInfo`: Thông tin user cơ bản

#### 2. **Repository Layer** (`internal/repository/auth_repo.go`)
- `GetUserByEmail()`: Tìm user theo email từ database

#### 3. **Service Layer** (`internal/service/auth_service.go`)
- `Login()`: Logic xử lý login:
  - Tìm user theo email
  - Verify password bằng bcrypt
  - Tạo JWT tokens (access + refresh)
  - Trả về response

#### 4. **Controller Layer** (`internal/controller/auth_controller.go`)
- `Login()`: Handle HTTP request/response
- Validation input
- Call service layer
- Return JSON response

#### 5. **Router** (`router/router.go`)
- Route: `POST /auth/login`

### 🛠️ Các Package được sử dụng

- **JWT**: `pkg/jwt` - Generate và verify JWT tokens
- **Hash**: `pkg/utils` - Hash và compare passwords
- **Error**: `constant/error.go` - Error constants
- **Handler**: `pkg/handler` - Response formatting

### ⚠️ TODO Items

1. **JWT Configuration**: Cần load RSA keys từ config file hoặc environment variables
2. **Error Handling**: Có thể cải thiện error messages cụ thể hơn
3. **Validation**: Có thể thêm validation cho email format
4. **Logging**: Có thể thêm logging cho security events

### 🚀 Cách test API

```bash
# Start server
go run cmd/main.go

# Test login API
curl -X POST http://localhost:8088/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### 📝 Notes

- API đã được implement theo đúng clean architecture pattern của dự án
- Sử dụng các package có sẵn trong dự án
- Tuân thủ naming convention và structure
- Ready để integrate với frontend
