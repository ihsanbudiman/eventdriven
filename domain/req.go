package domain

type WebRequest struct {
	ID int `json:"id"`
}

type WebResponse struct {
	User    User   `json:"user"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}
