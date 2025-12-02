package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Dial struct {
	currentPosition int
	minimumPosition int
	maximumPosition int
	numMinimum      uint
}

func (d *Dial) MoveRight(amount int) {
	for i := 0; i < amount; i++ {
		if d.currentPosition == d.maximumPosition {
			d.currentPosition = d.minimumPosition
		} else {
			d.currentPosition++
		}
	}
}
func (d *Dial) MoveLeft(amount int) {
	for i := 0; i < amount; i++ {
		if d.currentPosition == d.minimumPosition {
			d.currentPosition = d.maximumPosition
		} else {
			d.currentPosition--
		}
	}
}

func (d *Dial) Move(direction string, amount int) error {
	switch direction {
	case "L":
		d.MoveLeft(amount)
	case "R":
		d.MoveRight(amount)
	default:
		return fmt.Errorf("%s is not a valid direction", direction)
	}

	if d.currentPosition == d.minimumPosition {
		d.numMinimum++
	}
	return nil
}

func getInput() string {
	_, err := os.Stat("input/data.txt")
	if err != nil {
		client := http.Client{
			Timeout: 5 * time.Second,
		}

		req, err := http.NewRequest(
			"GET",
			"https://adventofcode.com/2025/day/1/input",
			nil,
		)

		session, err := os.ReadFile("../cookie.txt")

		req.AddCookie(&http.Cookie{
			Name:  "session",
			Value: string(session),
		})
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		contents, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create("input/data.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, err = file.Write(contents)
		if err != nil {
			log.Fatal(err)
		}
		return string(contents)
	}
	inputFile, err := os.Open("input/data.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()
	inputScanner := bufio.NewScanner(inputFile)
	inputScanner.Split(bufio.ScanLines)
	var lines []string
	for inputScanner.Scan() {
		lines = append(lines, inputScanner.Text())
	}
	return strings.Join(lines, "\n")
}

func main() {
	dial := Dial{
		currentPosition: 50,
		minimumPosition: 0,
		maximumPosition: 99,
		numMinimum:      0,
	}

	contents := getInput()
	for lineNum, line := range strings.Split(contents, "\n") {
		direction := string(line[0])
		amount, err := strconv.Atoi(line[1:])
		if err != nil {
			slog.Error(
				"failed to parse amount from input",
				slog.Int("line", lineNum),
				slog.String("contents", line),
				slog.Any("err", err),
			)
			os.Exit(1)
		}
		err = dial.Move(direction, amount)
		if err != nil {
			slog.Warn(
				"failed to move direction",
				slog.Int("line", lineNum),
				slog.String("contents", line),
				slog.Any("err", err),
			)
		}
	}
	slog.Info("done rotating dial", "realPass", dial.numMinimum)
}
