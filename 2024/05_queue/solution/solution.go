package solution

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CompareResult int8

const (
	NotStrict CompareResult = iota
	Yes
	No
)

func (cr CompareResult) String() string {
	switch cr {
	case NotStrict:
		return "NotStrict"
	case Yes:
		return "Yes"
	case No:
		return "No"
	default:
		panic("unexpected case")
	}
}

type Page struct {
	Val    int
	Before map[int]struct{}
	After  map[int]struct{}
}

func NewPageOrdering(val int) Page {
	return Page{
		Val:    val,
		Before: make(map[int]struct{}),
		After:  make(map[int]struct{}),
	}
}

func (p *Page) AddBefore(val int) {
	p.Before[val] = struct{}{}
}

func (p *Page) AddAfter(val int) {
	p.After[val] = struct{}{}
}

func (p *Page) ShouldBeBefore(other Page) CompareResult {
	if _, ok := p.Before[other.Val]; ok {
		return Yes
	}
	if _, ok := p.After[other.Val]; ok {
		return No
	}
	return NotStrict
}

func (p *Page) PrintRules() {
	fmt.Printf("Page %d\n", p.Val)
	for num := range p.Before {
		fmt.Printf("%d|%d\n", p.Val, num)
	}
	for num := range p.After {
		fmt.Printf("%d|%d\n", num, p.Val)
	}
	fmt.Println()
}

type PagePrinter struct {
	pages   [100]Page
	updates [][]int
}

func NewPagePrinter() PagePrinter {
	return PagePrinter{pages: [100]Page{}, updates: make([][]int, 0)}
}

func (pp *PagePrinter) CheckUpdateOrdered(update []int) bool {
	for i := 1; i < len(update); i++ {
		cur, ok := pp.GetPage(update[i])
		if !ok {
			continue
		}
	BeforeLoop:
		for j := 0; j < i; j++ {
			prev, ok := pp.GetPage(update[j])
			if !ok {
				continue BeforeLoop
			}

			orderResult := prev.ShouldBeBefore(cur)

			switch orderResult {
			case Yes, NotStrict:
				continue BeforeLoop
			case No:
				return false
			default:
				panic(fmt.Sprintf("unexpected case %s", orderResult))
			}
		}

	}
	return true
}

func (pp *PagePrinter) OrderUpdate(update []int) ([]int, bool) {
	updCopy := make([]int, len(update))
	copy(updCopy, update)
	update = updCopy
	alreadyOrdered := true
	for i := 1; i < len(update); i++ {
		cur, ok := pp.GetPage(update[i])
		if !ok {
			continue
		}
	BeforeLoop:
		for j := 0; j < i; j++ {
			prev, ok := pp.GetPage(update[j])
			if !ok {
				continue BeforeLoop
			}

			orderResult := prev.ShouldBeBefore(cur)

			switch orderResult {
			case Yes, NotStrict:
				continue BeforeLoop
			case No:
				alreadyOrdered = false
				update[i], update[j] = update[j], update[i]
			default:
				panic(fmt.Sprintf("unexpected case %s", orderResult))
			}
		}

	}
	return update, alreadyOrdered
}

func (pp *PagePrinter) GetPage(val int) (Page, bool) {
	page := pp.pages[val]
	if page.Val == 0 {
		return page, false
	}

	return page, true
}

func (pp *PagePrinter) AddRule(left, right int) {
	pp.ensurePage(left)
	pp.ensurePage(right)
	pp.pages[left].AddBefore(right)
	pp.pages[right].AddAfter(left)
}

func (pp *PagePrinter) ensurePage(val int) {
	if val < 10 || val > 99 {
		panic(fmt.Sprintf("invalid page value %d", val))
	}
	if pp.pages[val].Val == 0 {
		pp.pages[val] = NewPageOrdering(val)
	}
}

func (pp *PagePrinter) AddUpdate(update []int) {
	pp.updates = append(pp.updates, update)
}

func parsePageRule(raw string) (int, int) {
	var left, right int

	_, err := fmt.Sscanf(raw, "%d|%d", &left, &right)
	if err != nil {
		panic(err)
	}
	return left, right
}

func parseUpdate(raw string) []int {
	split := strings.Split(raw, ",")
	update := make([]int, len(split))
	for i, s := range split {
		val, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Sprintf("parsing %s: %s", s, err))
		}
		update[i] = val
	}

	return update
}

type ParsedInput = PagePrinter

func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	pp := NewPagePrinter()
	i := 0
	for {
		if i >= len(lines) {
			panic("end of input at rules parse")
		}

		left, right := parsePageRule(lines[i])
		pp.AddRule(left, right)
		i++

		if lines[i] == "" {
			break
		}
	}

	for i := i + 1; i < len(lines); i++ {
		update := parseUpdate(lines[i])
		pp.AddUpdate(update)
	}

	return pp, nil
}

func PartOne(inp ParsedInput) int {
	defer Track(time.Now(), "PartOne")

	midSum := 0
	for _, update := range inp.updates {
		if !inp.CheckUpdateOrdered(update) {
			continue
		}
		midNum := update[len(update)/2]
		midSum += midNum
	}

	return midSum
}

func PartTwo(inp ParsedInput) int {
	defer Track(time.Now(), "PartTwo")
	midSum := 0
	for _, update := range inp.updates {
		orderedUpdate, alreadyOrdered := inp.OrderUpdate(update)
		if alreadyOrdered {
			continue
		}
		midNum := orderedUpdate[len(orderedUpdate)/2]
		midSum += midNum
	}

	return midSum
}
