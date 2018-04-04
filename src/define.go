package main

type geoRequest struct {
	key      string
	address  string
	city     string
	batch    string
	sig      string
	output   string
	callback string
}

type geoRespone struct {
	Status   string
	Count    string
	Info     string
	Geocodes []geocoder
}

type geocoder struct {
	Formatted_address string
	Province          string
	City              string
	Citycode          string
	District          string
	Township          string
	Street            string
	Number            string
	Adcode            string
	Location          string
	Level             string
}
