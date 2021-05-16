package src

import (
	"fmt"
	"testing"
)

func TestReadFile(t *testing.T) {

	network, err := Parse("../2PC.txt")

	if err != nil {
		t.Fatalf("oh no %v", err)
	}

	fmt.Println(network)

}
