package types

type RegisterPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginPayload struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RestorePassword struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"password" binding:"required"`
}
