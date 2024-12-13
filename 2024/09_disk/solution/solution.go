package solution

import (
	"fmt"
	"math"
	"slices"
	"strings"
	"time"
	"unicode"
)

type FileID int

const emptyFID FileID = -1

type DiskInfo struct {
	Blocks []FileID
	// There could be only 1-9 and 0 is here to not drown in OBOE
	EmptyIndciesBySize [10][]int
}

type ParsedInput = DiskInfo

func NewBlockSequence(val FileID, count int) []FileID {
	blocks := make([]FileID, count)
	for i := range blocks {
		blocks[i] = val
	}
	return blocks
}

// 12345
// File ID: 0 -> 1 block
// Empty -> 2 blocks
// File ID: 1 -> 3 blocks
// Empty -> 4 blocks
// File ID: 2 -> 5 blocks
// 0..111....2222
func ParseInput(lines []string) (ParsedInput, error) {
	defer Track(time.Now(), "ParseInput")
	if len(lines) != 1 {
		return DiskInfo{}, fmt.Errorf("expected 1 line, got %d", len(lines))
	}
	line := lines[0]

	diskBlocks := make([]FileID, 0)
	emptyIndciesBySize := [10][]int{}

	curID := FileID(0)

	for i, c := range line {
		if !unicode.IsDigit(c) {
			return DiskInfo{}, fmt.Errorf("invalid character at position %d: %c", i, c)
		}
		blockCount := int(c - '0')

		var fileID FileID

		if i%2 == 0 {
			fileID = curID
			curID++
		} else {
			fileID = emptyFID
		}

		blocks := NewBlockSequence(fileID, blockCount)
		// add empty indices by size, len diskBlocks is the index of the index of the
		// first block of the current fileID added
		if fileID == emptyFID && blockCount != 0 {
			emptyIndciesBySize[blockCount] = append(
				emptyIndciesBySize[blockCount],
				len(diskBlocks),
			)
		}

		diskBlocks = append(diskBlocks, blocks...)

	}

	return DiskInfo{Blocks: diskBlocks, EmptyIndciesBySize: emptyIndciesBySize}, nil
}

// Disk is represented by contiguous slice of "blocks" where each block is a FileID value
// -1 is a special value of FileID for an empty block
// so 12345:
// i: 0    1  2   3 4 5   6  7  8  9   10 11 12 13 14
//
//	0   -1 -1   1 1 1  -1 -1 -1 -1   2  2  2  2  2
//
// EmptyCatalog contains indices pointing to start of the empty blocks of a certain size
// so for the example above it would be:
// 2: [1]
// 4: [6]
type Disk struct {
	Blocks       []FileID
	EmptyCatalog [10][]int
}

func (d Disk) String() string {
	var sb strings.Builder
	prevID := FileID(0)
	sb.WriteRune('{')
	for _, fileID := range d.Blocks {
		if fileID != prevID {
			prevID = fileID
			if fileID == emptyFID {
				sb.WriteString(fmt.Sprintf("}{.-"))
				continue
			}
			sb.WriteString(fmt.Sprintf("}{%d-", fileID))
			continue
		}
		if fileID == emptyFID {
			sb.WriteString(".-")
			continue
		}
		sb.WriteString(fmt.Sprintf("%d-", fileID))
	}
	sb.WriteRune('}')
	return sb.String()
}

func (d *Disk) Compactify() {
	// Simple two-pointers approach
	li, ri := 0, len(d.Blocks)-1
	// Move left pointer and right pointer towards each other until they meet
	for li < ri {
		// Seek the first empty block from the left
		if d.Blocks[li] != emptyFID {
			li++
			continue
		}
		// Seek the first non-empty block from the right
		if d.Blocks[ri] == emptyFID {
			ri--
			continue
		}
		// Swap the blocks, when both pointers are in the right place
		d.Blocks[li], d.Blocks[ri] = d.Blocks[ri], emptyFID
		li++
		ri--
	}
}

