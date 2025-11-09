package pkg

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/thecomputerm/localbox/internal/utils"
)

type Sandbox int

// directory path where the sandboxes are stored, specified in isolate.conf
const sandbox_base = "/var/lib/localbox"

func (s Sandbox) ID() int {
	return int(s)
}

func (s Sandbox) String() string {
	return strconv.Itoa(s.ID())
}

func (s Sandbox) BoxPath() string {
	return filepath.Join(sandbox_base, s.String(), "box")
}

func (s Sandbox) metadataFilePath() string {
	return "/tmp/box" + s.String() + "-meta.txt"
}

func (s Sandbox) Init() error {
	if err := exec.Command(
		Globals.IsolateBin,
		"--cg",
		"--box-id="+s.String(),
		"--init",
	).Run(); err != nil {
		return errors.Join(fmt.Errorf("could not init sandbox %d", s.ID()), err)
	}
	return nil
}

func (s Sandbox) Cleanup() error {
	if err := exec.Command(
		Globals.IsolateBin,
		"--cg",
		"--box-id="+s.String(),
		"--cleanup",
	).Run(); err != nil {
		return errors.Join(fmt.Errorf("could not cleanup sandbox %d", s.ID()), err)
	}
	return os.RemoveAll(s.metadataFilePath())
}

type SandboxFile struct {
	Name     string `json:"name" doc:"Path of the file within the sandbox" example:"hello.txt"`
	Content  string `json:"content" doc:"Content of the file" example:"Hello World"`
	Encoding string `json:"encoding,omitempty" doc:"Encoding of the content field" enum:"utf8,base64,hex" default:"utf8" `
}

func (s Sandbox) Mount(files []SandboxFile) error {
	for _, file := range files {
		if !filepath.IsLocal(file.Name) {
			return fmt.Errorf("file %s tries to escape from sandbox", file.Name)
		}
	}

	for _, file := range files {
		location := filepath.Join(s.BoxPath(), file.Name)
		err := os.MkdirAll(filepath.Dir(location), 0755)
		if err != nil {
			return err
		}
		f, err := os.Create(location)
		if err != nil {
			return err
		}
		defer f.Close()

		content := file.Content
		switch file.Encoding {
		case "utf8", "":
		case "base64":
			decoded, err := base64.StdEncoding.DecodeString(content)
			if err != nil {
				return errors.Join(fmt.Errorf("could not decode content for %s as %s", file.Name, file.Encoding), err)
			}
			content = string(decoded)
		case "hex":
			decoded, err := hex.DecodeString(content)
			if err != nil {
				return errors.Join(fmt.Errorf("could not decode content for %s as %s", file.Name, file.Encoding), err)
			}
			content = string(decoded)
		default:
			return fmt.Errorf("unknown encoding %s for file %s", file.Encoding, file.Name)
		}

		_, err = f.WriteString(content)
		if err != nil {
			return err
		}
	}

	return nil
}

type SandboxPhaseMetadata struct {
	Time     int    `json:"time" doc:"Run time of the program in milliseconds" example:"500"`
	WallTime int    `json:"wall_time" doc:"Wall time of the program in milliseconds" example:"1000"`
	Memory   int    `json:"memory" doc:"Total memory use by the whole control group in KB" example:"256"`
	MaxRSS   int    `json:"max_rss" doc:"Maximum resident set size of the program in KB" example:"128"`
	Status   string `json:"status" doc:"Two-letter status code" example:"OK" enum:"OK,RE,SG,TO,XX,OE,CE"`
	Message  string `json:"message" doc:"Human-readable message" example:"Executed"`
	ExitCode int    `json:"exit_code" doc:"Exit code/signal from the program" example:"0"`
}

// Helper to parse metadata file created by isolate
func (s Sandbox) parseMetadata() (*SandboxPhaseMetadata, error) {
	file, err := os.ReadFile(s.metadataFilePath())
	if err != nil {
		return nil, err
	}

	output := &SandboxPhaseMetadata{
		Status:  "OK",
		Message: "Executed",
	}

	for _, line := range strings.Split(string(file), "\n") {
		key, value, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		switch key {
		case "status":
			output.Status = value
		case "message":
			output.Message = value
		case "time":
			time, _ := strconv.ParseFloat(value, 64)
			output.Time = int(time * 1000)
		case "time-wall":
			wallTime, _ := strconv.ParseFloat(value, 64)
			output.WallTime = int(wallTime * 1000)
		case "cg-mem":
			output.Memory, _ = strconv.Atoi(value)
		case "max-rss":
			output.MaxRSS, _ = strconv.Atoi(value)
		case "exitcode", "exitsig":
			output.ExitCode, _ = strconv.Atoi(value)
		}
	}
	return output, nil
}

type SandboxPhaseResults struct {
	SandboxPhaseMetadata
	Stdout string `json:"stdout" doc:"stdout of the program" example:"program output"`
	Stderr string `json:"stderr" doc:"stderr of the program" example:""`
}

type SandboxCommand struct {
	Command  string   `json:"command" doc:"Command to execute in the sandbox" example:"cat hello.txt"`
	Packages []string `json:"packages,omitempty" doc:"Nix packages to install in the sandbox" example:"nixpkgs#cowsay,nixpkgs/nixos-25.05#busybox"`
}

type SandboxPhase struct {
	SandboxCommand
	SkipShell bool `json:"skip_shell,omitempty" doc:"Doesn't use a shell to run the command to if true, can be used to get more accurate results" default:"false"`
}

