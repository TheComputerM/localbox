POST: /engine/custom/execute

```jsonc
{
    "engine": {
        "compile": {
            "packages": ["nixpkgs/nixos-25.05#zig"],
            "command": "ZIG_GLOBAL_CACHE_DIR=\"$PWD/cache\" zig cc -target wasm32-wasi main.c -o output.wasm"
        },
        "execute": {
            "packages": ["nixpkgs/nixos-25.05#wasmtime"],
            "command": "wasmtime output.wasm"
        },
        "meta": {
            "run_file": "main.c"
        }
    },
    "files": [{
        "name": "@",
        "content": "#include <stdio.h>\nint main() {printf(\"Hello from C compiled to WASM\");return 0;}"
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