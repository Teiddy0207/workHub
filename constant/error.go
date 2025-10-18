package constant

import "errors"

var (
	ErrTakenCredential     = errors.New("credential already token")
	ErrInternalServer      = errors.New("InternalServerError")
	ErrVerificationExpired = errors.New("verification expired")
	ErrUnAuthentication    = errors.New("un authentication")
	ErrNotFound            = errors.New("not found")
	ErrUnprocessableEntity = errors.New("UnprocessableEntity")
	ErrBadRequestException = errors.New("BadRequestException")
	ErrUUIDInvalid         = errors.New("uuid invalid")

	//auth
	ErrUsernameOrPasswordIncorrect = errors.New("sai username")
	ErrPasswordIncorrect           = errors.New("mật khẩu không chính xác")
	ErrUserInactive                = errors.New("người dùng không hoạt động")
	ErrCantNotChangePassword       = errors.New("không được đổi mật khẩu")
	ErrEmailOrPhoneIncorrect       = errors.New("email hoặc số điện thoại không chính xác")
	ErrUnconfirmedPhoneNumber      = errors.New("số điện thoại chưa xác minh")
	ErrCantNotChangeUsername       = errors.New("không thể đổi được mật khẩu")

	ErrEmailConflictException    = errors.New("EmailConflictException")
	ErrPhoneConflictException    = errors.New("PhoneConflictException")
	ErrUsernameConflictException = errors.New("UsernameConflictException")

	ErrInvalidInput        = errors.New("thông tin nhập không hợp lệ")
	ErrInvalidFilterFormat = errors.New("thông tin lọc không hợp lệ")
	ErrInvalidOTP        = errors.New("mã xác nhận không hợp lệ hoặc đã hết hạn")
	ErrInvalidResetToken = errors.New("reset token không hợp lệ hoặc đã hết hạn")
	ErrPasswordMismatch  = errors.New("mật khẩu xác nhận không khớp")
	ErrUserEmailNotFound = errors.New("email không tồn tại trong hệ thống")

	// Registration specific errors
	ErrNoRolesAvailable      = errors.New("no roles available in system")
	ErrRoleNotFound          = errors.New("role not found")
	ErrPasswordHashingFailed = errors.New("password hashing failed")
	ErrDatabaseCreateFailed  = errors.New("database create operation failed")
	ErrInvalidFieldName      = errors.New("invalid field name")
	ErrDatabaseConnection    = errors.New("database connection error")
	
	// Storage errors
	ErrUploadToAWS      = errors.New("upload to AWS failed")
	ErrSaveStorageDb    = errors.New("save storage to database failed")
	ErrDeleteFromAWS    = errors.New("delete from AWS failed")
)
