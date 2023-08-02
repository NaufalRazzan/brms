package controllers

import (
	"brms/pkg/db"
	"brms/services/rules_management/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func InsertSpecificRule(ruleSetName string, insertedRule models.Rule) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{"rulesetname": ruleSetName}

	var checkRuleSet models.RuleSet

	if err := collectionName.FindOne(ctx, filter).Decode(&checkRuleSet); err != nil{
		return err
	}

	newRuleSet := models.RuleSet{
		RuleSetName: ruleSetName,
		Rules: append(checkRuleSet.Rules, insertedRule),
	}

	if _, err := collectionName.ReplaceOne(ctx, filter, newRuleSet); err != nil{
		return err
	}

	return nil
}

func UpdateSpecificRule(ruleSetName string, rulename string, updatedRule models.Rule) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil{
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"rulesetname": ruleSetName,
		"rules.name": rulename,
	}

	var checkRuleSet models.RuleSet

	if err := collectionName.FindOne(ctx, bson.M{"rulesetname": ruleSetName}).Decode(&checkRuleSet); err != nil{
		return err
	}

	updateFilter := bson.M{
		"$set": bson.M{
			"rules.$.conditions": updatedRule.Conditions,
			"rules.$.actions": updatedRule.Actions,
		},
	} 

	if _, err := collectionName.UpdateOne(ctx, filter, updateFilter); err != nil{
		return err
	}

	return nil
}

func DeleteSpecificRule(ruleSetName, ruleName string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil{
		return err
	}
	defer client.Disconnect(ctx)

	var checkRuleSet models.RuleSet

	if err := collectionName.FindOne(ctx, bson.M{"rulesetname": ruleSetName}).Decode(&checkRuleSet); err != nil{
		return err
	}

	filter := bson.M{
		"rulesetname": ruleSetName,
	}

	updateFilter := bson.M{
		"$pull": bson.M{
			"rules": bson.M{
				"name": ruleName,
			},
		},
	}

	if _, err := collectionName.UpdateOne(ctx, filter, updateFilter); err != nil{
		return err
	}

	return nil
}