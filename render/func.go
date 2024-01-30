package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/png"
	"io"
	"net/http"

	// "github.com/gofrs/uuid"
	fdk "github.com/fnproject/fdk-go"
	fauxgl "github.com/hawl1/brickgl"
)

const (
	scale = 3
	fovy  = 50
	near  = 0.1
	far   = 1000
)

var (
	eye    = fauxgl.V(-0.75, 0.85, -2)
	center = fauxgl.V(0, 0.06, 0)
	up     = fauxgl.V(0, 1, 0)
	light  = fauxgl.V(0, 6, -4).Normalize()

	def = "{\"user_id\":13,\"items\":{\"face\":0,\"hats\":[20121,0,0,0,0],\"head\":0,\"tool\":0,\"pants\":0,\"shirt\":0,\"figure\":0,\"tshirt\":0},\"colors\":{\"head\":\"eab372\",\"torso\":\"85ad00\",\"left_arm\":\"eab372\",\"left_leg\":\"37302c\",\"right_arm\":\"eab372\",\"right_leg\":\"37302c\"}}"
)

// RenderEvent input data to lambda to return an ImageResponse
type RenderEvent struct {
	AvatarJSON string `json:"avatarJSON"`
	Size       int    `json:"size"`
}

// ImageResponse lambda response for a base64 encoded render
type ImageResponse struct {
	// gonna fix this, can stay for now UUID  string `json:"uuid"`
	Image string `json:"image"`
}

// LoadMeshFromURL loads mesh from url
func LoadMeshFromURL(url string) *fauxgl.Mesh {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	mesh, _ := fauxgl.LoadOBJFromReader(resp.Body)

	return mesh
}

func main() {
	fdk.Handle(fdk.HandlerFunc(HandleRenderEvent))
}

// HandleRenderEvent function to process the rendering
func HandleRenderEvent(ctx context.Context, in io.Reader, out io.Writer) {
	e := RenderEvent{}
	err := json.NewDecoder(in).Decode(&e)
	if err != nil {
		fmt.Fprintln(out, "Error:", err)
		return
	}

	if e.AvatarJSON == "" {
		e.AvatarJSON = def
	}

	aspect := float64(e.Size) / float64(e.Size)
	matrix := fauxgl.LookAt(eye, center, up).Perspective(fovy, aspect, near, far)
	shader := fauxgl.NewPhongShader(matrix, light, eye)
	context := fauxgl.NewContext(e.Size, e.Size, scale, shader)
	scene := fauxgl.NewScene(context)

	mesh := LoadMeshFromURL("https://hawli.pages.dev/obj/Torso.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/Head.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/LeftArm.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/LeftLeg.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/RightArm.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	mesh = LoadMeshFromURL("https://hawli.pages.dev/obj/RightLeg.obj")
	mesh.SmoothNormals()
	scene.AddObject(&fauxgl.Object{
		Mesh:  mesh,
		Color: fauxgl.HexColor("777"),
	})

	shader.AmbientColor = fauxgl.HexColor("AAA")
	shader.DiffuseColor = fauxgl.HexColor("777")
	shader.SpecularPower = 0

	newMatrix := scene.FitObjectsToScene(eye, center, up, fovy, aspect, near, far)
	shader.Matrix = newMatrix
	scene.Draw()

	outImg := context.Image()
	buf := new(bytes.Buffer)
	err = png.Encode(buf, outImg)
	if err != nil {
		fmt.Fprintln(out, "Error:", err)
	}

	resp := ImageResponse{
		Image: base64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	json.NewEncoder(out).Encode(resp)
}