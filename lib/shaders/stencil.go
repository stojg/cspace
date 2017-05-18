package shaders

func NewStencil() *Stencil {
	s := &Stencil{}
	s.Program = buildShader("null", "null")
	s.LocModel = loc(s.Program, "model")
	s.LocView = loc(s.Program, "view")
	s.LocProjection = loc(s.Program, "projection")
	return s
}

type Stencil struct {
	Program       uint32
	LocModel      int32
	LocView       int32
	LocProjection int32
}
