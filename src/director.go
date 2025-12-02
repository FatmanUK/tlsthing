package main

import (
	"fmt"
	"log"
	"net"
	"errors"
)

const DocoptString = DocoptAppVer + `

Usage:
  {{ .Arg0 }} run [--verbosity=DBG]
  {{ .Arg0 }} --help
  {{ .Arg0 }} --version

Options:
  --help -h        Show this screen
  --version -v     Show version
  --verbosity=DBG  Log level 0-3; [default: {{ .Debug }}]`

func main() {
	bd := getBaseDir(getArg0())
	dotv := DocOptInfo{APP_NAME, VERSION, bd, DEBUG_LEVEL}
	err := SharedMainFunc(DocoptAppVer, DocoptString, dotv)
	if err != nil {
		log.Println(err)
	}
}

// If help or version passed, args will be empty.
func ThreadMain(args map[string]interface{}, logs chan Log) {
	defer close(logs)
	var err error
	logs <- Info("Starting up." )
	switch true {
		case args["run"].(bool): err = actRun(args, logs)
	}
	if err != nil {
		logs <- Error(err.Error())
	}
}

func makeDsnData(u string, p string) DsnData {
	dbPort := getEnvDefault("DATABASE_PORT", "5432")
	dbHost := getEnvDefault("DATABASE_HOST", "localhost")
	dbName := getEnvDefault("DATABASE_NAME", "tlsthing")
	dbTlsMode := getEnvDefault("DATABASE_TLSMODE", "disable")
	return DsnData{dbHost, u, p, dbName, dbPort, dbTlsMode}
}

func actRun(args map[string]interface{}, logs chan Log) error {
	port := getEnvDefault("TLSTHING_PORT", APP_PORT)
	tl, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	defer tl.Close()
	u, p, err := GetDynamicCreds()
	if err != nil {
		return err
	}
	dsn, err := MakePostgresDsn(makeDsnData(u, p))
	if err != nil {
		return err
	}
	db, err := GetDatabaseConnection(dsn, logs)
	if err != nil {
		return err
	}
	for {
		conn, err := tl.Accept()
		if err != nil {
			return err
		}
		go ThreadAcceptHandler(conn, logs, db)
	}
	return nil
}

func SwitchMessage(
		conn net.Conn,
		msg Message,
		data interface{},
		logs chan Log) error {
	var err error
	switch msg.Name {
		case "create": err = msgCreate(conn, msg, data, logs)
		case "read":   err =   msgRead(conn, msg, data, logs)
		case "update": err = msgUpdate(conn, msg, data, logs)
		case "delete": err = msgDelete(conn, msg, data, logs)
		default: {
			f := "Unknown command: %s"
			err = errors.New(fmt.Sprintf(f, msg.Name))
		}
	}
	return err
}

func msgCreate(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	var rde RegistryDatum
	err := DecodeData(msg.Data, &rde)
	if err != nil {
		return err
	}
	rde.Host = PopPort(conn.RemoteAddr().String())
	err = rde.Create(db)
	return nil
}

func msgRead(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
/*
	// decode
	// get missed messages while offline - match on Host
	// manual check for expired certs
	// query db, requests outstanding for client?
*/
	return nil
}

func msgUpdate(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	var rde RegistryDatum
	err := DecodeData(msg.Data, &rde)
	if err != nil {
		return err
	}
	rde.Host = PopPort(conn.RemoteAddr().String())
	err = rde.Update(db)
	return nil
}

func msgDelete(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	var rde RegistryDatum
	err := DecodeData(msg.Data, &rde)
	if err != nil {
		return err
	}
	rde.Host = PopPort(conn.RemoteAddr().String())
	err = rde.Delete(db)
	return err
}
