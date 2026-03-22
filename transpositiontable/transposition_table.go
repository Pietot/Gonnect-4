package transpositiontable

// A prime number close but under to 2^23 for better distribution.
const TT_SIZE = 8388593

type TranspositionTable struct {
	table []Entry
}

type Entry struct {
	key   uint64
	value uint8
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{

		table: make([]Entry, TT_SIZE),
	}
}

func (trans_table *TranspositionTable) Reset() {
	for i := range trans_table.table {
		trans_table.table[i] = Entry{}
	}
}

func (trans_table *TranspositionTable) Put(key uint64, value uint8) {
	entryIndex := trans_table.index(key)
	trans_table.table[entryIndex] = Entry{
		key:   key,
		value: value,
	}
}

func (trans_table *TranspositionTable) Get(key uint64) uint8 {
	entryIndex := trans_table.index(key)
	entry := trans_table.table[entryIndex]
	if entry.key == key {
		return entry.value
	}
	return 0
}

func (trans_table *TranspositionTable) index(key uint64) uint64 {
	return key % uint64(len(trans_table.table))
}
