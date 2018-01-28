package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gusga/ip-updater/domainer"
	"github.com/gusga/ip-updater/storage"
)

const IP_URL_SERVICE = "http://ident.me/"

var DOMAIN string

func lookupRemoteIp() string {
	response, err := http.Get(IP_URL_SERVICE)
	if err != nil {
		log.Fatalf("Error ", err)
	}
	ip, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Fatalf("Error parsing ip ", err)
	}
	return string(ip)
}
func init() {
	DOMAIN = os.Getenv("DOMAIN")
	if DOMAIN == "" {
		log.Fatalln("DOMAIN env must be setted")
	}
	domainer.SetDomain(DOMAIN)
}

func main() {
	var domainData *domainer.OwnDomain
	savedDomain, err := storage.GetDomainData(DOMAIN)
	if err != nil {
		log.Fatalln(err)
	}
	if savedDomain == "" {
		domainData, err = domainer.GetDomains()
		if err != nil {
			log.Fatalf("Could not retrieve domain records ", err)
		}
		err = domainData.Save()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		domainData, err = domainer.OwnDomainFromDB(savedDomain)
		if err != nil {
			log.Fatalln(err)
		}
	}

	ip := lookupRemoteIp()
	if domainData.Ip == ip {
		log.Println("The IP address has not changed")
		return
	}
	err = domainData.UpdateIp(ip)
	if err != nil {
		log.Fatalln(err)
	}
	return
}
