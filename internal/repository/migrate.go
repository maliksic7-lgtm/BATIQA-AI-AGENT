package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations reads all *.sql files in dir (sorted) and executes them in order.
// It is idempotent due to CREATE TABLE IF NOT EXISTS and INSERT IGNORE in seed.
// Uses simple file read + Exec; for Phase 2 this is sufficient vs golang-migrate.
func RunMigrations(db *sql.DB, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read migrations dir %s: %w", dir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", dir)
	}

	for _, f := range files {
		path := filepath.Join(dir, f)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", f, err)
		}
		sqlStr := strings.TrimSpace(string(content))
		if sqlStr == "" {
			continue
		}
		// Exec as a whole; MySQL driver supports multi-statements only if enabled.
		// We split by semicolon for safety and exec each statement.
		stmts := splitSQLStatements(sqlStr)
		for i, stmt := range stmts {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("migrate %s stmt %d failed: %w\nSQL: %s", f, i+1, err, stmt)
			}
		}
	}
	return nil
}

// splitSQLStatements splits SQL by semicolon and strips comment lines.
// Handles -- comments and empty lines; does not need to handle semicolons inside quotes for our schema.
func splitSQLStatements(s string) []string {
	lines := strings.Split(s, "\n")
	var cleaned []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "--") || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "/*") {
			continue
		}
		// Remove inline -- comments (simple, not inside quotes)
		if idx := strings.Index(trimmed, "--"); idx != -1 {
			// Keep only before comment, but avoid breaking strings containing --
			// For our files, inline -- only appears as comment start, so safe to strip
			trimmed = strings.TrimSpace(trimmed[:idx])
			if trimmed == "" {
				continue
			}
		}
		cleaned = append(cleaned, trimmed)
	}
	joined := strings.Join(cleaned, " ")
	parts := strings.Split(joined, ";")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if strings.TrimSpace(p) != "" {
			out = append(out, strings.TrimSpace(p))
		}
	}
	return out
}

// EnsureMigrationsTable is not needed for Phase 2 simple runner (no version tracking).
// We rely on IF NOT EXISTS for idempotency.
