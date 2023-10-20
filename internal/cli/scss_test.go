package cli

import "testing"

func TestCompileScssToCss(t *testing.T) {
	css, err := NewScssCompilerPool().CompileString(`
	@use 'bla/test';
	ul, ol {
		text-align: left;
	  
		& & {
		  padding: {
			bottom: 0;
			left: 0;
		  }
		}
	  }
	`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(css)
}

func TestCompileScssToCssFileImports(t *testing.T) {
	css, err := NewScssCompilerPool().CompileString(`
	@use 'test/dummy.css';
	ul, ol {
		text-align: left;
	  
		& & {
		  padding: {
			bottom: 0;
			left: 0;
		  }
		}
	  }
	`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(css)
}

func TestCompileScssToCssWebImports(t *testing.T) {
	css, err := NewScssCompilerPool().CompileString(`
	@use 'https://cdnjs.cloudflare.com/ajax/libs/picocss/1.5.2/pico.min.css';
	ul, ol {
		text-align: left;
	  
		& & {
		  padding: {
			bottom: 0;
			left: 0;
		  }
		}
	  }
	`)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(css)
}
