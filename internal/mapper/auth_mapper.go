package mapper

import (
    "time"

    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "workHub/internal/dto"
    "workHub/internal/entity"
)

func HashPassword(plain string) (string, error) {
    b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(b), nil
}


func ToRegisterResponse(user *entity.User) dto.RegisterResponse {
	return dto.RegisterResponse{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
	}
}


func ToUserEntity(req dto.RegisterRequest) (entity.User, error) {
    hashed, err := HashPassword(req.Password)
    if err != nil {
        return entity.User{}, err
    }
    return entity.User{
        ID:        uuid.NewString(),
        Email:     req.Email,
        Username:  req.Username,
        Password:  hashed,
        CreatedAt: time.Now(),
    }, nil
}