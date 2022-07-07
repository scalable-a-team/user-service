package forms

type UserSignUp struct {
	Username  string `form:"username" json:"username" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
}

type UserSignIn struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `form:"refresh_token" json:"refresh_token" binding:"required"`
}

type UserProfileResponse struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserGroupResponse struct {
	Name string `json:"name"`
}

type UserResponse struct {
	ID       uint                `json:"id"`
	Username string              `json:"username"`
	Profile  UserProfileResponse `json:"profile"`
	Group    UserGroupResponse   `json:"group"`
}

type LoginResponse struct {
	User    UserResponse `json:"user"`
	Token   string       `json:"access_token"`
	Refresh string       `json:"refresh_token"`
}
