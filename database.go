package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

var DB_PATHS = []string{}

type Config struct {
	Moziot *Moziot `json:"moziot"`
}

type Moziot struct {
	Config map[string]string `json:"config"`
}

func init() {
	moziotHome := os.Getenv("MOZIOT_HOME")
	if len(moziotHome) > 0 {
		DB_PATHS = append([]string{moziotHome}, DB_PATHS...)
	}

	moziotDatabase := os.Getenv("MOZIOT_DATABASE")
	if len(moziotDatabase) > 0 {
		DB_PATHS = append([]string{moziotDatabase}, DB_PATHS...)
	}
}

type Database struct {
	PackageName string
	Path        string
	Conn        *sql.DB
}

func NewDatabase(packageName string, path string) (*Database, error) {
	d := new(Database)
	d.PackageName = packageName
	d.Path = path
	d.Conn = nil

	if len(d.Path) == 0 {
		for _, p := range DB_PATHS {
			_, err := os.Stat(p)

			if err == nil {
				d.Path = p
				break
			}
		}
	}

	return d, nil
}

func (d *Database) Open() error {
	if len(d.Path) == 0 {
		return errors.New("Database path unknown")
	}

	connOptions := "?_busy_timeout=10000"
	db, err := sql.Open("sqlite3", strings.Join([]string{d.Path}, connOptions))
	defer db.Close()

	if err != nil {
		return err
	}

	d.Conn = db

	return nil
}

func (d *Database) Close() error {
	d.Conn.Close()

	return nil
}

func (d *Database) LoadConfig() (string, error) {
	key := fmt.Sprintf("addons.%s", d.PackageName)
	selectStmt, err := d.Conn.Prepare("SELECT value FROM settings WHERE key = ?")
	if err != nil {
		return "", err
	}

	selectStmt.Close()
	rows, err := selectStmt.Query(key)
	if err != nil {
		return "", err
	}

	var value string
	for rows.Next() {
		err = rows.Scan(&value)
		return value, nil
	}

	return "", nil
}

func (d *Database) SaveConfig(config map[string]string) error {
	key := fmt.Sprintf("addons.%s", d.PackageName)
	selectStmt, err := d.Conn.Prepare("SELECT value FROM settings WHERE key = ?")
	if err != nil {
		return err
	}

	selectStmt.Close()
	rows, err := selectStmt.Query(key)
	if err != nil {
		return err
	}

	var value string
	for rows.Next() {
		err = rows.Scan(&value)
	}

	if len(value) > 0 {
		var data *Config
		err := json.Unmarshal([]byte(value), &config)
		if err != nil {
			return err
		}

		data.Moziot.Config = config

		insertStmt, err := d.Conn.Prepare("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)")
		if err != nil {
			return err
		}

		insertStmt.Close()

		configValue, err := json.Marshal(data)
		if err != nil {
			return err
		}

		selectStmt.Exec(key, configValue)
	}

	return nil
}
