package dllm


type Mobile struct {
	buf []byte
}

func NewMobile() *Mobile {
	return &Mobile{}
}

func (m *Mobile) Query() string {
	return "on_mobile"
}
