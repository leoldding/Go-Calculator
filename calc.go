package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

func basic() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter first number: ")

	scanner.Scan()

	if _, err := strconv.ParseFloat(scanner.Text(), 32); err != nil {
		fmt.Println("Invalid number!")
		return
	}
	firstNum, _ := strconv.ParseFloat(scanner.Text(), 32)

	fmt.Print("Enter operator [+ - * /]: ")
	scanner.Scan()
	operator := scanner.Text()
	if len(operator) != 1 || !strings.Contains("+-*/%", operator) {
		fmt.Println("Invalid operator!")
		return
	}

	fmt.Print("Enter second number: ")
	scanner.Scan()
	if _, err := strconv.ParseFloat(scanner.Text(), 32); err != nil {
		fmt.Println("Invalid number!")
		return
	}
	secondNum, _ := strconv.ParseFloat(scanner.Text(), 32)

	fmt.Print("Result: ")
	fmt.Println(operate(operator, float32(firstNum), float32(secondNum)))
}

func operate(operator string, firstNum float32, secondNum float32) float32 {
	switch operator {
	case "+":
		return firstNum + secondNum
	case "-":
		return firstNum - secondNum
	case "*":
		return firstNum * secondNum
	case "/":
		return firstNum / secondNum
	}
	fmt.Println("Not a valid operator: " + operator)
	return -1.0
}

func solve() {
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
				return
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
	if variable == "" || (leftVar-rightVar) == 0 {
		fmt.Println("Equation not solvable.")
	} else {
		fmt.Print(variable + " = ")
		fmt.Println(float32(rightNum-leftNum) / float32(leftVar-rightVar))
	}
}

func simplify() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter equation: ")

	scanner.Scan()

	equation := strings.ReplaceAll(strings.ToLower(scanner.Text()), " ", "") // remove whitespace

	equalPos := strings.Index(equation, "=")                               // split equation
	if equalPos == -1 || strings.Index(equation[equalPos+1:], "=") != -1 { // check if input is valid
		fmt.Println("Not an equation!")
		return
	}

	leftStr := equation[:equalPos]
	rightStr := equation[equalPos+1:]

	if leftStr != "" && rightStr != "" {
		leftChan := make(chan map[int]float32, 1)
		rightChan := make(chan map[int]float32, 1)
		leftErrChan := make(chan string, 1)
		rightErrChan := make(chan string, 1)

		wg.Add(1)

		problem := false
		go func() {
			calculate(leftStr+" ", leftChan, leftErrChan)

			close(leftChan)
			close(leftErrChan)

			leftErr := <-leftErrChan
			if leftErr != "" {
				fmt.Println("Error on left side of equation: " + leftErr)
				problem = true
			}
		}()

		wg.Add(1)
		go func() {
			calculate(rightStr+" ", rightChan, rightErrChan)

			close(rightChan)
			close(rightErrChan)

			rightErr := <-rightErrChan
			if rightErr != "" {
				fmt.Println("Error on right side of equation: " + rightErr)
				problem = true
			}
		}()

		wg.Wait()

		if !problem {
			leftVals := <-leftChan
			rightVals := <-rightChan

			for key, val := range rightVals {
				leftVals[key] -= val
			}

			var keys []int
			for key := range leftVals {
				keys = append(keys, key)
			}

			sort.Ints(keys)

			fmt.Print("Simplified Equation: ")

			for i := len(keys) - 1; i >= 0; i-- {
				if leftVals[keys[i]] < 0 {
					fmt.Print(leftVals[keys[i]])
					if keys[i] > 1 {
						fmt.Print("x^(")
						fmt.Print(keys[i])
						fmt.Print(")")
					} else if keys[i] == 1 {
						fmt.Print("x")
					}
				} else if leftVals[keys[i]] > 0 {
					fmt.Print("+")
					fmt.Print(leftVals[keys[i]])
					if keys[i] > 1 {
						fmt.Print("x^(")
						fmt.Print(keys[i])
						fmt.Print(")")
					} else if keys[i] == 1 {
						fmt.Print("x")
					}
				}
			}

			fmt.Println()
		}
	} else {
		fmt.Println("At least one side of equation is empty.")
	}
}

