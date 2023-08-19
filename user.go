package todo

type User struct {
	Id       int    `json:"-" db:"id"` // id - not use
	Name     string `json:"name" binding:"required"` // necessarily
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
