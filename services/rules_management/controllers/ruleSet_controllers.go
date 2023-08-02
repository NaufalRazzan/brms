package controllers

import (
	"brms/pkg/db"
	"brms/services/rules_management/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InsertRuleSet(ruleSet models.RuleSet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collectionName.Indexes().CreateOne(ctx, indexModel); err != nil {
		return err
	}

	ruleSetIndex := mongo.IndexModel{
		Keys: bson.M{"rulesetname": 1},
		Options: options.Index().SetUnique(true),
	}

	if _, err := collectionName.Indexes().CreateOne(ctx, ruleSetIndex); err != nil{
		return err
	}

	if _, err := collectionName.InsertOne(ctx, ruleSet); err != nil{
		return err
	}

	return nil
}

func DeleteRuleSet(ruleSetName string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil{
		return err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"rulesetname": ruleSetName,
	}

	if _, err := collectionName.DeleteOne(ctx, filter); err != nil{
		return err
	}

	return nil
}

func FindOneRuleSet(ruleSetName string) (*models.RuleSet, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil{
		return nil, err
	}
	defer client.Disconnect(ctx)

	filter := bson.M{
		"rulesetname": ruleSetName,
	}

	var result models.RuleSet

	if err := collectionName.FindOne(ctx, filter).Decode(&result); err != nil{
		return nil, err
	}

	return &result, nil
}

func FindAllRuleSet() ([]models.RuleSet, error){
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	client, collectionName, err := db.ConnectDB("rule_engine")
	if err != nil{
		return nil, err
	}
	defer client.Disconnect(ctx)

	var results []models.RuleSet

	cursor, err := collectionName.Find(ctx, bson.M{})
	if err != nil{
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &results); err != nil{
		return nil, err
	}

	return results, nil
}
