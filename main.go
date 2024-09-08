package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	user    string
	balance = make(map[string]int)
	debt    = make(map[string]map[string]int)
	credit  = make(map[string]map[string]int)
)

func main() {
	snr := bufio.NewScanner(os.Stdin)

	for snr.Scan() {
		line := snr.Text()
		if len(line) == 0 {
			break
		}

		readCommand(line)
	}

	if err := snr.Err(); err != nil {
		if err != io.EOF {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func readCommand(str string) {
	input := strings.Split(str, " ")

	switch input[0] {
	// login
	case "login":
		if user != "" {
			fmt.Print("Please logout first!\n\n")
			return
		}

		if len(input) < 2 || input[1] == "" {
			fmt.Print("Please input an account name!\n\n")
			return
		}

		login(input[1])
		printBalance()

	// deposit
	case "deposit":
		if user == "" {
			fmt.Print("Please login first!\n\n")
			return
		}

		if len(input) < 2 {
			fmt.Print("Please input amount to deposit!\n\n")
			return
		}

		deposit(input[1])
		printBalance()
	// withdraw
	case "withdraw":
		if user == "" {
			fmt.Print("Please login first!\n\n")
			return
		}

		if len(input) < 2 {
			fmt.Print("Please input amount to withdraw!\n\n")
			return
		}

		withdraw(input[1])
		printBalance()

	// transfer
	case "transfer":
		if user == "" {
			fmt.Print("Please login first!\n\n")
			return
		}

		if input[1] == user {
			fmt.Print("Please input other account name!\n\n")
			return
		}

		if len(input) < 3 {
			fmt.Print("Please input account name and amount to transfer, separated by space!\n\n")
			return
		}

		transfer(input[1], input[2])
		printBalance()

	// logout
	case "logout":
		if user == "" {
			fmt.Print("Please login first!\n\n")
			return
		}
		logout()

	// default
	default:
		fmt.Print("Command unrecognized\n\n")
	}
}

func login(input string) {
	user = input
	if _, isExist := balance[user]; !isExist {
		balance[user] = 0
	}

	fmt.Println("Hello, " + user + "!")
}

func deposit(input string) {
	amount, _ := strconv.Atoi(input)

	// since there are no clarity in a case of which other user should be transferred into when a user has multiple debt,
	// then this will be just randomized based on whichever order the map set
	for i, v := range debt[user] {
		transferred := 0
		if v > amount {
			transferred = amount
			amount = 0
			debt[user][i] -= transferred
			credit[i][user] -= transferred
			balance[i] += transferred
		} else {
			transferred = v
			amount -= v
			delete(debt[user], i)
			delete(credit[i], user)
			balance[i] += transferred
		}

		fmt.Println("Transferred $" + strconv.Itoa(transferred) + " to " + i)
	}

	balance[user] += amount
}

func withdraw(input string) {
	amount, _ := strconv.Atoi(input)

	_, isExist := balance[user]
	if isExist {
		if amount > balance[user] {
			fmt.Println("Your balance is less than amount you want to draw")
			return
		}

		balance[user] -= amount
	}
}

func transfer(input1, input2 string) {
	amount, _ := strconv.Atoi(input2)
	transferred := 0

	// check if user has credit to other user being transferred
	if _, isExistCredit := credit[user]; isExistCredit {
		if _, isExistCredit2 := credit[user][input1]; isExistCredit2 {
			if amount > credit[user][input1] {
				delete(credit[user], input1)
				delete(debt[input1], user)
				amount -= credit[user][input1]
			} else {
				credit[user][input1] -= amount
				debt[input1][user] -= amount
				amount = 0
			}
		}
	}

	// check if user being transferred has an account
	if _, isExistBalance := balance[input1]; isExistBalance {
		if amount > balance[user] {
			// check if user has debt to other user being transferred
			if _, isExistDebt := debt[user]; isExistDebt {
				debt[user][input1] += amount - balance[user]
				if _, isExistCredit := credit[input1]; isExistCredit {
					credit[input1][user] += amount - balance[user]
				} else {
					credit[input1] = map[string]int{user: amount - balance[user]}
				}
			} else {
				debt[user] = map[string]int{input1: amount - balance[user]}
				credit[input1] = map[string]int{user: amount - balance[user]}
			}

			transferred = balance[user]
			balance[input1] += transferred
			balance[user] = 0
		} else {
			transferred = amount
			balance[input1] += transferred
			balance[user] -= transferred
		}
	} else {
		fmt.Println("There are no account with that name")
		return
	}

	if transferred > 0 {
		fmt.Println("Transferred $" + strconv.Itoa(transferred) + " to " + input1)
	}
}

func logout() {
	fmt.Println("Goodbye, " + user + "!\n")
	user = ""
}

func printBalance() {
	fmt.Println("Your balance is $" + strconv.Itoa(balance[user]))

	for i, v := range debt[user] {
		fmt.Println("Owed $" + strconv.Itoa(v) + " to " + i)
	}

	for i, v := range credit[user] {
		fmt.Println("Owed $" + strconv.Itoa(v) + " from " + i)
	}

	fmt.Println("")
}
