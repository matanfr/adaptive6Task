package main

type Data struct {
	countriesMap   map[string]int
	browsersMap    map[string]int
	osMap          map[OSType]int
	countriesCount int
	browsersCount  int
	osCount        int
}

func main() {
	logFilePath := "data/logs/apache_log.txt"
	geoLocationsFile := "data/GeoLite2-Country-CSV_20260224/GeoLite2-Country-Blocks-IPv4.csv"
	countriesFile := "data/GeoLite2-Country-CSV_20260224/GeoLite2-Country-Locations-en.csv"

	gp := &GeolocationParser_EN_IPV4{}
	gp.Init(countriesFile, geoLocationsFile)

	data := &Data{
		countriesMap:   make(map[string]int),
		browsersMap:    make(map[string]int),
		osMap:          make(map[OSType]int),
		countriesCount: 0,
		browsersCount:  0,
		osCount:        0,
	}
	ProcessLogFile(logFilePath, data, gp)
	PrintStats(data)
}
