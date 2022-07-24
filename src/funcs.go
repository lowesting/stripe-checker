package src

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Jeffail/gabs"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/refund"
	"gopkg.in/yaml.v2"
)

// card split
func CardSplit(rawCard, separator string) CardStruct {
	var cardStruct CardStruct
	Values := strings.Split(rawCard, separator)

	cardStruct.Number = Values[0]
	cardStruct.ExpMonth = Values[1]
	cardStruct.ExpYear = Values[2]
	cardStruct.Cvc = Values[3]

	return cardStruct
}

// open card list
func OpenCardList(filepath string, callback OpenCardListCallback) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for {
		if scanner.Scan() {
			callback(scanner.Text())
			continue
		}
		break
	}
	return nil
}

// load settings
func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func LoadSettings(configPath string) (Config, error) {
	var config Config
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if err := config.Parse(data); err != nil {
		return config, err
	}

	Logf("ammount: %s", 0, config.Ammount)
	Logf("currency: %s", 0, config.Currency)
	Logf("publish key: %s...", 0, config.PublishKey[0:15])
	Logf("private key: %s...", 0, config.PrivateKey[0:15])
	return config, nil
}

// check card
func CheckCard(cc CardStruct, config Config) (CheckerResult, error) {
	var result CheckerResult

	// sets
	result.cc = cc
	stripe.Key = config.PrivateKey

	// payment ident parameters
	amm, _ := strconv.ParseInt(config.Ammount, 10, 64)
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amm),
		Currency: stripe.String(config.Currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// create a payment indent
	pi, _ := paymentintent.New(params)

	// requests stripe api
	url := fmt.Sprintf("https://api.stripe.com/v1/payment_intents/%s/confirm", pi.ID)
	method := "POST"
	payload := strings.NewReader(fmt.Sprintf("return_url=https://stripe.com&payment_method_data[type]=card&payment_method_data[card][number]=%s&payment_method_data[card][cvc]=%s&payment_method_data[card][exp_year]=%s&payment_method_data[card][exp_month]=%s&payment_method_data[billing_details][address][country]=BR&payment_method_data[payment_user_agent]=stripe.js/eb14574ae;+stripe-js-v3/eb14574ae;+payment-element&payment_method_data[time_on_page]=145828&payment_method_data[guid]=a6cb16e6-1eee-420c-bdd7-49271c53ee9537ac30&payment_method_data[muid]=6a5417b9-177e-4293-af9b-200ef3fdac60ef8bc5&payment_method_data[sid]=c0d2f286-3f58-4e9b-95a5-d5bef895e91b377d41&expected_payment_method_type=card&use_stripe_sdk=true&key=%s&client_secret=%s", cc.Number, cc.Cvc, cc.ExpYear, cc.ExpMonth, config.PublishKey, pi.ClientSecret))

	client := &http.Client{}
	req, _ := http.NewRequest(method, url, payload)

	// set headers
	req.Header.Add("authority", "api.stripe.com")
	req.Header.Add("accept", "application/json")
	req.Header.Add("accept-language", "en-US,en;q=0.9,pt;q=0.8")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("origin", "https://js.stripe.com")
	req.Header.Add("referer", "https://js.stripe.com/")
	req.Header.Add("sec-ch-ua", "\" Not A;Brand\";v=\"99\", \"Chromium\";v=\"101\", \"Google Chrome\";v=\"101\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Windows\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-site")
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/101.0.4951.54 Safari/537.36")

	res, _ := client.Do(req)

	body, _ := ioutil.ReadAll(res.Body)

	// pass json response to gabs
	jsonParsed, _ := gabs.ParseJSON(body)

	// check
	if jsonParsed.ExistsP("error") {
		ErrorCode, ok := jsonParsed.Path("error.code").Data().(string)
		if ok {
			result.Code = ErrorCode
			if ErrorCode == "card_declined" {
				declinedReason, ok := jsonParsed.Path("error.decline_code").Data().(string)
				if ok {
					result.DeclinedReason = declinedReason
					switch declinedReason {
					case "currency_not_supported":
						result.Valid = true
					case "insufficient-funds":
						result.Valid = true
					case "amount-too-large":
						result.Valid = true
					case "balance-insufficient":
						result.Valid = true
					case "incorrect-cvc":
						result.Valid = true
					case "not_permitted":
						result.Valid = true
					}
				}
			}
		}
	} else {
		result.Valid = true

		// refund payment using payment ident id
		params := &stripe.RefundParams{
			PaymentIntent: stripe.String(pi.ID),
		}
		refund.New(params)
	}

	return result, nil
}

// process result
func ProcessResult(result CheckerResult) {
	if result.Valid {
		Logf("%s, %s/%s, %s", 1, result.cc.Number, result.cc.ExpMonth, result.cc.ExpYear, result.cc.Cvc)
	} else {
		Logf("%s, %s/%s, %s (Reason: %s)(Code: %s)", 2, result.cc.Number, result.cc.ExpMonth, result.cc.ExpYear, result.cc.Cvc, result.DeclinedReason, result.Code)
	}
}
