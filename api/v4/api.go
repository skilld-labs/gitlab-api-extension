package apiv4

import "../"

type ApiAPI struct {
	Api api.Api
}

func NewApiAPI(a api.Api) ApiAPI {
	return ApiAPI{Api: a}
}
