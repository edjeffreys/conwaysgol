package main

import (
	"conwaysgol/cmd/conwaysgol"

	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth = 800
	windowHeight = 800
)

func main() {
	window := conwaysgol.InitWindow("Conway's Game of Life", windowWidth, windowHeight)
	defer glfw.Terminate()

	window.Draw(conwaysgol.Square)
	//for !window.ShouldClose() {
	//}

}
