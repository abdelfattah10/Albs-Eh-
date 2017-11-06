package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
  "strings"
	"time"
	"fmt"
	"log"
	"os"

	"github.com/ramin0/chatbot"
)

import _ "github.com/joho/godotenv/autoload"

//------------------------ API KEY ------------------------------
const APIkey = "8630f0dade487bf96b54b18e5177f3c6"
const GoogleAPIkey = "AIzaSyDi1xOnlrDDKwTMFeO2eh3lgQyCUdeY2RY"
//------------------------ Structs ------------------------------
type weatherStruc struct{
    Main struct { Temp float64 `json:"temp"` }
}

// struct for JSON for the country validation
type coutryerror struct {
	Status float64 `json:"status"`
}

//structs for JSON for the city and country validation
type matchedSubstrings struct {
	Length float64 `json:"length"`
	Offset float64 `json:"offset"`
}
type structuredFormatting struct {
	Main_text                    string              `json:"main_text"`
	Main_text_matched_substrings []matchedSubstrings `json:"main_text_matched_substrings"`
	Secondary_text               string              `json:"secondary_text"`
}
type cityCountry struct {
	Offset float64 `json:"offset"`
	Value  string  `json:"value"`
}

type forEachCity struct {
	Description           string               `json:"description"`
	Id                    string               `json:"id"`
	Matched_substrings    []matchedSubstrings  `json:"matched_substrings"`
	Place_id              string               `json:"place_id"`
	Reference             string               `json:"reference"`
	Structured_formatting structuredFormatting `json:"structured_formatting"`
	Terms                 []cityCountry        `json:"terms"`
	Types                 []string             `json:"types"`
}

type cityerror struct {
	Predictions []forEachCity
}

//----------------------- Functions -----------------------------
func FloatToString(input_num float64) string {

    return strconv.FormatFloat(input_num, 'f', 2, 64)
}

func weather(country string, city string) float64 {
	weatherObj := weatherStruc{}
	url := "http://api.openweathermap.org/data/2.5/weather?q="+ country +","+ city +"&APPID=" + APIkey
	spaceClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
			log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Weather-API")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
			log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
			log.Fatal(readErr)
	}

	jsonErr := json.Unmarshal(body, &weatherObj)
	if jsonErr != nil {
			log.Fatal(jsonErr)
	}

 Celsius := weatherObj.Main.Temp - 273.15
 fmt.Println(Celsius)
 return Celsius
}

func countryCheck(country string) bool {
	countryname := country

		countryurl := "https://restcountries.eu/rest/v2/name/" + countryname + "?fullText=true"

		coutryerrorObj := coutryerror{}

		spaceClient := http.Client{
			Timeout: time.Second * 10, // Maximum of 10 secs
		}

		req, err := http.NewRequest(http.MethodGet, countryurl, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("User-Agent", "Country-API")

		res, getErr := spaceClient.Do(req)
		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}
		jsonErr := json.Unmarshal(body, &coutryerrorObj)
		if jsonErr != nil {
			fmt.Println("the country is valid")
			return true
		}

		if coutryerrorObj.Status == 404 {
			fmt.Println("the country is not valid")
			return false
		}
		return false
}

