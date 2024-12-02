package main

import (
	"MinecraftGolang/config"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

var chunks map[chunkPosition]chunkData = make(map[chunkPosition]chunkData)

func rebuildChunk(chunk chunkData, chunkPos chunkPosition) {
	propagateSunLight(chunkPos, chunk.blocksData)
	vao, trisCount := createChunkVAO(chunk.blocksData, chunkPos)
	chunks[chunkPos] = chunkData{
		blocksData: chunk.blocksData,
		vao:        vao,
		trisCount:  trisCount,
	}

}

func createChunks() {
	for x := int32(0); x < config.NumOfChunks; x++ {
		for z := int32(0); z < config.NumOfChunks; z++ {
			chunks[chunkPosition{x, z}] = chunk(chunkPosition{x, z})
		}
	}
	for chunkPos, _chunkData := range chunks {
		propagateSunLight(chunkPos, _chunkData.blocksData)

	}
	for chunkPos, _chunkData := range chunks {
		vao, trisCount := createChunkVAO(_chunkData.blocksData, chunkPos)
		chunks[chunkPos] = chunkData{
			blocksData: _chunkData.blocksData,
			vao:        vao,
			trisCount:  trisCount,
		}
	}
}

func loadTextureAtlas(textureFilePath string) uint32 {

	file, err := os.Open(textureFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var textureID uint32
	gl.GenTextures(1, &textureID)
	gl.BindTexture(gl.TEXTURE_2D, textureID)
	imageFile, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(imageFile.Bounds())
	draw.Draw(rgba, rgba.Bounds(), imageFile, image.Point{}, draw.Over)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Bounds().Dx()), int32(rgba.Bounds().Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	var maxAnisotropy int32
	gl.GetIntegerv(gl.MAX_TEXTURE_MAX_ANISOTROPY, &maxAnisotropy)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_ANISOTROPY, maxAnisotropy)
	return textureID
}
func getTextureCoords(blockID uint8, faceIndex uint8) []float32 {

	// Calculate UV coordinates
	u1 := float32(faceIndex*16) / float32(96)
	v1 := float32(blockID*16) / float32(48)
	u2 := float32((faceIndex+1)*16) / float32(96)
	v2 := float32((blockID+1)*16) / float32(48)

	return []float32{u1, v1, u2, v2}

}

/*
For each horizontal coordinate, iterate downwards until a block is hit or yLevel < limit.

If block hit, set light = 15.
If no block, check left right forward backward blocks. If block, light that one up as well.
*/
func propagateSunLight(chunkPos chunkPosition, blocksData map[blockPosition]blockData) {
	visited := make(map[blockPosition]bool)
	for x := uint8(0); x < 16; x++ {
		for z := uint8(0); z < 16; z++ {
			//height limit

			skyLight := uint8(15)

			for y := int16(128); y > -128; y-- {
				blockPos := blockPosition{x, y, z}
				if skyLight == 0 {
					break
				}
				if _, exists := blocksData[blockPos]; exists {
					if visited[blockPos] {
						continue
					}
					visited[blockPos] = true
					blocksData[blockPos] = blockData{
						blockType:  blocksData[blockPos].blockType,
						lightLevel: skyLight,
					}

					if _, exists := blocksData[blockPosition{x, y - 1, z}]; exists {

						skyLight--

					}

				} else {

					//left
					blockPos := blockPosition{x - 1, y, z}
					if x == 0 {
						//check neighbor chunk
						if _, exists := chunks[chunkPosition{chunkPos.x - 1, chunkPos.z}].blocksData[blockPosition{15, blockPos.y, blockPos.z}]; exists {

							chunks[chunkPosition{chunkPos.x - 1, chunkPos.z}].blocksData[blockPosition{15, blockPos.y, blockPos.z}] = blockData{
								blockType:  chunks[chunkPosition{chunkPos.x - 1, chunkPos.z}].blocksData[blockPosition{15, blockPos.y, blockPos.z}].blockType,
								lightLevel: skyLight,
							}

						}
					}

					if _, exists := blocksData[blockPos]; exists {
						if visited[blockPos] {
							continue
						}
						visited[blockPos] = true
						blocksData[blockPos] = blockData{
							blockType:  blocksData[blockPos].blockType,
							lightLevel: skyLight,
						}

					}

					//right
					blockPos = blockPosition{x + 1, y, z}
					if x == 15 {
						//check neighbor chunk
						if _, exists := chunks[chunkPosition{chunkPos.x + 1, chunkPos.z}].blocksData[blockPosition{0, y, z}]; exists {

							chunks[chunkPosition{chunkPos.x + 1, chunkPos.z}].blocksData[blockPosition{0, y, z}] = blockData{
								blockType:  chunks[chunkPosition{chunkPos.x + 1, chunkPos.z}].blocksData[blockPosition{0, y, z}].blockType,
								lightLevel: skyLight,
							}

						}
					}
					if _, exists := blocksData[blockPos]; exists {
						if visited[blockPos] {
							continue
						}
						visited[blockPos] = true
						blocksData[blockPos] = blockData{
							blockType:  blocksData[blockPos].blockType,
							lightLevel: skyLight,
						}

					}

					//forward
					blockPos = blockPosition{x, y, z + 1}
					if z == 15 {
						//check neighbor chunk
						if _, exists := chunks[chunkPosition{chunkPos.x, chunkPos.z + 1}].blocksData[blockPosition{x, y, 0}]; exists {

							chunks[chunkPosition{chunkPos.x, chunkPos.z + 1}].blocksData[blockPosition{x, y, 0}] = blockData{
								blockType:  chunks[chunkPosition{chunkPos.x, chunkPos.z + 1}].blocksData[blockPosition{x, y, 0}].blockType,
								lightLevel: skyLight,
							}

						}
					}
					if _, exists := blocksData[blockPos]; exists {
						if visited[blockPos] {
							continue
						}
						visited[blockPos] = true
						blocksData[blockPos] = blockData{
							blockType:  blocksData[blockPos].blockType,
							lightLevel: skyLight,
						}

					}
					//backward
					blockPos = blockPosition{x, y, z - 1}
					if z == 0 {
						if _, exists := chunks[chunkPosition{chunkPos.x, chunkPos.z - 1}].blocksData[blockPosition{x, y, 15}]; exists {

							chunks[chunkPosition{chunkPos.x, chunkPos.z - 1}].blocksData[blockPosition{x, y, 15}] = blockData{
								blockType:  chunks[chunkPosition{chunkPos.x, chunkPos.z - 1}].blocksData[blockPosition{x, y, 15}].blockType,
								lightLevel: skyLight,
							}

						}
					}
					if _, exists := blocksData[blockPos]; exists {
						if visited[blockPos] {
							continue
						}
						visited[blockPos] = true
						blocksData[blockPos] = blockData{
							blockType:  blocksData[blockPos].blockType,
							lightLevel: skyLight,
						}

					}

				}
			}
		}
	}

}
func propagateLight(chunkPos chunkPosition, blocksData map[blockPosition]blockData, startPos blockPosition, initialLight uint8) {
	/*
		if initialLight <= 4 {
			return
		}

		type queueEntry struct {
			pos        blockPosition
			lightLevel uint8
		}

		directions := []blockPosition{
			{1, 0, 0}, {-1, 0, 0}, // X-axis
			{0, 1, 0}, {0, -1, 0}, // Y-axis
			{0, 0, 1}, {0, 0, -1}, // Z-axis

		}

		queue := []queueEntry{{startPos, initialLight}}
		visited := make(map[blockPosition]bool)

		for len(queue) > 0 {
			current := queue[len(queue)-1]
			queue = queue[:len(queue)-1]

			if visited[current.pos] {
				continue
			}
			visited[current.pos] = true

			blocksData[current.pos] = blockData{
				blockType:  blocksData[current.pos].blockType,
				lightLevel: current.lightLevel,
			}

			for _, dir := range directions {
				neighborPos := blockPosition{
					x: current.pos.x + dir.x,
					y: current.pos.y + dir.y,
					z: current.pos.z + dir.z,
				}

				if neighborData, exists := blocksData[neighborPos]; exists {
					newLightLevel := uint8(float32(current.lightLevel) * 0.8)
					if newLightLevel > neighborData.lightLevel {
						queue = append(queue, queueEntry{neighborPos, newLightLevel})
					}
				} else {

							rebuildChunk(chunks[ChunkPos], ChunkPos)
							isBordering, borderingChunks := ReturnBorderingChunks(pos, ChunkPos)
							if isBordering {
								for i := range borderingChunks {
									rebuildChunk(chunks[borderingChunks[i]], borderingChunks[i])
								}
							}

						isBordering, borderingChunks := ReturnBorderingChunks(neighborPos, chunkPos)
						if isBordering {
							for i := range borderingChunks {
								propagateLight(borderingChunks[i], chunks[borderingChunks[i]].blocksData, neighborPos, uint8(float32(current.lightLevel)*0.8))

							}
						}

				}
			}

			}*/

}
func fractalNoise(x int32, z int32, amplitude float32, octaves int, lacunarity float32, persistence float32, scale float32) int16 {
	val := int16(0)
	x1 := float32(x)
	z1 := float32(z)

	for i := 0; i < octaves; i++ {
		val += int16(noise.Eval2(x1/scale, z1/scale) * amplitude)
		z1 *= lacunarity
		x1 *= lacunarity
		amplitude *= persistence
	}
	if val < -128 {
		return -128
	}
	if val > 128 {
		return 128
	}
	return val

}
func fractalNoise3D(x int32, y int32, z int32, amplitude float32, scale float32) float32 {
	val := float32(0)
	x1 := float32(x)
	y1 := float32(y)
	z1 := float32(z)

	val += noise.Eval3(x1/scale, y1/scale, z1/scale) * amplitude

	if val < -1 {
		return -1
	}
	if val > 1 {
		return 1
	}
	return val

}

func createChunkVAO(chunkData map[blockPosition]blockData, chunkPos chunkPosition) (uint32, uint32) {

	var chunkVertices []float32
	grassTint := mgl32.Vec3{0.486, 0.741, 0.419}
	noTint := mgl32.Vec3{1.0, 1.0, 1.0}
	for key := range chunkData {
		self := chunkData[blockPosition{key.x, key.y, key.z}]
		_, top := chunkData[blockPosition{key.x, key.y + 1, key.z}]
		_, bot := chunkData[blockPosition{key.x, key.y - 1, key.z}]
		_, l := chunkData[blockPosition{key.x - 1, key.y, key.z}]
		_, r := chunkData[blockPosition{key.x + 1, key.y, key.z}]
		_, b := chunkData[blockPosition{key.x, key.y, key.z - 1}]
		_, f := chunkData[blockPosition{key.x, key.y, key.z + 1}]

		//block touching blocks on each side, won't be visible
		if top && bot && l && r && b && f {
			continue
		}

		for i := 0; i < len(CubeVertices); i += 3 {
			curTint := noTint
			x := CubeVertices[i] + float32(key.x)
			y := CubeVertices[i+1] + float32(key.y)
			z := CubeVertices[i+2] + float32(key.z)
			uv := (i / 3) * 2
			var u, v uint8 = CubeUVs[uv], CubeUVs[uv+1]

			//FRONT FACE
			if i >= (0*18) && i <= (0*18)+15 {

				if !f {

					if key.z == 15 {
						////rowFront := col + 1
						//adjustedRow := (config.NumOfChunks * row)

						_, blockAdjChunk := chunks[chunkPosition{chunkPos.x, chunkPos.z + 1}].blocksData[blockPosition{key.x, key.y, 0}]
						if blockAdjChunk {
							continue
						}
					}
					textureUV := getTextureCoords(chunkData[key].blockType, 2)
					if self.blockType == GrassID {
						curTint = grassTint
						textureUVOverlay := getTextureCoords(chunkData[key].blockType, 5)
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.6, curTint[0], curTint[1], curTint[2], textureUVOverlay[u], textureUVOverlay[v])
					} else {

						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.6, curTint[0], curTint[1], curTint[2], 0, 0)
					}
				}
				continue
			}
			//BACK FACE
			if i >= (1*18) && i <= (1*18)+15 {

				if !b {
					if key.z == 0 {
						//rowFront := col - 1
						//adjustedRow := (config.NumOfChunks * row)
						_, blockAdjChunk := chunks[chunkPosition{chunkPos.x, chunkPos.z - 1}].blocksData[blockPosition{key.x, key.y, 15}]
						if blockAdjChunk {
							continue
						}
					}
					textureUV := getTextureCoords(chunkData[key].blockType, 3)
					if self.blockType == GrassID {
						curTint = grassTint
						textureUVOverlay := getTextureCoords(chunkData[key].blockType, 5)
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.6, curTint[0], curTint[1], curTint[2], textureUVOverlay[u], textureUVOverlay[v])
					} else {

						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.6, curTint[0], curTint[1], curTint[2], 0, 0)
					}
				}
				continue
			}
			//LEFT FACE
			if i >= (2*18) && i <= (2*18)+15 {
				if !l {
					if key.x == 0 {
						//rowFront := row - 1
						//adjustedRow := (config.NumOfChunks * rowFront)
						_, blockAdjChunk := chunks[chunkPosition{chunkPos.x - 1, chunkPos.z}].blocksData[blockPosition{15, key.y, key.z}]
						if blockAdjChunk {
							continue
						}
					}
					textureUV := getTextureCoords(chunkData[key].blockType, 4)
					if self.blockType == GrassID {
						curTint = grassTint
						textureUVOverlay := getTextureCoords(chunkData[key].blockType, 5)
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.8, curTint[0], curTint[1], curTint[2], textureUVOverlay[u], textureUVOverlay[v])
					} else {
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.8, curTint[0], curTint[1], curTint[2], 0, 0)
					}

				}
				continue
			}
			//RIGHT FACE
			if i >= (3*18) && i <= (3*18)+15 {

				if !r {
					if key.x == 15 {
						//rowFront := row + 1
						//adjustedRow := (config.NumOfChunks * rowFront)
						_, blockAdjChunk := chunks[chunkPosition{chunkPos.x + 1, chunkPos.z}].blocksData[blockPosition{0, key.y, key.z}]
						if blockAdjChunk {
							continue
						}
					}

					textureUV := getTextureCoords(chunkData[key].blockType, 5)

					if self.blockType == GrassID {
						curTint = grassTint
						textureUV := getTextureCoords(chunkData[key].blockType, 2)
						textureUVOverlay := getTextureCoords(chunkData[key].blockType, 5)
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.8, curTint[0], curTint[1], curTint[2], textureUVOverlay[u], textureUVOverlay[v])
					} else {
						chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.8, curTint[0], curTint[1], curTint[2], 0, 0)
					}
				}

				continue
			}
			//TOP FACE
			if i >= (4*18) && i <= (4*18)+15 {
				if !top {
					if self.blockType == GrassID {
						curTint = grassTint
					}
					textureUV := getTextureCoords(chunkData[key].blockType, 0)

					chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel), curTint[0], curTint[1], curTint[2], 0, 0)
				}
				continue
			}
			//BOTTOM FACE
			if i >= (5*18) && i <= (5*18)+15 {
				if !bot && key.y != -128 {
					textureUV := getTextureCoords(chunkData[key].blockType, 1)
					chunkVertices = append(chunkVertices, x, y, z, textureUV[u], textureUV[v], float32(self.lightLevel)*0.5, curTint[0], curTint[1], curTint[2], 0, 0)
				}
				continue
			}

		}
	}

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(chunkVertices), gl.Ptr(chunkVertices), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//position
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 11*4, nil)

	// Enable vertex attribute array for texture coordinates (location 1)
	gl.EnableVertexAttribArray(1)
	// Define the texture coordinate data layout: 2 components (u, v)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 11*4, uintptr(3*4))

	//light level
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointerWithOffset(2, 1, gl.FLOAT, false, 11*4, uintptr(5*4))

	//texture tint
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointerWithOffset(3, 3, gl.FLOAT, false, 11*4, uintptr(6*4))

	//overlay texture
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointerWithOffset(4, 2, gl.FLOAT, false, 11*4, uintptr(9*4))

	return vao, uint32(len(chunkVertices) / 5)
}
