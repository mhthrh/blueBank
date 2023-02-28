package Entity

type Gateway struct {
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	Ips         string `json:"ips"`
	GatewayName string `json:"gatewayName"`
	Status      bool   `json:"status"`
}
type GatewayLogin struct {
	UserName string `json:"userName" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum"`
}
