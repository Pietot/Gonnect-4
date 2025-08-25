package transposition_table

const MASK_49 uint64 = (1 << 49) - 1

type TranspositionTable struct {
	table []Entry
}

type Entry struct {
	key_49_bit uint64
	value      uint8
}

func NewTranspositionTable() *TranspositionTable {
	return &TranspositionTable{
		table: make([]Entry, 8388593),
	}
}

func (trans_table *TranspositionTable) Reset() {
	for i := range trans_table.table {
		trans_table.table[i] = Entry{}
	}
}

func (trans_table *TranspositionTable) Put(key uint64, value uint8) {
	if key >= (1 << 49) {
		panic("Key out of range (must be < 2^49)")
	}
	entryIndex := trans_table.index(key)
	trans_table.table[entryIndex] = Entry{
		key_49_bit: key & MASK_49,
		value:      value,
	}
}

func (trans_table *TranspositionTable) Get(key uint64) uint8 {
	if key >= (1 << 49) {
		panic("Key out of range (must be < 2^49)")
	}
	entryIndex := trans_table.index(key)
	entry := trans_table.table[entryIndex]
	if entry.key_49_bit == key {
		return entry.value
	}
	return 0
}

func (trans_table *TranspositionTable) index(key uint64) uint64 {
	return key % uint64(len(trans_table.table))
}
