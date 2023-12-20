package models

type User struct {
	Id          int32  `json:"id"`
	Name        string `json:"name"`
	MiddleName  string `json:"middlename"`
	Rut         string `json:"rut"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
type SignUpUser struct {
	Name        string `json:"name"`
	MiddleName  string `json:"middlename"`
	Rut         string `json:"rut"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}
type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}