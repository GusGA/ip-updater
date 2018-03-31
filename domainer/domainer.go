package domainer

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/digitalocean/godo"
	"github.com/gusga/ip-updater/storage"
	"golang.org/x/oauth2"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type OwnDomain struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
	Ip   string `json:"ip"`
}

func NewOwnDomain(domain string, record godo.DomainRecord) *OwnDomain {
	return &OwnDomain{
		Name: domain,
		Id:   record.ID,
		Ip:   record.Data,
	}
}

func OwnDomainFromDB(jsonString string) (*OwnDomain, error) {
	ownDomain := OwnDomain{}
	err := json.Unmarshal([]byte(jsonString), &ownDomain)
	return &ownDomain, err
}

func (od *OwnDomain) ToJSON() string {
	data, err := json.Marshal(od)
	if err != nil {
		return ""
	}
	return string(data)
}

func (od *OwnDomain) UpdateIp(ip string) error {
	editRequest := &godo.DomainRecordEditRequest{
		Data: ip,
	}
	log.Printf("Updating ip address with ip %s", ip)
	domainRecord, _, err := client.Domains.EditRecord(ctx, od.Name, od.Id, editRequest)
	if err != nil {
		return nil
	}
	if domainRecord.Data != ip {
		od.Ip = ip
	}
	return storage.SaveDomainData(od.Name, od.ToJSON())
}

func (od *OwnDomain) Save() error {
	return storage.SaveDomainData(od.Name, od.ToJSON())
}

var DO_TOKEN = os.Getenv("DO_TOKEN")
var workingDomain string

var ctx = context.TODO()
var oauthClient = oauth2.NewClient(oauth2.NoContext, &tokenSource{AccessToken: DO_TOKEN})
var client = godo.NewClient(oauthClient)
var opt = &godo.ListOptions{
	Page:    1,
	PerPage: 100,
}

func SetDomain(domain string) {
	workingDomain = domain
}

func GetDomains() (*OwnDomain, error) {
	var nDomain *OwnDomain
	records, _, err := client.Domains.Records(ctx, workingDomain, opt)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record.Type == "A" && record.Name != "" {
			nDomain = NewOwnDomain(workingDomain, record)
			log.Printf("Fetching domain %s", nDomain.Name)
		}
	}
	return nDomain, nil
}
