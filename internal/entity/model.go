package entity

// Request представляет структуру входного запроса
type Request struct {
	Fill int `json:"fill"`
}

// Response представляет структуру ответа от сервера
type Response struct {
	Fill int `json:"fill"`
}
