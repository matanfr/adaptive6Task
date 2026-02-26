package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func ParseEntry(entry string, data *Data) {
	re := regexp.MustCompile(`^(\S+) .*?"[^"]*" \d+ \d+ "[^"]*" "([^"]+)"`)
	matches := re.FindStringSubmatch(entry)

	if len(matches) < 3 {
		fmt.Println("Failed to parse log line")
		return
	}

	ip := matches[1]
	userAgent := matches[2]

	// Parse country
	countryName := GetCountryFromIP(ip)
	data.countriesMap[countryName]++
	data.countriesCount++

	// Parse browser
	browser := ""
	uaParts := strings.Split(userAgent, " ")
	if len(uaParts) > 0 {
		browser = strings.Split(uaParts[0], "/")[0]
	}
	data.browsersMap[browser]++
	data.browsersCount++

	// Parse OS
	detectedOS := OSUnknown
	for _, os := range allOS {
		if strings.Contains(userAgent, string(os)) {
			detectedOS = os
			break
		}
	}

	// Special case for Darwin -> MacOS
	if detectedOS == OSUnknown && strings.Contains(userAgent, "Darwin") {
		detectedOS = OSMac
	}
	data.osMap[detectedOS]++
	data.osCount++
}

func ProcessLogFile(logFilePath string, data *Data) {
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		entry := scanner.Text()
		ParseEntry(entry, data)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

}
