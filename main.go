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

		unit := generateUnit(info, composeDir)
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

// scanComposeFiles returns all *.yml and *.yaml files in dir.
func scanComposeFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		if ext == ".yml" || ext == ".yaml" {
			files = append(files, filepath.Join(dir, e.Name()))
		}
	}
	return files, nil
}
