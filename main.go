package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/bdenning/go-pushover"
	"github.com/spf13/pflag"
)

type (
	counties struct {
		XMLName xml.Name    `xml:"Counties"`
		County  []xmlCounty `xml:"County"`
	}

	xmlCounty struct {
		Name      string `xml:"CountyName"`
		Condition string `xml:"ConditionEnglish"`
	}

	options struct {
		userKey  string
		appKey   string
		counties []string
	}
)

func main() {
	var options options
	pflag.StringVarP(&options.userKey, "userkey", "u", "", "Pushover Userkey. REQUIRED.")
	pflag.StringVarP(&options.appKey, "appkey", "a", "", "Pushover Appkey. REQUIRED.")
	pflag.StringSliceVarP(&options.counties, "counties", "c", []string{}, "Counties to check. REQUIRED.")
	pflag.Parse()

	url := "http://novascotia.ca/natr/forestprotection/wildfire/burnsafe/counties.xml"
	web := "http://novascotia.ca/natr/forestprotection/wildfire/burnsafe/"
	client := http.Client{}
	res, err := client.Get(url)

	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		fmt.Println("Couldn't fetch XML from: ", url)
		return
	}

	b := new(bytes.Buffer)
	b.ReadFrom(res.Body)
	buffer := b.Bytes()

	var Counties counties

	err = xml.Unmarshal(buffer, &Counties)

	if err != nil {
		fmt.Println("Failed to unmarshal data: ", err.Error())
	}

	if len(Counties.County) == 0 {
		fmt.Println("No XML data yet, try after 14:00.")
		return
	}

	conditions := make(map[string]string)
	conditions["Burn"] = "Burning is permitted between the hours of 2pm and 8am the following day."
	conditions["Restricted"] = "Burning is permitted between the hours of 7pm and 8am the following day."

	for _, county := range Counties.County {
		for _, userCounty := range options.counties {
			if county.Name == userCounty {
				message := pushover.NewMessage(options.appKey, options.userKey)
				report := fmt.Sprintf("County: %s, Status: %s\n%s\n\n%s", county.Name, county.Condition, conditions[county.Condition], web)
				_, err := message.Push(report)

				if err != nil {
					fmt.Printf("Error sending notification: %s\n", err.Error())
				}
			}
		}
	}
}
