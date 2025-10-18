package constant

type NodeEnv string

const (
	NODE_ENV_DEVELOPMENT NodeEnv = "development"
	NODE_ENV_STAGING     NodeEnv = "staging"
	NODE_ENV_PRODUCTION  NodeEnv = "production"
)

func (s NodeEnv) Pointer() *NodeEnv {
	return &s
}

func (s *NodeEnv) Value() NodeEnv {
	if s == nil {
		return ""
	}

	return *s
}

type OrderBy int8

const (
	ORDER_BY_ASC  OrderBy = 1
	ORDER_BY_DESC OrderBy = -1
)

func (s OrderBy) Pointer() *OrderBy {
	return &s
}

func (s *OrderBy) Value() OrderBy {
	if s == nil {
		return 0
	}

	return *s
}

const (
	TRUE  = "true"
	FALSE = "false"

	// Context keys
	USER_TOKEN_CONTEXT_KEY = "user_token"
)

var (
	_true  bool = true
	True        = &_true
	_false bool = false
	False       = &_false
)

type Operator string

const (
	OPERATOR_EQUAL      = "equal_to"
	OPERATOR_NOT_EQUAL  = "not_equal_to"
	OPERATOR_IN         = "in"
	OPERATOR_NOT_IN     = "not_in"
	OPERATOR_LIKE       = "like"
	OPERATOR_ILIKE      = "ilike"
	OPERATOR_GREATER    = "gt"
	OPERATOR_GREATER_EQ = "gte"
	OPERATOR_LESS       = "lt"
	OPERATOR_LESS_EQ    = "lte"
)
