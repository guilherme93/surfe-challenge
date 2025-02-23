package dto

type CountResponse struct {
	Count int `json:"count"`
}

type Prediction struct {
	Action      string  `json:"action"`
	Probability float64 `json:"probability"`
}
