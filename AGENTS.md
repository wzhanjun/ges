# GES — Go Echo Skeleton

Single Go module, Cobra CLI scaffolding tool. One subcommand: `ges new <name>`.

## Commands

```bash
go install github.com/wzhanjun/ges@latest    # install
ges new myproject                              # scaffold a project
```

## Key facts

- **Template source**: cloned at runtime from `https://github.com/wzhanjun/go-echo-skeleton.git`. Override via `GES_LAYOUT_REPO` env var. `--branch` flag defaults to `main`.
- **Template cache**: stored in `~/.ges/repo/` — subsequent runs `git pull` instead of cloning fresh.
- **Module path replacement**: naive `bytes.ReplaceAll` on every copied file (old module path → project name). No Go template engine.
- **Prompt behavior**: if target directory exists, user is prompted to override. If `ges new` is called without a name argument, an interactive prompt (`survey/v2`) asks for a project name.
- **No tests** exist in the repository.
- **No Makefile** or CI workflows in the repo (the *generated* project has a `make run`).
