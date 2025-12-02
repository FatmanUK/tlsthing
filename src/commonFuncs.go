package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"bytes"
	"strings"
	"strconv"
	"text/template"
	"encoding/json"
	"path/filepath"
	"github.com/FatmanUK/fatgo/doh"
)

const DocoptAppVer = "{{ .Name }} v{{ .Version }}"

type DocOptInfo struct {
	Name string
	Version string
	Arg0 string
	Debug string
}

type Message struct {
	Name string  `json:"Name,omitempty"`
	Data []byte  `json:"Data,omitempty"`
}

type DsnData struct {
	Host string
	User string
	Password string
	Database string
	Port string
	TlsMode string
}

func NaiveTwoToThePowerOf(p int) int {
	var rv int = 1
	for ; p > 0; p-- {
		rv *= 2
	}
	return rv
}

// Implement something better later and call it here. Saves changing
// working caller code later.
func TwoToThePowerOf(p int) int {
	return NaiveTwoToThePowerOf(p)
}

func DoublingSleep(count int, scale float32) {
	offset := float32(TwoToThePowerOf(count)) * scale
	time.Sleep(time.Second * time.Duration(offset))
}

// get from Vault
func GetDynamicCreds() (string, string, error) {
	return "postgres", "temppw", nil
}

func MakePostgresDsn(dsnStruct DsnData) (string, error) {
	dsnFormatArray := []string{
		"host={{ .Host }}",
		"user={{ .User }}",
		"password={{ .Password }}",
		"dbname={{ .Database }}",
		"port={{ .Port }}",
		"sslmode={{ .TlsMode }}",
	}
	return TPrintf(
		strings.Join(dsnFormatArray, " "),
		dsnStruct,
	)
}

func EncodeNetworkMessages(
		msgs []Message,
		logs chan Log) ([]byte, error) {
	buffer, err := json.Marshal(msgs)
	if err != nil {
		logs <- Error(err.Error())
		return nil, err
	}
	return buffer, nil
}

func DecodeNetworkMessages(
		bytes []byte,
		logs chan Log) ([]Message, error) {
	var msgs []Message
	err := json.Unmarshal(bytes, &msgs)
	if err != nil {
		logs <- Error(err.Error())
		return nil, err
	}
	return msgs, nil
}

// I tried to find out why we need template names, but I got confused
// over the dichotomy of templates within a template (???).
// Also, you can't Delete a template once it's New'd. (???)
func TPrintf(
		templateString string,
		object interface{}) (string, error) {
	mt := template.New("TempTemplate") // ???
	mt, err := mt.Parse(templateString)
	if err != nil {
		return "", err
	}
	var wr bytes.Buffer
	err = mt.Execute(&wr, object)
	if err != nil {
		return "", err
	}
	return string(wr.Bytes()), nil
}

// data should be a *struct
func DecodeData(bytes []byte, data any) error {
	err := json.Unmarshal(bytes, data)
	if err != nil {
		return err
	}
	return nil
}

func EncodeData(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func IntFromString(s string) (int, error) {
	threshold, err := strconv.Atoi(s)
	if err != nil {
		return -1, err
	}
	return threshold, nil
}

func SharedMainFunc(appVer string, usage string, dotv any) error {
	args, err := doh.RunDocopt(appVer, usage, dotv)
	if err != nil {
		return err
	}
	if len(args) > 0 {
		verb := args["--verbosity"].(string)
		threshold, err := IntFromString(verb)
		if err != nil {
			return err
		}
		logs := MakeLogChan()
		go ThreadMain(args, logs)
		ThreadLoggingLoop(uint(threshold), logs)
	}
	return nil
}

func ResolvePath(s string) (string, error) {
	return filepath.Abs(s)
}

func getBaseDir(s string) string {
	return filepath.Base(s)
}

func getArg0() string {
	return os.Args[0]
}

func getPairType(args map[string]interface{}) string {
	switch true {
		case args["pem"].(bool): return "pem"
		case args["der"].(bool): return "der"
		case args["jks"].(bool): return "jks"
	}
	return ""
}

func getEnvDefault(name string, def string) string {
	env := os.Getenv(name)
	if env == "" {
		return def
	}
	return env
}

func PopPort(s string) string {
	arr := strings.Split(s, ":")
	arr = arr[:len(arr) - 1]
	s = strings.Join(arr, ":")
	return s
}

func ProcessMessages(
		conn net.Conn,
		msgs []Message,
		logs chan Log,
		data interface{}) error {
	var err error
	for _, msg := range msgs {
		logs <- Debug(fmt.Sprintf("Received: %s", msg.Name))
		err = SwitchMessage(conn, msg, data, logs)
		if err != nil {
			logs <- Warn(err.Error())
		}
	}
	return err
}
