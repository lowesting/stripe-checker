package src

// struct of a card
type CardStruct struct {
	Number   string
	ExpMonth string
	ExpYear  string
	Cvc      string
}

// callback to function @OpenCardList
type OpenCardListCallback func(rawCard string)

// struct of a config
type Config struct {
	PublishKey string `yaml:"stripePublishKey"`
	PrivateKey string `yaml:"stripePrivateKey"`
	Ammount    string `yaml:"amount"`
	Currency   string `yaml:"currency"`
}

// struct of checker result
type CheckerResult struct {
	Code           string
	DeclinedReason string
	Valid          bool
	cc             CardStruct
}
