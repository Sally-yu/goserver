package model

type User struct {
	Userid   string `bson:"userid"`
	Username string `bson:"username"`
	Phone    string `bson:"phone"`
	Email    string `bson:"email"`
	Status   string `bson:"status"`
}

func (user *User) Save() {

}
