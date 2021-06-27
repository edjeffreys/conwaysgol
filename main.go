package main

import (
	"conwaysgol/pkg/gled"
	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth = 800
	windowHeight = 800
)

func main() {
	window := gled.InitWindow("Conway's Game of Life", windowWidth, windowHeight)
	defer glfw.Terminate()

	shape := gled.Square

	gled.BindBuffer(gl.ARRAY_BUFFER)
	gled.BindBuffer(gl.ELEMENT_ARRAY_BUFFER)

	vao := gled.BindNewVertexArray()

	// populate buffers
	gled.BufferDataFloat(gl.ARRAY_BUFFER, shape.Vertices, gl.STATIC_DRAW)
	gled.BufferDataUInt(gl.ELEMENT_ARRAY_BUFFER, shape.Indices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()

		gled.UseProgram(window.ShaderProgram)
		gled.BindVertexArray(vao)

		//gl.DrawElements(gl.TRIANGLES, 2, gl.UNSIGNED_INT, gl.Ptr(shape.Indices))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Vertices) / 3))

		gled.UnbindVertexArray()

		window.SwapBuffers()
	}

}
