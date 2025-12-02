package postgres

import (
	"context"
	"embed"
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// sqlFs contains the SQL files for the migrations.
//
//go:embed sql
var sqlFs embed.FS

const (
	// sqlFolder is the folder containing the SQL files. Should match the folder name in sqlFs.
	sqlFolder = "sql"
)

// New creates a new pgxpool.Pool.
func New(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// sqlFilePattern matches the file name pattern for SQL files.
var sqlFilePattern = regexp.MustCompile(`^(\d+)_.*\.sql$`)

// getSqlFileNames returns the list of SQL file names.
func getSqlFileNames() (fileNames []string, err error) {
	files, err := sqlFs.ReadDir(sqlFolder)
	if err != nil {
		return nil, err
	}

	fileNames = make([]string, 0, len(files))
	for _, file := range files {
		if !sqlFilePattern.MatchString(file.Name()) {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}

	return fileNames, nil
}

// getSortedSqlFileNames returns the sorted list of SQL file names.
func getSortedSqlFileNames() (fileNames []string, err error) {
	fileNames, err = getSqlFileNames()
	if err != nil {
		return nil, err
	}

	slices.Sort(fileNames)
	if err != nil {
		return nil, err
	}

	return fileNames, nil
}

// Init initialises the database. Sets up the schema and functions. Uses
// the migration table to track the applied migrations.
func Init(pool *pgxpool.Pool, ctx context.Context) error {
	fileNames, err := getSortedSqlFileNames()
	if err != nil {
		return err
	}

	if len(fileNames) == 0 {
		panic("no SQL files found")
	}

	// Note that this assumes that the first file sorted by name is the migration table schema.
	err = createMigrationsTable(pool, ctx, filepath.Join(sqlFolder, fileNames[0]))
	if err != nil {
		return err
	}
	fileNames = fileNames[1:]

	count, err := getLastMigration(pool, ctx)
	if err != nil {
		return err
	} else if count == nil {
		return applyMigrations(pool, ctx, fileNames)
	}

	{ // Check if we're already at the latest migration
		lastKey, err := getKey(fileNames[len(fileNames)-1])
		if err != nil {
			return err
		}

		if *lastKey == *count {
			return nil
		}
	}

	return applyMigrations(pool, ctx, fileNames[*count:])
}

// getKey returns the migration key from the file name.
func getKey(filePath string) (*int, error) {
	key := sqlFilePattern.FindStringSubmatch(filePath)[1]
	intKey, err := strconv.Atoi(key)
	if err != nil {
		return nil, err
	}

	return &intKey, nil
}

// createMigrationsTable creates the migration table if it doesn't exist.
func createMigrationsTable(pool *pgxpool.Pool, ctx context.Context, filePath string) error {
	b, err := sqlFs.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, string(b))
	return err
}

// getLastMigration returns the last migration applied to the database.
func getLastMigration(pool *pgxpool.Pool, ctx context.Context) (count *int, err error) {
	const stmt = "SELECT MAX(id) FROM schema_migrations"
	err = pool.QueryRow(ctx, stmt).Scan(&count)
	if err != nil {
		return nil, err
	}

	return count, nil
}

// applyMigrations applies the migrations to the database.
func applyMigrations(pool *pgxpool.Pool, ctx context.Context, fileNames []string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, fileName := range fileNames {
		filePath := filepath.Join(sqlFolder, fileName)
		b, err := sqlFs.ReadFile(filePath)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, string(b))
		if err != nil {
			return fmt.Errorf(
				"error executing %s, error: %v",
				filePath,
				err,
			)
		}

		key, err := getKey(fileName)
		if err != nil {
			return err
		}

		{
			const stmt = "INSERT INTO schema_migrations (id) VALUES ($1)"
			_, err = tx.Exec(ctx, stmt, *key)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
