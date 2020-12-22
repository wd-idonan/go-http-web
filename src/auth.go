package main


import (
	"net/http"
	"strings"
	"log"
	"time"
	"strconv"
)

func APIAuth(r *http.Request, body string) bool {
	var qtimeInt int64

	if strings.HasPrefix(r.RemoteAddr, "127.0.0.1") {
		return true
	}

	identity := r.FormValue("identity")
	if "" == identity {
		log.Printf(" [CHECK] param error, no identity found")
		return false
	}
	salt, bid, err := getAuthInfo(identity)
	if err != nil {
		log.Printf("[CHECK] param error, no salt for identity: '%s', err: '%s'", identity, err.Error())
		return false
	}
	if salt == "" {
		log.Printf("[CHECK] param error, no default salt, identity: %s", identity)
		return false
	}
	reqBid := r.FormValue("bid")
	if !strings.EqualFold(bid, reqBid) {
		log.Printf("[CHECK] param error, bad bid: %s for identity: %s", reqBid, identity)
		return false
	}

	qtimeStr := r.FormValue("qtime")
	if qtimeStr != "" {
		qtimeInt, err = strconv.ParseInt(qtimeStr, 10, 64)
		if err != nil {
			log.Printf("[CHECK] APIAuth param error, bad qtime: %s", qtimeStr)
			return false
		}
	}else {
		log.Printf("[CHECK] APIAuth param error, qtime not exist")
		return false
	}

	nowInt := time.Now().Unix()
	diff := nowInt - qtimeInt
	if diff > cfg.AuthMaxInterval {
		log.Printf("[CHECK] APIAuth qtime diff: %v, max: %d, now:%d, qtime:%d", diff, cfg.AuthMaxInterval, nowInt, qtimeInt)
	}

	return true
}

func getAuthInfo(indetity string) (salt string, bid string, err error) {
	//you can do some auth in here, egg select db
	return salt, bid, nil
}