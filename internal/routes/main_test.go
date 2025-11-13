package routes_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"
	"github.com/thecomputerm/localbox/internal"
	"github.com/thecomputerm/localbox/internal/routes"
)

const EXAMPLES_DIR = "../../examples"

type apitest struct {
	Method   string
	Endpoint string
	Input    string
	Expected string
}

func exampleToAPITest(filename string) (*apitest, error) {
	fullpath, err := filepath.Abs(filepath.Join(EXAMPLES_DIR, filename))
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(fullpath)
	if err != nil {
		return nil, err
	}
	input, _, found := strings.Cut(string(data), "---")
	if !found {
		return nil, fmt.Errorf("invalid example format")
	}

	suite := &apitest{}

	codeblockRegex := "```json[c]?" + `\s*({[\s\S]*?})\s*` + "```"

	inputRegex := regexp.MustCompile(`([A-Z]+):\s*(\/\S+)\s+` + codeblockRegex)
	matches := inputRegex.FindStringSubmatch(input)
	suite.Method = matches[1]
	suite.Endpoint = matches[2]
	suite.Input = matches[3]

	return suite, nil
}

func TestE2E(t *testing.T) {
	t.Skip()
	huma.DefaultArrayNullable = false
	examples, err := os.ReadDir(EXAMPLES_DIR)
	require.NoError(t, err)

	for _, example := range examples {
		if example.IsDir() || !strings.HasSuffix(example.Name(), ".md") {
			continue
		}
		testdata, err := exampleToAPITest(example.Name())
		require.NoError(t, err)
		_, api := humatest.New(t)

		routes.AddRoutes(api)
		t.Run(strings.TrimSuffix(example.Name(), ".md"), func(t *testing.T) {
			resp := api.Do(testdata.Method, testdata.Endpoint, strings.NewReader(testdata.Input))
			require.Less(t, resp.Result().StatusCode, 300)
			require.GreaterOrEqual(t, resp.Result().StatusCode, 200)
		})
	}
}

func TestMain(m *testing.M) {
	if err := internal.InitCGroup(); err != nil {
		panic(errors.Join(errors.New("failed to init cgroup"), err))
	}
	code := m.Run()
	os.Exit(code)
}
