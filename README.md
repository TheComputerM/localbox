# LocalBox

**WIP: Work in Progress**

## What?

LocalBox is a **easy-to-host**, **general purpose** and **fast** code execution system for running **untrusted** code in sandboxes.

## Why?

[Piston](https://github.com/engineer-man/piston) and [Judge0](https://github.com/judge0/judge0) are already well established and tested code execution systems, so why the need to reinvent the wheel?

Both of these projects have some *issues* which did not fit my usecase, namely:

1. In both Piston and Judge0
    - You cannot install the latest versions of compilers and runtimes, you are only limited to the versions which the creators bothered to add and update.
    - Very difficult to run arbitrary commands for specific usecases.
2. In Piston
    - Isolate doesn't share cgroups with host system which can lead to inaccurate program limits and runtime metrics.
    - All programs are run in bash which adds additional overhead to measuring program runtime metrics.
    - You have to compile/build the runtime on your system for many languages instead of getting prebuilt binaries.
3. In Judge0
    - Only supports outdated linux cgroups v1, you will have to change kernel options in any machine with a recent linux kernel.
    - Cannot take multiple files as input.

## Getting Started

### With Docker

The image exposes port 2000 by default, the `--privileged --cgroupns=host` are important as they allow localjudge to manipulate cgroups.

```sh
docker run --rm -it --privileged --cgroupns=host -p 2000:2000 localbox
```

## Usage

You can visit http://localhost:2000/docs to get the full API documentation.

### Using an Engine

`GET /engines`: List all the available engines.

Here is a [list of available languages/runtimes](./engines/) packaged as engines, these provide preconfigured compile and execute steps so you don't have to set stuff up.

`POST /engine/{engine_name}`: Execute a predefined engine with an execution phase whose options can be overridden.

```jsonc
// Request Body
{
  "options": {
    "stdin": "hello world"
  },
  "files": [
    {
      "content": "print(input())",
      "encoding": "utf8",
      "name": "@"
    }
  ]
}

// Response
{
  "$schema": "http://localhost:9000/schemas/SandboxPhaseResults.json",
  "time": 20,
  "wall_time": 42,
  "memory": 4076,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "hello world",
  "stderr": ""
}
```

Here are the options you can provide the sandbox with:

```go
type SandboxPhaseOptions struct {
	MemoryLimit  int               `json:"memory_limit,omitempty" doc:"Maximum total memory usage allowed by the whole control group in KB, '-1' for no limit" default:"-1"`
	TimeLimit    int               `json:"time_limit,omitempty" doc:"Maximum CPU time of the program in milliseconds, '-1' for no limit" default:"5000"`
	FilesLimit   int               `json:"files_limit,omitempty" doc:"Maximum number of open files allowed in the sandbox, '-1' for no limit" default:"64"`
	ProcessLimit int               `json:"process_limit,omitempty" doc:"Maximum number of processes allowed in the sandbox" default:"64"`
	Network      bool              `json:"network,omitempty" doc:"Whether to enable network access in the sandbox" default:"false"`
	Stdin        string            `json:"stdin,omitempty" doc:"Text to pass into stdin of the program" default:""`
	BufferLimit  int               `json:"buffer_limit,omitempty" doc:"Maximum kilobytes to capture from stdout and stderr" default:"64"`
	Environment  map[string]string `json:"environment,omitempty" doc:"Environment variables to set in the sandbox" example:"{}"`
}
```