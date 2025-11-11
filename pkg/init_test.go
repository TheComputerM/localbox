package pkg_test

import (
	"errors"

	"github.com/thecomputerm/localbox/internal"
)

// initializes localbox for tests
func init() {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
}
