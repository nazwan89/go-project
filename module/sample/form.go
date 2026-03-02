package sample

type Request struct {
	Name string `json:"name" form:"name"`
}

type Response struct {
	Message string `json:"message"`
}
