package solution

import (
	"fmt"
	"slices"
)

const abcSize = 5

type AlphaIndex int

const (
	w AlphaIndex = iota
	u
	b
	r
	g
)

func alphaIndex(char rune) AlphaIndex {
	switch char {
	case 'w':
		return w
	case 'u':
		return u
	case 'b':
		return b
	case 'r':
		return r
	case 'g':
		return g
	default:
		panic("Invalid character")
	}
}

func byteAlphaIndex(char byte) AlphaIndex {
	switch char {
	case 'w':
		return w
	case 'u':
		return u
	case 'b':
		return b
	case 'r':
		return r
	case 'g':
		return g
	default:
		panic("Invalid character")
	}
}

type TrieNode struct {
	Char     rune
	Children [abcSize]*TrieNode
	Depth    int
	IsRoot   bool
}

func NewTrieNode(char rune, depth int) *TrieNode {
	return &TrieNode{
		Char:     char,
		Children: [abcSize]*TrieNode{},
		Depth:    depth,
		IsRoot:   char == 0,
	}
}

func (t *TrieNode) IterChildren() func(func(int, *TrieNode) bool) {
	return func(yield func(int, *TrieNode) bool) {
		for i, child := range t.Children {
			if child == nil {
				continue
			}
			if !yield(i, child) {
				return
			}
		}
	}
}

func (t *TrieNode) HasChild(idx AlphaIndex) bool {
	return t.Children[idx] != nil
}

func (t *TrieNode) GetChild(idx AlphaIndex) *TrieNode {
	return t.Children[idx]
}

func (t *TrieNode) AddChild(char rune, cIdx AlphaIndex) *TrieNode {
	child := NewTrieNode(char, t.Depth+1)
	t.Children[cIdx] = child
	return child
}

func (t *TrieNode) String() string {
	children := make([]string, 0, len(t.Children))
	for _, child := range t.IterChildren() {
		children = append(children, fmt.Sprintf("%c", child.Char))
	}
	slices.Sort(children)
	identifier := string(t.Char)
	if t.IsRoot {
		identifier = "<>"
	}
	return fmt.Sprintf("%s (%d) -> %v", identifier, t.Depth, children)
}

type AhoCorasick struct {
	Root      *TrieNode
	FailLinks map[*TrieNode]*TrieNode
	// Patterns are unique and added once, so we can use slice here
	Outputs    map[*TrieNode][]string
	TotalDepth int
}

func (ac *AhoCorasick) Match(word string) bool {
	state := ac.Root
	matches := make([]bool, len(word)+1)
	matches[0] = true
	for i, char := range word {
		cIdx := alphaIndex(char)
		for !state.IsRoot && !state.HasChild(cIdx) {
			state = ac.FailLinks[state]
		}
		if !state.HasChild(cIdx) {
			return false
		}
		state = state.GetChild(cIdx)
		end := i + 1
		for _, pattern := range ac.Outputs[state] {
			start := end - len(pattern)
			if matches[start] {
				matches[end] = true
				break
			}
		}
		// If match is not found after the length of the longest pattern
		// we can break the loop, because we can't create our string
		// with the patterns we have
		if i > ac.TotalDepth && ((i+1)%ac.TotalDepth) == 1 && !matches[i+1] {
			break
		}

	}
	return matches[len(word)]
}

func (ac *AhoCorasick) MatchCount(word string) int {
	state := ac.Root
	matches := make([]int, len(word)+1)
	matches[0] = 1
	for i, char := range word {
		cIdx := alphaIndex(char)
		for !state.IsRoot && !state.HasChild(cIdx) {
			state = ac.FailLinks[state]
		}
		if !state.HasChild(cIdx) {
			return 0
		}
		state = state.GetChild(cIdx)
		end := i + 1
		for _, pattern := range ac.Outputs[state] {
			start := end - len(pattern)
			if matches[start] > 0 {
				matches[end] = matches[end] + matches[start]
			}
		}
		if i > ac.TotalDepth && ((i+1)%ac.TotalDepth) == 1 && matches[i+1] == 0 {
			return 0
		}

	}
	return matches[len(word)]
}

