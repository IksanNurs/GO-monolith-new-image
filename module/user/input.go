package usercamp

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     int    `json:"role"`
}

type UpdateUserInput struct {
	ID       int    `json:"id"`
	Name     string `json:"username"`
	Email    string `json:"email" binding:"email"`
	Password string `json:"password"`
	Role     int    `json:"role"`
}

type LoginInput struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type FormCreateUserInput struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type FormUserInput struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
	Role     int    `form:"role"`
	Error    error
}
