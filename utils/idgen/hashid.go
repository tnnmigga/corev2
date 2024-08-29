package idgen

// BKDRHash BKDR hash algorithm for 32bit
func BKDRHash(d []byte) uint32 {
	s := uint32(31)
	v := uint32(0)
	for _, b := range d {
		v = v*s + uint32(b)
	}
	return v
}

// BKDRHash64 BKDR hash algorithm for 64bit
func BKDRHash64(d []byte) uint64 {
	s := uint64(31)
	v := uint64(0)
	for _, b := range d {
		v = v*s + uint64(b)
	}
	return v
}

// HashToID hash string to uint32
func HashToID(s string) uint32 {
	return BKDRHash([]byte(s))
}

// HashToID64 hash string to uint64
func HashToID64(s string) uint64 {
	return BKDRHash64([]byte(s))
}
