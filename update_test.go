package main

import "testing"

func Test_applyWordWrap(t *testing.T) {
	t.Run("easy mode", func(t *testing.T) {
		words := `happy duck paste`
		width := 40
		got := applyWordWrap(words, width)
		if got != words {
			t.Errorf("error, expected %s, but got %s", words, got)
		}
	})

	t.Run("hard mode", func(t *testing.T) {
		words := `happy duck paste mustache pottery potluckz`
		expected := `happy duck paste
mustache pottery
potluckz`
		width := 20
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", expected, got)
		}
	})

	t.Run("harder mode", func(t *testing.T) {
		words := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
		width := 80
		got := applyWordWrap(words, width)
		if got != words {
			t.Errorf("error, expected: %s, but got: %s", words, got)
		}
	})

	t.Run("hardest mode", func(t *testing.T) {
		words := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World there cruel world, I want to lett you know that I don't like it!")
}
`
		expected := `package main

import "fmt"

func main() {
    fmt.Println("Hello, World there cruel world, I want to lett you know that I
don't like it!")
}
`
		width := 80
		got := applyWordWrap(words, width)
		if got != expected {
			t.Errorf("error, expected: %s, but got: %s", expected, got)
		}
	})

	t.Run("indented list", func(t *testing.T) {
		words := `
List:
    1. Eggs
    2. Bacon
`
		width := 80
		got := applyWordWrap(words, width)
		if got != words {
			t.Errorf("error, expected: %s, but got: %s", words, got)
		}
	})
}
