package cli

import (
	"log"
	"os"

	"github.com/a-h/templ/cmd/templ/generatecmd"
)

func generateTempl() error {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	err := generatecmd.Run(generatecmd.Arguments{
		WorkerCount: 4,
	})
	null.Close()
	os.Stdout = sout
	os.Stderr = serr
	log.SetOutput(os.Stderr)
	return err
}
