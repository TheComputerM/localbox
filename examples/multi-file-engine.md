POST: /engine/python/execute

```jsonc
{
    "files": [{
        "content": "from utils import add\nprint(add(2,3))",
        "name": "@"
    }, {
        "content": "def add(a, b):\n\treturn a+b",
        "name": "utils.py"
    }]
}
```

---

```jsonc
{
  "$schema": "http://localhost:3000/schemas/SandboxPhaseResults.json",
  "time": 20,
  "wall_time": 23,
  "memory": 4600,
  "max_rss": 11008,
  "status": "OK",
  "message": "Executed",
  "exit_code": 0,
  "stdout": "5",
  "stderr": ""
}
```