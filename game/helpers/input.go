package helpers

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetAnswer() (string, error) {
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
	return strings.ToLower(answer), nil
}

func GetName() (string, error) {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name: ")
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
	fmt.Println("Enter your description: ")
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return answer, nil
}
