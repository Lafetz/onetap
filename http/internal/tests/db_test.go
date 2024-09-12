package integration_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Lafetz/loyalty_marketplace/internal/repository"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "loyality"
	DbUser = "user"
	DbPass = "password"
)

type TestDatabase struct {
	DbInstance *sql.DB
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	// remove test container
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *sql.DB, string, error) {

	var env = map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	//
	db, err := repository.OpenDB(fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName))
	if err != nil {
		log.Fatal(err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {

	// get location of test
	// _, path, _, ok := runtime.Caller(0)
	// if !ok {
	// 	return fmt.Errorf("failed to get path")
	// }
	goModRoot, err := findGoModRoot()
	if err != nil {

		return err
	}
	pathToMigrationFiles := fmt.Sprintf("%s/db/migrations", goModRoot)
	//filepath.Dir(path) + "/../../db/migrations"
	// log.Fatal(pathToMigrationFiles)
	databaseURL := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	log.Println("migration done")

	return nil
}
func findGoModRoot() (string, error) {
	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %v", err)
	}

	// Traverse upwards until go.mod is found
	for {
		// Check if go.mod exists in the current directory
		modPath := filepath.Join(workingDir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			return workingDir, nil
		}

		// Move up one directory level
		parentDir := filepath.Dir(workingDir)
		if parentDir == workingDir {
			break // Reached the root, stop searching
		}
		workingDir = parentDir
	}

	return "", fmt.Errorf("go.mod file not found")
}
