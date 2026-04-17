//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"warehouse/config"
	"warehouse/logger"
	"warehouse/modules/db"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDB struct {
	DB        *sqlx.DB
	Container testcontainers.Container
	Ctx       context.Context
}

func StartTestDB(t *testing.T) *TestDB {
	//сначала находим путь до енв и миграций, также запускаем логгер и конфиги и загружаем миграции.
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	ctx := context.Background()
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	envPath, err := findEnvPath(wd)
	if err != nil {
		panic(err)
	}
	migrationsDir, err := findDirUpwards(wd, "migrations")
	if err != nil {
		panic(err)
	}

	appConfig := config.MustLoadConfig(envPath)
	log := logger.Initlogger(appConfig.LogLevel, appConfig.Production)

	files, err := LoadMigrations(migrationsDir)
	if err != nil {
		log.Error("Failed to load migrations from /migrations, FATAL: ", slog.Any("err", err))
		panic(err)
	}
	//создаем контейнер
	/////////////////////////////////////////////////////////////
	req := testcontainers.ContainerRequest{
		Image: "postgres:15",
		Env: map[string]string{
			"POSTGRES_DB":       appConfig.DB.Name,
			"POSTGRES_USER":     appConfig.DB.User,
			"POSTGRES_PASSWORD": appConfig.DB.Password,
		},
		Files:        files,
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Error("Failed to start postgres container, FATAL: ", slog.Any("err", err))
		panic(err)
	}

	t.Cleanup(func() {
		err := pgContainer.Terminate(ctx)
		if err != nil {
			log.Error("Failed to terminate postgres container: ", slog.Any("err", err))
		}
	})

	//Подключаемся и добавляем все в testdb
	//////////////////////////////////////////////////////////////////////////////////////////
	newDBConf := appConfig.DB
	newDBConf.Host, err = pgContainer.Host(ctx)
	if err != nil {
		log.Error("Failed to change host on config, FATAL: ", slog.Any("err", err))
		panic(err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Error("Failed to get port from PGContainer, FATAL: ", slog.Any("err", err))
		panic(err)
	}
	newDBConf.Port = port.Port()

	testSqlxDB, err := db.NewSqlDB(log, &newDBConf)
	if err != nil {
		log.Error("failed to connect PG on tests, Fatal Error", slog.Any("err", err))
		panic(err)
	}

	return &TestDB{
		DB:        testSqlxDB,
		Container: pgContainer,
		Ctx:       ctx,
	}

}

//функция для поиска файлов миграций

func LoadMigrations(dir string) ([]testcontainers.ContainerFile, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []testcontainers.ContainerFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		if !strings.HasSuffix(name, ".sql") {
			continue
		}

		migrationFiles = append(migrationFiles, testcontainers.ContainerFile{
			HostFilePath:      filepath.Join(dir, name),
			ContainerFilePath: "/docker-entrypoint-initdb.d/" + name,
			FileMode:          0644,
		})
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].ContainerFilePath < migrationFiles[j].ContainerFilePath
	})

	return migrationFiles, nil

}

//функция для поиска env. поднимается вверх для поиска migrations.

func findEnvPath(startDir string) (string, error) {
	dir := startDir

	for {
		candidate := filepath.Join(dir, ".env")
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf(".env not found starting from %s", startDir)
		}
		dir = parent
	}
}

//функция для поиска директорий, используется для поиска миграций, поднимается пока не найдет нужную директорию.

func findDirUpwards(startDir, dirname string) (string, error) {
	dir := startDir

	for {
		candidate := filepath.Join(dir, dirname)
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s dir not found starting from %s", dirname, startDir)
		}
		dir = parent
	}
}
