package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/bdenning/go-pushover"
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
)

func main() {

	url := "http://novascotia.ca/natr/forestprotection/wildfire/burnsafe/counties.xml"
	client := http.Client{}
	res, err := client.Get(url)

	if err != nil {
		fmt.Println("Couldn't fetch XML from: ", url)
	}

	b := new(bytes.Buffer)
	b.ReadFrom(res.Body)
	buffer := b.Bytes()

	var Counties counties

	err = xml.Unmarshal(buffer, &Counties)

	if err != nil {
		fmt.Println("Failed to unmarshal data: ", err.Error())
	}

	for _, county := range Counties.County {
		if county.Name == "Lunenburg" {
			token := ""
			user := ""

			message := pushover.NewMessage(token, user)

			report := fmt.Sprintf("County: %s, Status: %s", county.Name, county.Condition)
			res, err := message.Push(report)
		}
	}

}
