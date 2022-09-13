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
	firstNum, _ := strconv.Atoi(scanner.Text())

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
	secondNum, _ := strconv.Atoi(scanner.Text())

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

	fmt.Print("Only linear equations.\nOnly use a single type of variable.\nVariable must come after number.\nOnly addition supported right now.\n")
	fmt.Print("Enter equation: ")
	scanner.Scan()

	variable := ""
	left := true
	leftNum := 0
	rightNum := 0
	curNum := ""
	leftVar := 0
	rightVar := 0
	operation := "" // var num

	for _, char := range strings.ToLower(scanner.Text()) {
		switch {
		case char >= 97 && char <= 122: // letters
			if variable != "" && variable != string(char) {
				fmt.Println("Variable already set to " + variable)
			}
			variable = string(char)
			if operation != "var" {
				operation = "var"
			} else {
				fmt.Println("Multiple variables in one term!")
				return
			}
		case char >= 48 && char <= 57: // digits
			curNum += string(char)
			if operation != "var" {
				operation = "num"
			}
		case char == 43: // plus sign
			addNum, _ := strconv.Atoi(curNum)
			curNum = ""
			switch {
			case left && operation == "num":
				leftNum += addNum
			case left && operation == "var":
				if addNum == 0 {
					leftVar += 1
				} else {
					leftVar += addNum
				}
			case !left && operation == "num":
				rightNum += addNum
			case !left && operation == "var":
				if addNum == 0 {
					rightVar += 1
				} else {
					rightVar += addNum
				}
			}
			operation = ""
		case char == 61: // equal sign
			addNum, _ := strconv.Atoi(curNum)
			curNum = ""
			switch {
			case left && operation == "num":
				leftNum += addNum
			case left && operation == "var":
				if addNum == 0 {
					leftVar += 1
				} else {
					leftVar += addNum
				}
			case !left && operation == "num":
				rightNum += addNum
			case !left && operation == "var":
				if addNum == 0 {
					rightVar += 1
				} else {
					rightVar += addNum
				}
			}
			operation = ""
			if left {
				left = false
			} else {
				fmt.Println("Not an equation (too many ='s).")
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
	addNum, _ := strconv.Atoi(curNum)
	if operation == "num" {
		rightNum += addNum
	} else {
		if addNum == 0 {
			rightVar += 1
		} else {
			rightVar += addNum
		}
	}
	fmt.Print("x = ")
	fmt.Println(float32(rightNum-leftNum) / float32(leftVar-rightVar))
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
