/*
Copyright © 2021 The Goat <we@thegoat.ir>
Copyright © 2021 Mahdyar Hasanpour <hi@mahdyar.me>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Json bool

// whoisCmd represents the whois command
var whoisCmd = &cobra.Command{
	Use:   "whois [domain]",
	Short: "Whois lookup a domain.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		lookup(args[0])
	},
}

func init() {
	rootCmd.AddCommand(whoisCmd)

	// --json, -j flag for printing as json
	whoisCmd.Flags().BoolVarP(&Json, "json", "j", false, "print as json")
}

func lookup(domain string) {
	responseBytes := getLookupData(viper.GetString("API_URL"), domain)
	if Json {
		fmt.Print(string(responseBytes))
	} else {
		var lookupData map[string]interface{}
		if err := json.Unmarshal(responseBytes, &lookupData); err != nil {
			log.Printf("Couldn't unmarshal response: %s", err)
		}
		for index, datum := range lookupData {
			fmt.Printf("%v: %v\n", index, datum)
		}
	}
}

func getLookupData(baseAPI string, domain string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		baseAPI+"?token="+viper.GetString("token")+"&domain="+domain,
		nil,
	)
	if err != nil {
		log.Printf("Couldn't request the whois lookup data: %s", err)
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("User-Agent", "thegoat-cli")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Couldn't make a request: %s", err)
	}

	responseBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Couldn't read response body: %s", err)
	}
	return responseBytes
}
