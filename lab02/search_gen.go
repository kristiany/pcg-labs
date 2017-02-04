package main

import (
	"fmt"
	"math/rand"
	"time"
)

const XSIZE = 50
const YSIZE = 50
const SIZE = YSIZE * XSIZE
const WALL = "#"
const START = "S"
const EXIT = "E"
const MONSTER = "M"
const TREASURE = "*"
const EMPTY = " "
const ROOM_GEN_ITERATIONS = 800
const MAX_ROOM_SIZE = 15
const MIN_ROOM_SIZE = 3

type room struct {
	x int
	y int
	width int
	height int
}

func main() {
	const WALL_RATIO = 0.36
	const MONSTER_RATIO = 0.02
	const TREASURE_RATIO = 0.02
	var grid [SIZE] string
	for i := 0; i < SIZE; i++ {
		grid[i] = WALL
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	generateRooms(&grid, r)
	//retryingRandomAdd(&grid, START, r)
	//retryingRandomAdd(&grid, EXIT, r)
	/*for i := 0; i < WALL_RATIO * SIZE; i++ {
		retryingRandomAdd(&grid, WALL, r)
	}*/
	/*for i := 0; i < MONSTER_RATIO * SIZE; i++ {
		retryingRandomAdd(&grid, MONSTER, r)
	}
	for i := 0; i < TREASURE_RATIO * SIZE; i++ {
		retryingRandomAdd(&grid, TREASURE, r)
	}*/
	print(&grid)
}

func generateRooms(g *[SIZE]string, r *rand.Rand) {
	for i := 0; i < ROOM_GEN_ITERATIONS; i++ {
		var x = r.Intn(XSIZE)
		var y = r.Intn(YSIZE)
		var width = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE) + MIN_ROOM_SIZE
		var height = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE) + MIN_ROOM_SIZE
		var random = room{x, y, Min(width, XSIZE - x), Min(height, YSIZE - y)}
		if(areaFree(g, random)) {
			for y := random.y + 1; y < random.y + random.height - 1 && y < YSIZE; y++ {
				for x := random.x + 1; x < random.x + random.width - 1 && x < XSIZE; x++ {
					g[position1d(x, y)] = EMPTY
				}
			}
		}
	}
}

func areaFree(g *[SIZE]string, room room) bool {
	for y := room.y; y < room.y + room.height && y < YSIZE; y++ {
		for x := room.x; x < room.x + room.width && x < XSIZE; x++ {
			if(g[position1d(x, y)] != WALL) {
				return false;
			}
		}
	}
	return true
}

func print(g *[SIZE]string) {
	for y := 0; y < YSIZE; y++ {
		for x := 0; x < YSIZE; x++ {
			fmt.Printf(g[position1d(x, y)] + " ")
		}
		fmt.Printf("\n")
	}
}

func retryingRandomAdd(g *[SIZE]string, value string, r *rand.Rand) {
	for i := 0; !addSafe(g, randomPosition(r), value) && i < 200; i++ {
		//fmt.Println("Index taking, trying another")
	}
}

func addSafe(g *[SIZE]string, i int, value string) bool {
	if(g[i] == EMPTY) {
		g[i] = value
		return true
	}
	return false
}

func randomPosition(r *rand.Rand) int {
	return position1d(r.Intn(XSIZE), r.Intn(YSIZE))
}

func position1d(x int, y int) int {
	return y * XSIZE + x
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}