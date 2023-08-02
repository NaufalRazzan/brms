package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/PaesslerAG/gval"
)

type Rule struct {
	Name       string
	Conditions []string
	Actions    []string
}

func Exec(customer map[string]interface{}, rules []Rule) error {
	totalDiscount := 0.0

	for _, rule := range rules {
		var parsedExpression []gval.Evaluable

		for _, condition := range rule.Conditions {
			expression, err := gval.Full().NewEvaluable(condition)
			if err != nil {
				return err
			}
			parsedExpression = append(parsedExpression, expression)
		}

		ruleResult := true
		for _, expression := range parsedExpression {
			result, err := expression.EvalBool(context.Background(), customer)
			if err != nil {
				return err
			}

			ruleResult = ruleResult && result
		}

		if ruleResult {
			for _, action := range rule.Actions {
				if action != ""{
					if err := PerformAction(action, &totalDiscount, customer); err != nil{
						return err
					}
				}
			}
		}
	}

	fmt.Printf("Total discount for customer: %.2f%%\n", totalDiscount)

	return nil
}

func PerformAction(action string, floatVariable *float64, data map[string]interface{}) error{
	actionParts := strings.Fields(action)
	if len(actionParts) < 3{
		return fmt.Errorf("invalid action expression: %s", action)
	}

	switch actionParts[0]{
	case "Apply":
		if actionParts[1] == "discount"{
			discount, err := gval.Evaluate(actionParts[2], data)
			if err != nil{
				return err
			}
			if discountValue, ok := discount.(float64); ok{
				*floatVariable += discountValue
			}
		}
	case "Print":
		fmt.Println(actionParts[1])
	default:
		return fmt.Errorf("invalid action keyword: %s", actionParts[0])
	}

	return nil
}

func InsertRule() []Rule{
	return []Rule{
		{
			Name:       "Rule1",
			Conditions: []string{"age > 30", `dept == "electronics"`},
			Actions:    []string{"totalDiscount = 0.1"},
		},
		{
			Name:       "Rule2",
			Conditions: []string{"age < 40", `dept == "clothing"`},
			Actions:    []string{"totalDiscount = 0.05"},
		},
	}
} 

func main(){
	customer := map[string]interface{}{
		"age": 35,
		"dept": "electronics",
	}

	if err := Exec(customer, InsertRule()); err != nil{
		fmt.Println("Error executing rules: ", err)
	}
}
 