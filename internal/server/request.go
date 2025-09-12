package server

type EmailLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneLoginRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type EmailRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneRegisterRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}
