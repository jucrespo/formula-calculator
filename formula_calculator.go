package formula_calculator

import (
	"fmt"
	"github.com/antonmedv/expr"
	"strings"
)

type FormulaGetter interface {
	GetFormula(code string) (string, error)
}

func CalculateFormula(formulaCode string, parameters interface{}, formulaGetter FormulaGetter) (interface{}, error) {
	codeWithoutBrackets := strings.Replace(formulaCode, "(", "", -1)
	codeWithoutBrackets = strings.Replace(codeWithoutBrackets, ")", "", -1)
	words := strings.Split(codeWithoutBrackets, " ")

	for _, word := range words {
		if len(word) >= 7 && word[:7] == "formula" {
			innerCode, err := formulaGetter.GetFormula(word[8:])
			if err != nil {
				return nil, err
			}

			result, err := CalculateFormula(innerCode, parameters, formulaGetter)
			if err != nil {
				return nil, err
			}

			formulaCode = strings.Replace(formulaCode, word, fmt.Sprintf("%v", result), -1)
		}
	}

	program, err := expr.Compile(formulaCode, expr.Env(parameters))
	if err != nil {
		return nil, err
	}

	out, err := expr.Run(program, parameters)
	if err != nil {
		return nil, err
	}

	return out, err
}
