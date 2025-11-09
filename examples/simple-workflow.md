POST: /execute

```jsonc
{
    "files": [{
        "content": "@TheComputerM on GitHub seems like a chill guy",
        "name": "truth.txt"
    }],
    "phases": [
        {
            "command": "cat truth.txt | cowsay",
            "packages": ["nixpkgs#cowsay","nixpkgs/nixos-25.05#busybox"]
        }
    ]
}
```

---

```jsonc
[
  {
    "time": 26,
    "wall_time": 65,
    "memory": 4148,
    "max_rss": 8192,
    "status": "OK",
    "message": "Executed",
    "exit_code": 0,
    "stdout": "______________________________________\n/ @TheComputerM on GitHub seems like a \\\n\\ chill guy                            /\n --------------------------------------\n        \\   ^__^\n         \\  (oo)\\_______\n            (__)\\       )\\/\\\n                ||----w |\n                ||     ||",
    "stderr": ""
  }
]
```