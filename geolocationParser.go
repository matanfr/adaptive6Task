package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"

	"inet.af/netaddr"
)

type Block struct {
	Prefix netaddr.IPPrefix
	Line   string
}

var blocks []Block
var countriesToGeoMap = make(map[string]string)

func LoadCountries(countriesFile string) error {
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
			countriesToGeoMap[id] = countryName
		}
	}

	// No errors
	return nil
}

func LoadGeolocations(filename string) error {
	file, err := os.Open(filename)
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

		blocks = append(blocks, Block{
			Prefix: prefix,
			Line:   strings.Join(record, ","),
		})
	}

	// Sort by starting IP of the prefix
	sort.Slice(blocks, func(i, j int) bool {
		return blocks[i].Prefix.Masked().IP().Compare(blocks[j].Prefix.Masked().IP()) < 0
	})

	// No errors
	return nil
}

func GetGeonameIDFromIP(ipStr string) string {
	block := &Block{}
	var geonameId string

	ip, err := netaddr.ParseIP(ipStr)
	if err != nil {
		fmt.Println("could not parse IP")
		return ""
	}

	low, high := 0, len(blocks)-1
	for low <= high {
		mid := (low + high) / 2
		p := blocks[mid].Prefix

		if p.Contains(ip) {
			block = &blocks[mid]
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

func GetCountryFromIP(ip string) string {
	geonameID := GetGeonameIDFromIP(ip)
	countryName := GetCountryFromGeonameID(geonameID)
	return countryName
}

func GetCountryFromGeonameID(geonameID string) string {
	if name, ok := countriesToGeoMap[geonameID]; ok {
		return name
	}
	return "Unknown"
}
