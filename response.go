package form

type customResponse struct {
	response map[string]any
}

func (c customResponse) String(key string) string {
	return c.response[key].(string)
}

func (c customResponse) Int(key string) int {
	return c.response[key].(int)
}

func (c customResponse) Float(key string) float64 {
	return c.response[key].(float64)
}

func (c customResponse) Bool(key string) bool {
	return c.response[key].(bool)
}

type CustomResponse interface {
	String(key string) string
	Int(key string) int
	Float(key string) float64
	Bool(key string) bool
}
