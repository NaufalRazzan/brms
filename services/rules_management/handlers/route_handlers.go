package handlers

import (
	"brms/pkg/response"
	"brms/services/rules_management/controllers"
	"brms/services/rules_management/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routes(app *fiber.App) {
	ruleEngine := app.Group("/rulesEngine")

	// bagian ruleset
	ruleEngine.Post("/addRuleSet", InsertRuleSet)
	ruleEngine.Get("/listAllRuleSet", ListAllRuleSet)
	ruleEngine.Delete("/removeRuleSet", DeleteOneRuleSet)

	// bagian rule
	ruleEngine.Post("/addSpecificRule", InsertOneRule)
	ruleEngine.Patch("/updateSpecificRule", UpdateOneRule)
	ruleEngine.Delete("/removeSpecificRule", DeleteOneRule)
}

func InsertOneRule(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	ruleSetname := c.Query("ruleSetName")

	if ruleSetname == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty query field")
	}

	var insertedRule models.Rule

	if err := c.BodyParser(&insertedRule); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "The request entity contains invalid or missing data")
	}

	if insertedRule.Name == "" || len(insertedRule.Conditions) == 0 || len(insertedRule.Actions) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	if err := controllers.InsertSpecificRule(ruleSetname, insertedRule); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "rule set not found")
		}
		if err.Error() == "rule already exists" {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.StatusCreated("new rule appended", &fiber.Map{
		"details": fmt.Sprintf("rule appended for rule set %s", ruleSetname),
	}))
}

func UpdateOneRule(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPatch {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	ruleSetName := c.Query("ruleSetName")
	ruleName := c.Query("ruleName")

	if ruleSetName == "" || ruleName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty query fields")
	}

	var updatedRule models.Rule

	if err := c.BodyParser(&updatedRule); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "The request entity contains invalid or missing data")
	}

	if updatedRule.Name == "" || len(updatedRule.Conditions) == 0 || len(updatedRule.Actions) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	if err := controllers.UpdateSpecificRule(ruleSetName, ruleName, updatedRule); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "rule set or rule name does not exist")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.StatusOK("updated one rule", &fiber.Map{
		"details": fmt.Sprintf("rule %s has been updated on rule set %s", ruleName, ruleSetName),
	}))
}

func DeleteOneRule(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodDelete {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	ruleSetName := c.Query("ruleSetName")
	ruleName := c.Query("ruleName")

	if ruleSetName == "" || ruleName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "empty query fields")
	}

	if err := controllers.DeleteSpecificRule(ruleSetName, ruleName); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, "rule set or rule name does not exists")
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.StatusOK("one rule deleted", &fiber.Map{
		"details": fmt.Sprintf("rule %s has been deleted from rule set %s", ruleName, ruleSetName),
	}))
}

func InsertRuleSet(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodPost {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	var ruleSet models.RuleSet

	if err := c.BodyParser(&ruleSet); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "The request entity contains invalid or missing data")
	}

	if ruleSet.RuleSetName == "" || len(ruleSet.Rules) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "empty fields")
	}

	// insert semua ruleset
	if err := controllers.InsertRuleSet(ruleSet); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(response.StatusCreated("new rule set inserted", &fiber.Map{"details": "new rule set created"}))
}

func ListAllRuleSet(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodGet {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	// fetch all rule set
	listRuleSet, err := controllers.FindAllRuleSet()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(listRuleSet) == 0 {
		return c.Status(fiber.StatusOK).JSON(response.StatusOK("list empty", &fiber.Map{"details": "list empty"}))
	}

	return c.Status(fiber.StatusOK).JSON(response.StatusOK("list all rule sets", &fiber.Map{
		"details": listRuleSet,
	}))
}

func DeleteOneRuleSet(c *fiber.Ctx) error {
	if c.Method() != fiber.MethodDelete {
		return fiber.NewError(fiber.StatusMethodNotAllowed, "invalid http method")
	}

	ruleSetName := c.Query("ruleSetName")
	if ruleSetName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "invalid query parameter")
	}

	if err := controllers.DeleteRuleSet(ruleSetName); err != nil {
		if err == mongo.ErrNoDocuments {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response.StatusOK("Rule set deleted", &fiber.Map{
		"details": "One rule set has been deleted",
	}))
}
