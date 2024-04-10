package confiq

import (
	"strconv"
	"strings"
)

const (
	segmentDividerChar = "."
	openBraceChar      = "["
	closeBraceChar     = "]"
)

type segment interface {
	String() string
}

type keySegment string

func (kS keySegment) String() string {
	return string(kS)
}

func (kS keySegment) asString() string {
	return string(kS)
}

type indexSegment int

func (iS indexSegment) String() string {
	return "[" + strconv.Itoa(int(iS)) + "]"
}

func (iS indexSegment) asInt() int {
	return int(iS)
}

func getNextSegment(path string) (segment, string) {
	var (
		nextSegment, remainingPath string
		segmentDividerIndex        = strings.Index(path, segmentDividerChar)
	)

	if segmentDividerIndex == -1 {
		nextSegment = path
		remainingPath = ""
	} else {
		nextSegment = path[:segmentDividerIndex]
		remainingPath = path[segmentDividerIndex+1:]
	}

	openBraceIndex := strings.Index(path, openBraceChar)
	if openBraceIndex == -1 {
		return keySegment(nextSegment), remainingPath
	}

	closeBraceIndex := strings.Index(nextSegment, closeBraceChar)
	if closeBraceIndex == -1 || closeBraceIndex < openBraceIndex {
		return keySegment(nextSegment), remainingPath
	}

	if openBraceIndex != 0 {
		return keySegment(nextSegment[:openBraceIndex]), nextSegment[openBraceIndex:closeBraceIndex+1] + segmentDividerChar + remainingPath
	}

	index := nextSegment[openBraceIndex+1 : closeBraceIndex]

	if indexInt, err := strconv.Atoi(index); err == nil {
		return indexSegment(indexInt), remainingPath
	}

	return keySegment(index), remainingPath
}
