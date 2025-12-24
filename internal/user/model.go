package user

type User struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}
