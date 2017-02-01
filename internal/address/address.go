package address

import (
	"context"
	"errors"
	"fmt"
	"googlemaps.github.io/maps"
	"log"
	"strings"
)

var (
	ErrNoAddress = errors.New("no results for address")
	ErrNotUS     = errors.New("not US")
)

type Address struct {
	Latitude  float64
	Longitude float64
	Address   []string
}

func Get(ctx context.Context, name, address string, client *maps.Client) (*Address, error) {
	geocoded, err := client.Geocode(ctx, &maps.GeocodingRequest{
		Address: address,
	})
	if err != nil {
		log.Printf("unable to geocode: %v\n", err)
		return nil, err
	}

	if len(geocoded) == 0 {
		log.Printf("no results for address %v\n", address)
		return nil, ErrNoAddress
	}

	formatted := &Address{
		Latitude:  geocoded[0].Geometry.Location.Lat,
		Longitude: geocoded[0].Geometry.Location.Lng,
	}

	var USok bool
	var administrative_area_level_1 []string
	var postal_code []string
	var locality []string
	var route []string
	var street_number []string
	var subpremise []string
	for _, v := range geocoded[0].AddressComponents {
		for _, t := range v.Types {
			switch t {
			case "country":
				if v.ShortName == "US" {
					USok = true
				}
			case "administrative_area_level_1":
				administrative_area_level_1 = append(administrative_area_level_1, v.ShortName)
			case "postal_code":
				postal_code = append(postal_code, v.LongName)
			case "locality":
				locality = append(locality, v.LongName)
			case "route":
				route = append(route, v.LongName)
			case "street_number":
				street_number = append(street_number, v.LongName)
			case "subpremise":
				subpremise = append(subpremise, v.LongName)
			}
		}
	}

	if !USok {
		log.Printf("address must be in the US\n")
		return nil, ErrNotUS
	}

	var addr []string
	if len(name) > 0 {
		addr = append(addr, name)
	}
	addr = append(addr, fmt.Sprintf("%s %s", strings.Join(street_number, " "), strings.Join(route, " ")))
	if len(subpremise) > 0 {
		addr = append(addr, fmt.Sprintf("#%s", strings.Join(subpremise, " ")))
	}
	addr = append(addr, fmt.Sprintf("%s, %s %s", strings.Join(locality, " "), strings.Join(administrative_area_level_1, " "), strings.Join(postal_code, " ")))

	formatted.Address = addr
	return formatted, nil
}
