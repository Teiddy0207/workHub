package dto

type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	FullName    string  `json:"full_name" binding:"required"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Role        string  `json:"role" binding:"omitempty,oneof=student teacher admin"`
}

type RegisterResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Role        string  `json:"role"`
	CreatedAt   string  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
	User         UserInfo `json:"user"`
}

type UserInfo struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Role        string  `json:"role"`
}

type UpdateUserRequest struct {
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Role        *string `json:"role" binding:"omitempty,oneof=student teacher admin"`
}

type UserResponse struct {
	ID          string  `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Role        string  `json:"role"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type UserItem struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type PaginatedUserResponse = Pagination[RegisterResponse]
