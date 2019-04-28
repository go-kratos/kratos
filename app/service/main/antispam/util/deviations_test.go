package util

import "testing"

func TestMax(t *testing.T) {
	cases := []struct {
		input []int64
	}{
		{
			input: []int64{},
		},
		{
			input: []int64{
				268826797,
			},
		},
		{
			input: []int64{
				1,
				2,
				3,
				4,
				5,
				272668542,
				87759075,
				272670403,
				192148035,
				189361023,
				88269430,
				89400131,
				272690653,
				269278391,
				268823477,
				268826797,
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rs := Max(c.input)
			t.Logf("input %v \noutputs:%v\n", c.input, rs)
		})
	}
}

func TestStdDeviation(t *testing.T) {
	cases := []struct {
		input []float64
	}{
		{
			input: []float64{},
		},
		{
			input: []float64{
				268826797,
			},
		},
		{
			input: []float64{
				1,
				2,
				3,
				4,
				5,
				272668542,
				87759075,
				272670403,
				192148035,
				189361023,
				88269430,
				89400131,
				272690653,
				269278391,
				268823477,
				268826797,
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rs := StdDeviation(c.input)
			t.Logf("input %v \noutputs:%v\n", c.input, rs)
		})
	}
}

func TestDeviation(t *testing.T) {
	cases := []struct {
		input []float64
	}{
		{
			input: []float64{
				1,
				2,
				3,
				4,
				5,
				272668542,
				87759075,
				272670403,
				192148035,
				189361023,
				88269430,
				89400131,
				272690653,
				269278391,
				268823477,
				268826797,
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rs := Deviation(c.input)
			t.Logf("input %v \noutputs:%v\n", c.input, rs)
		})
	}
}

func TestNormallization(t *testing.T) {
	cases := []struct {
		input []int64
	}{
		{
			input: []int64{
				1,
				2,
				3,
				4,
				5,
				272668542,
				87759075,
				272670403,
				192148035,
				189361023,
				88269430,
				89400131,
				272690653,
				269278391,
				268823477,
				268826797,
			},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			rs := Normallization(c.input)
			t.Logf("input %v\n, outputs:%v\n", c.input, rs)
		})
	}
}

func TestExpectation(t *testing.T) {
	cases := []struct {
		input []float64
	}{
		{
			input: []float64{
				272668542,
				87759075,
				272670403,
				192148035,
				189361023,
				88269430,
				89400131,
				272690653,
				269278391,
				268823477,
				268826797,
				107338074,
				272635619,
				272655050,
				272280850,
				274403561,
				274404369,
				274408851,
				33561986,
				274440108,
				274440766,
				274440926,
				273730691,
				274443040,
				273744469,
				274443230,
				273749519,
				274837710,
				273748847,
				274837920,
				273814911,
				274838433,
				269278030,
				273815488,
				273819536,
				274838710,
				274838821,
				274875985,
				274876083,
			},
		},
		{
			input: []float64{
				23062012,
				32199188,
				114238752,
				35134612,
				18289010,
				29669239,
				25702393,
				19201815,
				38191035,
				24237961,
				10155123,
				36970516,
				14015382,
				89085165,
				28324458,
				38405237,
				20190219,
				175006499,
				178984876,
				34990873,
				158026517,
				111053972,
				25519948,
				39061494,
				155286833,
				87469728,
				62849283,
				210174070,
				13841839,
				65905090,
				8219588,
				37192235,
				274356431,
				44363565,
				22899631,
				43582749,
				101217680,
				25011431,
				33447081,
				222278335,
				20785287,
				38448378,
				40153047,
			},
		},
	}

	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			output := Expectation(c.input)
			t.Logf("input %v, output:%v", c.input, output)
		})
	}
}
