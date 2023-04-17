package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getCoord() (string, error) {
	var coord string
	fmt.Print("Enter coordinates: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		coord = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}

	var coordList []string
	for _, i := range coord {
		// Trim whitespaces from coord
		if i := strings.ToUpper(string(i)); strings.Contains("ABCDEFGHIJ0123456789", i) {
			coordList = append(coordList, i)
		}
	}

	if len(coordList) < 2 && len(coordList) > 3 {
		return "", errors.New("please enter a valid coordinate (too short or too long, ex. A4, G10)")
	}

	if !strings.Contains("ABCDEFGHIJ", coordList[0]) {
		return "", errors.New("please enter a valid coordinate (use letters from A to J)")
	}

	if number, err := strconv.Atoi(strings.Join(coordList[1:], "")); number > 10 || number < 1 {
		if err != nil {
			return "", err
		}
		return "", errors.New("please enter a valid coordinate (use number from 1 to 10")
	}
	return strings.Join(coordList, ""), nil
}
