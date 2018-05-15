package colon

import (
	"encoding/base32"
)

// MSeed2Bytes converts a seed into a [32]byte array
func MSeed2Bytes(seed string) (seedBytes [32]byte, err error) {
	// convert the seed to []byte
	data, err := base32.StdEncoding.DecodeString(seed)
	if err != nil {
		return seedBytes, err
	}
	// copy the bytes from 1 to 32
	for i := 0; i < 32; i++ {
		seedBytes[i] = data[i+1]
	}

	return seedBytes, err
}
