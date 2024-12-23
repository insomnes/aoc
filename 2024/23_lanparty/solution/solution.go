package solution

import (
	"fmt"
	"slices"
	"sort"
	"strings"
	"time"
)

type Computer struct {
	Name        string
	Connections map[string]*Computer
	Degree      int
}

func NewComputer(name string) *Computer {
	return &Computer{Name: name, Connections: make(map[string]*Computer)}
}

func (c *Computer) Connect(other *Computer) {
	c.Connections[other.Name] = other
	other.Connections[c.Name] = c
	c.Degree++
}

func (c *Computer) FindCommons(other *Computer) []string {
	commons := make([]string, 0, len(c.Connections))
	checked := make(map[string]struct{})
	for name := range c.Connections {
		if name == c.Name || name == other.Name {
			continue
		}
		if _, ok := checked[name]; ok {
			continue
		}
		if _, ok := other.Connections[name]; !ok {
			continue
		}
		commons = append(commons, name)
		checked[name] = struct{}{}
	}
	return commons
}

type Set map[string]struct{}

func NewSet(n int) Set {
	return make(map[string]struct{}, n)
}

func NewSetFromValues(values ...string) Set {
	s := NewSet(len(values))
	for _, value := range values {
		s[value] = struct{}{}
	}
	return s
}

func (s Set) Add(name string) {
	s[name] = struct{}{}
}

func (s Set) Remove(name string) {
	delete(s, name)
}

func (s Set) Contains(name string) bool {
	_, ok := s[name]
	return ok
}

func (s Set) Empty() bool {
	return s.Size() == 0
}

func (s Set) Size() int {
	return len(s)
}

func (s Set) Intersection(other Set) Set {
	minimal, toCheck := s, other
	if other.Size() < s.Size() {
		minimal, toCheck = other, s
	}
	intersect := NewSet(minimal.Size())
	for name := range minimal {
		if toCheck.Contains(name) {
			intersect.Add(name)
		}
	}
	return intersect
}

func (s Set) Union(other Set) Set {
	union := NewSet(s.Size() + other.Size())
	for name := range s {
		union.Add(name)
	}
	for name := range other {
		union.Add(name)
	}
	return union
}

func (s Set) ToSortedSlice() []string {
	slice := make([]string, 0, s.Size())
	for name := range s {
		slice = append(slice, name)
	}
	slices.Sort(slice)
	return slice
}

type LANMap struct {
	LAN       map[string]Set
	Computers map[string]*Computer
}

func NewLANMap(lan map[string]Set) LANMap {
	computers := make(map[string]*Computer)
	lm := LANMap{LAN: lan, Computers: computers}
	for name, connections := range lan {
		for connection := range connections {
			lm.Connect(name, connection)
		}
	}
	return lm
}

func (l *LANMap) Connect(node1, node2 string) {
	computer1 := l.EnsureComputer(node1)
	computer2 := l.EnsureComputer(node2)
	computer1.Connect(computer2)
}

func (l *LANMap) EnsureComputer(name string) *Computer {
	if computer, ok := l.Computers[name]; ok {
		return computer
	}
	computer := NewComputer(name)
	l.Computers[name] = computer
	return computer
}

func (l *LANMap) SearchTTriangles() map[Triangle]struct{} {
	triangles := map[Triangle]struct{}{}
	checked := make(map[string]struct{})
	for cn, computer := range l.Computers {
		if _, ok := checked[cn]; ok {
			continue
		}
		cnStartsT := cn[0] == 't'
		for nn, neighbor := range computer.Connections {
			nnStartsT := nn[0] == 't'
			commons := computer.FindCommons(neighbor)
			for _, common := range commons {
				commonStartsT := common[0] == 't'
				if cnStartsT || nnStartsT || commonStartsT {
					triangle := NewTriangle([3]string{cn, nn, common})
					triangles[triangle] = struct{}{}
				}
			}
			checked[cn] = struct{}{}
		}
	}

	return triangles
}

func (lm *LANMap) CheckRegularity() int {
	knownDegree := -1
	for _, computer := range lm.Computers {
		if knownDegree == -1 {
			knownDegree = computer.Degree
			continue
		}
		if knownDegree != computer.Degree {
			return -1
		}
	}
	return knownDegree
}

type Triangle struct {
	Nodes [3]string
}

func (t Triangle) String() string {
	return fmt.Sprintf("%s-%s-%s", t.Nodes[0], t.Nodes[1], t.Nodes[2])
}

func NewTriangle(nodes [3]string) Triangle {
	// sort the node by alphabetical order
	nSlice := nodes[:]
	sort.Strings(nSlice)
	nodes = [3]string{nSlice[0], nSlice[1], nSlice[2]}
	return Triangle{Nodes: nodes}
}

type ParsedInput = LANMap

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	lan := make(map[string]Set)
	for _, line := range lines {
		split := strings.Split(line, "-")
		if len(split) != 2 {
			return ParsedInput{}, fmt.Errorf("line %s", line)
		}
		node1, node2 := split[0], split[1]
		if _, ok := lan[node1]; !ok {
			lan[node1] = NewSet(1)
		}
		if _, ok := lan[node2]; !ok {
			lan[node2] = NewSet(1)
		}
		lan[node1].Add(node2)
		lan[node2].Add(node1)
	}
	return NewLANMap(lan), nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")
	total := 0
	lm := inp
	for _, computer := range lm.Computers {
		connections := make([]string, 0, len(computer.Connections))
		for name := range computer.Connections {
			connections = append(connections, name)
			total++
		}
	}

	triangles := lm.SearchTTriangles()

	total = len(triangles)

	return total
}

func BronKerbosh(clique Set, toCheck Set, excluded Set, graph map[string]Set, knownMax Set) Set {
	if toCheck.Empty() && excluded.Empty() {
		return clique
	}

	for vertex := range toCheck {
		vertexSet := NewSetFromValues(vertex)
		vertexConnections := graph[vertex]

		// Add vertex to clique
		newClique := clique.Union(vertexSet)
		// Vertex neighbors in candidates set
		newCandidates := toCheck.Intersection(vertexConnections)
		// Do not continue if the new clique size is not enough to beat known max
		if newClique.Size()+newCandidates.Size() <= knownMax.Size() {
			continue
		}
		// Vertex neighbors in excluded set
		newExcluded := excluded.Intersection(vertexConnections)
		foundClique := BronKerbosh(newClique, newCandidates, newExcluded, graph, knownMax)
		if foundClique.Size() > knownMax.Size() {
			knownMax = foundClique
		}

		// Move vertex from candidates to excluded
		toCheck.Remove(vertex)
		excluded.Add(vertex)

	}

	return knownMax
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	lm := inp

	vc := len(lm.Computers)
	candidates := NewSet(vc)
	for name := range lm.Computers {
		candidates.Add(name)
	}
	clique := BronKerbosh(NewSet(vc+1), candidates, NewSet(vc+1), lm.LAN, NewSet(0))
	fmt.Println("Clique size:", clique.Size())
	password := strings.Join(clique.ToSortedSlice(), ",")
	fmt.Println("Password:", password)

	return clique.Size()
}
