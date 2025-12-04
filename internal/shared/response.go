package shared

type (
	ResponseJSON struct {
		Message string      `json:"message"`
		Count   int         `json:"count"`
		Data    interface{} `json:"data"`
	}
)
