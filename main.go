package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: hitch-generator <normal-dir> [early-dir] [late-dir]\n")
		os.Exit(1)
	}

	normalDir := os.Args[1]

	composeDir := os.Getenv("HITCH_COMPOSE_DIR")
	if composeDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "hitch-generator: cannot determine home directory: %v\n", err)
			os.Exit(1)
		}
		composeDir = filepath.Join(home, ".config", "containers", "compose")
	}

	if err := run(composeDir, normalDir); err != nil {
		fmt.Fprintf(os.Stderr, "hitch-generator: %v\n", err)
		os.Exit(1)
	}
}

// run scans composeDir for compose files and generates systemd units into normalDir.
func run(composeDir, normalDir string) error {
	files, err := scanComposeFiles(composeDir)
	if err != nil {
		return fmt.Errorf("scanning %s: %w", composeDir, err)
	}

	if len(files) == 0 {
		return nil
	}

	// Ensure the wants directory exists.
	wantsDir := filepath.Join(normalDir, "default.target.wants")
	if err := os.MkdirAll(wantsDir, 0755); err != nil {
		return fmt.Errorf("creating wants dir: %w", err)
	}

	for _, f := range files {
		info, err := parseComposeFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "hitch-generator: skipping %s: %v\n", f, err)
			continue
		}

		unit := generateUnit(info)
		unitName := info.Name + ".service"
		unitPath := filepath.Join(normalDir, unitName)

		if err := os.WriteFile(unitPath, []byte(unit), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", unitPath, err)
		}

		// Create wants symlink: default.target.wants/<name>.service → ../<name>.service
		linkPath := filepath.Join(wantsDir, unitName)
		target := filepath.Join("..", unitName)
		// Remove any existing symlink (ignore error if not present).
		os.Remove(linkPath)
		if err := os.Symlink(target, linkPath); err != nil {
			return fmt.Errorf("creating symlink %s: %w", linkPath, err)
		}
	}

	return nil
}

// scanComposeFiles returns all *.yml and *.yaml files in dir, recursively.
func scanComposeFiles(dir string) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}

	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext == ".yml" || ext == ".yaml" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
