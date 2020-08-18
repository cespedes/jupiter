Jupiter is a network storage system that stores data blocks.  A 256-bit
hash of the data (called score) acts as the address of the data.
This enforces a write-once policy since no other data block can be
found with the same address: the addresses of multiple writes of the
same data are identical, so duplicate data is easily identified and
the data block is stored only once.

Jupiter is heavily inspired by Venti, a network storage system
designed and implemented at Bell Labs for the Plan 9 distribution.

Details
-------

Jupiter uses the SHA-512/256 hash function, designed by the United
States NSA and first published in 2001.  As of 2020, it is considered
secure.

Jupiter is a user space daemon.  Clients connect to Jupiter over TCP and
communicate using a simple RPC protocol.  The most important messages of
the protocol are listed below:

+ read(score) returns the data identified by score.
+ write(data) stores data at the address calculated by its hash (score),
  and returns this score.

Storage
-------

The implementation uses an append-only log of data blocks and an index
that maps fingerprints to locations in this log.

The simplicity of the append-only log structure eliminates many possible
software errors that might cause data corruption and facilitates a
variety of additional integrity strategies. A separate index structure
allows a block to be efficiently located in the log; however, the index
can be regenerated from the data log if required and thus does not have
the same reliability constraints as the log itself.

### The data log

The data log contains a series of data blocks, that are filled in an
append-only manner. In Jupiter, data blocks are variable sized, but
since blocks are immutable they can be densely packed without
fragmentation.

Each block is prefixed by a header that describes the contents of the
block.  The header contains the score, the type, the compression
algorithm, the compressed size and uncompressed size.

### The index

We implement the index using a hash table.  The index is divided into
fixed-sized buckets, each of which is stored as a single disk block.
Each bucket contains the index map for a small section of the
score space, as a list of entries sorted by score.  A hash function is
used to map scores to index buckets in a roughly uniform manner, and
then the bucket is examined using binary search. In case there is a
buffer overflow when adding an entry to a bucket, a new bucket will be
allocated and the hash function will change to reflect this.

Each bucket has a small header and some fixed-size entries, which have
part of the score of one block, and the address of that block in
the data log.

We store in the index buckets part of the score, not all of it, for
space constraints.  Thus, if there is a match in a bucket it is
still necessary to access the data in order to know for sure if a block
is present in the archive.  Additionally, the part of the score we store
could be shared by several blocks in the index, and we will need to check
all of them to look for a specific block.

#### Bucket header

The bucket header contains:

- the number of bytes in the score used to identify each entry
- the number and value of the initial bits in the score which are
  common to all the entries in this bucket
- the number of bytes used to address each entry in the data log

So, each bucket can contain a different number of entries.

#### Entries in a bucket

The hash function to determine in which bucket can we find a given block
is a combination of the score of that block and a binary tree, stored as
a binary heap.  The first bits in the score determines the position in
the binary tree, and the value of that node, if not empty, indicates the
bucket where the block should be looked up.

As an example, let's suppose we begin with an index with 4 buckets, from
1 to 4, where the first bits in the score determine the bucket:

- 00 bucket 1
- 01 bucket 2
- 10 bucket 3
- 11 bucket 4

If bucket number 3 overflows, a new bucket (5) will be allocated, and
the new table will be:

- 00  bucket 1
- 01  bucket 2
- 100 bucket 3
- 101 bucket 5
- 11  bucket 4

After a new overflow of bucket 5, the table will be:

- 00   bucket 1
- 01   bucket 2
- 100  bucket 3
- 1010 bucket 5
- 1011 bucket 6
- 11   bucket 4

And, after a new overflow of bucket 1:

- 000  bucket 1
- 001  bucket 7
- 01   bucket 2
- 100  bucket 3
- 1010 bucket 5
- 1011 bucket 6
- 11   bucket 4

The table index will always reside in memory.  The index buckets will be
in memory after they are used for the first time.

### Index initialization

Disk storage
------------

### Configuration

For a given Jupiter instance, we need to know a few things.
This configuration is stored as a text file, with one line for each
configuration option:

* the server listen address
  listen *address*
* the HTTP port for stats
  http *address*
* location of the binary heap in disk
  heap *location*
* size of each bucket in bytes
  bucketsize *size*
* number of bytes used for fingerprint in each bucket entry
  fpsize *size*
* location of the index buckets in disk
  buckets *location*
* location of the data log in disk
  data *location*


### Binary heap

Each entry in the tree has the address of a bucket, expressed as a
big-endian 32-bit unsigned int.  Thus, the maximum number of buckets
is 2^32.i, which will give us:
* 32TB maximum memory for buckets, assuming 8K buckets
* 2T blocks of data, assuming 16 bytes per entry.
* 4 PB of storage, assuming 2K bytes per data block.

### Index buckets
### Data log


References
----------

* Venti: a new approach to archival storage: `http://doc.cat-v.org/plan_9/4th_edition/papers/venti/`
* Venti analysis and memventi implementation: `http://essay.utwente.nl/694/1/scriptie_Lukkien.pdf`
* `https://github.com/mjl-/memventi`
* `https://github.com/mjl-/ventisrv`
