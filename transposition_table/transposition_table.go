package transposition_table

const MASK_49 uint64 = (1 << 49) - 1

type TranspositionTable struct {
	table []Entry
}

type Entry struct {
	key_56_bit uint64
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
	if key >= (1 << 56) {
		panic("Key out of range (must be < 2^56)")
	}
	entryIndex := trans_table.index(key)
	trans_table.table[entryIndex].key_56_bit = key & MASK_56
	trans_table.table[entryIndex].value = value
}

func (trans_table *TranspositionTable) Get(key uint64) uint8 {
	if key >= (1 << 56) {
		panic("Key out of range (must be < 2^56)")
	}
	entryIndex := trans_table.index(key)
	if trans_table.table[entryIndex].key_56_bit == key {
		return trans_table.table[entryIndex].value
	}
	return 0
}

func (trans_table *TranspositionTable) index(key uint64) uint64 {
	return key % uint64(len(trans_table.table))
}
