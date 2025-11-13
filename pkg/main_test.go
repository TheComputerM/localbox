package pkg_test

import (
	"errors"
	"os"
	"testing"

	"github.com/thecomputerm/localbox/internal"
)

func TestMain(m *testing.M) {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
	code := m.Run()
	os.Exit(code)
}
