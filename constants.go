package main

var (
	seed           int64   = 1
	tickUpdateRate float32 = float32(1.0 / 30.0) //tick rate (30 TPS) ticks per second
	runningSpeed   float32 = 1.3
	flyingSpeed    float32 = 5
	walkingSpeed   float32 = 2.3
	jumpHeight     float32 = 0.25
)
var (
	DirtID  uint16 = 0
	GrassID uint16 = 1
	StoneID uint16 = 2
)

var CubeVertices []float32 = []float32{

	// Front face
	-0.5, -0.5, 0.5, // Bottom-left
	0.5, -0.5, 0.5, // Bottom-right
	0.5, 0.5, 0.5, // Top-right
	-0.5, -0.5, 0.5, // Bottom-left
	0.5, 0.5, 0.5, // Top-right
	-0.5, 0.5, 0.5, // Top-left

	// Back face
	-0.5, -0.5, -0.5, // Bottom-left
	-0.5, 0.5, -0.5, // Top-left
	0.5, 0.5, -0.5, // Top-right
	-0.5, -0.5, -0.5, // Bottom-left
	0.5, 0.5, -0.5, // Top-right
	0.5, -0.5, -0.5, // Bottom-right

	// Left face
	-0.5, -0.5, -0.5, // Bottom-left
	-0.5, -0.5, 0.5, // Bottom-right
	-0.5, 0.5, 0.5, // Top-right
	-0.5, -0.5, -0.5, // Bottom-left
	-0.5, 0.5, 0.5, // Top-right
	-0.5, 0.5, -0.5, // Top-left

	// Right face
	0.5, -0.5, -0.5, // Bottom-left
	0.5, 0.5, -0.5, // Top-left
	0.5, 0.5, 0.5, // Top-right
	0.5, -0.5, -0.5, // Bottom-left
	0.5, 0.5, 0.5, // Top-right
	0.5, -0.5, 0.5, // Bottom-right

	// Top face
	-0.5, 0.5, -0.5, // Bottom-left
	-0.5, 0.5, 0.5, // Bottom-right
	0.5, 0.5, 0.5, // Top-right
	-0.5, 0.5, -0.5, // Bottom-left
	0.5, 0.5, 0.5, // Top-right
	0.5, 0.5, -0.5, // Top-left

	// Bottom face
	-0.5, -0.5, -0.5, // Bottom-left
	0.5, -0.5, -0.5, // Bottom-right
	0.5, -0.5, 0.5, // Top-right
	-0.5, -0.5, -0.5, // Bottom-left
	0.5, -0.5, 0.5, // Top-right
	-0.5, -0.5, 0.5, // Top-left

}
var CubeUVs []uint8 = []uint8{
	// Front face
	0, 3, // Bottom-left
	2, 3, // Bottom-right
	2, 1, // Top-right
	0, 3, // Bottom-left
	2, 1, // Top-right
	0, 1, // Top-left

	// Back face
	2, 3, // Bottom-left
	2, 1, // Top-left
	0, 1, // Top-right
	2, 3, // Bottom-left
	0, 1, // Top-right
	0, 3, // Bottom-right

	// Left face
	0, 3, // Bottom-left
	2, 3, // Bottom-right
	2, 1, // Top-right
	0, 3, // Bottom-left
	2, 1, // Top-right
	0, 1, // Top-left

	// Right face
	0, 3, // Bottom-left
	0, 1, // Top-left
	2, 1, // Top-right
	0, 3, // Bottom-left
	2, 1, // Top-right
	2, 3, // Bottom-right

	// Top face
	0, 3, // Bottom-left
	0, 1, // Bottom-right
	2, 1, // Top-right
	0, 3, // Bottom-left
	2, 1, // Top-right
	2, 3, // Top-left

	// Bottom face
	0, 3, // Bottom-left
	2, 3, // Bottom-right
	2, 1, // Top-right
	0, 3, // Bottom-left
	2, 1, // Top-right
	0, 1, // Top-left
}
