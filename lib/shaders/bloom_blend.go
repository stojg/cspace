package shaders

type BloomBlend struct {
	Program          uint32
	LocScreenTexture int32
	LocBloomTexture  int32
}

func NewBloomBlend() *BloomBlend {
	c := buildShader("fx", "fx_bloom_blender")

	return &BloomBlend{
		Program:          c,
		LocScreenTexture: loc(c, "screenTexture"),
		LocBloomTexture:  loc(c, "bloomTexture"),
	}
}
