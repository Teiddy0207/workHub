package validate

type Validate struct {
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
}

