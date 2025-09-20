# Developing

Use the provided devcontainer and run `sudo su` to switch to root so that localbox can manipulate cgroups. [Air](https://github.com/air-verse/air) is also installed in the devcontainer for faster workflow.

# Publishing

GitHub action will automatically publish container image when new tag is pushed

```
git tag v${VERSION}
git push origin v${VERSION}
```