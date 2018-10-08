package models

// Message Define our message object
type Message struct {
	Username  string `json:"username"`
	Message   string `json:"message"`
	UpdatedAt string `json:"updated_at"`
}
