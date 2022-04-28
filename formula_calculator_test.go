package formulacalculator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _code = "(Output.Six - Input.Two) * Input.Inner.Three / formula.inner_formula"

type InnerStruct struct {
	Three float64
	Four  float64
}

type Input struct {
	One   float64
	Two   float64
	Inner InnerStruct
}

type Output struct {
	Five float64
	Six  float64
	Text string
}

var (
	_formulas map[string]string
	_input    Input
	_output   Output
	_env      map[string]interface{}
)

func setup() {
	_formulas = make(map[string]string)

	_input = Input{
		One: 1,
		Two: 2,
		Inner: InnerStruct{
			Three: 3,
			Four:  4,
		},
	}

	_output = Output{
		Five: 5,
		Six:  6,
		Text: "asd",
	}

	_env = map[string]interface{}{
		"Input":  &_input,
		"Output": &_output,
	}
}

type FormulaGetterMock struct{}

func (f FormulaGetterMock) GetFormula(key string) (string, error) {
	formula, ok := _formulas[key]
	if !ok {
		return "", fmt.Errorf("formula '%s' does not exist", key)
	}

	return formula, nil
}

func TestCalculateFormula_ResultNoDecimals(t *testing.T) {
	setup()

	_formulas["inner_formula"] = "Output.Five - Input.One"
	code := "(Output.Six - Input.Two) * Input.Inner.Three / (formula.inner_formula + Input.Two)"

	f := NewFormulaCalculator()
	f.AddParameters(_env)

	out, err := f.CalculateFormula(code, FormulaGetterMock{})
	assert.Nil(t, err)
	assert.Equal(t, float64(2), out.(float64))
}

func TestCalculateFormula_ResultWithDecimals(t *testing.T) {
	setup()

	_formulas["inner_formula"] = "Output.Six - Input.One"

	f := NewFormulaCalculator()
	f.AddParameters(_env)

	out, err := f.CalculateFormula(_code, FormulaGetterMock{})
	assert.Nil(t, err)
	assert.Equal(t, 2.4, out.(float64))
}

func TestCalculateFormula_InexistentInnerFormula(t *testing.T) {
	setup()

	f := NewFormulaCalculator()
	f.AddParameters(_env)

	_, err := f.CalculateFormula(_code, FormulaGetterMock{})
	assert.NotNil(t, err)
	assert.Equal(t, "formula 'inner_formula' does not exist", err.Error())
}

func TestCalculateFormula_InvalidFormula(t *testing.T) {
	setup()

	code := "name + age"

	f := NewFormulaCalculator()
	f.AddParameters(_env)

	_, err := f.CalculateFormula(code, FormulaGetterMock{})
	assert.NotNil(t, err)
}

func TestCalculateFormula_InvalidInnerFormula(t *testing.T) {
	setup()

	_formulas["inner_formula"] = "name + age"

	f := NewFormulaCalculator()
	f.AddParameters(_env)

	_, err := f.CalculateFormula(_code, FormulaGetterMock{})
	assert.NotNil(t, err)
}

func TestCalculateFormula_CustomFormulaReturnError(t *testing.T) {
	setup()

	code := "error()"

	env := map[string]interface{}{
		"error": func() (int, error) {
			return 0, fmt.Errorf("custom error")
		},
	}

	f := NewFormulaCalculator()
	f.AddParameters(env)

	_, err := f.CalculateFormula(code, FormulaGetterMock{})
	assert.NotNil(t, err)
	assert.Equal(t, "custom error", err.Error())
}

func TestCalculateFormula_CustomFunctionRoundDownTwoDecimals(t *testing.T) {
	setup()

	code := "round((Output.Six - Input.Two) * Input.Inner.Three / 3.4)"

	f := NewFormulaCalculator()
	f.AddParameters(_env)
	f.AddParameter("round", RoundDown(2))

	out, err := f.CalculateFormula(code, FormulaGetterMock{})
	assert.Nil(t, err)
	assert.Equal(t, 3.52, out.(float64))
}

func TestCalculateFormula_CustomFunctionRoundNearestZeroDecimals(t *testing.T) {
	setup()

	code := "round((Output.Six - Input.Two) * Input.Inner.Three / 3.4)"

	f := NewFormulaCalculator()
	f.AddParameters(_env)
	f.AddParameter("round", RoundNearest(0))

	out, err := f.CalculateFormula(code, FormulaGetterMock{})
	assert.Nil(t, err)
	assert.Equal(t, float64(4), out)
}