func calculate(equation string, channel chan map[int]float32, errChan chan string) {
	defer wg.Done()
	values := make(map[int]float32)
	var powers []int64
	var numbers []float32
	var operations []string
	var power int64 = 0
	var variable = ""
	var val int64
	num := ""
	pointer := 0

	for i := 0; i < len(equation); i++ {
		char := equation[i]
		switch {
		case char >= 48 && char <= 57: // digits
			num += string(char)
		case char >= 97 && char <= 122: // letters
			if variable == "" {
				variable = string(char)
			}
			if string(char) != variable {
				errChan <- "Too many different variables!"
				goto END
			} else if power != 0 {
				errChan <- "Two variables in one term!"
				goto END
			} else {
				if num == "" {
					num = "1"
				}
				if string(equation[i+1]) != "^" {
					power = 1
					if equation[i+1] >= 48 && equation[i+1] <= 57 {
						errChan <- "Missing '^'"
						goto END
					}
				} else {
					count := 0
					for j := i + 2; j < len(equation); j++ {
						if equation[j] < 48 || equation[j] > 57 {
							break
						}
						count += 1
					}
					power, _ = strconv.ParseInt(equation[i+2:i+2+count], 10, 32)
					i = i + count + 1
				}
			}
		case char == 43: // addition
			if num == "" {
				errChan <- "Missing number."
				goto END
			}
			val, _ = strconv.ParseInt(num, 10, 32)
			numbers = append(numbers, float32(val))
			powers = append(powers, power)
			operations = append(operations, "+")
			num = ""
			power = 0
		case char == 45: // subtraction
			if num != "" {
				val, _ = strconv.ParseInt(num, 10, 32)
				numbers = append(numbers, float32(val))
				powers = append(powers, power)
				operations = append(operations, "+")
				num = "-"
				power = 0
			} else {
				num += "-"
			}
		case char == 42: // multiplication
			if num == "" {
				errChan <- "Missing number."
				goto END
			}
			val, _ = strconv.ParseInt(num, 10, 32)
			numbers = append(numbers, float32(val))
			powers = append(powers, power)
			operations = append(operations, "*")
			num = ""
			power = 0
		case char == 47: // division
			if num == "" {
				errChan <- "Missing number."
				goto END
			}
			val, _ = strconv.ParseInt(num, 10, 32)
			numbers = append(numbers, float32(val))
			powers = append(powers, power)
			operations = append(operations, "/")
			num = ""
			power = 0
		case char == 32: // space at end
			continue
		default:
			errChan <- "Invalid character: " + string(char)
			goto END
		}
	}
	val, _ = strconv.ParseInt(num, 10, 32)
	numbers = append(numbers, float32(val))
	powers = append(powers, power)

	for {
		if pointer >= len(operations) {
			break
		}
		if operations[pointer] == "*" || operations[pointer] == "/" {
			if operations[pointer] == "*" {
				numbers[pointer] *= numbers[pointer+1]
				powers[pointer] += powers[pointer+1]
			} else {
				numbers[pointer] /= numbers[pointer+1]
				powers[pointer] -= powers[pointer+1]
			}
			if pointer < len(operations)-1 {
				numbers = append(numbers[:pointer+1], numbers[pointer+2:]...)
				powers = append(powers[:pointer+1], powers[pointer+2:]...)
				operations = append(operations[:pointer], operations[pointer+1:]...)
			} else {
				numbers = numbers[:pointer+1]
				powers = powers[:pointer+1]
				operations = operations[:pointer]
			}
		} else {
			pointer += 1
		}
	}

	for i := 0; i < len(numbers); i++ {
		values[int(powers[i])] += numbers[i]
	}

	channel <- values
END:
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("[Basic], [Simplify], or [Solve]: ")
	scanner.Scan()

	switch strings.ToLower(scanner.Text()) {
	case "basic":
		basic()
	case "simplify":
		simplify()
	case "solve":
		solve()
	default:
		fmt.Println("Not a valid input!")
	}
}
