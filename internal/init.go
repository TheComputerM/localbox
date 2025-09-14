package internal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func initCGroup() error {
	cgroups := "/sys/fs/cgroup"
	name := "isolate"

	if info, err := os.Stat(filepath.Join(cgroups, name)); !errors.Is(err, os.ErrNotExist) && info.IsDir() {
		return nil
	}

	if err := os.Mkdir(filepath.Join(cgroups, name), 0755); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(cgroups, name, "cgroup.subtree_control"), []byte("+cpuset +memory"), 0644); err != nil {
		return err
	}
	return nil
}

func init() {
	if os.Getuid() != 0 {
		log.Fatal("LocalBox must be run as root")
	}
	if err := initCGroup(); err != nil {
		log.Fatal(errors.Join(fmt.Errorf("couldn't init cgroup"), err))
	}
}
