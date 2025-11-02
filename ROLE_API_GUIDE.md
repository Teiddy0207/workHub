# Role CRUD API Test Guide

## 🚀 Server đang chạy tại: http://localhost:8088

## 📋 Role API Endpoints

### 1. **Tạo Role mới**
```bash
POST /roles
Content-Type: application/json

{
    "name": "Administrator",
    "code": "ADMIN",
    "description": "Quản trị viên hệ thống",
    "is_active": true
}
```

### 2. **Lấy danh sách Roles**
```bash
GET /roles?page=1&size=10&search=admin
```

### 3. **Lấy Role theo ID**
```bash
GET /roles/{id}
```

### 4. **Cập nhật Role**
```bash
PUT /roles/{id}
Content-Type: application/json

{
    "name": "Super Administrator",
    "description": "Quản trị viên cấp cao",
    "is_active": true
}
```

### 5. **Xóa Role**
```bash
DELETE /roles/{id}
```

## 🧪 Test Cases

### Test 1: Tạo Role
```bash
curl -X POST http://localhost:8088/roles \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Manager",
    "code": "MANAGER",
    "description": "Quản lý",
    "is_active": true
  }'
```

### Test 2: Lấy danh sách Roles
```bash
curl -X GET "http://localhost:8088/roles?page=1&size=10"
```

### Test 3: Lấy Role theo ID (thay {id} bằng ID thực tế)
```bash
curl -X GET http://localhost:8088/roles/{id}
```

### Test 4: Cập nhật Role
```bash
curl -X PUT http://localhost:8088/roles/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Senior Manager",
    "description": "Quản lý cấp cao"
  }'
```

### Test 5: Xóa Role
```bash
curl -X DELETE http://localhost:8088/roles/{id}
```

## 📊 Response Examples

### Success Response (Create/Update/Get)
```json
{
  "status": "success",
  "message": "Tạo mới quyền hạn thành công",
  "data": {
    "id": "uuid-here",
    "name": "Manager",
    "code": "MANAGER",
    "description": "Quản lý",
    "is_active": true,
    "created_at": "2025-10-20T15:30:00Z",
    "updated_at": "2025-10-20T15:30:00Z"
  }
}
```

### Success Response (List)
```json
{
  "status": "success",
  "message": "Lấy danh sách quyền hạn thành công",
  "data": {
    "items": [
      {
        "id": "uuid-here",
        "name": "Manager",
        "code": "MANAGER",
        "description": "Quản lý",
        "is_active": true,
        "created_at": "2025-10-20T15:30:00Z",
        "updated_at": "2025-10-20T15:30:00Z"
      }
    ],
    "total_items": 1,
    "page_number": 1,
    "page_size": 10
  }
}
```

### Success Response (Delete)
```json
{
  "status": "success",
  "message": "Xoá quyền hạn thành công"
}
```

### Error Response
```json
{
  "status": "error",
  "code": 400,
  "message": "credential already token"
}
```

## 🔍 Logs

Tất cả các API đều có logs chi tiết với emoji để dễ theo dõi:
- 🎯 Controller được gọi
- 📝 Request data
- 🔍 Repository operations
- ✅ Success operations
- ❌ Error operations

