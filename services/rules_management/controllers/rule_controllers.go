package controllers

import (
	"brms/pkg/db"
	"brms/services/rules_management/models"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func InsertSpecificRule(ruleSetName string, insertedRule models.Rule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{"rulesetname": ruleSetName}

	var checkRuleSet models.RuleSet

	if err := collectionName.FindOne(ctx, filter).Decode(&checkRuleSet); err != nil {
		return err
	}

	for _, rule := range checkRuleSet.Rules {
		if rule.RuleID == insertedRule.RuleID {
			return fmt.Errorf("rule already exists")
		}
	}

	insertedRule.RuleID = uuid.NewString()

	newRuleSet := models.RuleSet{
		RuleSetName: ruleSetName,
		Rules:       append(checkRuleSet.Rules, insertedRule),
	}

	if _, err := collectionName.ReplaceOne(ctx, filter, newRuleSet); err != nil {
		return err
	}

	return nil
}

func UpdateSpecificRule(ruleSetName string, ruleID string, updatedRule models.Rule) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"rulesetname":  ruleSetName,
		"rules.ruleid": ruleID,
	}

	var checkRuleSet models.RuleSet

	if err := collectionName.FindOne(ctx, filter).Decode(&checkRuleSet); err != nil {
		return err
	}

	updateFilter := bson.M{
		"$set": bson.M{
			"rules.$.conditions": updatedRule.Conditions,
			"rules.$.action":     updatedRule.Action,
		},
	}

	if _, err := collectionName.UpdateOne(ctx, filter, updateFilter); err != nil {
		return err
	}

	return nil
}

func DeleteSpecificRule(ruleSetName, ruleID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	var checkRuleSet models.RuleSet

	filter := bson.M{
		"rulesetname": ruleSetName,
		"rules.ruleid": ruleID,
	}

	if err := collectionName.FindOne(ctx, filter).Decode(&checkRuleSet); err != nil {
		return err
	}

	updateFilter := bson.M{
		"$pull": bson.M{
			"rules": bson.M{
				"ruleid": ruleID,
			},
		},
	}

	if _, err := collectionName.UpdateOne(ctx, filter, updateFilter); err != nil {
		return err
	}

	return nil
}
