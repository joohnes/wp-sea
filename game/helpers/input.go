package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetAnswer(caseSensitive bool) (string, error) {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(answer) == 0 {
		return "", nil
	}
	if caseSensitive {
		return answer, nil
	}
	return strings.ToLower(answer), nil
}

func GetName() (string, error) {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your name: ")
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return answer, nil
}
func GetDesc() (string, error) {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter your description: ")
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return answer, nil
}
