package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Line struct {
	Id     string
	Values interface{}
}

func main() {

	f, _ := os.Open("./import.txt")

	scanner := bufio.NewScanner(f)
	var lines []Line

	for scanner.Scan() {

		line := scanner.Text()
		split_line := strings.Split(line, " ")

		number := strings.Split(string(split_line[0]), "N")

		lines = append(lines, Line{
			Id:     string(number[1]),
			Values: split_line[1:],
		})

	}

	result, err := worker(lines)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)

}

func worker(lines []Line) ([]Line, error) {
	return lines, nil
}
