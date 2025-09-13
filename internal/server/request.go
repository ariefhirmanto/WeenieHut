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

type UpdateUserProfileRequest struct {
	FileID            string `json:"fileId"`
	BankAccountName   string `json:"bankAccountName" validate:"omitempty,min=4,max=32"`
	BankAccountHolder string `json:"bankAccountHolder" validate:"omitempty,min=4,max=32"`
	BankAccountNumber string `json:"bankAccountNumber" validate:"omitempty,min=4,max=32"`
}

type UpdateUserContactRequest struct {
	Email string `json:"email" validate:"omitempty"`
	Phone string `json:"phone" validate:"omitempty"`
}
