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
    - Isolate is not run as root and doesn't share cgroup with host system which can lead to inaccurate program runtime metrics.
    - All scripts are run in bash which adds additional overhead to measuring program runtime metrics.
3. In Judge0
    - Only supports outdated linux cgroups v1, you will have to change kernel options in any machine with a recent linux kernel.
    - Cannot take multiple files as input.