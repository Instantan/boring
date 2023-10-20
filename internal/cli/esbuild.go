package cli

import (
	"errors"
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

func buildJsFile(entrypoint string) ([]api.OutputFile, error) {
	var err error
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{entrypoint},
		Bundle:            true,
		MinifySyntax:      true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		Target:            api.DefaultTarget,
		Outfile:           entrypoint,
	})
	if len(result.Errors) > 0 {
		for i := range result.Errors {
			err = errors.Join(fmt.Errorf("%v at %v:%v; ", result.Errors[i].Text, result.Errors[i].Location.File, result.Errors[i].Location.Line))
		}
	}
	return result.OutputFiles, err
}
