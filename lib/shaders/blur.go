package shaders

type Blur struct {
	Program          uint32
	LocScreenTexture int32
}

func NewBlur() *Blur {
	c := buildShader("fx", "fx_blur")

	return &Blur{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
	}
}
