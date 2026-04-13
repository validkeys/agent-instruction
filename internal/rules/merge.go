package rules

// MergeInstructions merges instruction arrays from multiple sources in order.
// It preserves exact order and does not deduplicate - user controls structure
// through their rule files.
func MergeInstructions(sources [][]Instruction) []Instruction {
	// Calculate total size to pre-allocate slice
	totalSize := 0
	for _, source := range sources {
		totalSize += len(source)
	}

	// Pre-allocate result slice for efficiency
	result := make([]Instruction, 0, totalSize)

	// Flatten arrays in order
	for _, source := range sources {
		result = append(result, source...)
	}

	return result
}
