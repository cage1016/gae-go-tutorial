package responses

type DataRes struct {
	APIVersion string      `json:"apiVersion"`
	Data       interface{} `json:"data"`
}

type Responser interface {
	Response() interface{}
}

type Paging struct {
	Total  uint64 `json:"total"`
	Offset uint64 `json:"offset"`
	Limit  uint64 `json:"limit"`
}
