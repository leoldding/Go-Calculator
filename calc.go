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
var variable string

// basic calculator function
func basic() {
	scanner := bufio.NewScanner(os.Stdin)

	// read in first number
	fmt.Print("Enter first number: ")
	scanner.Scan()
	firstNum, err := strconv.ParseFloat(scanner.Text(), 32)
	if err != nil { // break if input is not a number
		fmt.Println("Invalid number!")
		return
	}

	// read in operator
	fmt.Print("Enter operator [+ - * /]: ")
	scanner.Scan()
	operator := scanner.Text()
	if len(operator) != 1 || !strings.Contains("+-*/", operator) { // break if input is not a valid operator
		fmt.Println("Invalid operator!")
		return
	}

	// read in second number
	fmt.Print("Enter second number: ")
	scanner.Scan()
	secondNum, err := strconv.ParseFloat(scanner.Text(), 32)
	if err != nil { // break if input is not a number
		fmt.Println("Invalid number!")
		return
	}

	// output result
	fmt.Print("Result: ")
	fmt.Println(operate(operator, float32(firstNum), float32(secondNum)))
}

// function to calculate basic operations
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

// simplifies any polynomial or polynomial equation regardless of power
func simplify() {
	scanner := bufio.NewScanner(os.Stdin)

	// read in equation
	fmt.Print("Enter equation: ")
	scanner.Scan()
	equation := strings.ReplaceAll(strings.ToLower(scanner.Text()), " ", "") // remove whitespace

	// split equation into two halves
	equalPos := strings.Index(equation, "=")
	if equalPos == -1 || strings.Index(equation[equalPos+1:], "=") != -1 { // break if input is not an equation
		fmt.Println("Not an equation!")
		return
	}
	leftStr := equation[:equalPos]
	rightStr := equation[equalPos+1:]

	if leftStr == "" || rightStr == "" { // break if one side of equation is empty
		fmt.Println("At least one side of equation is empty.")
		return
	}

	// create channels
	leftChan := make(chan map[int]float32, 1)
	rightChan := make(chan map[int]float32, 1)
	leftErrChan := make(chan string, 1)
	rightErrChan := make(chan string, 1)

	// error checker
	problem := false

	// add to waitgroup for syncing
	wg.Add(1)

	// go routine that processes left side
	go func() {
		calculate(leftStr+" ", leftChan, leftErrChan)

		// close channels
		close(leftChan)
		close(leftErrChan)

		// check if there was an error
		leftErr := <-leftErrChan
		if leftErr != "" {
			fmt.Println("Error on left side of equation: " + leftErr)
			problem = true
		}
	}()

	// add to waitgroup for syncing
	wg.Add(1)

	// go routine that processes right side
	go func() {
		calculate(rightStr+" ", rightChan, rightErrChan)

		// close channels
		close(rightChan)
		close(rightErrChan)

		// check if there was an erro
		rightErr := <-rightErrChan
		if rightErr != "" {
			fmt.Println("Error on right side of equation: " + rightErr)
			problem = true
		}
	}()

	// wait for go routines to finish
	wg.Wait()

	// continue if there were no problems with processing
	if !problem {
		// retrieve maps
		leftVals := <-leftChan
		rightVals := <-rightChan

		// bring all values from "right" to "left"
		for key, val := range rightVals {
			leftVals[key] -= val
		}

		// retrieve and sort map keys
		var keys []int
		for key := range leftVals {
			keys = append(keys, key)
		}
		sort.Ints(keys)

		fmt.Print("Simplified Equation: ")

		// output simplified equation
		for i := len(keys) - 1; i >= 0; i-- {
			if leftVals[keys[i]] < 0 { // negative values
				fmt.Print(" - ")
				if keys[i]-keys[0] == 0 { // constants
					fmt.Print(leftVals[keys[i]] * -1)
				} else if keys[i]-keys[0] == 1 { // power of one
					if leftVals[keys[i]] != -1 {
						fmt.Print(leftVals[keys[i]] * -1)
					}
					fmt.Print(variable)
				} else if keys[i]-keys[0] != 0 { // non-zero powers
					if leftVals[keys[i]] != -1 {
						fmt.Print(leftVals[keys[i]] * -1)
					}
					fmt.Print(variable + "^(")
					fmt.Print(keys[i] - keys[0])
					fmt.Print(")")
				}
			} else if leftVals[keys[i]] > 0 { // positive values
				fmt.Print(" + ")
				if keys[i]-keys[0] == 0 { // constants
					fmt.Print(leftVals[keys[i]])
				} else if keys[i]-keys[0] == 1 { // power of one
					if leftVals[keys[i]] != 1 {
						fmt.Print(leftVals[keys[i]])
					}
					fmt.Print(variable)
				} else if keys[i]-keys[0] != 0 { // non-zero powers
					if leftVals[keys[i]] != 1 {
						fmt.Print(leftVals[keys[i]])
					}
					fmt.Print(variable + "^(")
					fmt.Print(keys[i] - keys[0])
					fmt.Print(")")
				}
			}
		}
		fmt.Println()
	}
}

