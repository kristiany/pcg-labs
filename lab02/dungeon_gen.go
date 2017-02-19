package main

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/kristiany/pcg-labs/lab02/utils"
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
const EXTRA_DOORS_ONE_IN = 50
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
	const MONSTER_RATIO = 0.02
	const TREASURE_RATIO = 0.03
	var grid [SIZE] string
	var regions [SIZE] int
	for i := 0; i < SIZE; i++ {
		grid[i] = WALL
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var currentRegionId = 1
	var rooms [] room
	currentRegionId, rooms = generateRooms(&grid, &regions, r, currentRegionId)
	currentRegionId = fillWithMazes(&grid, &regions, r, currentRegionId)
	connect(&grid, &regions, r)
	sparsify(&grid)
	retryingRandomAdd(&grid, rooms, START, r)
	retryingRandomAdd(&grid, rooms, EXIT, r)
	for i := 0; i < MONSTER_RATIO * SIZE; i++ {
		retryingRandomAdd(&grid, rooms, MONSTER, r)
	}
	for i := 0; i < TREASURE_RATIO * SIZE; i++ {
		retryingRandomAdd(&grid, rooms, TREASURE, r)
	}
	fmt.Println("Tiles:")
	print(&grid)
}
func sparsify(g *[SIZE]string) {
	for y := 1; y < YSIZE - 1; y++ {
		for x := 1; x < XSIZE - 1; x++ {
			removeDeadend(g, x, y)
		}
	}
}

func removeDeadend(g *[SIZE]string, x int, y int) {
	if g[position1d(x, y)] != EMPTY || x < 1 || x >= XSIZE - 1 || y < 1 || y >= YSIZE - 1 {
		return
	}
	closeableDirection := closeable(g, x, y)
	if closeableDirection == NO_DIR {
		return
	}
	g[position1d(x, y)] = WALL
	if closeableDirection == LEFT {
		removeDeadend(g, x - 1, y)
	}
	if closeableDirection == RIGHT {
		removeDeadend(g, x + 1, y)
	}
	if closeableDirection == UP {
		removeDeadend(g, x, y - 1)
	}
	if closeableDirection == DOWN {
		removeDeadend(g, x, y + 1)
	}
}
func closeable(g *[SIZE]string, x int, y int) int {
	if g[position1d(x - 1, y)] == EMPTY &&
			g[position1d(x, y - 1)] == WALL &&
			g[position1d(x + 1, y)] == WALL &&
			g[position1d(x, y + 1)] == WALL {
		return LEFT
	}
	if g[position1d(x + 1, y)] == EMPTY &&
			g[position1d(x, y - 1)] == WALL &&
			g[position1d(x - 1, y)] == WALL &&
			g[position1d(x, y + 1)] == WALL {
		return RIGHT
	}
	if g[position1d(x, y - 1)] == EMPTY &&
			g[position1d(x + 1, y)] == WALL &&
			g[position1d(x, y + 1)] == WALL &&
			g[position1d(x - 1, y)] == WALL {
		return UP
	}
	if g[position1d(x, y + 1)] == EMPTY &&
			g[position1d(x + 1, y)] == WALL &&
			g[position1d(x, y - 1)] == WALL &&
			g[position1d(x - 1, y)] == WALL {
		return DOWN
	}
	return NO_DIR
}

func connect(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand) {
	edges := findConnectors(g, regions)
	spanningTree := utils.NewIntSet()
	open(g, &edges[r.Intn(len(edges))], spanningTree)
	doors := 1
	edges = validConnectors(edges, spanningTree)
	for len(edges) > 0 {
		unitedEdges := filter(edges, func(v connector) bool {
			return spanningTree.Contains(v.region1) || spanningTree.Contains(v.region2)
		})
		open(g, &unitedEdges[r.Intn(len(unitedEdges))], spanningTree)
		doors = doors + 1
		edges = validConnectors(edges, spanningTree)
	}
	// Extra doors
	edges = findConnectors(g, regions)
	spanningTree = utils.NewIntSet()
	max := int(1.0 / EXTRA_DOORS_ONE_IN * float64(doors))
	for i := 0; i < max; i++ {
		open(g, &edges[r.Intn(len(edges))], spanningTree)
	}
}

func validConnectors(edges []connector, spanningTree *utils.IntSet) []connector {
	return filter(edges, func(v connector) bool {
		return !(spanningTree.Contains(v.region1) && spanningTree.Contains(v.region2))
	})
}

func open(g *[SIZE]string, edge *connector, spanningTree *utils.IntSet) {
	g[position1d(edge.x, edge.y)] = EMPTY
	spanningTree.Add(edge.region1)
	spanningTree.Add(edge.region2)
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

func generateRooms(g *[SIZE]string, regions *[SIZE]int, r *rand.Rand, currentRegionId int) (int, []room) {
	rooms := make([]room, 0)
	for i := 0; i < ROOM_GEN_ITERATIONS; i++ {
		var x = r.Intn(XSIZE - MIN_ROOM_SIZE)
		var y = r.Intn(YSIZE - MIN_ROOM_SIZE)
		var width = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE + 1) + MIN_ROOM_SIZE
		var height = r.Intn(MAX_ROOM_SIZE - MIN_ROOM_SIZE + 1) + MIN_ROOM_SIZE
		var random = room{x, y, utils.Min(width, XSIZE - x - 1), utils.Min(height, YSIZE - y - 1)}
		if(areaFree(g, random)) {
			rooms = append(rooms, random)
			for y := random.y + 1; y < random.y + random.height && y < YSIZE; y++ {
				for x := random.x + 1; x < random.x + random.width && x < XSIZE; x++ {
					g[position1d(x, y)] = EMPTY
					regions[position1d(x, y)] = currentRegionId
				}
			}
			currentRegionId = currentRegionId + 1
		}
	}
	return currentRegionId, rooms
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
	for y := room.y; y <= room.y + room.height && y < YSIZE; y++ {
		for x := room.x; x <= room.x + room.width && x < XSIZE; x++ {
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

func retryingRandomAdd(g *[SIZE]string, rooms [] room, value string, r *rand.Rand) {
	for i := 0; !addSafe(g, randomRoomPosition(r, rooms[r.Intn(len(rooms))]), value) && i < 100; i++ {
		//fmt.Println("Index taking, trying another")
	}
}

func addSafe(g *[SIZE]string, i int, value string) bool {
	if(g[i] == EMPTY) {
		g[i] = value
		fmt.Printf("Placing %s at %d\n", value, i)
		return true
	}
	return false
}

func randomRoomPosition(r *rand.Rand, room room) int {
	roomx := r.Intn(room.width) + 1
	roomy := r.Intn(room.height) + 1
	fmt.Printf("Trying random room position (%d, %d) at position (%d, %d)\n", roomx, roomy, room.x, room.y)
	return position1d(room.x + roomx, room.y + roomy)
}

func position1d(x int, y int) int {
	return y * XSIZE + x
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