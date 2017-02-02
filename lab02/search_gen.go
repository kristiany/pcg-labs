package main

import (
	"fmt"
	"math/rand"
	"time"
)

const XSIZE = 50
const YSIZE = 50
const WALL = "#"
const START = "S"
const EXIT = "E"
const MONSTER = "M"
const TREASURE = "*"
const EMPTY = "_"

func main() {
	const WALL_RATIO = 0.46
	const MONSTER_RATIO = 0.02
	const TREASURE_RATIO = 0.02
	var grid [YSIZE * XSIZE] string
	for i := 0; i < YSIZE * XSIZE; i++ {
		grid[i] = EMPTY
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	retryingRandomAdd(&grid, START, r)
	retryingRandomAdd(&grid, EXIT, r)
	for i := 0; i < WALL_RATIO * XSIZE * YSIZE; i++ {
		retryingRandomAdd(&grid, WALL, r)
	}
	for i := 0; i < MONSTER_RATIO * XSIZE * YSIZE; i++ {
		retryingRandomAdd(&grid, MONSTER, r)
	}
	for i := 0; i < TREASURE_RATIO * XSIZE * YSIZE; i++ {
		retryingRandomAdd(&grid, TREASURE, r)
	}
	print(&grid)
}

func print(g *[YSIZE * XSIZE]string) {
	for y := 0; y < YSIZE; y++ {
		for x := 0; x < YSIZE; x++ {
			fmt.Printf(g[position1d(x, y)] + " ")
		}
		fmt.Printf("\n")
	}
}

func retryingRandomAdd(g *[YSIZE * XSIZE]string, value string, r *rand.Rand) {
	for !addSafe(g, randomPosition(r), value) {
		fmt.Printf("Index taking, trying another")
	}
}

func addSafe(g *[YSIZE * XSIZE]string, i int, value string) bool {
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
