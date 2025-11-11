POST: /engine/cpp/execute

```jsonc
{
    "files": [{
        "content": "#include <iostream>\nint main() {std::cout << \"Hello World\\n\";};",
        "name": "@"
    }]
}
```

---

```jsonc
{
  "$schema": "http://localhost:3000/schemas/SandboxPhaseResults.json",
  "time": 2,
  "wall_time": 6,
  "memory": 572,
  "max_rss": 3456,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "Hello World",
  "stderr": ""
}
```
