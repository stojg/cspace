package shaders

type BloomSeparator struct {
	Program          uint32
	LocScreenTexture int32
}

func NewBloomSeparator() *BloomSeparator {
	c := buildShader("fx", "fx_brigthness_sep")

	return &BloomSeparator{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
	}
}
