package jupiter

type Jupiter struct {
	config   *Config
	binheap  *BinHeap
	indexes  []*Index
	datalogs []*DataLog
}

type Type byte

func Open(c *Config) *Jupiter {
	return nil
}

func (j *Jupiter) Read(s Score) (Type, []byte) {
	return 0, nil
}

func (j *Jupiter) Write(t Type, b []byte) Score {
	return ZeroScore
}
