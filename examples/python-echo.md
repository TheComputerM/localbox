POST: /engine/python/execute

```jsonc
{
    "files": [{
        "content": "print(input())",
        "name": "@"
    }],
    "options": {
        "stdin": "change da world. my final message. goodbye"
    }
}
```

---

```jsonc
{
  "$schema": "http://localhost:3000/schemas/SandboxPhaseResults.json",
  "time": 24,
  "wall_time": 31,
  "memory": 14476,
  "max_rss": 10880,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "change da world. my final message. goodbye",
  "stderr": ""
}
```
