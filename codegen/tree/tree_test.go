package tree_test

import (
	"github.com/jt0610/scaf/codegen/tree"
	"os"
	"testing"
)

func TestTree_Render(t *testing.T) {
	err := os.Remove("testResult")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}
	dev := tree.New("test_dev", true)
	err = dev.Render("testResult")
	if err != nil {
		t.Fatal(err)
	}
}
