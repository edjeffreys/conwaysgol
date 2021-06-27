package gled

import (
	"io/ioutil"
	"log"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	vertexShaderPath = "assets/basic.vert"
	fragmentShaderPath = "assets/basic_orange.frag"
)

type Window struct {
	title         string
	width         int
	height        int
	glWindow      *glfw.Window
	ShaderProgram ProgramId
}

type ShaderId uint32
type ProgramId uint32
type BufferId uint32
type VertexArrayId uint32

func init() {
	runtime.LockOSThread()
}

// InitWindow initialises a new window with the given title, width, and height; returning the `Window` struct
func InitWindow(title string, width int, height int) Window {
	if err := glfw.Init(); err != nil {
		log.Panicln("Failed to initialise GLFW", err)
	}
	log.Println("Initialised GLFW successfully")

	glfw.WindowHint(glfw.Resizable, glfw.False)

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, title, nil, nil)

	if err != nil {
		log.Panicln("Failed to initialise window", err)
	}
	log.Println("Initialised window successfully")

	window.MakeContextCurrent()

	return Window{
		title:         title,
		width:         width,
		height:        height,
		glWindow:      window,
		ShaderProgram: initOpenGL(),
	}
}

// CompileShaderFromString compiles a shader from a GLSL string with a given shader type e.g. `gl.FRAGMENT_SHADER` and returns the generated gl shader id.
func CompileShaderFromString(shaderSource string, shaderType uint32) ShaderId {
	cShaderSource, freeShaderFn := gl.Strs(shaderSource + "\x00")
	defer freeShaderFn()

	shaderId := gl.CreateShader(shaderType)
	gl.ShaderSource(shaderId, 1, cShaderSource, nil)
	gl.CompileShader(shaderId)
	getShaderStatus(shaderId)

	return ShaderId(shaderId)
}

// CompileShaderFromFile compiles a shader from a GLSL file with a given shader type e.g. `gl.FRAGMENT_SHADER` and returns the generated gl shader id.
func CompileShaderFromFile(filePath string, shaderType uint32) ShaderId {
	file, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Panicln("Failed to load shader file:\n", err)
	}

	return CompileShaderFromString(string(file), shaderType)
}

// CreateProgram links a new program with given shaders, returning the program id.
func CreateProgram(shaderIds ...ShaderId) ProgramId {
	programId := gl.CreateProgram()

	for _, shader := range shaderIds {
		gl.AttachShader(programId, uint32(shader))
		// free resource
		gl.DeleteShader(uint32(shader))
	}

	gl.LinkProgram(programId)

	getProgramLinkStatus(programId)

	return ProgramId(programId)
}

// BindBuffer generates a buffer id and binds the id to given buffer type, returning the buffer id.
func BindBuffer(bufferType uint32) BufferId {
	var bufferId uint32

	gl.GenBuffers(1, &bufferId)
	gl.BindBuffer(bufferType, bufferId)

	return BufferId(bufferId)
}

// BindNewVertexArray generates a vertex array id and binds it, returning the id.
func BindNewVertexArray() VertexArrayId {
	var vertexArrayId uint32

	gl.GenVertexArrays(1, &vertexArrayId)
	gl.BindVertexArray(vertexArrayId)

	return VertexArrayId(vertexArrayId)
}

// BindVertexArray binds the given vertex array id. Use `BindNewVertexArray` instead to generate a new id and bind it.
func BindVertexArray(vertexArrayId VertexArrayId) {
	gl.BindVertexArray(uint32(vertexArrayId))
}

// UnbindVertexArray unbinds the currently bound vertex array
func UnbindVertexArray() {
	gl.BindVertexArray(0)
}

// BufferDataFloat writes into the given buffer type, and usage, the given slice
func BufferDataFloat(bufferType uint32, data []float32, usage uint32) {
	gl.BufferData(bufferType, len(data) * 4, gl.Ptr(data), usage)
}

// BufferDataUInt writes into the given buffer type, and usage, the given slice
func BufferDataUInt(bufferType uint32, data []uint32, usage uint32) {
	gl.BufferData(bufferType, len(data) * 4, gl.Ptr(data), usage)
}

// UseProgram installs the given program as part of the current rendering state
func UseProgram(programId ProgramId) {
	gl.UseProgram(uint32(programId))
}

// ShouldClose returns value of close flag
func (window Window) ShouldClose() bool {
	return window.glWindow.ShouldClose()
}

// SwapBuffers swaps the front and back buffers of the window
func (window Window) SwapBuffers() {
	window.glWindow.SwapBuffers()
}

func initOpenGL() ProgramId {
	if err := gl.Init(); err != nil {
		log.Panicln("Failed to initialise OpenGL", err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("Initialised OpenGL version", version)

	vertexShader := CompileShaderFromFile(vertexShaderPath, gl.VERTEX_SHADER)
	fragmentShader := CompileShaderFromFile(fragmentShaderPath, gl.FRAGMENT_SHADER)

	shaderProgram := CreateProgram(vertexShader, fragmentShader)

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

		log.Panicln("Failed to compile shader", shaderType, ":\n", shaderLog)
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
