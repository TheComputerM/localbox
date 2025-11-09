POST: /engine/python/execute

```jsonc
{
    "files": [{
        "content": "print('hello world')",
        "name": "@"
    }]
}
```

---

```jsonc
{
  "$schema": "http://localhost:3000/schemas/SandboxPhaseResults.json",
  "time": 31,
  "wall_time": 1287,
  "memory": 15656,
  "max_rss": 10880,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "hello world",
  "stderr": ""
}
```