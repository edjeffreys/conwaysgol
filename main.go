package main

import (
	"conwaysgol/pkg/gled"
	"github.com/go-gl/gl/v3.3-core/gl"
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

	vao := gled.BindNewVertexArray()

	gled.BindBuffer(gl.ARRAY_BUFFER)
	gled.BufferDataFloat(gl.ARRAY_BUFFER, shape.Vertices, gl.STATIC_DRAW)

	gled.BindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	gled.BufferDataUInt(gl.ELEMENT_ARRAY_BUFFER, shape.Indices, gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()

		gled.UseProgram(window.ShaderProgram)
		gled.BindVertexArray(vao)

		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Vertices) / 5))

		window.SwapBuffers()
	}

}
