package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
)

var (
	ES *elasticsearch.Client
)

type csp_raport struct {
	date               string
	document_uri       string
	referrer           string
	blocked_uri        string
	violated_directive string
	original_policy    string
}

func ConnectES() {
	var err error
	ES, err = elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
}

func ConnectStatus() bool {
	return true
}

func Report() {

}
