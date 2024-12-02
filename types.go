package main

import "github.com/go-gl/mathgl/mgl32"

type blockPosition struct {
	x uint8
	y int16
	z uint8
}
type blockPositionHoriz struct {
	x uint8
	z uint8
}

type chunkPosition struct {
	x int32
	z int32
}

type blockData struct {
	blockType  uint8
	lightLevel uint8
}

type chunkData struct {
	blocksData map[blockPosition]blockData
	vao        uint32
	trisCount  uint32
}

func ReturnBorderingChunks(pos blockPosition, chunkPos chunkPosition) (bool, []chunkPosition) {

	var borderingChunks []chunkPosition

	if _, exists := chunks[chunkPosition{chunkPos.x + 1, chunkPos.z}]; exists {

		if pos.x == 15 {
			borderingChunks = append(borderingChunks, chunkPosition{chunkPos.x + 1, chunkPos.z})
		}
	}
	if _, exists := chunks[chunkPosition{chunkPos.x - 1, chunkPos.z}]; exists {
		if pos.x == 0 {
			borderingChunks = append(borderingChunks, chunkPosition{chunkPos.x - 1, chunkPos.z})
		}
	}
	if _, exists := chunks[chunkPosition{chunkPos.x, chunkPos.z + 1}]; exists {
		if pos.z == 15 {
			borderingChunks = append(borderingChunks, chunkPosition{chunkPos.x, chunkPos.z + 1})
		}
	}
	if _, exists := chunks[chunkPosition{chunkPos.x, chunkPos.z - 1}]; exists {
		if pos.z == 0 {
			borderingChunks = append(borderingChunks, chunkPosition{chunkPos.x, chunkPos.z - 1})
		}
	}
	if len(borderingChunks) > 0 {
		return true, borderingChunks
	}
	return false, borderingChunks

}

func chunk(pos chunkPosition) chunkData {
	var blocksData map[blockPosition]blockData = make(map[blockPosition]blockData)
	//blocks that are touching sky
	//var exposedBlocks map[blockPositionHoriz]int16 = make(map[blockPositionHoriz]int16)
	var scale float32 = 100 // Adjust as needed for terrain detail
	var amplitude float32 = 30
	for x := uint8(0); x < 16; x++ {

		for z := uint8(0); z < 16; z++ {

			noiseValue := fractalNoise(int32(x)+(pos.x*16), int32(z)+(pos.z*16), amplitude, 4, 1.5, 0.5, scale)
			//maxValue := noiseValue
			for y := noiseValue; y >= int16(-128); y-- {

				//determine block type
				blockType := DirtID
				fluctuation := int16(random.Float32() * 5)

				if y < ((noiseValue - 6) + fluctuation) {
					blockType = DirtID
				}
				if y < ((noiseValue - 10) + fluctuation) {
					blockType = StoneID
				}

				//top most layer
				if y == noiseValue {
					blocksData[blockPosition{x, y, z}] = blockData{
						blockType:  GrassID,
						lightLevel: 0,
					}
				} else {
					blocksData[blockPosition{x, y, z}] = blockData{
						blockType:  blockType,
						lightLevel: 0,
					}
				}

				if y < 0 {
					isCave := fractalNoise3D(int32(x)+(pos.x*16), int32(y), int32(z)+(pos.z*16), 0.7, 8)

					if isCave > 0.1 {
						delete(blocksData, blockPosition{x, y, z})
						/*
							if y == maxValue {
								maxValue = y - 1
							}
						*/
					}
				}

			}
			/*
				if block, exists := blocksData[blockPosition{x, maxValue, z}]; exists {
					block.lightLevel = 15

					blocksData[blockPosition{x, maxValue, z}] = block
					exposedBlocks[blockPositionHoriz{x, z}] = maxValue

				}
			*/
		}
	}
	/*
		for blockXZ, blockY := range exposedBlocks {
			propagateSunLight(pos, blocksData, blockPosition{blockXZ.x, blockY, blockXZ.z})
		}
	*/

	return chunkData{
		blocksData: blocksData,
		vao:        0,
		trisCount:  0,
	}
}

type aabb struct {
	Min, Max mgl32.Vec3
}

func AABB(min, max mgl32.Vec3) aabb {
	return aabb{Min: min, Max: max}
}
func Intersects(a, b aabb) bool {
	return (a.Min.X() <= b.Max.X() && a.Max.X() >= b.Min.X()) &&
		(a.Min.Y() <= b.Max.Y() && a.Max.Y() >= b.Min.Y()) &&
		(a.Min.Z() <= b.Max.Z() && a.Max.Z() >= b.Min.Z())
}

type text struct {
	VAO      uint32
	Texture  uint32
	Position mgl32.Vec2
	Update   bool
	FontSize float64
	Content  interface{}
}
type collider struct {
	Time   float32
	Normal []int
}
