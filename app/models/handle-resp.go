package models

type GlobalErrorHandlerResp struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type GlobalSuccessHandlerResp struct {
	Status  bool `json:"status"`
	Message any  `json:"message"`
	Data    any  `json:"data,omitempty"`
	Page    any  `json:"page,omitempty"`
	Limit   any  `json:"limit,omitempty"`
	Total   any  `json:"total,omitempty"`
}
