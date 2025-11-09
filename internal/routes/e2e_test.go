package routes_test

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const EXAMPLES_DIR = "../../examples"

type apitest struct {
	Method   string
	Endpoint string
	Input    any
	Expected any
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
	input, output, found := strings.Cut(string(data), "---")
	if !found {
		return nil, fmt.Errorf("invalid example format")
	}

	suite := &apitest{}

	codeRegex := "```json[c]?" + `\s*({[\s\S]*?})\s*` + "```"

	inputRegex := regexp.MustCompile(`([A-Z]+):\s*(\/\S+)\s+` + codeRegex)
	matches := inputRegex.FindStringSubmatch(input)
	suite.Method = matches[1]
	suite.Endpoint = matches[2]
	suite.Input = matches[3]

	outputRegex := regexp.MustCompile(codeRegex)
	outputMatches := outputRegex.FindStringSubmatch(output)
	suite.Expected = outputMatches[1]

	return suite, nil
}

func TestE2E(t *testing.T) {
	examples, err := os.ReadDir(EXAMPLES_DIR)
	require.NoError(t, err)

	for _, example := range examples {
		if example.IsDir() || !strings.HasSuffix(example.Name(), ".md") {
			continue
		}
		data, err := exampleToAPITest(example.Name())
		require.NoError(t, err)
		fmt.Println(data)
	}
}
