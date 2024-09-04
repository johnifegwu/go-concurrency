package main

import (
	"fmt"
	"strings"
)

func wordCount(text string) string {
	words := strings.Fields(text)
	counts := map[string]int{}

	for _, word := range words {
		counts[strings.ToLower(word)]++
	}

	return fmt.Sprintln(counts)
}

func main() {

	text := `Obil was the former Governor 
	of Anambra State, he was also the former 
	presidential candidate for the Labor party.
	`
	fmt.Println(wordCount(text))

	crypto := map[string]float64{
		"BTC":  64000.25,
		"ETH":  3000,
		"SHIB": 0.00055478,
	}

	// print the length of the ma
	fmt.Printf("lenght: %v, \n", len(crypto))

	//print all
	for key, value := range crypto {

		fmt.Printf("%v", key)
		fmt.Printf(" : %v, \n", value)
	}
	/*
		// Slices
		nameSlice := []string{"John", "Paul", "James"}

		nameSlice = append(nameSlice, "Kalu")

		// Say helo
		for _, name := range nameSlice {
			fmt.Printf("helo %v,\n", name)
		}
	*/

	/*
		count := 0

		// Even ended numbers
		for a := 1000; a <= 9999; a++ {

			for b := 1000; b <= 9999; b++ {
				n := a * b

				// if a*b is even ended
				s := fmt.Sprintf("%v", n)

				if s[0] == s[len(s)-1] {

					// increment count
					count++
				}

			}
		}

		fmt.Println(count)
	*/

	/*
		x, y := 3.4, 6.8

		r := y * x

		// Using fmt.Printf for formatted output
		fmt.Printf("y=%v, type of %T\n", y, y)
		fmt.Printf("x=%v, type of %T\n", x, x)
		fmt.Printf("r=%v, type of %T\n", r, r)

		for i := 1; i <= 20; i++ {
			if i%3 == 0 && i%5 == 0 {
				// if number is divisible by bith 3 and 5 print fizz buzz
				fmt.Println("fizz buzz")
			} else if i%3 == 0 {
				// if number is divisible by 3 print fizz
				fmt.Println("fizz")
			} else if i%5 == 0 {
				// if number is divisible by 5 print buzz
				fmt.Println("buzz")
			} else {
				// print the number
				fmt.Println(i)
			}
		}
	*/
}
