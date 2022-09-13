package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func simple() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter first number: ")
	scanner.Scan()
	if _, err := strconv.Atoi(scanner.Text()); err != nil {
		fmt.Println("Invalid number!")
		return
	}
	firstNum,_ := strconv.Atoi(scanner.Text())

	fmt.Print("Enter operator [+ - * / %]: ")
	scanner.Scan()
	operator := scanner.Text()
	if len(operator) != 1 || !strings.Contains("+-*/%", operator) {
		fmt.Println("Invalid operator!")
		return
	}

	fmt.Print("Enter second number: ")
	scanner.Scan()
	if _, err := strconv.Atoi(scanner.Text()); err != nil {
		fmt.Println("Invalid number!")
		return
	}
	secondNum,_ := strconv.Atoi(scanner.Text())

	fmt.Print("Result: ")
	switch operator {
	case "+":
		fmt.Println(firstNum + secondNum)
	case "-":
		fmt.Println(firstNum - secondNum)
	case "*":
		fmt.Println(firstNum * secondNum)
	case "/":
		fmt.Println(float32(firstNum) / float32(secondNum))
	case "%":
		fmt.Println(firstNum % secondNum)
	}
}

func equation() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Only linear equations.\nOnly use a single type of variable.\nOnly numbers and [+ =].\n")
	fmt.Print("Enter equation: ")
	scanner.Scan()

	variable := ""
	left := true
	leftNum := 0
	rightNum := 0
	curNum := ""
	waitingForNum := true
	operation := ""	// var num

	for _, char := range strings.ToLower(scanner.Text()){
		switch {
		case char >= 97 && char <= 122: // letters
			if variable != "" && variable != string(char) {
				fmt.Println("Variable already set to " + variable)
			}
			variable = string(char)
		case char >= 48 && char <= 57: // digits

		case char == 43: // plus sign

		case char == 61: // equal sign
			if left {
				left = false
			} else {
				fmt.Println("Not an equation (too many ='s).")
				return
			}
			if waitingForNum {
				fmt.Print(string(char))
				fmt.Println("is not a number!")
				return
			}
		case char == 32: // space
			continue
		default:
			fmt.Print("Invalid value: ")
			fmt.Println(string(char))
			return
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("[Simple] or [Equation]: ")
	scanner.Scan()

	switch strings.ToLower(scanner.Text()) {
	case "simple":
		simple()
	case "equation":
		equation()
	default:
		fmt.Println("Not a valid input!")
	}
}