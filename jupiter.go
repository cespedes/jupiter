package jupiter

type Jupiter struct {
	index *Index
}

func Open(config string) *Jupiter {
	return nil
}

func (j *Jupiter) Read(s Score) []byte {
	return nil
}

func (j *Jupiter) Write(b []byte) Score {
	return ZeroScore
}
