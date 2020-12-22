package main


const (
	StatusSuccess				= 200
	ErrorReadBodyFailed			= -300100
	ErrorMarshalJSON			= -300101
	ErrorAuthentication			= -300102
)

var StatusText = map[int]string{
	StatusSuccess:			"successful",
}


type APIResponse struct {
	Code		int			`json:"code"`	
	Success 	bool		`json:"success"`	
	Msg 		string		`json:"message"`	
	Module		string		`json:"module,omitempty"`	
	Data 		interface{}	`json:"data,omitempty"`	
}


type HttpPara struct {
	
}

type AppCfg struct {
	NoAuth 				bool 		`json:"no_auth"`
	AuthMaxInterval		int64		`json:"auth_max_interval"`

}