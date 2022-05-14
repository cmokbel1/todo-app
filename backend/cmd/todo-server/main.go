package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"

	"github.com/cmokbel1/todo-app/backend/crypto"
	"github.com/cmokbel1/todo-app/backend/postgres"
	"github.com/cmokbel1/todo-app/backend/todo"
)

var (
	version = "NA"
	commit  = "NA"
	date    = "NA"
)

func main() {
	todo.Build.Version = version
	todo.Build.Commit = commit
	todo.Build.Date = date

	logger := todo.NewLogger()
	logger.SetOutput(os.Stderr)
	logger.Info(todo.BuildDetails())
	os.Exit(realMain(logger))
}

func realMain(logger todo.Logger) int {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() { <-c; cancel() }()

	app := NewApp()
	if err := app.ParseFlagsAndLoadConfig(ctx, os.Args[1:]); err != nil {
		logger.Error(err.Error())
		return 1
	}

	if err := app.Run(ctx); err != nil {
		logger.Error(err.Error())
		if err = app.Close(); err != nil {
			logger.Error(err.Error())
		}
		return 1
	}
	<-ctx.Done()

	if err := app.Close(); err != nil {
		logger.Error(err.Error())
		return 1
	}
	return 0
}

func NewApp() *App {
	return &App{
		Logger: todo.NewLogger(),
		DB:     postgres.New(""),
	}
}

type App struct {
	Config Config

	Logger todo.Logger
	DB     *postgres.DB
}

func (app *App) Run(ctx context.Context) error {
	logger := todo.NewLogger()
	if app.Config.Log.Enabled {
		logger.SetLevel(app.Config.Log.Level)
		logger.SetOutput(os.Stderr)
	}
	app.Logger = logger

	app.DB = postgres.New(app.Config.DB.DSN)
	app.DB.EnableQueryLogging = app.Config.DB.EnableQueryLogging
	app.DB.Logger = app.Logger
	if err := app.DB.Open(ctx); err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}

	if err := app.DB.Migrate(); err != nil {
		return fmt.Errorf("failed to migrate db: %v", err)
	}

	return nil
}

func (app *App) ParseFlagsAndLoadConfig(ctx context.Context, args []string) error {
	var configFile string

	fs := flag.NewFlagSet("todo", flag.ContinueOnError)
	fs.StringVar(&configFile, "config", os.Getenv("TODO_CONFIG"), "path to the config file")

	if err := fs.Parse(args); err != nil {
		return err
	} else if app.Config, err = LoadConfig(configFile); err != nil {
		return err
	}

	return nil
}

func (app *App) Close() error {
	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			return err
		}
	}

	return nil
}

type Config struct {
	DB struct {
		DSN                string `json:"dsn"`
		EnableQueryLogging bool   `json:"query_logging_enabled"`
	} `json:"db"`

	HTTP struct {
		Addr   string  `json:"addr"`
		APIKey *string `json:"api_key,omitempty"`
		Domain string  `json:"domain"`
		TLS    bool    `json:"tls"`
	} `json:"http"`

	Log struct {
		// If enabled the application logs to stderr.
		Enabled bool   `json:"enabled"`
		Level   string `json:"level"`
	} `json:"log"`
}

func DefaultConfig() Config {
	var c Config
	c.DB.DSN = ""
	c.HTTP.Addr = "0.0.0.0:8058"
	c.HTTP.Domain = "localhost"
	return c
}

func LoadConfig(filename string) (Config, error) {
	config := DefaultConfig()
	if filename == "" {
		return config, errors.New("must specify a config file path using either TODO_CONFIG environment variable or the --config flag")
	} else if b, err := ioutil.ReadFile(filename); err != nil {
		return config, err
	} else if err = json.Unmarshal(b, &config); err != nil {
		return config, err
	}

	if config.HTTP.APIKey == nil {
		apiKey := crypto.RandomString()
		config.HTTP.APIKey = &apiKey
	}
	return config, nil
}
