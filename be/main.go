package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Point definition
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Request body with start and end points
type PathRequest struct {
	Start Point `json:"start"`
	End   Point `json:"end"`
}

// Response body containing the array of points after path calculations
type PathResponse struct {
	Path []Point `json:"result_path"`
}

func FindPathHandler(w http.ResponseWriter, r *http.Request) {
	var req PathRequest
	fmt.Printf("Got request for find-path\n")
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}
	fmt.Printf("start point : %v, end point : %v\n", req.Start, req.End)

	// Call findPath function
	path := findPath(req.Start, req.End)
	res := PathResponse{Path: path}

	json.NewEncoder(w).Encode(res)
}

func findPath(start, end Point) []Point {
	// Use a 4x4 grid for a 4x4 matrix
	const gridSize = 4

	visited := make([][]bool, gridSize)
	for i := range visited {
		visited[i] = make([]bool, gridSize)
	}

	var path []Point

	// If start and end are the same, return the start point
	if start == end {
		fmt.Println("Start and end points are the same.")
		return []Point{start}
	}

	fmt.Println("Starting DFS from start to end:", start, end)
	// Call DFS to find the path
	if dfs(start, end, visited, &path, gridSize) {
		fmt.Println("Path found: ", path)
		return path
	}
	fmt.Println("No Path found.")

	return []Point{}
}

func dfs(current, end Point, visited [][]bool, path *[]Point, gridSize int) bool {
	fmt.Printf("Visiting point : (%d,%d)\n", current.X, current.Y)

	// Base case: If the current point is the end point, add to path and return true
	if current.X == end.X && current.Y == end.Y {
		fmt.Printf("Reached the endpoint: (%d,%d)\n", end.X, end.Y)
		*path = append(*path, current) // Append end point to the path
		return true
	}

	visited[current.X][current.Y] = true
	fmt.Printf("Marking (%d,%d) as visited.\n", current.X, current.Y)

	// Define possible directions (up, down, left, right)
	directions := []Point{
		{X: 0, Y: 1},  //right
		{X: 1, Y: 0},  //down
		{X: 0, Y: -1}, //left
		{X: -1, Y: 0}, //up
	}

	// Explore each direction
	for _, dir := range directions {
		next := Point{X: current.X + dir.X, Y: current.Y + dir.Y}

		if isValid(next, gridSize, gridSize) && !visited[next.X][next.Y] {
			fmt.Printf("Moving to next point: (%d,%d)\n", next.X, next.Y)

			if dfs(next, end, visited, path, gridSize) {

				fmt.Printf("Appending point: (%d,%d) to the path.\n", current.X, current.Y)
				*path = append([]Point{current}, *path...) // Prepend current point
				return true
			}
		} else {
			fmt.Printf("Skipping point: (%d,%d) - Out of bounds or visited.\n", next.X, next.Y)
		}
	}

	visited[current.X][current.Y] = false
	fmt.Printf("Backtracking from point: (%d,%d).\n", current.X, current.Y)
	return false
}

func isValid(point Point, maxX int, maxY int) bool {
	return point.X >= 0 && point.X < maxX && point.Y >= 0 && point.Y < maxY
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allow specific HTTP methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // Allow specific headers

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	http.Handle("/findpath", enableCORS(http.HandlerFunc(FindPathHandler)))

	fmt.Println("Server running on port 3333")
	http.ListenAndServe(":3333", nil)
}

