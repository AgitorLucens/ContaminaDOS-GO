package types

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
	Others  interface{} `json:"others"`
}
