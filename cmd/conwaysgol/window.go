package conwaysgol

import (
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
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
	shaderProgram  uint32
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
	glfw.WindowHint(glfw.ContextVersionMinor, 6)
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
		shaderProgram:  initOpenGL(),
	}
}

func (window Window) Draw(shape Shape) {

	// vertex buffer object
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	// vertex array object
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	// element buffer object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)

	// populate buffers
	gl.BufferData(gl.ARRAY_BUFFER, len(shape.Vertices)*4, gl.Ptr(shape.Vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(shape.Indices)*4, gl.Ptr(shape.Indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)

	for !window.glWindow.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		glfw.PollEvents()

		gl.UseProgram(window.shaderProgram)

		gl.BindVertexArray(vao)

		//gl.DrawElements(gl.TRIANGLES, 2, gl.UNSIGNED_INT, gl.Ptr(shape.Indices))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Vertices) / 3))

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
	vertexShaderSource, freeVertexShaderFn := gl.Strs(vertexShaderSource + "\x00")
	defer freeVertexShaderFn()

	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexShader, 1, vertexShaderSource, nil)
	gl.CompileShader(vertexShader)
	getShaderStatus(vertexShader)

	fragmentShaderSource, freeFragmentShaderFn := gl.Strs(fragmentShaderSource + "\x00")
	defer freeFragmentShaderFn()

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentShader, 1, fragmentShaderSource, nil)
	gl.CompileShader(fragmentShader)
	getShaderStatus(fragmentShader)

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	getProgramLinkStatus(shaderProgram)

	// free resources
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram
}

func getShaderStatus(shader uint32) {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	var shaderTypeId int32
	gl.GetShaderiv(shader, gl.SHADER_TYPE, &shaderTypeId)

	shaderType := "UNKNOWN"
	if shaderTypeId == gl.FRAGMENT_SHADER {
		shaderType = "FRAGMENT_SHADER"
	} else if shaderTypeId == gl.VERTEX_SHADER {
		shaderType = "VERTEX_SHADER"
	}

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		shaderLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(shaderLog))

		log.Panicln("Failed to compile shader:\n", shaderType)
	} else {
		log.Println("Compiled shader", shaderType, "successfully")
	}
}

func getProgramLinkStatus(program uint32) {
	var status int32

	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		programLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(programLog))

		log.Panicln("Failed to link program:\n", programLog)
	} else {
		log.Println("Linked program successfully")
	}
}