func cityCheck(city string, country string) bool {
	cityname := city

	var finalcityname string
 	words := strings.Fields(cityname)
 	for i := 0; i < len(words); i++ {
 		finalcityname = finalcityname + words[i]
 	}
 	cityurl := "https://maps.googleapis.com/maps/api/place/autocomplete/json?input=" + finalcityname + "&types=(cities)&key=" + GoogleAPIkey

	cityerrorObj := cityerror{}

	spaceClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
	}

	req, err := http.NewRequest(http.MethodGet, cityurl, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Country-API")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	jsonErr := json.Unmarshal(body, &cityerrorObj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	if strings.EqualFold(country, "united states of america") {
		country = "united states"
	}
	citIsValid := false
	for i := 0; i < len(cityerrorObj.Predictions); i++ {
		if strings.EqualFold(cityerrorObj.Predictions[i].Terms[len(cityerrorObj.Predictions[i].Terms)-1].Value, country) {
			citIsValid = true
		}
	}
	return citIsValid
	}


 func chatbotProcess(session chatbot.Session, message string) (string, error) {
 	if strings.EqualFold(message, "chatbot") {
 		return "", fmt.Errorf("This can't be, I'm the one and only %s!", message)
 	}
fmt.Println(message)

 	var questionMarksCount int
	resMsg := ""


 	// Try fetching the count of question marks
 	count, found := session["questionMarksCount"]
 	// If a count is saved in the session
 	if found {
 		// Cast it into an int (since sessions values are generic)
 		questionMarksCount = count.(int)
 	} else {
 		// Otherwise, initialize the count to 1
 		questionMarksCount = 1
 	}

	if questionMarksCount == 1 {
		session["Country"] = message
		country := session["Country"].(string)

		countryCheck := countryCheck(country)

		if countryCheck {
		resMsg = " ,Great! Now What Is Your City ?"
			} else {
				resMsg = "Invalid Country, Pls try again"
				message = ""
				questionMarksCount = 0
			 }
	} else {
		session["City"] = message

		country := session["Country"].(string)
		city := session["City"].(string)

		CelsiusF := weather(country,city)
		Celsius  := FloatToString(CelsiusF)

		cityCheck := cityCheck(city,country)

		if CelsiusF == -273.15 || !cityCheck {
		resMsg = "City Not Found, Pls Enter Your Country Again"
		message = ""
		questionMarksCount = 0
		} else {

		clothes := ""
		switch(true){

		case CelsiusF >= 40 : clothes = ", It is melting outside you better stay at home otherwise put on the swimsuit"
		case CelsiusF >= 31 && CelsiusF  < 40: clothes = ", It is very sunny and hot today we recommend you to wear sleeveless tops and shorts also try avoid wearing clothes made of polyester, nylon, or silk"
		case CelsiusF >= 24 && CelsiusF  < 31: clothes = ", It seems to be hot today we recommend you to wear light clothes like T-shirts and sweatpants would be nice also do not forget your hat and sunglasses"
		case CelsiusF >= 16 && CelsiusF  < 24: clothes = ", It seems to be clear comfortable today we recommend you to wear light clothes like T-shirts or shirts and jeans"
		case CelsiusF >= 9  && CelsiusF  < 16: clothes = ", It is seems to be cold today we recommend you to wear heavy clothes like pull over or hoodies and jeans"
		case CelsiusF >= -5 && CelsiusF  < 9: clothes = ", It is very cold outside we recommend you to wear heavy clothes like jackets and denim jeans also do not forget your scarfs and icecap"
		case CelsiusF <  -5 : clothes = ", It is freezing outside you better stay at home"

		}

    endMsg := "<br> <br> If You Want To Check For Another Country Pls Inform Me, Otherwise GoodBye :))"
		message = ""
		resMsg = "Weather in <span style='color:red;'>"+city+"</span> is <span style='color:CornflowerBlue;'>"+Celsius+" C </span>"+clothes+endMsg
		questionMarksCount = 0

		}
	}

 	// Save the updated question marks count to the session
 	session["questionMarksCount"] = questionMarksCount + 1

fmt.Println(message)
fmt.Println(resMsg)

 	return fmt.Sprintf(" <b style='color:red'>%s</b> %s", message,resMsg), nil
 }

func main() {

	 chatbot.WelcomeMessage = "Hello, Pls Enter Your Country ..."
	 chatbot.ProcessFunc(chatbotProcess)

	// Use the PORT environment variable
	port := os.Getenv("PORT")
	// Default to 3000 if no PORT environment variable was defined
	if port == "" {
		port = "3000"
	}

	// Start the server
	fmt.Printf("Listening on port %s...\n", port)
	log.Fatalln(chatbot.Engage(":" + port))
}
