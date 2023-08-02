package models

type Rule struct {
	Name       string   `bson:"name" json:"name"`
	Conditions []string `bson:"condition" json:"condition"`
	Actions    []string `bson:"actions" json:"actions"`
}

type RuleSet struct {
	RuleSetName string `bson:"rulesetname" json:"rulesetname"`
	Rules       []Rule `bson:"rules" json:"rules"`
}
