package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pborman/getopt/v2"
	"github.com/yeluolei/dnsmasq-ha-web/frontend"
	"github.com/yeluolei/dnsmasq-ha-web/migrations"
	"go.uber.org/zap"

	_ "github.com/mattn/go-sqlite3"
)

type (
	host struct {
		ID      int    `json:"id"`
		IP      string `json:"ip"`
		FQDN    string `json:"fqdn"`
		COMMENT string `json:"comment"`
	}
)

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(frontend.EmbededFiles, "dist")
	if err != nil {
		panic(err)
	}
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("path=%q, isDir=%v\n", path, d.IsDir())
		return nil
	})
	return http.FS(fsys)
}

var dbFile = "./hosts.db"
var hostFile = "/etc/dnsmasq-ha-web/hosts"
var port = 80

func init() {
	getopt.Flag(&dbFile, 'f', "The sqlite3 db file path")
	getopt.Flag(&hostFile, 'h', "The hosts file used for dnsmasq")
	getopt.Flag(&port, 'p', "The server port to listening")
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Parse the program arguments
	getopt.Parse()

	absDbFile, err := filepath.Abs(dbFile)
	if err != nil {
		sugar.Fatal(err)
	}

	sugar.Infof("Start with SQL lite file path: %s", absDbFile)

	// Create and use an existing database instance.
	db, err := sql.Open("sqlite3", absDbFile)
	if err != nil {
		sugar.Fatal(err)
	}
	defer db.Close()

	migrationFiles, err := iofs.New(migrations.Migrations, ".")
	if err != nil {
		sugar.Fatal(err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", migrationFiles, fmt.Sprintf("sqlite3://%s", absDbFile))
	if err != nil {
		sugar.Fatal(err)
	}

	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		sugar.Fatal(err)
	}

	// Create DB layer
	dbAPI := HostAPI{
		db: db,
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Handle hosts requests
	e.POST("/hosts", func(c echo.Context) error {
		return createHost(c, &dbAPI)
	})

	e.PUT("/hosts/:id", func(c echo.Context) error {
		return updateHost(c, &dbAPI)
	})

	e.GET("/hosts", func(c echo.Context) error {
		return getAllHosts(c, &dbAPI)
	})

	e.DELETE("/hosts/:id", func(c echo.Context) error {
		return deleteHost(c, &dbAPI)
	})

	e.POST("/generate", func(c echo.Context) error {
		return generateHosts(c, &dbAPI, hostFile)
	})

	// Handle the static web assets here
	assetHandler := http.FileServer(getFileSystem())
	e.GET("/*", echo.WrapHandler(assetHandler))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
