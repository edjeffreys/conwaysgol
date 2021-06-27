package conwaysgol

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
)

var (
	vertexShaderSource = `
			#version 330 core
			layout (location = 0) in vec3 aPos;
			
			void main()
			{
				gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
			}`
	fragmentShaderSource = `
			#version 330 core
			out vec4 FragColor;
			
			void main()
			{
    			FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
			}`
)

type Window struct {
	title    string
	width    int
	height   int
	glWindow *glfw.Window
	program  uint32
}

func init() {
	runtime.LockOSThread()
}

func InitWindow(title string, width int, height int) Window {
	if err := glfw.Init(); err != nil {
		log.Panicln("Failed to initialise GLFW", err)
	}
	log.Println("Initialised GLFW successfully")

	glfw.WindowHint(glfw.Resizable, glfw.False)

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		log.Panicln("Failed to initialise window", err)
	}
	log.Println("Initialised window successfully")

	window.MakeContextCurrent()

	return Window{
		title:    title,
		width:    width,
		height:   height,
		glWindow: window,
		program:  initOpenGL(),
	}
}

func (window Window) Draw(shape Shape) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// vertex array object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// vertex buffer object
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(shape.Vertices), gl.Ptr(shape.Vertices), gl.STATIC_DRAW)

	// element buffer object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(shape.Indices), gl.Ptr(shape.Indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 6, gl.FLOAT, false, 4*3, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	for !window.glWindow.ShouldClose() {
		glfw.PollEvents()

		gl.BindVertexArray(vao)
		gl.UseProgram(window.program)

		gl.DrawElements(gl.TRIANGLES, 2, gl.UNSIGNED_INT, gl.Ptr(shape.Indices))
		//gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Vertices) / 3))

		gl.BindVertexArray(0)
		gl.UseProgram(0)

		window.glWindow.SwapBuffers()
	}
}

func (window Window) ShouldClose() bool {
	return window.glWindow.ShouldClose()
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		log.Panicln("Failed to initialise OpenGL", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Initialised OpenGL version", version)

	// compile shaders
	vertexShaderSource, freeVertexShaderFn := gl.Strs(vertexShaderSource, "\x00")
	defer freeVertexShaderFn()

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, 1, vertexShaderSource, nil)
	gl.CompileShader(vertexShader)

	fragmentShaderSource, freeFragmentShaderFn := gl.Strs(fragmentShaderSource, "\x00")
	defer freeFragmentShaderFn()

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, 1, fragmentShaderSource, nil)
	gl.CompileShader(fragmentShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	return program
}
