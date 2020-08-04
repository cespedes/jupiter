package jupiter

const (
	BlockSize       = 8192
	DirtyBucket     = 1 << 31 // this bucket has not been synced to disk
	EntriesInBucket = 682
)

//type Bucket [BlockSize]byte

//type Bucket int

type Bucket struct {
	entries [EntriesInBucket]Entry
}

//type Bucket [EntriesInBucket]Score

// GetAddress returns the possible addresses for a given Score, if found in the bucket
func (b Bucket) GetAddress(s Score) []uint64 {
	return nil
}

// Add adds an entry in a bucket, identified by a score, and points it to a given address.
// It returns true if successful, false if there is no space in the bucket.
func (b Bucket) Add(s Score, a uint64) bool {
	return false
}

func (b Bucket) score(i int) Score {
	return b.entries[i].s
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
