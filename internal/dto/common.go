package dto

type ListUID struct {
	UIDS []string `json:"uids" validate:"required"`
}

type ErrorResponse struct {
	Err string `json:"error"`
}
