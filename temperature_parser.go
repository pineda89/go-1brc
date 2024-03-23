package main

func parseTemperature(temperature []byte) int {
	if len(temperature) == 0 {
		return 0
	}

	val := int(temperature[len(temperature)-1] - '0' + (temperature[len(temperature)-3]-'0')*10)

	switch len(temperature) {
	case 3: // "9.9"
		return val
	case 4: // "99.9" , "-9.9"
		if temperature[0] == '-' {
			return val * -1
		} else {
			return val + int((temperature[0]-'0')*100)
		}
	case 5: // "-99.9"
		return (val + (int(temperature[1]-'0') * 100)) * -1
	}

	return val
}
