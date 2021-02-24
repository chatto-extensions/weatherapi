package ext

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jaimeteb/chatto/extension"
	"github.com/jaimeteb/chatto/fsm"
	"github.com/jaimeteb/chatto/query"
)

// RegisteredFuncs is a list of command to function mappings
var RegisteredFuncs = extension.RegisteredFuncs{
	"weather": Weather,
}

var weatherKey = os.Getenv("WEATHER_API_KEY")
var weatherURL = "http://api.weatherapi.com/v1/current.json?key=%s&q=%s"

type weatherResponse struct {
	Location weatherResponseLocation `json:"location"`
	Current  weatherResponseCurrent  `json:"current"`
}

type weatherResponseLocation struct {
	Name    string `json:"name"`
	Region  string `json:"region"`
	Country string `json:"country"`
}

type weatherResponseCurrent struct {
	Condition  weatherResponseCurrentCondition `json:"condition"`
	TempC      float32                         `json:"temp_c"`
	TempF      float32                         `json:"temp_f"`
	FeelsLikeC float32                         `json:"feelslike_c"`
	FeelsLikeF float32                         `json:"feelslike_f"`
	Humidity   int                             `json:"humidity"`
}

type weatherResponseCurrentCondition struct {
	Text string `json:"text"`
}

// Weather retrieves weather information from api.weatherapi.com and returns the response
func Weather(req *extension.Request) (res *extension.Response) {
	res, err := weather(req)
	if err != nil {
		log.Println("ERROR: %s", err)
	}

	return res
}

func weather(req *extension.Request) (res *extension.Response, err error) {
	location := url.QueryEscape(req.Question.Text)

	resp, err := http.Get(fmt.Sprintf(weatherURL, weatherKey, location))
	if err != nil {
		return &extension.Response{
			FSM:     req.FSM,
			Answers: []query.Answer{{Text: req.Domain.DefaultMessages.Error}},
		}, err
	}

	defer resp.Body.Close()
	var weatherResp weatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return &extension.Response{
			FSM:     req.FSM,
			Answers: []query.Answer{{Text: req.Domain.DefaultMessages.Error}},
		}, err
	}

	var message string
	switch resp.StatusCode {
	case 200:
		weatherMessage := "In %s, %s, it is %s. The temperature is %2.1f °C (%2.1f °F) and feels like %2.1f °C (%2.1f °F)."
		message = fmt.Sprintf(
			weatherMessage,
			weatherResp.Location.Name,
			weatherResp.Location.Region,
			strings.ToLower(weatherResp.Current.Condition.Text),
			weatherResp.Current.TempC,
			weatherResp.Current.TempF,
			weatherResp.Current.FeelsLikeC,
			weatherResp.Current.FeelsLikeC,
		)
	case 400:
		message = "Sorry, I couldn't find your location, try with another one please."
		return &extension.Response{
			FSM: &fsm.FSM{
				State: req.Domain.StateTable["ask_location"],
				Slots: req.FSM.Slots,
			},
			Answers: []query.Answer{{Text: message}}}, nil
	default:
		return &extension.Response{
			FSM:     req.FSM,
			Answers: []query.Answer{{Text: req.Domain.DefaultMessages.Error}},
		}, errors.New(resp.Status)
	}

	return &extension.Response{
		FSM:     req.FSM,
		Answers: []query.Answer{{Text: message}},
	}, nil
}
