package formulacalculator

import (
	"fmt"
	"math"
	"strings"

	"github.com/antonmedv/expr"
)

type FormulaGetter interface {
	GetFormula(code string) (string, error)
}

type FormulaCalculator struct {
	parameters map[string]interface{}
}

func NewFormulaCalculator() FormulaCalculator {
	f := FormulaCalculator{}
	f.parameters = make(map[string]interface{})
	return f
}

func (f FormulaCalculator) AddParameter(key string, value interface{}) {
	f.parameters[key] = value
}

func (f FormulaCalculator) AddParameters(parameters map[string]interface{}) {
	for key, value := range parameters {
		f.parameters[key] = value
	}
}

func (f FormulaCalculator) CalculateFormula(formulaCode string, formulaGetter FormulaGetter) (interface{}, error) {
	codeWithoutBrackets := strings.ReplaceAll(formulaCode, "(", "")
	codeWithoutBrackets = strings.ReplaceAll(codeWithoutBrackets, ")", "")
	words := strings.Split(codeWithoutBrackets, " ")

	for _, word := range words {
		if len(word) >= 7 && word[:7] == "formula" {
			innerCode, err := formulaGetter.GetFormula(word[8:])
			if err != nil {
				return nil, err
			}

			result, err := f.CalculateFormula(innerCode, formulaGetter)
			if err != nil {
				return nil, err
			}

			formulaCode = strings.ReplaceAll(formulaCode, word, fmt.Sprintf("%v", result))
		}
	}

	program, err := expr.Compile(formulaCode, expr.Env(f.parameters))
	if err != nil {
		return nil, err
	}

	out, err := expr.Run(program, f.parameters)
	if err != nil {
		return nil, err
	}

	return out, err
}

func RoundDown(decimals int) func(number float64) float64 {
	return func(number float64) float64 {
		multiplier := math.Pow10(decimals)
		return math.Floor(number*multiplier) / multiplier
	}
}

func RoundUp(decimals int) func(number float64) float64 {
	return func(number float64) float64 {
		multiplier := math.Pow10(decimals)
		return math.Ceil(number*multiplier) / multiplier
	}
}

func RoundNearest(decimals int) func(number float64) float64 {
	return func(number float64) float64 {
		multiplier := math.Pow10(decimals)
		return math.Round(number*multiplier) / multiplier
	}
}
