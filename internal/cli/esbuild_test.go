package cli

import (
	"fmt"
	"testing"
)

func TestEsbuild(t *testing.T) {
	result, err := buildJsFile("./../../examples/boring-cli/my.js")
	if err != nil {
		t.Fatal(err)
	}
	for i := range result {
		fmt.Printf("%v", string(result[i].Contents))
	}
}
