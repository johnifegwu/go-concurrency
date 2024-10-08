package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pkgerrors "github.com/pkg/errors"
)

// generics
type Ordered interface {
	int | float64 | string
}

func min[T Ordered](values []T) (T, error) {
	if len(values) == 0 {
		var zero T
		return zero, errors.New("min of empty slice")
	}
	m := values[0]

	for _, v := range values {
		if v < m {
			m = v
		}
	}

	return m, nil
}

// Capper implements io.Writer and turn everything to uppercase.
type Capper struct {
	wrt io.Writer
}

// Write implements io.Writer.
func (c *Capper) Write(p []byte) (n int, err error) {
	diff := byte('a' - 'A')

	out := make([]byte, len(p))

	for i, c := range p {

		if c >= 'a' && c <= 'z' {
			c -= diff
		}

		out[i] = c
	}

	return c.wrt.Write(out)
}

type Squre struct {
	x      int
	y      int
	lenght int
}

func NewSqure(x int, y int, length int) (*Squre, error) {
	if length <= 0 {
		return nil, errors.New("length must be greater then zero")
	}

	s := Squre{
		x:      x,
		y:      y,
		lenght: length,
	}

	return &s, nil
}

func (s *Squre) Move(dx int, dy int) {
	s.x = dx
	s.y = dy
}

func (s Squre) Area() int {
	return s.lenght * s.lenght
}

func wordCount(text string) string {
	words := strings.Fields(text)
	counts := map[string]int{}

	for _, word := range words {
		counts[strings.ToLower(word)]++
	}

	return fmt.Sprintln(counts)
}

func contentType(url string) (string, error) {
	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	conType := resp.Header.Get("Content-Type")

	if conType == "" {
		return "", errors.New("comtent-type header not found")
	}

	return conType, nil
}

func siteSerial(urls []string) {
	for _, url := range urls {
		ctype, err := contentType(url)
		if err == nil {
			fmt.Printf("%s -> %s, \n", url, ctype)
		}
	}
}

// Goroutines example
func siteSerialConcurrent(urls []string) {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			ctype, err := contentType(url)
			if err == nil {
				fmt.Printf("%s -> %s, \n", url, ctype)
			}
			wg.Done()
		}(url)
	}

	wg.Wait()
}

// sending data through a channel that has no listener will get the system blocked
// rather we use goroutine to send the data concurrently
// example
func sendViaGoroutine(data int) {
	// create the channel
	ch := make(chan int) // we created a channel of type int
	// ch <- data // sending from here will block

	// send from goroutine
	go func() {
		// print sending data
		fmt.Printf("sending %d\n", data)
		// send data to the channel
		ch <- data
	}()

	// receive data from the channel
	val := <-ch

	// print received data
	fmt.Printf("got %d\n", val)
}

// Gaurding against panic
func safevalue(vals []int, index int) (n int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	return vals[index], nil
}

// Opening files safely in go
func killServer(pidFilePath string) error {
	file, err := os.Open(pidFilePath)
	if err != nil {
		return err
	}

	defer file.Close()

	var pid int
	if _, err := fmt.Fscanf(file, "%d", &pid); err != nil {
		return pkgerrors.Wrap(err, "bad proccess ID")
	}

	// Simulate killing the server
	fmt.Println("Killing the server")

	if err := os.Remove(pidFilePath); err != nil {
		// we can go on if we fail here
		log.Printf("warning: can't remove pid file - %s", err)
	}

	return nil
}

// Algorithm
// ==========================================================================
var (
	urlTemplate = "https://s3.amazonaws.com/nyc-tlc/trip+data/%s_tripdata_2020-%02d.csv"
	colors      = []string{"green", "yellow"}
)

// Calls an http head request for the provided url and returns the content size
func downloadSize(url string) (int, error) {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status %s : url %s", resp.Status, url)
	}

	// get the content length from the header
	return strconv.Atoi(resp.Header.Get("Content-Length"))
}

// to make the Algorithm perform better
// we will add a worker
type result struct {
	url  string
	size int
	err  error
}

func sizeWorker(url string, ch chan result) {
	fmt.Println(url)
	res := result{url: url}
	res.size, res.err = downloadSize(url)
	ch <- res
}

//=======================================================================

func main() {

	//Concurrent Algorithm
	//====================================================
	start := time.Now()
	size := 0
	// creat channel
	ch := make(chan result)
	for month := 1; month <= 12; month++ {
		for _, color := range colors {
			url := fmt.Sprintf(urlTemplate, color, month)
			go sizeWorker(url, ch)
		}
	}

	// lets collect the result
	for i := 0; i < len(colors)*12; i++ {
		// we get the result from channel
		r := <-ch
		if r.err != nil {
			// we exit the program by logging fatal
			log.Fatal(r.err)
		}
		// we increase the size here
		size += r.size
	}

	duration := time.Since(start)
	fmt.Println(size, duration)
	//====================================================

	// Goroutines
	//====================================================
	datas := []int{101, 201, 301, 401, 501}
	for index := range datas {
		sendViaGoroutine(datas[index])
	}
	//====================================================

	urls := []string{
		"https://golang.org",
		"https://api.github.com",
		"https://httpbin.org/ip",
	}

	start = time.Now()
	siteSerial(urls)
	fmt.Println(time.Since(start))

	start = time.Now()
	siteSerialConcurrent(urls)
	fmt.Println(time.Since(start))

	//=========================================
	if err := killServer("server.pid"); err != nil {
		fmt.Fprintf(os.Stderr, "errors: %s\n", err)
		os.Exit(1)
	}
	//==========================================

	vals := []int{1, 2, 3, 4, 5}
	v, err := safevalue(vals, 10)

	if err != nil {
		fmt.Printf("Error: %v, \n", err)
	} else {
		fmt.Printf("Value: %v, \n", v)
	}
	// min value of a slice
	//=================================================
	fmt.Println(min([]float64{10, 1, 8, 4}))
	fmt.Println(min([]string{"H", "C", "A", "Z"}))
	//=================================================

	// convert to uppercase
	//==================================================

	c := &Capper{os.Stdout}
	fmt.Fprintln(c, "Helo world")

	//==================================================

	s, err := NewSqure(1, 1, 10)

	if err != nil {
		log.Fatalf("ERROR: an error occured while creating squre")
	}

	// move squre
	s.Move(4, 5)
	// use %#v to print VALUE and TYPE in one go
	fmt.Printf("%#v", s)
	fmt.Println(s.Area())

	fmt.Println(contentType("https://linkedin.com"))

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

	// Slices
	nameSlice := []string{"John", "Paul", "James"}

	nameSlice = append(nameSlice, "Kalu")

	// Say helo
	for _, name := range nameSlice {
		fmt.Printf("helo %v,\n", name)
	}

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
}
