package main

import (
	"fmt"
	"strings"
)

// type Employee struct {
// 	Name   string
// 	Salary float64
// }

// func (e *Employee) increaseSalary(percentage float64) {
// 	e.Salary += e.Salary * percentage / 100
// }

// func main() {
// 	emp := Employee{Name: "John", Salary: 5000}
// 	emp.increaseSalary(10)
// 	fmt.Println(emp.Salary)
// }

// func printNumbers() {
// 	for i := 1; i <= 5; i++ {
// 		fmt.Println(i)
// 		time.Sleep(100 * time.Millisecond)
// 	}
// }

// func main() {
// 	start := time.Now()
// 	go printNumbers()
// 	fmt.Println("selesai")

// 	elapsed := time.Since(start)
// 	fmt.Printf("Sum function took %s", elapsed)
// }

func main() {
	s := "Go is amazing"
	words := strings.Split(s, " ")
	fmt.Println(len(words))
}

// c,a,b,b,c,b,c,e,a,b,a,a,a, b, e
