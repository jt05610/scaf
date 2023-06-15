package builder

import (
	"github.com/jt05610/scaf/core"
	"os"
	"testing"
	"text/template"
)

func TestNewDir(t *testing.T) {
	dir := NewDir(func(*core.Module) string { return "testDir" }, []*File{})
	if dir.Name(nil) != "testDir" {
		t.Errorf("Expected 'testDir', got '%s'", dir.Name(nil))
	}
	if len(dir.Files) != 0 {
		t.Errorf("Unexpected files in directory")
	}
}

func TestDirBuilder(t *testing.T) {
	dir := NewDir(func(*core.Module) string { return "testDir" }, []*File{
		{Name: func(*core.Module) string { return "file1" }, Template: template.Must(template.New("test1").Parse(""))},
		{Name: func(*core.Module) string { return "file2" }, Template: template.Must(template.New("test2").Parse(""))},
	})
	builder := NewDirBuilder(dir).(*dirBuilder)
	builder.parent = "testParent"
	err := builder.VisitModule(&core.Module{Name: "testModule"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Now, verify the directory and files were created as expected
	path := "testParent/testModule/testDir"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Directory not created: %s", path)
	}
	for _, file := range dir.Files {
		if _, err := os.Stat(path + "/" + file.Name(nil)); os.IsNotExist(err) {
			t.Errorf("File not created: %s", path+"/"+file.Name(nil))
		}
	}

	// Clean up
	_ = os.RemoveAll("testParent")
}

func TestDirBuilderWithCircularDependency(t *testing.T) {
	dir1 := NewDir(func(*core.Module) string { return "dir1" }, []*File{
		{Name: func(*core.Module) string { return "file1" }, Template: template.Must(template.New("test1").Parse(""))},
	})
	dir2 := NewDir(func(*core.Module) string { return "dir2" }, []*File{
		{Name: func(*core.Module) string { return "file2" }, Template: template.Must(template.New("test2").Parse(""))},
	})
	dir1.AddChild(dir2)
	dir2.AddChild(dir1) // This creates a circular dependency

	builder := NewDirBuilder(dir1).(*dirBuilder)
	builder.parent = "testParent"
	err := builder.VisitModule(&core.Module{Name: "testModule"})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	// Now, verify the directories and files were created as expected
	paths := []string{
		"testParent/testModule/dir1",
		"testParent/testModule/dir1/file1",
		"testParent/testModule/dir1/dir2",
		"testParent/testModule/dir1/dir2/file2",
	}

	for _, dirPath := range paths {
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Not created: %s", dirPath)
		}
	}
	_ = os.RemoveAll("testParent")
}
