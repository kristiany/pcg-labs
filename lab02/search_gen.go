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
const MAX_ROOM_SIZE = 11
const MIN_ROOM_SIZE = 4
const LEFT = 0
const UP = 1
const RIGHT = 2
const DOWN = 3
const NO_DIR = -1

type room struct {
	x int
	y int
	width int
	height int
}

type pos struct {
	x int
	y int
}

// Inspired by http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/
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
	fillWithMazes(&grid, r)


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
		var width = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE + 1) + MIN_ROOM_SIZE
		var height = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE + 1) + MIN_ROOM_SIZE
		var random = room{x, y, min(width, XSIZE - x), min(height, YSIZE - y)}
		if(areaFree(g, random)) {
			for y := random.y + 1; y < random.y + random.height - 1 && y < YSIZE; y++ {
				for x := random.x + 1; x < random.x + random.width - 1 && x < XSIZE; x++ {
					g[position1d(x, y)] = EMPTY
				}
			}
		}
	}
}

// Depth-first search https://en.wikipedia.org/wiki/Maze_generation_algorithm
func fillWithMazes(g *[SIZE]string, r *rand.Rand) {
	for y := 1; y < YSIZE - 1; y++ {
		for x := 1; x < XSIZE - 1; x++ {
			generateMaze(g, r, x, y, NO_DIR)
		}
	}
}

func generateMaze(g *[SIZE]string, r *rand.Rand, x int, y int, dir int) {
	if x < 1 || x >= XSIZE - 1 || y < 1 || y >= YSIZE - 1 || !possibleDirection(g, x, y, dir) {
		return
	}
	position := position1d(x, y)
	g[position] = EMPTY
	var unvisited = toUnvisitedDirections(dir)
	for len(unvisited) > 0 {
		var i = r.Intn(len(unvisited))
		if unvisited[i] == LEFT {
			generateMaze(g, r, x - 1, y, LEFT)
		} else if unvisited[i] == UP {
			generateMaze(g, r, x, y - 1, UP)
		} else if unvisited[i] == RIGHT {
			generateMaze(g, r, x + 1, y, RIGHT)
		} else if unvisited[i] == DOWN {
			generateMaze(g, r, x, y + 1, DOWN)
		}
		unvisited = remove(unvisited, i)
	}
}

func toUnvisitedDirections(dir int) []int {
	if dir == LEFT {
		return []int {LEFT, UP, DOWN}
	}
	if dir == UP {
		return []int {LEFT, UP, RIGHT}
	}
	if dir == RIGHT {
		return []int {UP, RIGHT, DOWN}
	}
	if dir == DOWN {
		return []int {LEFT, RIGHT, DOWN}
	}
	return []int {LEFT, UP, RIGHT, DOWN}
}

func remove(a []int, i int) []int {
	return append(a[:i], a[i + 1:]...)
}

func openableSpace(g *[SIZE]string, x int, y int) bool {
	return g[position1d(x - 1, y - 1)] == WALL &&
			g[position1d(x, y - 1)] == WALL &&
			g[position1d(x + 1, y - 1)] == WALL &&
			g[position1d(x - 1, y)] == WALL &&
			g[position1d(x, y)] == WALL &&
			g[position1d(x + 1, y)] == WALL &&
			g[position1d(x - 1, y + 1)] == WALL &&
			g[position1d(x, y + 1)] == WALL &&
			g[position1d(x + 1, y + 1)] == WALL
}

func possibleDirection(g *[SIZE]string, x int, y int, dir int) bool {
	if dir == LEFT {
		return g[position1d(x, y)] == WALL &&
				g[position1d(x, y - 1)] == WALL &&
				g[position1d(x - 1, y - 1)] == WALL &&
				g[position1d(x - 1, y)] == WALL &&
				g[position1d(x - 1, y + 1)] == WALL &&
				g[position1d(x, y + 1)] == WALL
	}
	if dir == UP {
		return g[position1d(x, y)] == WALL &&
				g[position1d(x - 1, y)] == WALL &&
				g[position1d(x - 1, y - 1)] == WALL &&
				g[position1d(x, y - 1)] == WALL &&
				g[position1d(x + 1, y - 1)] == WALL &&
				g[position1d(x + 1, y)] == WALL
	}
	if dir == RIGHT {
		return g[position1d(x, y)] == WALL &&
				g[position1d(x, y - 1)] == WALL &&
				g[position1d(x + 1, y - 1)] == WALL &&
				g[position1d(x + 1, y)] == WALL &&
				g[position1d(x + 1, y + 1)] == WALL &&
				g[position1d(x, y + 1)] == WALL
	}
	if dir == DOWN {
		return g[position1d(x, y)] == WALL &&
				g[position1d(x - 1, y)] == WALL &&
				g[position1d(x - 1, y + 1)] == WALL &&
				g[position1d(x, y + 1)] == WALL &&
				g[position1d(x + 1, y + 1)] == WALL &&
				g[position1d(x + 1, y)] == WALL
	}
	return openableSpace(g, x, y)
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
		for x := 0; x < XSIZE; x++ {
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}