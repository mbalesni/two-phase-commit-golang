package src

import (
	"testing"
)

// func GetProcessIdxById(id int, processes []*Process) int {
// 	for i, process := range processes {
// 		if process.Id == id {
// 			return i
// 		}
// 	}
// 	return -1
// }

func TestSpawn(t *testing.T) {
	//verbose := true

	network, err := Parse("../2PC.txt")

	if err != nil {
		t.Fatalf("%v", err)
	}

	network.ListHistory()

}