// function to process polynomials
func calculate(equation string, channel chan map[int]float32, errChan chan string) {
	defer wg.Done()

	values := make(map[int]float32) // map for constants
	var powers []int64              // array to track power of each term
	var numbers []float32           // array to track constants of each term
	var operations []string         // array to track operations between terms
	var power int64 = 0
	variable = ""
	var val int64
	num := ""
	pointer := 0 // indexing variable for processing later

	// iterate through characters for processing
	for i := 0; i < len(equation); i++ {
		char := equation[i] // current character being processed
		switch {

		case char >= 48 && char <= 57: // digits
			num += string(char) // append digits

		case char >= 97 && char <= 122: // letters
			if variable == "" { // set variable
				variable = string(char)
			}
			if string(char) != variable { // letter is different from set variable
				errChan <- "Too many different variables!"
				goto END
			} else if power != 0 { // variable appears at least twice in one term (i.e. 2xx; not supported)
				errChan <- "Two variables in one term!"
				goto END
			} else {
				if num == "" { // variable without constant
					num = "1"
				} else if num == "-" {
					num = "-1"
				}
				if string(equation[i+1]) != "^" { // check if power follows variable
					power = 1
					if equation[i+1] >= 48 && equation[i+1] <= 57 { // number follows variable without ^ denoting power
						errChan <- "Missing '^'"
						goto END
					}
				} else {
					count := 0 // moves index i to correct spot after finding power
					// find the power of the variable
					for j := i + 2; j < len(equation); j++ {
						if (equation[j] < 48 || equation[j] > 57) && equation[j] != 45 {
							break
						}
						count += 1
					}
					power, _ = strconv.ParseInt(equation[i+2:i+2+count], 10, 32)
					i = i + count + 1
				}
			}

		case char == 42 || char == 43 || char == 47: // multiplication / addition / division
			if num == "" { // break if no number precedes operator
				errChan <- "Missing number before operator."
				goto END
			}
			val, _ = strconv.ParseInt(num, 10, 32) // parse number value
			numbers = append(numbers, float32(val))
			powers = append(powers, power)
			// add appropriate operator
			if char == 42 {
				operations = append(operations, "*")
			} else if char == 43 {
				operations = append(operations, "+")
			} else {
				operations = append(operations, "/")
			}
			num = ""
			power = 0

		case char == 45: // subtraction
			if num != "" {
				val, _ = strconv.ParseInt(num, 10, 32) // parse number value
				numbers = append(numbers, float32(val))
				powers = append(powers, power)
				operations = append(operations, "+") // subtraction is just addition of a negative number
				num = "-"
				power = 0
			} else { // leading negative sign
				num += "-"
			}

		case char == 32: // skip the additional space at end
			continue

		default: // not a valid character input
			errChan <- "Invalid character: " + string(char)
			goto END
		}
	}

	// add last number and power
	val, _ = strconv.ParseInt(num, 10, 32)
	numbers = append(numbers, float32(val))
	powers = append(powers, power)

	// do all multiplications and divisions
	for {
		if pointer >= len(operations) { // breakpoint
			break
		}
		if operations[pointer] == "*" || operations[pointer] == "/" {
			if operations[pointer] == "*" { // multiplication
				numbers[pointer] *= numbers[pointer+1]
				powers[pointer] += powers[pointer+1]
			} else { // division
				numbers[pointer] /= numbers[pointer+1]
				powers[pointer] -= powers[pointer+1]
			}
			//remove processed values
			if pointer < len(operations)-1 {
				numbers = append(numbers[:pointer+1], numbers[pointer+2:]...)
				powers = append(powers[:pointer+1], powers[pointer+2:]...)
				operations = append(operations[:pointer], operations[pointer+1:]...)
			} else {
				numbers = numbers[:pointer+1]
				powers = powers[:pointer+1]
				operations = operations[:pointer]
			}
		} else { // increment pointer; skip addition
			pointer += 1
		}
	}

	// add all same-powered terms together
	for i := 0; i < len(numbers); i++ {
		values[int(powers[i])] += numbers[i]
	}
	// send values back through channel
	channel <- values

END:
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// input of calculator type
	fmt.Print("[Basic] or [Simplify]")
	scanner.Scan()

	switch strings.ToLower(scanner.Text()) {
	case "basic":
		basic()
	case "simplify":
		simplify()
	default:
		fmt.Println("Not a valid input!")
	}
}
