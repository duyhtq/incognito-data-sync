package serializers

type PagingResp struct {
	Records     interface{}
	Page        *int
	Limit       *int
	TotalRecord *int
	TotalPage   *int
}

type Resp struct {
	Result interface{} `json:"Result"`
	Error  interface{} `json:"Error"`
}
