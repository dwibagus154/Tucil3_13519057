package graph

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// GRAPH
// Atribut Graph
type Graph struct {
	adjacencyMatrix [][]float64
	totalNodes      int
	nodes           []Info
}

// Getter Graph
func (graf *Graph) GetTotalNodes() int {
	return graf.totalNodes
}

// mendapatkan matriks ketetanggaan yang berisi jarak antara A dan B
func (graf *Graph) GetDistance(A string, B string) float64 {
	idx1, idx2 := Search(A, B, graf.nodes)
	if idx1 == -1 || idx2 == -1 {
		log.Fatalf("Error : Nama Tidak Ditemukan!")
	}
	return graf.adjacencyMatrix[idx1][idx2]
}

// mendapatkan nodes yang ada
func (graf *Graph) GetNodes() []Info {
	return graf.nodes
}

// mendapatkan nodes yang indexnya ada disebuah string visited
func (graf *Graph) GetNodeswithName(visited []string, distance float64) []Info {
	answer := make([]Info, len(visited)+1)

	for i := 0; i < len(visited); i++ {
		for j := 0; j < graf.totalNodes; j++ {
			if visited[i] == graf.nodes[j].Name {
				answer[i] = graf.nodes[j]
				break
			}
		}
	}
	answer[len(visited)] = Info{distance, 0, "a"}
	return answer
}

// mendapatkan index
func (graf *Graph) GetIndex(search string) int {
	for i := 0; i < graf.totalNodes; i++ {
		if graf.nodes[i].Name == search {
			return i
		}
	}
	return 0
}

// mencari index dari 2 nodes di array nodes
func Search(A string, B string, nodes []Info) (int, int) {
	idx1 := -1
	idx2 := -1
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Name == A {
			idx1 = i
		}
		if nodes[i].Name == B {
			idx2 = i
		}
	}
	return idx1, idx2
}

// Print List Nodes
func PrintListNodes(nodes []Info) {
	for _, v := range nodes {
		fmt.Println("→" + v.Name)
	}
}

// Info Nodes
// Atribut Info
type Info struct {
	Latitude  float64
	Longitude float64
	Name      string
}

// Getter Info
func (info Info) GetLatitude() float64 {
	return info.Latitude
}
func (info Info) GetLongitude() float64 {
	return info.Longitude
}
func (info Info) GetName() string {
	return info.Name
}

// Item yang membantu jalannya program A*
type Item struct {
	current string
	goal    string
	gn      float64
	fn      float64
	visited []string
	index   int
}

// Mendefinisikan tipe data dari Item yang brnama PriorityQueue
type PriorityQueue []*Item

// Mendapatkan panjang array
func (pq PriorityQueue) Len() int { return len(pq) }

//membuat operasi perbandingan antar index dan mereturn true apabila fn lebih kecil
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].fn < pq[j].fn
}

// swap index
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// memasukkan ke dalam array
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

// menghapus data di dalam array object
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Error input
func Check(fileName string, e error) {
	if e != nil {
		log.Fatalf(fileName + ": The system cannot find the file specified.")
	}
}

// Euclidan Distance
// mengubah derajat ke radian
func degreesToRadians(d float64) float64 {
	return d * math.Pi / 180
}

// menggunakan haversine formula untuk mendapatkan euclidean distance
func GetEuclidanDistance(A Info, B Info) float64 {
	// Coordinate A in radians
	latitude1 := degreesToRadians(A.Latitude)
	longitude1 := degreesToRadians(A.Longitude)

	// Coordinate B in radians
	latitude2 := degreesToRadians(B.Latitude)
	longitude2 := degreesToRadians(B.Longitude)

	// Haversine Formula
	differencelongitude := longitude2 - longitude1
	differencelatitude := latitude2 - latitude1

	a := math.Pow(math.Sin(differencelatitude/2), 2) + math.Cos(latitude1)*math.Cos(latitude2)*math.Pow(math.Sin(differencelongitude/2), 2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// earth radius in KM
	km := c * 6371
	return km
}

// Read External File
func ReadFile(fileName string) *Graph {
	graf := new(Graph)
	fileName = "../../test/" + fileName
	file, err := os.Open(fileName)
	Check(fileName, err)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	// read total nodes
	scanner.Scan()
	n, _ := strconv.Atoi(scanner.Text())
	graf.totalNodes = n
	// read detail nodes
	graf.nodes = make([]Info, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		line := strings.Fields(scanner.Text())
		simpang := line[0]
		for j := 1; j < len(line)-2; j++ {
			simpang = simpang + " " + line[j]
		}
		latitude, _ := strconv.ParseFloat(line[len(line)-2], 64)
		longitude, _ := strconv.ParseFloat(line[len(line)-1], 64)
		graf.nodes[i] = Info{latitude, longitude, simpang}
	}

	// Read & Split adjacency matrix
	graf.adjacencyMatrix = make([][]float64, n)
	for i := 0; i < n; i++ {
		graf.adjacencyMatrix[i] = make([]float64, n)
		scanner.Scan()
		line := strings.Fields(scanner.Text())
		for j := 0; j < n; j++ {
			if line[j] == "0" {
				graf.adjacencyMatrix[i][j] = 0
			} else {
				a := graf.nodes[i]
				b := graf.nodes[j]
				// convert array to real euclidean distance
				graf.adjacencyMatrix[i][j] = GetEuclidanDistance(a, b)
			}
		}
	}
	return graf
}

// A* Algorithm

// cek apabila sudah pernah dikunjungi
func isVisited(visited []string, idx int, nodes []Info) bool {
	for i := 0; i < len(visited); i++ {
		if visited[i] == nodes[idx].Name {
			return true
		}
	}
	return false
}

//fungsi Astar
func (graf *Graph) Astar(A string, B string) (float64, []string) {
	pq := make(PriorityQueue, 1)
	pq[0] = &Item{
		current: A,
		goal:    B,
		gn:      0,
		fn:      graf.GetDistance(A, B),
		visited: []string{A},
		index:   0,
	}
	heap.Init(&pq)
	// while pq is not empty
	for pq.Len() > 0 {
		now := heap.Pop(&pq).(*Item)
		// apabila node sekarang == tujuan maka return
		fmt.Println(now.visited)
		if now.current == now.goal {
			return now.fn, now.visited
		}
		a := graf.GetIndex(now.current)

		// mengecek semua tetangga yang belum pernah dikunjungi
		for i := 0; i < graf.totalNodes; i++ {
			if graf.adjacencyMatrix[a][i] > 0 && !isVisited(now.visited, i, graf.nodes) {
				updategn := now.gn + graf.adjacencyMatrix[a][i]
				updatehn := graf.GetDistance(graf.nodes[i].Name, now.goal)
				updatefn := updategn + updatehn
				updatevisited := append(now.visited, graf.nodes[i].Name)
				item := &Item{
					current: graf.nodes[i].Name,
					goal:    now.goal,
					gn:      updategn,
					fn:      updatefn,
					visited: updatevisited,
				}
				heap.Push(&pq, item)
			}
		}

	}
	return 0., []string{A,B}
}
