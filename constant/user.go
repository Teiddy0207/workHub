package constant

type Gender int8

const (
	GENDER_MALE Gender = 1 + iota
	GENDER_FEMALE
	GENDER_OTHER
)

func (s Gender) Pointer() *Gender {
	return &s
}

func (s *Gender) Value() Gender {
	if s == nil {
		return 0
	}

	return *s
}

type UserStatus string

const (
	USER_STATUS_ACTIVED  UserStatus = "actived"
	USER_STATUS_INACTIVE UserStatus = "inactive"
	USER_STATUS_PENDING  UserStatus = "pending"
)

func (s UserStatus) Pointer() *UserStatus {
	return &s
}

func (s *UserStatus) Value() UserStatus {
	if s == nil {
		return "UNKNOWN"
	}

	return *s
}

var (
	USER_LOGGED_IN_SESSION_ID_CONTEXT_KEY = "USER_SESSION_ID"
	USER_LOGGED_INFO                      = "USER_INFO"
	USER_ROLE_INFO                        = "USER_ROLE_INFO"
	USER_PERMISSIONS                      = "USER_PERMISSIONS"
	USER_PASSWORD_CHANGED_SUCCESSFULLY    = "Mật khẩu đã được đổi thành công"
	USER_UPDATED_SUCCESSFULLY             = "Cập nhật người dùng thành công"
	USER_DELETED                          = "Xoá người dùng thành công"
	USER_GET_LIST_SUCCESSFULLY            = "Lấy danh sách người dùng thành công"
	USER_GET_SUCCESSFULLY                 = "Lấy thông tin người dùng thành công"

	USER_PASSWORD_CHANGED_FAIL = "Mật khẩu đã được đổi thất bại"
	USER_UPDATED_FAIL          = "Cập nhật người dùng thất bại"
	USER_DELETED_FAIL          = "Xoá người dùng thất bại"
	USER_GET_LIST_FAIL         = "Lấy danh sách người dùng thất bại"
	USER_GET_FAIL              = "Lấy thông tin người dùng thất bại"
)