type SandboxPhaseOptions struct {
	MemoryLimit   int               `json:"memory_limit,omitempty" doc:"Maximum total memory usage allowed by the whole control group in KB, '-1' for no limit" default:"-1"`
	TimeLimit     int               `json:"time_limit,omitempty" doc:"Maximum CPU time of the program in milliseconds, '-1' for no limit" default:"5000"`
	WallTimeLimit int               `json:"wall_time_limit,omitempty" doc:"Maximum wall time of the program in milliseconds, '-1' for no limit" default:"10000"`
	FilesLimit    int               `json:"files_limit,omitempty" doc:"Maximum number of open files allowed in the sandbox, '-1' for no limit" default:"64"`
	FileSizeLimit int               `json:"file_size_limit,omitempty" doc:"Maximum size a file created/modified in the sandbox in KB, -1 for no limit" default:"10000"`
	ProcessLimit  int               `json:"process_limit,omitempty" doc:"Maximum number of processes allowed in the sandbox" default:"64"`
	Network       bool              `json:"network,omitempty" doc:"Whether to enable network access in the sandbox" default:"false"`
	Stdin         string            `json:"stdin,omitempty" doc:"Text to pass into stdin of the program" default:""`
	BufferLimit   int               `json:"buffer_limit,omitempty" doc:"Maximum kilobytes to capture from stdout and stderr" default:"8"`
	Environment   map[string]string `json:"environment,omitempty" doc:"Environment variables to set in the sandbox" example:"{}"`
}

// Run a command inside the isolated sandbox which the listed nix packages installed
func (s Sandbox) Run(
	phase *SandboxPhase,
	options *SandboxPhaseOptions,
) (*SandboxPhaseResults, error) {
	if options.Stdin != "" {
		if err := os.WriteFile(
			filepath.Join(s.BoxPath(), "stdin.txt"),
			[]byte(options.Stdin),
			0600,
		); err != nil {
			return nil, err
		}
	}

	var cmd *exec.Cmd
	isolateCommand := buildIsolateCommand(s, phase, options)
	if len(phase.Packages) > 0 {
		cmd = exec.Command("nix", buildNixShell(phase.Packages, isolateCommand)...)
	} else {
		cmd = exec.Command(isolateCommand[0], isolateCommand[1:]...)
	}

	stdout := utils.NewLimitedWriter(uint64(options.BufferLimit) * 1024)
	cmd.Stdout = stdout

	stderr := utils.NewLimitedWriter(uint64(options.BufferLimit) * 1024)
	cmd.Stderr = stderr

	isOutputBufferExceeded := false

	if err := cmd.Run(); err != nil {
		if _, found := os.Stat(s.metadataFilePath()); errors.Is(found, os.ErrNotExist) {
			// throw error if metadata file doesn't exist
			return nil, errors.Join(fmt.Errorf("sandbox error: %s", cmd.String()), errors.New(stderr.String()), err)
		} else if errors.Is(err, utils.ErrWriteLimitExceeded) {
			isOutputBufferExceeded = true
		}
	}

	results := &SandboxPhaseResults{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	meta, err := s.parseMetadata()
	if err != nil {
		return nil, err
	}

	results.SandboxPhaseMetadata = *meta

	if isOutputBufferExceeded {
		results.Status = "OE"
		results.Message = "Output buffer limit exceeded"
	}

	return results, nil
}

// Run a command on the sandbox in a shell without any isolation or limits
func (s Sandbox) UnsafeRun(command *SandboxCommand) error {
	wrapper := []string{Globals.ShellBin, "-c", command.Command}
	cmd := exec.Command("nix", buildNixShell(command.Packages, wrapper)...)
	cmd.Dir = s.BoxPath()
	return cmd.Run()
}

func buildNixShell(packages, run []string) []string {
	args := []string{"shell", "--quiet", "-i", "-k", "ISOLATE_CONFIG_FILE"}
	args = append(args, packages...)
	args = append(args, "--command")
	args = append(args, run...)
	return args
}

// Helper to build isolate command string with options
func buildIsolateCommand(
	s Sandbox,
	phase *SandboxPhase,
	options *SandboxPhaseOptions,
) []string {

	filesLimit := options.FilesLimit
	if filesLimit == -1 {
		filesLimit = 0 // 0 means no limit in isolate
	}

	command := []string{
		Globals.IsolateBin,
		"--cg",
		"-s",
		"--meta=" + s.metadataFilePath(),
		"--dir=/nix=/nix",
		"--dir=/etc=/etc:noexec",
		"--box-id=" + s.String(),
		"--open-files=" + strconv.Itoa(filesLimit),
		"--processes=" + strconv.Itoa(options.ProcessLimit),
		"-e",
		"--env=HOME=/tmp",
	}

	for key, value := range options.Environment {
		command = append(command, fmt.Sprintf("--env=%s=%s", key, value))
	}

	if options.Stdin != "" {
		command = append(command, "--stdin=/box/stdin.txt")
	}

	if options.TimeLimit != -1 {
		command = append(command,
			"--time="+strconv.FormatFloat(float64(options.TimeLimit)/1000, 'f', 3, 64),
		)
	}

	if options.WallTimeLimit != -1 {
		command = append(command,
			"--wall-time="+strconv.FormatFloat(float64(options.WallTimeLimit)/1000, 'f', 3, 64),
		)
	}

	if options.MemoryLimit != -1 {
		command = append(command,
			"--cg-mem="+strconv.Itoa(options.MemoryLimit),
		)
	}

	if options.FileSizeLimit != -1 {
		command = append(command,
			"--fsize="+strconv.Itoa(options.FileSizeLimit),
		)
	}

	if options.Network {
		command = append(command, "--share-net")
	}

	command = append(command, "--run", "--")
	if phase.SkipShell {
		command = append(command, phase.Command)
	} else {
		command = append(command, Globals.ShellBin, "-c", phase.Command)
	}

	return command
}
