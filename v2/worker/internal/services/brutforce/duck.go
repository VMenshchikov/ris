package brutforce

import (
	"bytes"
	"crypto/md5"
	utils "hash_worker/internal/services/brutforce/util"
)

const (
	ALPHABET = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func md5Check(word string, target []byte) bool {
	hash := md5.Sum([]byte(word))
	return bytes.Equal(hash[:], []byte(target))
}

func CheckRange(target [16]byte, maxLen, blockSize uint, startIndex uint64) []string {
	results := make([]string, 0)

	generator := utils.NewGenerator(ALPHABET, startIndex, maxLen)

	for i := 0; i < int(blockSize); i++ {
		word, err := generator.Next()
		if err != nil {
			break
		}

		if md5Check(word, target[:]) {
			results = append(results, word)
		}
	}

	return results
}
