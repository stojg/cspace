package main

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type ShaderType int

const (
	NoRenderMesh ShaderType = iota
	TextureMesh
	MaterialMesh
)

type Node struct {
	children   []*Node
	shaderType ShaderType
	transform  *mgl32.Mat4
	mesh       *Mesh
}

func (n *Node) Render(projection, view mgl32.Mat4, tShader TextureShader, mShader MaterialShader) {

	for _, child := range n.children {
		if child.transform != nil {
			transform := child.transform.Mul4(*n.transform)
			if child.mesh.MeshType == TextureMesh {
				tShader.UsePV(projection, view)
				gl.UniformMatrix4fv(tShader.ModelUniform(), 1, false, &transform[0])
			} else {
				mShader.UsePV(projection, view)
				gl.UniformMatrix4fv(mShader.ModelUniform(), 1, false, &transform[0])
			}
			child.mesh.Render(tShader, mShader)
		}
		child.Render(projection, view, tShader, mShader)
	}
}

func (n *Node) Destroy() {
	n.children = make([]*Node, 0)
}

func (n *Node) Add(mesh []*Mesh, s ShaderType, transform mgl32.Mat4) {

	for i := range mesh {
		child := &Node{
			shaderType: s,
			mesh:       mesh[i],
			transform:  &transform,
		}
		n.children = append(n.children, child)
	}
}
