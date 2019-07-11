package caasiu

type Rules struct {
	QueryString Rule
	Body        Rule
}

type Rule = map[string][]string
