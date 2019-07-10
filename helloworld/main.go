package main

import (
	"log"
	"math"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"kiwanoengine.com/kiwano"
	"kiwanoengine.com/kiwano/core"
	"kiwanoengine.com/kiwano/render"
)

var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;

out vec3 ourColor;

void main()
{
	gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
	ourColor = aColor;
}
`

var fragmentShaderSource = `
#version 330 core
in vec3 ourColor;
out vec4 FragColor;

void main()
{
    FragColor = vec4(ourColor, 1.0);
}
`

var vertices = []float32{
	// 位置          // 颜色
	0.5, -0.5, 0.0, 1.0, 0.0, 0.0, // 右下
	-0.5, -0.5, 0.0, 0.0, 1.0, 0.0, // 左下
	0.0, 0.5, 0.0, 0.0, 0.0, 1.0, // 顶部
}

var indices = []uint32{ // 注意索引从0开始!
	0, 1, 3, // 第一个三角形
	1, 2, 3, // 第二个三角形
}

type MainScene struct {
	VAO    uint32
	VBO    uint32
	EBO    uint32
	shader *render.Shader
}

func (s *MainScene) OnEnter() {
	// Create shader program
	var err error
	s.shader, err = render.CreateShader(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatalln(err)
	}

	// Set up vertex array
	gl.GenVertexArrays(1, &s.VAO)
	gl.BindVertexArray(s.VAO)

	// Set up vertex buffer
	gl.GenBuffers(1, &s.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, s.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &s.EBO)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.EBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*len(indices), gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)
}

func (s *MainScene) OnUpdate(dt time.Duration) {
	if kiwano.Window.GLFWWindow.GetKey(glfw.KeyEscape) == glfw.Press {
		kiwano.Window.GLFWWindow.SetShouldClose(true)
	}

	timeValue := glfw.GetTime()
	greenValue := (math.Sin(timeValue) / 2.0) + 0.5
	vertexColorLocation := gl.GetUniformLocation(s.shader.ID, gl.Str("ourColor\x00"))

	s.shader.Use()
	gl.Uniform4f(vertexColorLocation, 0.0, float32(greenValue), 0.0, 1.0)

	//gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (s *MainScene) OnExit() {
	gl.DeleteVertexArrays(1, &s.VAO)
	gl.DeleteBuffers(1, &s.VBO)
}

func main() {
	// Init kiwano engine
	if err := kiwano.Init(&core.Option{
		Width:      640,
		Height:     480,
		Title:      "LearnOpenGL",
		ClearColor: core.ColorRGB(0.2, 0.3, 0.3),
		NoTitleBar: false,
		Resizable:  true,
		Fullscreen: false,
		Vsync:      true,
	}); err != nil {
		log.Fatalln(err)
	}
	defer kiwano.Destroy()

	// Enter scene
	kiwano.EnterScene(&MainScene{})

	// Start game
	kiwano.Run()
}
