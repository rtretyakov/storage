package app

type createItemRequest struct {
	Value interface{}
	Ttl   int
}

type getItemResponse struct {
	Value interface{} `json:"value"`
}

type incrItemResponse struct {
	Value float64 `json:"value"`
}
