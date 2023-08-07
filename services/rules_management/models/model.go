package models

type Rule struct {
	RuleID     string      `bson:"ruleid" json:"ruleid"`
	RuleType   string      `bson:"ruletype" json:"ruletype"`
	Conditions []string    `bson:"conditions" json:"conditions"`
	Action     interface{} `bson:"action" json:"action"`
}

type RuleSet struct {
	RuleSetName string `bson:"rulesetname" json:"rulesetname"`
	Rules       []Rule `bson:"rules" json:"rules"`
}
