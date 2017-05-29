package shaders

type IBLBrdf struct {
	Program uint32
}

func NewIBLBrdf() *IBLBrdf {
	c := buildShader("ibl_brdf", "ibl_brdf")
	return &IBLBrdf{
		Program: c,
	}
}