func NewAhoCorasick(words []string) *AhoCorasick {
	root := NewTrieNode(0, 0)
	ac := &AhoCorasick{
		Root:      root,
		FailLinks: make(map[*TrieNode]*TrieNode, 100),
		Outputs:   make(map[*TrieNode][]string, 100),
	}

	for _, word := range words {
		ac.AddWord(word)
	}
	ac.BuildFailLinks()
	return ac
}

func (ac *AhoCorasick) AddWord(word string) {
	node := ac.Root
	for _, char := range word {
		cIdx := alphaIndex(char)
		if node.HasChild(cIdx) {
			node = node.GetChild(cIdx)
		} else {
			node = node.AddChild(char, cIdx)
		}
	}
	ac.Outputs[node] = append(ac.Outputs[node], word)
	if len(word) > ac.TotalDepth {
		ac.TotalDepth = len(word)
	}
}

func (ac *AhoCorasick) BuildFailLinks() {
	queue := NewQueue()
	node := ac.Root
	// Init with the depth 1 nodes
	for _, child := range node.IterChildren() {
		ac.FailLinks[child] = node
		queue.Push(child)
	}

	for !queue.Empty() {
		node = queue.Pop()
		ac.buildFailLink(node)
		// Here we proceed with BFS and also initialize fail links for children
		// as the fail link of it's parent (current node)
		failLink := ac.FailLinks[node]
		for _, child := range node.IterChildren() {
			ac.FailLinks[child] = failLink
			queue.Push(child)
		}
	}
}

func (ac *AhoCorasick) buildFailLink(node *TrieNode) {
	// Our init fail link is the fail link of the parent node
	fail := ac.FailLinks[node]
	for {
		for _, child := range fail.IterChildren() {
			if child.Char != node.Char || child == node {
				continue
			}
			// If child was found and it's not the same as the current node
			// we can set the fail link of the current node to that child
			ac.FailLinks[node] = child
			ac.Outputs[node] = append(ac.Outputs[node], ac.Outputs[child]...)
			// If we are at the longest pattern
			return
		}
		// If we are at the root and we didn't find the child, we should stop
		// and set the fail link of the current node to the root
		if fail.IsRoot {
			ac.FailLinks[node] = fail
			return
		}

		// Otherwise we should continue with the fail link of the current fail link
		fail = ac.FailLinks[fail]
	}
}

// This method is used for checking the build of the trie and fail links
func (ac *AhoCorasick) Search(word string) []ACSResult {
	state := ac.Root // This is the current state of the automaton
	results := make([]ACSResult, 0, 10)
	i := 0 // This is the position in the input string
	cIdx := byteAlphaIndex(word[i])

	// Work till the end of the input string
	for i < len(word) {
		// This signals that we can't proceed deeper in the trie
		if !state.HasChild(cIdx) {
			// If we are at the root, we should just move to the next character
			// of the input word, otherwise we should go to the fail link and try again
			if state.IsRoot {
				i++
				cIdx = byteAlphaIndex(word[i])
			} else {
				state = ac.FailLinks[state]
			}
			continue
		}
		state = state.GetChild(cIdx)
		// If we have outputs for the current state, we should add them to the results,
		// all the fail links' outputs are also outputs of the current state and
		// have been added during the building of the fail links
		for _, output := range ac.Outputs[state] {
			results = append(results, NewACSResult(i, output))
		}
		i++
		cIdx = byteAlphaIndex(word[i])
	}

	return results
}

func patternPositions(pos int, pattern string) (int, int) {
	end := pos + 1
	start := end - len(pattern)
	return start, end
}

type ACSResult struct {
	Start, End int
	Pattern    string
}

func NewACSResult(pos int, val string) ACSResult {
	start, end := patternPositions(pos, val)
	return ACSResult{start, end, val}
}

type Queue []*TrieNode

func NewQueue() *Queue {
	q := make(Queue, 0, 100)
	return &q
}

func (q *Queue) Push(node *TrieNode) {
	*q = append(*q, node)
}

func (q *Queue) Pop() *TrieNode {
	node := (*q)[0]
	*q = (*q)[1:]
	return node
}

func (q *Queue) Empty() bool {
	return len(*q) == 0
}
