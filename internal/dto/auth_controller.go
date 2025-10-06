package dto
import "workHub/utils"
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3,max=32"`
    Password string `json:"password" binding:"required,min=6"`
}

type RegisterResponse struct {
    ID       string `json:"id"`
    Email    string `json:"email"`
    Username string `json:"username"`
}

type UserItem struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Username  string `json:"username"`
    CreatedAt string `json:"created_at"`
 
}

type UserListResponse struct {
    Items []UserItem     `json:"items"`
    Meta  utils.Pagination `json:"meta"`
}