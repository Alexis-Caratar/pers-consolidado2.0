package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/promartingranja/SISNOVA-PERS-CONSOLIDADO/APP/models"
	//Para conexion con SQLSERVER
	_ "github.com/denisenkom/go-mssqldb"
)

// CONSOLIDADO

func getConfiguracion() (config models.ServerDB, err error) {
	js, err := os.Open("../config/bd_cauca.json")
	if err != nil {
		return config, err
	}
	defer js.Close()
	deco := json.NewDecoder(js)
	err = deco.Decode(&config)
	return config, err
}

func GetSQLServer() (db *sql.DB, err error) {
	config, err := getConfiguracion()
	if err != nil {
		return db, err
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", config.Servidor, config.Usuario, config.Clave, config.Puerto, config.BD)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	return db, err
}

// SUBIENDO ARCHIVOS Y MIGRANDO A BD

func getConfiguracionUpload() (config models.ServerDB, err error) {
	js, err := os.Open("../config/bdUpload.json")
	if err != nil {
		return config, err
	}
	defer js.Close()
	deco := json.NewDecoder(js)
	err = deco.Decode(&config)
	return config, err
}

func GetSQLServerUpload() (db *sql.DB, err error) {
	config, err := getConfiguracionUpload()
	if err != nil {
		return db, err
	}

	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", config.Servidor, config.Usuario, config.Clave, config.Puerto, config.BD)

	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	return db, err
}