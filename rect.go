package main

import (
	"context"
	"fmt"
	"time"
)

type rect struct {
	largura, altura int
}
type Result struct {
	Area      int
	Perimetro int
	Error     error
}
func (r *rect) area() int {
	
	return r.largura * r.altura
}
func (r *rect) perimetro() int {
	l, a := 2*r.largura, 2*r.altura
	fmt.Println("Largura: ", l, "\nAltura: ", a)
	return l + a
}

func (r rect) metodos() (int,int, error){
	area := r.area()
	perimetro := r.perimetro()
	time.Sleep(1 * time.Second)
	return area, perimetro, nil
}
// receptor executa o cálculo de área com valores fornecidos (largura/altura)
// e retorna via canal com suporte a timeout.
// Retorna (resultado, nil) no sucesso ou (0, context.DeadlineExceeded) em caso de timeout.
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
