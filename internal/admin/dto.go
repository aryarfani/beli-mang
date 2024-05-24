package admin

type RegisterAdminRequest struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=5,max=33"`
}

type LoginAdminRequest struct {
	Username string `json:"username" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15"`
}
