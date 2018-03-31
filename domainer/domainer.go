package domainer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/gusga/ip-updater/storage"
	"golang.org/x/oauth2"
)

var doToken = os.Getenv("DO_TOKEN")
var workingDomain string

var ctx = context.TODO()
var oauthClient = oauth2.NewClient(oauth2.NoContext, &tokenSource{AccessToken: doToken})
var client = godo.NewClient(oauthClient)
var opt = &godo.ListOptions{
	Page:    1,
	PerPage: 100,
}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type ownDomain struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	IP   string `json:"ip"`
}

// DomainList Strut
type DomainList struct {
	Name       string       `json:"name"`
	IP         string       `json:"ip"`
	SubDomains []*ownDomain `json:"subdomains"`
}

func newOwnDomain(domain string, record godo.DomainRecord) *ownDomain {
	return &ownDomain{
		Name: record.Name,
		ID:   record.ID,
		IP:   record.Data,
	}
}

// DomainListFromDB fetch domain list from db
func DomainListFromDB(jsonString string) (*DomainList, error) {
	ownDomainList := DomainList{}
	err := json.Unmarshal([]byte(jsonString), &ownDomainList)
	return &ownDomainList, err
}

// ToJSON convert DomainList to string to save on db
func (dl *DomainList) ToJSON() string {
	data, err := json.Marshal(dl)
	if err != nil {
		return ""
	}
	return string(data)
}

// SetIP set the domain list ip
func (dl *DomainList) SetIP(ip string) {
	dl.IP = ip
}
func (od *ownDomain) updateIP(ip string) error {
	editRequest := &godo.DomainRecordEditRequest{
		Data: ip,
	}
	log.Printf("Updating ip address with ip %s", ip)
	domainRecord, _, err := client.Domains.EditRecord(ctx, workingDomain, od.ID, editRequest)
	if err != nil {
		return err
	}
	if domainRecord.Data != ip {
		od.IP = ip
	}
	return nil
}

// Save DomainList to DB
func (dl *DomainList) Save() error {
	return storage.SaveDomainData(dl.Name, dl.ToJSON())
}

// UpdateDomainsIP Update all subdomains ip
func (dl *DomainList) UpdateDomainsIP(ip string) error {
	for _, od := range dl.SubDomains {
		log.Printf("Updating domain %s", od.Name)
		err := od.updateIP(ip)
		if err != nil {
			return err
		}
	}
	return dl.Save()
}

// CheckIP Check all subs ip
func (dl *DomainList) CheckIP(ip string) bool {
	for _, od := range dl.SubDomains {
		log.Printf("Checkin ip from domain %s", od.Name)
		if strings.Compare(od.IP, ip) != 0 {
			return false
		}
	}
	return true
}

// SetDomain set global variable
func SetDomain(domain string) {
	workingDomain = domain
}

// GetDomains from DO
func GetDomains() (*DomainList, error) {
	nDomains := &DomainList{Name: workingDomain}
	records, _, err := client.Domains.Records(ctx, workingDomain, opt)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Type == "A" && record.Name != "" {
			domainPlusSub := record.Name + "_" + workingDomain
			log.Printf("Fetching domain %s", domainPlusSub)
			subdomain := newOwnDomain(domainPlusSub, record)
			nDomains.SubDomains = append(nDomains.SubDomains, subdomain)
		}
	}
	return nDomains, nil
}
