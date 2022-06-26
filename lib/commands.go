package lib

import (
	"log"
	"os/exec"
	"strings"
)

func Command(cmd *exec.Cmd) error {
	var w strings.Builder
	cmd.Stderr = &w
	err := cmd.Run()
	if err != nil {
		log.Println(strings.TrimSpace(w.String()))
		return err
	}
	return nil
}
