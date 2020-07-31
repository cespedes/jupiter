package jupiter

type Index struct {
	filename string
	DataLog  *DataLog
}

func OpenIndex(filename string) *Index {
	return nil
}

func (i *Index) GetEntry(score Score) *Entry {
	return nil
}
