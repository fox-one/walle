package number

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"strconv"
)

func RandomPin() string {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panicln(err)
	}
	c := binary.LittleEndian.Uint64(b[:]) % 1000000
	if c < 100000 {
		c = 100000 + c
	}

	return strconv.FormatUint(c, 10)
}
