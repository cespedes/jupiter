package jupiter

const (
	BlockSize        = 8192
	DirtyBucket      = 1 << 31 // this bucket has not been synced to disk
	EntriesInBucket  = 682
)

//type Bucket [BlockSize]byte

//type Bucket int

type Bucket struct {
	scores [EntriesInBucket]Score
}

//type Bucket [EntriesInBucket]Score

func (b Bucket) score(i int) Score {
	return b.scores[i]
}

func (b *Bucket) setScore(i int, s Score) {
	b.scores[i] = s
}

func (b Bucket) size() int {
	var i int
	for i = 0; i < EntriesInBucket; i++ {
		if b.score(i).Equal(ZeroScore) {
			break
		}
	}
	return i
}

var buckets []Bucket

