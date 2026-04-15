package models

//
//type Login struct {
//	ID           uint64 `gorm:"column:id;primary_key;auto_increment" json:"id"`
//	Email        string `gorm:"column:email;size:45;not null;unique" validate:"required,email" json:"email" `
//	Password     string `gorm:"column:password;size:255;not null;" validate:"required,min=8,max=255" json:"password"`
//	RefreshToken string `gorm:"column:refresh_token;size:255;not null;" json:"refresh_token"`
//}
//
//type JWToken struct {
//	Data    string `gorm:"column:data;not null;" json:"data"`
//	Token   string `gorm:"column:access_token;not null;" json:"access_token"`
//	Refresh string `gorm:"column:refresh_token;not null;" json:"refresh_token"`
//}
