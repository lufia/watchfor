// +build: darwin dragonfly freebsd linux nacl netbsd openbsd

package main

import (
	"os/exec"
)

func System(cmd Command) error {
	c := exec.Command("bash", "-c", string(cmd))
	if err := c.Start(); err != nil {
		return err
	}
	if err := c.Wait(); err != nil {
		return err
	}
	return nil
}
