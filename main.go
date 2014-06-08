// shida2 project main.go
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"
)

var (
	count int     = 0
	N     int     = 20
	xm    float64 = 0.0
	ym    float64 = 0.5
	h     float64 = 0.6
)

var (
	width    int    = 500
	height   int    = 500
	filename string = "shida.png"
)

var (
	bgcolor   color.Color = color.RGBA{255, 255, 255, 255}
	linecolor color.Color = color.RGBA{0, 128, 0, 255}
)

func W1x(x, y float64) float64 {
	return 0.836*x + 0.044*y
}

func W1y(x, y float64) float64 {
	return -0.044*x + 0.836*y + 0.169
}

func W2x(x, y float64) float64 {
	return -0.141*x + 0.302*y
}

func W2y(x, y float64) float64 {
	return 0.302*x + 0.141*y + 0.127
}

func W3x(x, y float64) float64 {
	return 0.141*x - 0.302*y
}

func W3y(x, y float64) float64 {
	return 0.302*x + 0.141*y + 0.169
}

func W4x(x, y float64) float64 {
	return 0
}

func W4y(x, y float64) float64 {
	return 0.175337 * y
}

func f(m *image.RGBA, k int, x, y float64) {
	var wg sync.WaitGroup

	defer wg.Wait()

	if 0 < k {
		wg.Add(1)
		go func() {
			f(m, k-1, W1x(x, y), W1y(x, y))
			wg.Done()
		}()
		wg.Add(1)
		go func() {
			if rand.Float64() < 0.3 {
				go f(m, k-1, W2x(x, y), W2y(x, y))
			}
			if rand.Float64() < 0.3 {
				go f(m, k-1, W3x(x, y), W3y(x, y))
			}
			if rand.Float64() < 0.3 {
				go f(m, k-1, W4x(x, y), W4y(x, y))
			}
			wg.Done()
		}()
	} else {
		var s float64 = 490.0
		m.Set(int(x*s+float64(width)*0.5), int(float64(height)-y*s), linecolor)
		count += 1
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	rand.Seed(time.Now().UTC().UnixNano())

	start := time.Now()

	m := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), &image.Uniform{bgcolor}, image.ZP, draw.Src)

	f(m, N, 0, 0)

	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		file.Close()
		fmt.Printf("Count = %d\n", count)
		fmt.Printf("Time=%f\n", time.Now().Sub(start).Seconds())
	}()
	err = png.Encode(file, m)
	if err != nil {
		log.Fatal(err)
	}
}
