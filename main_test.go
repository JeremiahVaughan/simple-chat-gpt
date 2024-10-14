package main

import "testing"

func Test_isEven(t *testing.T) {

	t.Run("low numbers 0", func(t *testing.T) {
		with := 0
		got := isEven(with)
		if !got {
			t.Errorf("expected even for %d, but it was odd", with)
		}
	})

	t.Run("low numbers 1", func(t *testing.T) {
		with := 1
		got := isEven(with)
		if got {
			t.Errorf("expected got for %d, but it was even", with)
		}
	})

	t.Run("low numbers 2", func(t *testing.T) {
		with := 2
		got := isEven(with)
		if !got {
			t.Errorf("expected even for %d, but it was odd", with)
		}
	})

	t.Run("low numbers 3", func(t *testing.T) {
		with := 3
		got := isEven(with)
		if got {
			t.Errorf("expected odd for %d, but it was even", with)
		}
	})

}
