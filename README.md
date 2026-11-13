# quadlet-compose-generator

`quadlet-compose` manages Podman Compose projects through systemd user services.
Its generator, `quadlet-compose-generator`, discovers Compose YAML files and
generates a systemd service unit for each project, so you can manage Compose
projects with `systemctl --user`.

The project is maintained by [kuyacarlo](https://github.com/kuyacarlo).

## How it works

systemd invokes the generator during user manager startup and reload:

```text
quadlet-compose-generator <normal-dir> [early-dir] [late-dir]
```

The generator scans `~/.config/containers/compose/` (or the directory set by
`QUADLET_COMPOSE_DIR`) recursively for `*.yml` and `*.yaml` files. It writes a
`.service` unit for each Compose project and links each unit from
`default.target.wants` so it starts with the user's default target.

Each generated unit:

- Creates missing external Podman networks before startup
- Runs `podman compose up -d` on start
- Runs `podman compose down` on stop
- Supports reload through `podman compose up -d --force-recreate`

## Install

Build and install the generator into the systemd user generator directory:

```bash
make install
systemctl --user daemon-reload
```

`make install` places `quadlet-compose-generator` in
`~/.config/systemd/user-generators/`.

## Build and test

```bash
make build
make test
```

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `QUADLET_COMPOSE_DIR` | `$HOME/.config/containers/compose/` | Directory to scan for Compose YAML files |
| `HITCH_COMPOSE_DIR` | — | Legacy alias, used only when `QUADLET_COMPOSE_DIR` is unset |

## License

MIT
