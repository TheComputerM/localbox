package internal

import (
	"errors"
	"os"
	"path/filepath"
)

func InitCGroup() error {
	cgroups := "/sys/fs/cgroup"
	name := "localbox"

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
