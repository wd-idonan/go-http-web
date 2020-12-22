package main


import (
	"flag"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"net/http"
)

var (
	cfg 				AppCfg
	buildGitRevision	string
	buildTimeStamp		string
)


func readConfig(configFile string, cfg *AppCfg) error {
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buf, &cfg)
	if err != nil {
		return err
	}

	return nil
}

func registerAPI() {
	http.Handle("/test/getdata", AdvHandler{h: describeUserInfo})
}

func main()  {
	pConfigFile := flag.String("f", "./test.json", "config file path")
	pVersion := flag.Bool("v", false, "show version")
	flag.Parse()

	if *pVersion {
		fmt.Printf("Git Revision: %s\n Build Time: %s\n", buildGitRevision, buildTimeStamp)
		os.Exit(0)
	}

	err := readConfig(*pConfigFile, &cfg)
	if err != nil {
		log.Fatalf("read config file err: %s", err.Error())
	}
}