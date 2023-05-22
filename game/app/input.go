package app

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (a *App) getCoord() (string, error) { // dodac obsluge błędu
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
	if coord == "quit" {
		err := a.client.Resign()
		if err != nil {
			return "", err
		}
		os.Exit(1) //TODO: Zmienić na ładniejsze
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
			fmt.Println(err)
			os.Exit(1)
		}
		return "", errors.New("please enter a valid coordinate (use number from 1 to 10)")
	}
	return strings.ToLower(strings.Join(coordList, "")), nil
}

func (a *App) getAnswer() (string, error) {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(string(answer)) == 0 {
		return "", nil
	}
	return strings.ToLower(string(answer[0])), nil
}

func (a *App) getName() error {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your name: ")
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	a.nick = answer
	return nil
}
func (a *App) getDesc() error {
	var answer string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter your description: ")
	for scanner.Scan() {
		answer = scanner.Text()
		break
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	a.desc = answer
	return nil
}
