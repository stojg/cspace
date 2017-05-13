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

type SceneNode interface {
	Render(projection, view mgl32.Mat4, tShader TextureShader, mShader MaterialShader)
	Add(mesh []*Mesh, s ShaderType, transform mgl32.Mat4)
	Destroy()
}

func NewBaseNode() SceneNode {
	origin := mgl32.Translate3D(0, 0, 0)
	q := &BaseNode{
		Node: Node{
			transform: &origin,
		},
	}
	return q
}

type BaseNode struct {
	Node
}

func (n *BaseNode) Render(projection, view mgl32.Mat4, tShader TextureShader, mShader MaterialShader) {
	var tMeshes []*Node
	var mMeshes []*Node
	children := n.Node.Children()
	for _, child := range children {
		if child.mesh.MeshType == TextureMesh {
			tMeshes = append(tMeshes, child)
		} else if child.mesh.MeshType == MaterialMesh {
			mMeshes = append(mMeshes, child)
		}
	}

	tShader.UsePV(projection, view)
	for i := range tMeshes {
		tMeshes[i].Render(projection, view, tShader, mShader)
	}

	mShader.UsePV(projection, view)
	for i := range mMeshes {
		mMeshes[i].Render(projection, view, tShader, mShader)
	}
}

type Node struct {
	children   []*Node
	shaderType ShaderType
	transform  *mgl32.Mat4
	mesh       *Mesh
}

func (n *Node) Render(projection, view mgl32.Mat4, tShader TextureShader, mShader MaterialShader) {
	transform := *n.transform
	if n.mesh.MeshType == TextureMesh {
		gl.UniformMatrix4fv(tShader.ModelUniform(), 1, false, &transform[0])
	} else {
		gl.UniformMatrix4fv(mShader.ModelUniform(), 1, false, &transform[0])
	}
	n.mesh.Render(tShader, mShader)
}

func (n *Node) Destroy() {
	n.children = make([]*Node, 0)
}

func (n *Node) Children() []*Node {
	var children []*Node
	for _, child := range n.children {
		children = append(children, child)
		children = append(children, child.Children()...)
	}
	return children
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
