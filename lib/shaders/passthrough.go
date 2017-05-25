package shaders

type Passthrough struct {
	Program          uint32
	LocScreenTexture int32
}

func NewPassthrough() *Passthrough {
	c := buildShader("fx", "fx_pass")
	return &Passthrough{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
	}
}
