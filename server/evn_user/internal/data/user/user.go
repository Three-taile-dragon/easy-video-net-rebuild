package user

type User struct {
	Id            int64
	Email         string
	Account       string
	Password      string
	Name          string
	Mobile        string
	CreateTime    int64
	LastLoginTime int64
}

func (*User) TableName() string {
	return "user"
}
