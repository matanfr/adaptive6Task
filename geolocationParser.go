package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"inet.af/netaddr"
)

type GeolocationParser interface {
	Init(countriesFile string, geolocationsFile string)
	GetCountryFromIP(ip string) string
}

type Block struct {
	Prefix netaddr.IPPrefix
	Line   string
}

type GeolocationParser_EN_IPV4 struct {
	blocks            []Block
	countriesToGeoMap map[string]string
}

func (gp *GeolocationParser_EN_IPV4) Init(countriesFile string, geolocationsFile string) {
	gp.countriesToGeoMap = make(map[string]string)
	err := gp.loadCountries(countriesFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	gp.loadGeolocations(geolocationsFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
}

func (gp *GeolocationParser_EN_IPV4) loadCountries(countriesFile string) error {
	file, err := os.Open(countriesFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	firstLine := true

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if firstLine {
			firstLine = false
			continue // skip header
		}

		if len(record) < 6 {
			continue // invalid line
		}

		id := record[0]
		countryName := record[5]
		countryName = strings.TrimSpace(countryName)

		if id != "" && countryName != "" {
			gp.countriesToGeoMap[id] = countryName
		}
	}

	// No errors
	return nil
}

func (gp *GeolocationParser_EN_IPV4) loadGeolocations(geolocationsFile string) error {
	file, err := os.Open(geolocationsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.FieldsPerRecord = -1

	firstLine := true
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}

		if firstLine {
			firstLine = false
			continue // skip header
		}

		network := record[0]
		if network == "" {
			continue
		}

		prefix, err := netaddr.ParseIPPrefix(network)
		if err != nil {
			continue
		}

		gp.blocks = append(gp.blocks, Block{
			Prefix: prefix,
			Line:   strings.Join(record, ","),
		})
	}

	// Sort by starting IP of the prefix
	sort.Slice(gp.blocks, func(i, j int) bool {
		return gp.blocks[i].Prefix.Masked().IP().Compare(gp.blocks[j].Prefix.Masked().IP()) < 0
	})

	// No errors
	return nil
}

func (gp *GeolocationParser_EN_IPV4) GetCountryFromIP(ip string) string {
	geonameID := gp.getGeonameIDFromIP(ip)
	countryName := gp.getCountryFromGeonameID(geonameID)
	return countryName
}

func (gp *GeolocationParser_EN_IPV4) getGeonameIDFromIP(ipStr string) string {
	block := &Block{}
	var geonameId string

	ip, err := netaddr.ParseIP(ipStr)
	if err != nil {
		fmt.Println("could not parse IP")
		return ""
	}

	low, high := 0, len(gp.blocks)-1
	for low <= high {
		mid := (low + high) / 2
		p := gp.blocks[mid].Prefix

		if p.Contains(ip) {
			block = &gp.blocks[mid]
		}

		if ip.Less(p.Masked().IP()) {
			high = mid - 1
		} else {
			low = mid + 1
		}
	}

	parts := strings.Split(block.Line, ",")

	if len(parts) > 1 {
		geonameId = parts[1]
		return geonameId
	}

	fmt.Println("Error parsing geoname_id")
	return ""

}

func (gp *GeolocationParser_EN_IPV4) getCountryFromGeonameID(geonameID string) string {
	if name, ok := gp.countriesToGeoMap[geonameID]; ok {
		return name
	}
	return "Unknown"
}
