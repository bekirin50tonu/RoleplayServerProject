package dto

type SwaggerPayload[T any, TD any] struct {
	Meta T  `json:"meta"`
	Data TD `json:"data"`
}

type SwaggerErrors struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SwaggerErrorResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Errors  []SwaggerErrors `json:"Errors"`
}

type SwaggerSuccessResponse[T any, TD any] struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Payload SwaggerPayload[T, TD] `json:"payload"`
}
