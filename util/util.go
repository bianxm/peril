package util

import "math/rand"

// generate 4 letter room id
func GenerateRoomId() string {
	const l = "BCDFGHJKLMNPQRSTVWXYZ"
	len := len(l)
	rc := make([]byte, 4)
	for i := 0; i < 4; i++ {
		r := rand.Intn(len)
		rc[i] = l[r]
	}
	return string(rc)
}
