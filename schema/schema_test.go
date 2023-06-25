package schema

import (
	"github.com/jt05610/scaf/context"
	"github.com/jt05610/scaf/testData"
	"os"
	"testing"
)

func TestSchemer_VisitModule(t *testing.T) {
	df, err := os.Create("test.json")
	defer func() {
		_ = df.Close()
	}()
	if err != nil {
		t.Fatal(err)
	}
	v := NewSchemer(df)
	ctx := context.NewContext(nil)
	mod := testData.APIs()
	if err := v.VisitModule(ctx, mod); err != nil {
		t.Fatal(err)
	}
}
