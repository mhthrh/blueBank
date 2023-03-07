package Entity

import (
	"time"
)

type (
	Customer struct {
		ID       int       `json:"id"validate:"required"`
		FullName string    `json:"fullName" validate:"required,alphanum"`
		UserName string    `json:"userName" validate:"required,alphanum"`
		PassWord string    `json:"passWord" validate:"required,printascii"`
		Email    string    `json:"email" validate:"required,email"`
		CreateAt time.Time `json:"createAt"`
		ExpireAt time.Time `json:"expireAt"`
	}
	CustomerLogin struct {
		UserName string `json:"userName" validate:"required,alphanum"`
		PassWord string `json:"passWord" validate:"required,printascii"`
	}
	CustomerLoginResponse struct {
		UserName  string `json:"userName"`
		Token     string `json:"token"`
		ValidTill string `json:"valid-till"`
	}
)
