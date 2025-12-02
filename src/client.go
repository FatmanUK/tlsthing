package main

import (
	"fmt"
	"log"
	"net"
	"time"
	"errors"
)

const DocoptString = DocoptAppVer + `

Usage:
  {{ .Arg0 }} run [--verbosity=DBG]
  {{ .Arg0 }} add (pem | der | jks) --name=NAME
        --key=KEY --cert=CERT [--verbosity=DBG]
  {{ .Arg0 }} upd (pem | der | jks) --name=NAME
        --key=KEY --cert=CERT [--verbosity=DBG]
  {{ .Arg0 }} del --name=NAME [--verbosity=DBG]
  {{ .Arg0 }} --help
  {{ .Arg0 }} --version

Options:
  --help -h        Show this screen
  --version -v     Show version
  --name=NAME      Unique name for key pair
  --key=KEY        PEM key file
  --cert=CERT      PEM cert file
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
		case args["add"].(bool): err = actAdd(args, logs)
		case args["del"].(bool): err = actDel(args, logs)
		case args["upd"].(bool): err = actUpd(args, logs)
	}
	if err != nil {
		logs <- Error(err.Error())
	}
}

// Host: get from director
// NotBefore: get from cert update
// NotAfter: get from cert update
func MakeRegistration(
		args map[string]interface{},
		logs chan Log) (RegistryDatum, error) {
	keyPath, err := ResolvePath(args["--key"].(string))
	if err != nil {
		return RegistryDatum{}, err
	}
	certPath, err := ResolvePath(args["--cert"].(string))
	if err != nil {
		return RegistryDatum{}, err
	}
	return RegistryDatum{
		Name: args["--name"].(string),
		Host: "",
		Type: getPairType(args),
		CertPath: certPath,
		KeyPath: keyPath,
		RegDate: time.Now(),
		NotBefore: time.Time{},
		NotAfter: time.Time{},
	}, nil
}

func sendData(r RegistryDatum, t string, logs chan Log) error {
	s, err := EncodeData(r)
	if err != nil {
		return err
	}
	logs <- Debug(fmt.Sprintf("Encoded: %s", s))
	var c []Message
	c = append(c, Message{Name: t, Data: s})
	buffer, err := EncodeNetworkMessages(c, logs)
	if err != nil {
		return err
	}
	DebugOutputByteArray(buffer, logs)
	port := getEnvDefault("TLSTHING_PORT", APP_PORT)
	host := getEnvDefault("TLSTHING_HOST", "localhost")
	hostPort := fmt.Sprintf("%s:%s", host, port)
	logs <- Debug("Connecting director on " + hostPort)
	conn, err := GetDirectorConnection(hostPort, logs)
	if err != nil {
		return err
	}
	defer conn.Close()
	err = WriteTcpBytes(conn, logs, buffer)
	if err != nil {
		return err
	}
	return nil
}

func actRun(args map[string]interface{}, logs chan Log) error {
	port := getEnvDefault("TLSTHING_PORT", APP_PORT)
	host := getEnvDefault("TLSTHING_HOST", "localhost")
	hostPort := fmt.Sprintf("%s:%s", host, port)
	logs <- Debug("Connecting director on " + hostPort)
	conn, err := GetDirectorConnection(hostPort, logs)
	if err != nil {
		return err
	}
	defer conn.Close()
	ThreadAcceptHandler(conn, logs, nil)
	return nil
}

func actAdd(args map[string]interface{}, logs chan Log) error {
	r, err := MakeRegistration(args, logs)
	if err != nil {
		return err
	}
	return sendData(r, "create", logs)
}

func actUpd(args map[string]interface{}, logs chan Log) error {
	r, err := MakeRegistration(args, logs)
	if err != nil {
		return err
	}
	return sendData(r, "update", logs)
}

func actDel(args map[string]interface{}, logs chan Log) error {
	r := RegistryDatum{
		Name: args["--name"].(string),
	}
	return sendData(r, "delete", logs)
}

func SwitchMessage(
		conn net.Conn,
		msg Message,
		data interface{},
		logs chan Log) error {
	var err error
	switch msg.Name {
//		case "create": err = msgCreate(conn, msg, data, logs)
//		case "read":   err = msgRead(conn, msg, data, logs)
//		case "update": err = msgUpdate(conn, msg, data, logs)
//		case "delete": err = msgDelete(conn, msg, data, logs)
		default: {
			f := "Unknown command: %s"
			err = errors.New(fmt.Sprintf(f, msg.Name))
		}
	}
	return err
}

//////////////////////////////////////////////////////////////////////
/*
func msgCreate(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	return nil
}

func msgRead(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	return nil
}

func msgUpdate(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	return nil
}

func msgDelete(
		conn net.Conn,
		msg Message,
		db interface{},
		logs chan Log) error {
	return nil
}
*/

/*
func getCreateMessage(logs chan Log) ([]byte, error) {
	t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		return nil, err
	}
	r := RegistryDatum{
		Name: "TestCreate",
		Host: "192.168.11.145",
		Type: "CreateType",
		RegDate: t,
		NotBefore: t,
		NotAfter: t,
	}
	s, err := encodeData(r)
	logs <- Debug(fmt.Sprintf("Encoded: %s", s))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getUpdateMessage(logs chan Log) ([]byte, error) {
	t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		return nil, err
	}
	r := RegistryDatum{
		Name: "TestCreate",
		Host: "192.168.11.145",
		Type: "UpdateType",
		RegDate: t,
		NotBefore: t,
		NotAfter: t,
	}
	s, err := encodeData(r)
	logs <- Debug(fmt.Sprintf("Encoded: %s", s))
	if err != nil {
		return nil, err
	}
	return s, nil
}

func getDeleteMessage(logs chan Log) ([]byte, error) {
	t, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	if err != nil {
		return nil, err
	}
	r := RegistryDatum{
		Name: "TestCreate",
		Host: "192.168.11.145",
		Type: "DeleteType",
		RegDate: t,
		NotBefore: t,
		NotAfter: t,
	}
	s, err := encodeData(r)
	logs <- Debug(fmt.Sprintf("Encoded: %s", s))
	if err != nil {
		return nil, err
	}
	return s, nil
}
*/

/*
func setupStuff(logs chan Log, dbConn string, configFile string) {
	// record configs in conf file
	// confs overridden by env vars
	// set up database
	logs <- Debug("dbConn: " + dbConn)
	logs <- Debug("configFile: " + configFile)
	logs <- Warn("This will erase any previous configuration.")
	// ...
}
*/
