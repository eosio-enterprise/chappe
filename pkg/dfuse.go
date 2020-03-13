package pkg

import (
	"fmt"
	"log"

	"github.com/dfuse-io/eosws-go"
	"github.com/spf13/viper"
)

// func getToken(apiKey string) (token string, expiration time.Time, err error) {
// 	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"api_key":"%s"}`, apiKey)))
// 	resp, err := http.Post("https://auth.dfuse.io/v1/auth/issue", "application/json", reqBody)
// 	if err != nil {
// 		err = fmt.Errorf("unable to obtain token: %s", err)
// 		return
// 	}

// 	if resp.StatusCode != 200 {
// 		err = fmt.Errorf("unable to obtain token, status not 200, got %d: %s", resp.StatusCode, reqBody.String())
// 		return
// 	}

// 	if body, err := ioutil.ReadAll(resp.Body); err == nil {
// 		token = gjson.GetBytes(body, "token").String()
// 		expiration = time.Unix(gjson.GetBytes(body, "expires_at").Int(), 0)
// 	}
// 	return
// }

// GetClient ...
func GetClient() *eosws.Client {
	apiKey := viper.GetString("Dfuse.ApiKey")
	if apiKey == "" {
		log.Fatalf("Missing Dfuse.ApiKey in config")
	}

	jwt, _, err := eosws.Auth(apiKey)
	if err != nil {
		log.Fatalf("cannot get auth token: %s", err.Error())
	}

	// var dfuseEndpoint = viper.GetString("Dfuse.WSEndpoint")
	var origin = viper.GetString("Dfuse.Origin")

	client, err := eosws.New(viper.GetString("Dfuse.WSEndpoint"), jwt, origin)
	if err != nil {
		log.Fatalf("cannot connect to dfuse endpoint: %s", err.Error())
	}
	return client
}

// GetActionTraces ...
func GetActionTraces() *eosws.GetActionTraces {
	ga := &eosws.GetActionTraces{}
	ga.ReqID = "chappe"
	ga.StartBlock = -300
	ga.Listen = true
	ga.WithProgress = 3
	ga.IrreversibleOnly = false
	ga.Data.Accounts = viper.GetString("Eosio.PublishAccount")
	ga.Data.ActionNames = "pub"
	fmt.Printf("Connecting...  %s::%s\n", ga.Data.Accounts, ga.Data.ActionNames)
	ga.Data.WithInlineTraces = true
	return ga
}
