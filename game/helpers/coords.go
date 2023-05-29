package helpers

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// NumericCords change coords to numeric representation
func NumericCords(coord string) (map[string]uint8, error) {
	if len(coord) != 2 && len(coord) != 3 {
		return map[string]uint8{}, errors.New("param coord must be of length 2 or 3")
	}
	translatedCoords := make(map[string]uint8, 2)
	coord = strings.ToLower(coord)
	if !strings.Contains("abcdefghij", string(coord[0])) {
		return map[string]uint8{}, errors.New("first letter must be from a-j, example: A4, a2")
	}
	translatedCoords["x"] = coord[0] - 97

	numbers, err := strconv.Atoi(coord[1:])
	if err != nil {
		return map[string]uint8{}, err
	}
	if numbers < 1 || numbers > 10 {
		return map[string]uint8{}, errors.New("second coord (number) must be between 1 and 10")
	}
	translatedCoords["y"] = uint8(numbers) - 1

	return translatedCoords, nil
}

func AlphabeticCoords(x, y int) string {
	letters := "abcdefghij"
	return string(letters[x]) + fmt.Sprint(y+1)
}
