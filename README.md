# hitch-generator

A systemd user generator that creates service units for Podman Compose files.

## How it works

systemd invokes generators during boot/reload:

```
hitch-generator <normal-dir> [early-dir] [late-dir]
```

The generator scans `~/.config/containers/compose/` (or `$HITCH_COMPOSE_DIR`) for `*.yml`/`*.yaml` files and produces a `.service` unit for each one. Services are wired into `default.target.wants` so they start automatically.

Each generated unit:
- Ensures external Podman networks exist (`ExecStartPre`)
- Runs `podman compose up -d` on start
- Runs `podman compose down` on stop
- Supports reload via `podman compose up -d --force-recreate`

## Install

```bash
make install
# Installs to ~/.config/systemd/user-generators/hitch-generator
```

Then reload systemd to pick up the generator:

```bash
systemctl --user daemon-reload
```

## Build & Test

```bash
make build
make test
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `HITCH_COMPOSE_DIR` | `$HOME/.config/containers/compose/` | Directory to scan for compose files |

## License

MIT
