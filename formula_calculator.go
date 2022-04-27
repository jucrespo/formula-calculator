package formulacalculator

import (
	"fmt"
	"strings"

	"github.com/antonmedv/expr"
)

type FormulaGetter interface {
	GetFormula(code string) (string, error)
}

func CalculateFormula(formulaCode string, parameters interface{}, formulaGetter FormulaGetter) (interface{}, error) {
	codeWithoutBrackets := strings.ReplaceAll(formulaCode, "(", "")
	codeWithoutBrackets = strings.ReplaceAll(codeWithoutBrackets, ")", "")
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

			formulaCode = strings.ReplaceAll(formulaCode, word, fmt.Sprintf("%v", result))
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
