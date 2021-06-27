package conwaysgol

type Shape struct {
	Vertices []float32
	Indices []uint32
}

var (
	Square = Shape{
		Vertices: []float32{
			//0.5,  0.5, 0.0,  // top right
			0.5, -0.5, 0.0,  // bottom right
			-0.5, -0.5, 0.0, // bottom left
			-0.5,  0.5, 0.0, // top left
		},
		Indices: []uint32{
			0, 1, 3,   // first triangle
			1, 2, 3,   // second triangle
		},
	}
)
