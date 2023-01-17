package store

type baseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type listResponse struct {
	baseResponse
	Items []listItem `json:"items"`
}

type listItem struct {
	K string `json:"k"`
}

type getResponse struct {
	baseResponse
	Value string `json:"value"`
}

type setRequest struct {
	Value string `json:"value"`
}
