package model

type APIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type UploadItem struct {
	URL string `json:"url"`
	POS string `json:"pos"`
}

