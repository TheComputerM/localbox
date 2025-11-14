POST: /engine/custom/execute

```jsonc
{
    "engine": {
        "compile": {
            "packages": ["nixpkgs/nixos-25.05#go"],
            "command": "GOCACHE=/box GOOS=wasip1 GOARCH=wasm go build -o output.wasm main.go"
        },
        "execute": {
            "packages": ["nixpkgs/nixos-25.05#wasmtime"],
            "command": "wasmtime output.wasm"
        },
        "meta": {
            "run_file": "main.go"
        }
    },
    "files": [{
        "name": "@",
        "content": "package main\nimport \"fmt\"\nfunc main() {fmt.Println(\"Hello from Go compiled to WASM\")}"
    }]
}
```

---

```jsonc
{
  "$schema": "http://localhost:3000/schemas/SandboxPhaseResults.json",
  "time": 54,
  "wall_time": 37,
  "memory": 10952,
  "max_rss": 29084,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "Hello from C compiled to WASM",
  "stderr": ""
}
```