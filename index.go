package jupiter

type Index struct {
	filename string
}

func OpenIndex(filename string) *Index {
	return &Index{filename: filename}
}

func (i *Index) GetEntry(score Score) *Entry {
	return nil
}

func (i *Index) GetBucketFromScore(score Score) int {
	return 0
}
