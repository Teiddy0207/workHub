package constant

type QueueStatus string

const (
	QUEUE_STATUS_BOOKED    UserStatus = "booked"
	QUEUE_STATUS_WAITING   UserStatus = "waiting"
	QUEUE_STATUS_SELFCARE  UserStatus = "selfcare"
	QUEUE_STATUS_CALLED    UserStatus = "called"
	QUEUE_STATUS_COMPLETE  UserStatus = "completed"
	QUEUE_STATUS_CANCELLED UserStatus = "cancelled"
)

func (s QueueStatus) Pointer() *QueueStatus {
	return &s
}

func (s *QueueStatus) Value() QueueStatus {
	if s == nil {
		return "UNKNOWN"
	}

	return *s
}

var (
	QUEUE_CREATED_FAIL                 = "Tạo mới xếp số thất bại"
	QUEUE_CREATED_SERVICE_TYPE_ID_FAIL = "Tạo mới xếp số thiếu ID dịch vụ"
	QUEUE_CREATED_STORE_ID_FAIL        = "Tạo mới xếp số thiếu ID cửa hàng"
	QUEUE_CREATED_PHONE_NUMBER_FAIL    = "Tạo mới xếp số thiếu số điện thoại"
	QUEUE_CREATED_CUSTOMER_NAME_FAIL   = "Tạo mới xếp số thiếu tên khách hàng"
	QUEUE_UPDATED_FAIL                 = "Cập nhật xếp số thất bại"
	QUEUE_DELETED_FAIL                 = "Xoá xếp số thất bại"
	QUEUE_GET_LIST_FAIL                = "Lấy danh sách xếp số thất bại"
	QUEUE_GET_FAIL                     = "Lấy thông tin xếp số thất bại"
	QUEUE_GET_STORE_ID_FAIL            = "Lấy thông tin store_id thất bại"

	QUEUE_CREATED_SUCCESSFULLY  = "Tạo mới xếp số thành công"
	QUEUE_UPDATED_SUCCESSFULLY  = "Cập nhật xếp số thành công"
	QUEUE_DELETED_SUCCESSFULLY  = "Xoá xếp số thành công"
	QUEUE_GET_LIST_SUCCESSFULLY = "Lấy danh sách xếp số thành công"
	QUEUE_GET_SUCCESSFULLY      = "Lấy thông tin xếp số thành công"
	QUEUE_UPDATED_CANCELLED_SUCCESSFULLY  = "Cập nhật huỷ hàng chờ thành công"

)
