package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/cspace/lib/rgbe"
)

type TextureType string

const (
	Albedo    TextureType = "albedo"
	Metallic  TextureType = "metallic"
	Roughness TextureType = "roughness"
	Normal    TextureType = "normal"
)

type Texture struct {
	ID          uint32
	textureType TextureType // type of texture, like diffuse, specular or bump
}

// GetTexture will load and return a Texture and panic if the texture could not be loaded
func GetTexture(texType TextureType, file string, gammaCorrect bool) *Texture {
	texture, err := newTexture(texType, filepath.Join("textures", file), gammaCorrect)
	if err != nil {
		panic(err)
	}
	return texture
}

// GetHDRTexture will load and return a Texture and panic if the texture could not be loaded
func GetHDRTexture(file string) *Texture {
	texture, err := NewHDRTexture(filepath.Join("textures", file))
	if err != nil {
		panic(err)
	}
	return texture
}

// GetLDRCubeMap will load and return an Texture id for a OpenGL texture cube map
func GetLDRCubeMap() uint32 {
	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, textureID)

	var textures = []string{
		"textures/skybox1/right.jpg",
		"textures/skybox1/left.jpg",
		"textures/skybox1/top.jpg",
		"textures/skybox1/bottom.jpg",
		"textures/skybox1/back.jpg",
		"textures/skybox1/front.jpg",
	}

	for i, file := range textures {
		rgba, err := loadImage(file)
		if err != nil {
			panic(err)
		}
		gl.TexImage2D(gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.SRGB8_ALPHA8,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, 0)
	return textureID
}

func NewHDRTexture(file string) (*Texture, error) {

	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	width, height, data, err := rgbe.Decode(fi)
	if err != nil {
		return nil, err
	}

	data = flipImgData(width, height, data)

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)

	gl.TexImage2D(gl.TEXTURE_2D,
		0,
		gl.RGB32F,
		int32(width),
		int32(height),
		0,
		gl.RGB,
		gl.FLOAT,
		gl.Ptr(data),
	)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	return &Texture{
		ID:          textureID,
		textureType: Albedo,
	}, nil
}

func loadImage(file string) (*image.RGBA, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Texture %q not found on disk: %v", file, err)
	}

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)
	return rgba, nil
}

func newTexture(name TextureType, file string, gammaCorrect bool) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	// flip it into open GL format
	rgba = flip(rgba)

	var texture uint32
	gl.GenTextures(1, &texture)

	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.ActiveTexture(gl.TEXTURE0)

	var internalFormat int32 = gl.RGBA
	if gammaCorrect {
		internalFormat = gl.SRGB8_ALPHA8
	}

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		internalFormat,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return &Texture{
		ID:          texture,
		textureType: name,
	}, nil
}

// flip the image upside down so that opengl can use it as a texture properly
func flip(src *image.RGBA) *image.RGBA {
	srcW := src.Bounds().Max.X
	srcH := src.Bounds().Max.Y
	dstW := srcW
	dstH := srcH

	dst := image.NewRGBA(src.Bounds())

	for dstY := 0; dstY < dstH; dstY++ {
		for dstX := 0; dstX < dstW; dstX++ {
			srcX := dstX
			srcY := dstH - dstY - 1
			srcOff := srcY*src.Stride + srcX*4
			dstOff := dstY*dst.Stride + dstX*4
			copy(dst.Pix[dstOff:dstOff+4], src.Pix[srcOff:srcOff+4])
		}
	}

	return dst
}

func flipImgData(width, height int, src []float32) []float32 {
	dst := make([]float32, len(src))

	rowSize := width * 3

	for y := 0; y < height; y++ {
		srcStart := y * rowSize
		srcEnd := srcStart + rowSize

		dstStart := (height - y - 1) * rowSize
		dstEnd := dstStart + rowSize

		copy(dst[dstStart:dstEnd], src[srcStart:srcEnd])
	}
	return dst
}
