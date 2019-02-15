package schema

// The Error contains error relevant information.
type Error struct {
	// The general error message
	Error string `json:"error"`

	// The http error code.
	ErrorCode int `json:"error_code"`

	// The http error code.
	ErrorDescription string `json:"error_description"`
}
