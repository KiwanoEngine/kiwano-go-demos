package main

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"time"

	"kiwanoengine.com/kiwano"
	"kiwanoengine.com/kiwano/external/gl"
	"kiwanoengine.com/kiwano/render"
)

var vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 aColor;
layout (location = 2) in vec2 aTexCoord;

out vec3 ourColor;
out vec2 TexCoord;

void main()
{
    gl_Position = vec4(aPos, 1.0);
    ourColor = aColor;
    TexCoord = aTexCoord;
}
`

var fragmentShaderSource = `
#version 330 core
out vec4 FragColor;

in vec3 ourColor;
in vec2 TexCoord;

uniform sampler2D ourTexture;

void main()
{
    FragColor = texture(ourTexture, TexCoord);
}
`

var indices = []uint32{
	0, 1, 3, // first triangle
	1, 2, 3, // second triangle
}

var vertices = []float32{
	//     ---- 位置 ----       ---- 颜色 ----     - 纹理坐标 -
	0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // 右上
	0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // 右下
	-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // 左下
	-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // 左上
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/KiwanoEngine/kiwano-go-demos/texture")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}

type MainScene struct {
	textureID uint32
	VAO       uint32
	VBO       uint32
	EBO       uint32
	shader    *render.Shader
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	// Load the texture
	s.textureID, err = newTexture("assets/wall.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	// gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	// gl.TextureParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR_MIPMAP_LINEAR)

	// gl.GenTextures(1, &s.textureID)
	// gl.BindTexture(gl.TEXTURE_2D, s.textureID)
}

func (s *MainScene) OnExit() {
	gl.DeleteVertexArrays(1, &s.VAO)
	gl.DeleteBuffers(1, &s.VBO)
	gl.DeleteBuffers(1, &s.EBO)
}

func (s *MainScene) OnUpdate(dt time.Duration) {
	// bind Texture
	gl.BindTexture(gl.TEXTURE_2D, s.textureID)

	// render container
	s.shader.Use()
	gl.BindVertexArray(s.VAO)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return texture, nil
}

func setup() {
	kiwano.EnterScene(&MainScene{})
}

func main() {
	option := &kiwano.Option{
		Title:     "LearnTexture",
		Width:     640,
		Height:    480,
		Vsync:     true,
		Resizable: true,
	}
	if err := kiwano.Setup(option, setup); err != nil {
		log.Fatalln(err)
	}
}
