package grid

type MoveEntry struct {
	Move  uint64
	Score int
}

type MoveSorter struct {
	size    int
	entries [WIDTH]MoveEntry
}

func newMoveSorter() *MoveSorter {
	return &MoveSorter{
		size: 0,
	}
}

func (ms *MoveSorter) addMove(move uint64, score int) {
	pos := ms.size
	ms.size++
	for pos > 0 && ms.entries[pos-1].Score > score {
		ms.entries[pos] = ms.entries[pos-1]
		pos--
	}
	ms.entries[pos].Move = move
	ms.entries[pos].Score = score
}

func (ms *MoveSorter) getNextMove() uint64 {
	if ms.size == 0 {
		return 0
	}
	ms.size--
	return ms.entries[ms.size].Move
}
