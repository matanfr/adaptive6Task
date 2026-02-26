package main

type OSType string

const (
	OSWindows OSType = "Windows"
	OSMac     OSType = "Mac OS"
	OSLinux   OSType = "Linux"
	OSUnknown OSType = "Unknown"
)

var allOS = []OSType{
	OSWindows,
	OSMac,
	OSLinux,
}
