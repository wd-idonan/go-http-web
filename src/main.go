package main


import (
	"flag"
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"log"
	"log/syslog"
	"net/http"
	"net"
	"os/signal"
	"syscall"
	"time"
	"database/sql"
)

var (
	cfg 				AppCfg
	buildGitRevision	string
	buildTimeStamp		string
	vpc_db				*sql.DB
	app_db				*sql.DB
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

func initLog() {
	logWriter, err := syslog.New(syslog.LOG_ERR|syslog.LOG_USER, "wd_test")
	log.SetFlags(log.Lshortfile)
	if err != nil {
		fmt.Println("unable to dail syslog server, discarding all log")
		log.SetOutput(ioutil.Discard)
	} else {
		log.SetOutput(logWriter)
	}
}


func initMysql() {
	var err error
	vpc_db, err = initDB(cfg.VpcDB)
	if err != nil {
		log.Fatalf("init vpc_db failed, err: %s", err)
	}
	
	app_db, err = initDB(cfg.AppDB)
	if err != nil {
		log.Fatalf("init app_db failed, err: %s", err)
	}
	
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

	ln, err := net.Listen("tcp", ":"+cfg.ListenPort)
	if err != nil {
		log.Fatalf("listen port: %s failed, err: %s", cfg.ListenPort, err.Error())
	}

	signal.Ignore(syscall.SIGPIPE)

	initLog()
	initMysql()


	srv := &http.Server{
		Handler:		http.DefaultServeMux,
		ReadTimeout:	30 * time.Second,
		WriteTimeout:	30 * time.Second,
	}

	registerAPI()
	log.Println("[SMS] server is restarting.")
	err = srv.Serve(ln)
	if err != nil {
		log.Fatalf("server http server failed, err: %s", err.Error())
	}
}