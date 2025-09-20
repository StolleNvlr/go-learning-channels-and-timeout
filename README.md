# Go Project to Learn Channels, Structs and Timeout management
Practice project using math, channels, and timeouts in Go to calculate the area and perimeter of a rectangle. Includes functions to analyze execution time and verify if the runtime exceeds the response time limit.

## Importing the necessary modules
```go
package main

import (
	"context"
	"fmt"
	"time"
)
```

## Defining the rectangle struct

```go
type rect struct {
	largura, altura int
}
```

## Defining the Result struct

```go
type Result struct {
	Area      int
	Perimetro int
	Error     error
}
```

## Creating the functions to calculate the area and perimeter
### -> Area:
```go
func (r *rect) area() int {
	
	return r.largura * r.altura
}
```
### -> Perimeter
```go
func (r *rect) perimetro() int {
	l, a := 2*r.largura, 2*r.altura
	fmt.Println("Largura: ", l, "\nAltura: ", a)
	return l + a
}
```

## Create a methods() function to group the math calculations
```go
func (r rect) metodos() (int,int, error){
	area := r.area()
	perimetro := r.perimetro()
	time.Sleep(1 * time.Second)
	return area, perimetro, nil
}
```

## Create the receptor() function to receive the values and return the results via channel
### This function returns values as follows: if successful, (results, nil); if the timeout is exceeded, (0, context.DeadlineExceeded).
#### The argument context.DeadlineExceeded is used when the response time exceeds the predefined limit.

```go
func (r rect) receptor(largura, altura int) (int, int, error) {
	resCh := make(chan Result, 1) // canal de resultados, buffer = 1 (área e perímetro)
	timeout := 3 * time.Second
	// "Requisição": executar area() com os valores desejados em uma goroutine
	go func() {
		rr := rect{largura: largura, altura: altura}
		area, perimetro, err := rr.metodos()
		resCh <- Result{Area: area, Perimetro: perimetro, Error: err}
	}()

	select {
	case res := <-resCh:
		return res.Area, res.Perimetro, nil
	case <-time.After(timeout):
		return 0, 0, context.DeadlineExceeded
	}
}
```

## Main function to initialize variables, create the rectangle, and test timeout/math code
```go
func main() {
	start := time.Now()
	r := rect{largura: 2, altura: 4}
	resArea, resPerimetro, resErr := r.receptor(2, 4)
	duration := time.Since(start)
	fmt.Println("Area imediata (usa r atual):", r.area())

	// Demonstração do teste de canal com timeout
	// Caso 1: timeout limite fixo de (2s) -> area() recebe timeout de 1s, então dorme 1s, deve concluir com sucesso (2*4=8)
	if resErr != nil {
		fmt.Printf("receptor (%.2fs): timeout -> %v\n", duration.Seconds(), resErr)
	} else {
		fmt.Printf("receptor (%.2fs): area = %v\nperimetro = %v\n", duration.Seconds(), resArea, resPerimetro)
	}
}
```
