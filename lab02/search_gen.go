package main

import (
	"fmt"
	"math/rand"
	"time"
)

const XSIZE = 70
const YSIZE = 50
const SIZE = YSIZE * XSIZE
const WALL = "#"
const START = "S"
const EXIT = "E"
const MONSTER = "M"
const TREASURE = "*"
const EMPTY = " "
const ROOM_GEN_ITERATIONS = 400
const EXTRA_DOORS_RATE = 0.02
const MAX_ROOM_SIZE = 12
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

type connector struct {
	x int
	y int
	region1 int
	region2 int
	horizontal bool
}

// Inspired by http://journal.stuffwithstuff.com/2014/12/21/rooms-and-mazes/
func main() {
	const WALL_RATIO = 0.36
	const MONSTER_RATIO = 0.02
	const TREASURE_RATIO = 0.02
	var grid [SIZE] string
	var regions [SIZE] int
	for i := 0; i < SIZE; i++ {
		grid[i] = WALL
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var currentRegionId = 1
	currentRegionId = generateRooms(&grid, &regions, r, currentRegionId)
	currentRegionId = fillWithMazes(&grid, &regions, r, currentRegionId)

	connect(&grid, &regions, r)

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
	fmt.Println("Tiles:")
	print(&grid)
	//fmt.Println("Regions:")
	//printInt(&regions)
}

func connect(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand) {
	edges := findConnectors(g, regions)
	spanningTree := NewIntSet()
	open(g, &edges[r.Intn(len(edges))], spanningTree)
	edges = validConnectors(edges, spanningTree)
	for len(edges) > 0 {
		unitedEdges := filter(edges, func(v connector) bool {
			return spanningTree.contains(v.region1) || spanningTree.contains(v.region2)
		})
		open(g, &unitedEdges[r.Intn(len(unitedEdges))], spanningTree)
		edges = validConnectors(edges, spanningTree)
	}
}

func validConnectors(edges []connector, spanningTree *IntSet) []connector {
	return filter(edges, func(v connector) bool {
		return !(spanningTree.contains(v.region1) && spanningTree.contains(v.region2))
	})
}

func open(g *[SIZE]string, edge *connector, spanningTree *IntSet) {
	g[position1d(edge.x, edge.y)] = EMPTY
	spanningTree.add(edge.region1)
	spanningTree.add(edge.region2)
}

func findConnectors(g *[SIZE]string, regions *[SIZE]int) [] connector {
	var edges [SIZE] connector
	edgesIndex := 0
	for y := 0; y < YSIZE; y++ {
		for x := 0; x < XSIZE; x++ {
			edge := connectable(g, regions, x, y)
			if edge != nil {
				edges[edgesIndex] = *edge
				edgesIndex = edgesIndex + 1
			}
		}
	}
	return edges[:edgesIndex]
}

func connectable(g *[SIZE]string, regions *[SIZE]int, x int, y int) *connector {
	if x < 1 || x >= XSIZE - 1 || y < 1 || y >= YSIZE - 1 {
		return nil
	}
	left := position1d(x - 1, y)
	right := position1d(x + 1, y)
	up := position1d(x, y - 1)
	down := position1d(x, y + 1)
	if g[position1d(x, y)] != WALL {
		return nil
	}
	if regions[left] != 0 && regions[right] != 0 && regions[left] != regions[right] {
		return &connector{x, y, regions[left], regions[right], true}
	}
	if regions[up] != 0 && regions[down] != 0 && regions[up] != regions[down] {
		return &connector{x, y, regions[up], regions[down], false}
	}
	return nil
}

func generateRooms(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand, currentRegionId int) int {
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
					regions[position1d(x, y)] = currentRegionId
				}
			}
			currentRegionId = currentRegionId + 1
		}
	}
	return currentRegionId
}

// Depth-first search https://en.wikipedia.org/wiki/Maze_generation_algorithm
func fillWithMazes(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand, currentRegionId int) int {
	for y := 1; y < YSIZE - 1; y++ {
		for x := 1; x < XSIZE - 1; x++ {
			if generateMaze(g, regions, r, x, y, NO_DIR, currentRegionId) {
				currentRegionId = currentRegionId + 1
			}
		}
	}
	return currentRegionId
}

func generateMaze(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand, x int, y int, dir int, currentRegionId int) bool {
	if x < 1 || x >= XSIZE - 1 || y < 1 || y >= YSIZE - 1 || !possibleDirection(g, x, y, dir) {
		return false
	}
	position := position1d(x, y)
	g[position] = EMPTY
	regions[position] = currentRegionId
	var unvisited = toUnvisitedDirections(dir)
	for len(unvisited) > 0 {
		var i = r.Intn(len(unvisited))
		if unvisited[i] == LEFT {
			generateMaze(g, regions, r, x - 1, y, LEFT, currentRegionId)
		} else if unvisited[i] == UP {
			generateMaze(g, regions, r, x, y - 1, UP, currentRegionId)
		} else if unvisited[i] == RIGHT {
			generateMaze(g, regions, r, x + 1, y, RIGHT, currentRegionId)
		} else if unvisited[i] == DOWN {
			generateMaze(g, regions, r, x, y + 1, DOWN, currentRegionId)
		}
		unvisited = remove(unvisited, i)
	}
	return true
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

func printInt(g *[SIZE]int) {
	for y := 0; y < YSIZE; y++ {
		for x := 0; x < XSIZE; x++ {
			fmt.Printf("%d ", g[position1d(x, y)])
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

func filter(vs []connector, f func(connector) bool) []connector {
	result := make([]connector, 0)
	for _, v := range vs {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// Oh this fudging language, you have to do everything yourself
// https://play.golang.org/p/tDdutH672-
type IntSet struct {
	set map[int]bool
}

func NewIntSet() *IntSet {
	return &IntSet{make(map[int]bool)}
}

func (set *IntSet) add(i int) bool {
	_, found := set.set[i]
	set.set[i] = true
	return !found	//False if it existed already
}

func (set *IntSet) contains(i int) bool {
	_, found := set.set[i]
	return found	//true if it existed already
}