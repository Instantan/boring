package cli

import "testing"

func TestGetRequest(t *testing.T) {
	data, err := getRequest("https://cdnjs.cloudflare.com/ajax/libs/picocss/1.5.2/pico.min.css")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestReadFile(t *testing.T) {
	{
		data, err := readFile("../../examples/boring-cli/dummy.css")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(data)
	}
	{
		data, err := readFile("../../examples/boring-cli", ".css", ".scss", "/index.css", "/index.scss", "index.css", "index.scss")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(data)
	}
}

func TestFileExists(t *testing.T) {
	if fileExists("test") {
		t.Fatalf("Is not a file")
	}
	if fileExists("../../examples/boring-cli/dummy") {
		t.Fatalf("Should not exist")
	}
	if !fileExists("../../examples/boring-cli/dummy.css") {
		t.Fatalf("Should exsist")
	}
}

func TestGetRelevantFiles(t *testing.T) {
	files, err := getRelevantFiles("../../examples/boring-cli", []string{
		".css",
		".scss",
		".js",
		".ts",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(files)
}

func TestTurnFilepathIntoMinifiedVersion(t *testing.T) {
	t.Log(turnFilepathIntoMinifiedVersion("/bla/bla/file.css"))
}

func TestIsDirectoryEmpty(t *testing.T) {
	empty, err := isDirectoryEmpty(".", []string{"go.mod", "go.sum"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(empty)
}
