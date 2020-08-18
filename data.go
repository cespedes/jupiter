package jupiter

type DataLog struct {
	filename string
}

// Each block is prefixed by a header that describes the contents of the
// block.  The header contains the score, the type, the compression
// algorithm, the compressed size and uncompressed size.

func OpenDataLog(filename string) *DataLog {
	panic("not implemented")
	return nil
}

func NewDataLog() *DataLog {
	panic("not implemented")
	return nil
}

func (d *DataLog) NewChunk(score Score, t Type, b []byte) (addr uint64, err error) {
	panic("not implemented")
	return 0, nil
}

// PeekChunk is used to check if a given block is stored at an address
func (d *DataLog) PeekChunk(score Score, addr uint64) (t Type, err error) {
	panic("not implemented")
	return 0, nil
}

// GetChunk returns the block with a given score stored at an address
func (d *DataLog) GetChunk(score Score, addr uint64) (t Type, b []byte, err error) {
	panic("not implemented")
	return 0, nil, nil
}

// func (j *Jupiter) Write(t Type, b []byte) (Score, error) {
