package util

import "testing"

func Test_GenerateRoomID(t *testing.T) {
	id := GenerateRoomId()
	t.Logf("%s\n", id)
}