func (d *Disk) RunDefragmentation() {
	rightIndex := len(d.Blocks) - 1

	movedEnds := make(map[int]int)
	var fileStart, fileEnd, fileSize int

	// Here we are going to iterate from the right to the left
	for rightIndex >= 0 {
		// We don't care about empty blocks on the right
		if d.Blocks[rightIndex] == emptyFID {
			rightIndex--
			continue
		}
		if nextFileEnd, ok := movedEnds[rightIndex]; ok {
			rightIndex = nextFileEnd
			continue
		}

		// No need to proceed if we are at the beginning of the disk (last file)
		fileStart = d.SeekFileStart(rightIndex)
		if fileStart == 0 {
			break
		}
		fileEnd, rightIndex = rightIndex, fileStart-1
		fileSize = fileEnd - fileStart + 1

		emptySize := d.FindLeftestEmptyFit(fileSize)
		if emptySize == -1 {
			continue
		}
		emptyBlocks := d.EmptyCatalog[emptySize]

		// If the first fitting empty block is to the right of the file,
		// than we can't use it
		firstEmptyBlock := emptyBlocks[0]
		if firstEmptyBlock > fileStart {
			continue
		}
		d.EmptyCatalog[emptySize] = slices.Clone(emptyBlocks[1:])

		d.MoveFileToEmptyBlock(fileStart, fileEnd, firstEmptyBlock)
		// The skip map mentioned above is set here to the next left block from
		// found empty block
		movedEnds[fileEnd] = firstEmptyBlock - 1

		// If it was not perfect fit, we need to update the empty catalog
		if emptySize == fileSize {
			continue
		}
		newEmptySize := emptySize - fileSize
		newEmptyIndex := firstEmptyBlock + fileSize
		d.UpdateEmptyCatalog(newEmptySize, newEmptyIndex)
	}
}

func (d *Disk) SeekFileStart(fileEnd int) int {
	if fileEnd >= len(d.Blocks) {
		panic(fmt.Sprintf("disk index out of bounds: %d vs %d", fileEnd, len(d.Blocks)))
	}

	fileID := d.Blocks[fileEnd]

	fileStart := fileEnd
	for {
		if fileStart == 0 || d.Blocks[fileStart-1] != fileID {
			break
		}
		fileStart--
	}
	return fileStart
}

func (d *Disk) FindLeftestEmptyFit(fileSize int) int {
	minEmptyIndex := math.MaxInt32
	foundSize := -1
	// We check all buckets from the current file size
	// and choose the one with the smallest index
	for emptySize := fileSize; emptySize < 10; emptySize++ {
		if len(d.EmptyCatalog[emptySize]) == 0 {
			return -1
		}
		if d.EmptyCatalog[emptySize][0] < minEmptyIndex {
			minEmptyIndex = d.EmptyCatalog[emptySize][0]
			foundSize = emptySize
		}
	}
	return foundSize
}

func (d *Disk) MoveFileToEmptyBlock(fileStart, fileEnd int, emptyStart int) {
	fileID := d.Blocks[fileStart]
	size := fileEnd - fileStart
	d.ChangeSegmentValue(emptyStart, emptyStart+size, fileID)
	d.DeleteBlocks(fileStart, fileEnd)
}

func (d *Disk) DeleteBlocks(start, end int) {
	d.ChangeSegmentValue(start, end, emptyFID)
}

func (d *Disk) ChangeSegmentValue(start, end int, val FileID) {
	for i := start; i <= end; i++ {
		if val != emptyFID && d.Blocks[i] != emptyFID {
			panic(
				fmt.Sprintf(
					"Trying to write a non-empty block %d [%d:%d] with %d",
					d.Blocks[i],
					start,
					end,
					val,
				),
			)
		}
		d.Blocks[i] = val
	}
}

func (d *Disk) UpdateEmptyCatalog(size, index int) {
	emptyBlocks := d.EmptyCatalog[size]
	emptyBlocks = append(emptyBlocks, index)
	slices.Sort(emptyBlocks)
	d.EmptyCatalog[size] = emptyBlocks
}

func (d *Disk) Checksum() uint64 {
	var cs uint64 = 0

	for i, fileID := range d.Blocks {
		if fileID == emptyFID {
			continue
		}
		cs += uint64(i) * uint64(fileID)
	}

	return cs
}

func PartOne(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartOne")

	blocks := slices.Clone(inp.Blocks)
	disk := Disk{Blocks: blocks}
	disk.Compactify()

	return disk.Checksum()
}

func PartTwo(inp ParsedInput) uint64 {
	defer Track(time.Now(), "PartTwo")

	disk := Disk{
		Blocks:       slices.Clone(inp.Blocks),
		EmptyCatalog: inp.EmptyIndciesBySize,
	}
	disk.RunDefragmentation()

	return disk.Checksum()
}
