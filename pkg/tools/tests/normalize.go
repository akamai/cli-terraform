// Package tests contains methods and tools that should be used for tests only.
package tests

import (
	"sort"
	"strings"
)

// NormalizeWholeFile sorts all lines to produce a stable
// representation regardless of map iteration order.
func NormalizeWholeFile(content string) string {
	lines := strings.Split(content, "\n")
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// NormalizeFieldsInBlock sorts field lines within each block matched by blockPrefix
// in-place to handle non-deterministic map iteration in templates.
func NormalizeFieldsInBlock(content, blockPrefix string) string {
	lines := strings.Split(content, "\n")
	var result []string
	var openingLine string
	var fieldLines []string
	inBlock := false

	for _, line := range lines {
		switch {
		case !inBlock && strings.HasPrefix(line, blockPrefix):
			inBlock = true
			openingLine = line
			fieldLines = nil

		case inBlock && line == "}":
			sort.Strings(fieldLines)
			result = append(result, openingLine)
			result = append(result, fieldLines...)
			result = append(result, line)
			inBlock = false

		case inBlock:
			fieldLines = append(fieldLines, line)

		default:
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

// NormalizeBlocksInFile collects all blocks matched by blockPrefix, sorts them as units,
// and appends them after all non-block lines to handle non-deterministic ordering from map iteration.
func NormalizeBlocksInFile(content, blockPrefix string) string {
	lines := strings.Split(content, "\n")
	var nonBlockLines []string
	var completedBlocks []string
	var openingLine string
	var fieldLines []string
	inBlock := false

	for _, line := range lines {
		switch {
		case !inBlock && strings.HasPrefix(line, blockPrefix):
			inBlock = true
			openingLine = line
			fieldLines = nil

		case inBlock && line == "}":
			blockLines := append([]string{openingLine}, fieldLines...)
			blockLines = append(blockLines, line)
			completedBlocks = append(completedBlocks, strings.Join(blockLines, "\n"))
			inBlock = false

		case inBlock:
			fieldLines = append(fieldLines, line)

		default:
			nonBlockLines = append(nonBlockLines, line)
		}
	}

	sort.Strings(completedBlocks)
	out := strings.Join(nonBlockLines, "\n")
	for _, block := range completedBlocks {
		out += "\n" + block
	}
	return out
}
