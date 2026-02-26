package main

import "fmt"

func PrintStats(data *Data) {

	// fmt.Println(countriesMap)
	// fmt.Println(browsersMap)
	// fmt.Println(osMap)

	// fmt.Printf("Total countries: %d\n", countriesCount)
	// fmt.Printf("Total browsers: %d\n", browsersCount)
	// fmt.Printf("Total os: %d\n", osCount)

	// Print countries
	fmt.Println("Country:")
	for key, value := range data.countriesMap {
		percentage := float64(value) / float64(data.countriesCount) * 100
		fmt.Printf("%s: %.2f%%\n", key, percentage)
	}

	// Print OS
	fmt.Println("\nOS:")
	for key, value := range data.osMap {
		percentage := float64(value) / float64(data.osCount) * 100
		fmt.Printf("%s: %.2f%%\n", key, percentage)
	}

	// Print browsers
	fmt.Println("\nBrowser:")
	for key, value := range data.browsersMap {
		percentage := float64(value) / float64(data.countriesCount) * 100
		fmt.Printf("%s: %.2f%%\n", key, percentage)
	}

}
