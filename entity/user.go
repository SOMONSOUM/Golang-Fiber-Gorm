package entity

type User struct {
	ID       uint64 `gorm:"primary_key:autoIncrement" json:"id"`
	Username string `gorm:"type: varchar(255)" json:"username"`
	Email    string `gorm:"uniqueIndex;type: varchar(255)" json:"email"`
	Password string `gorm:"->;<-;not null" json:"password"`
	Token    string `gorm:"-" json:"token,omitempty"`
}
