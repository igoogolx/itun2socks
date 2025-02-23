package constants

var hubAddress = ""

func SetHubAddress(address string) {
	hubAddress = address
}

func HubAddress() string {
	return hubAddress
}
