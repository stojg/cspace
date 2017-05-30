package shaders

type Gaussian struct {
	Program          uint32
	LocScreenTexture int32
	LocHorizontal    int32
}

func NewGaussian() *Gaussian {
	c := buildShader("fx", "fx_guassian_blur")

	return &Gaussian{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
		LocHorizontal:    loc(c, "horizontal"),
	}
}
