# stripe-checker
Credit card checker using stripe payment gateway.

# how to use
to check only one card you can use
```bash
  # command once arrives only one card, you can abbreviate with: o, 0
  $ schecker once 5555555555555555|05|2025|555  
```
```bash
  # your card list should look like this:
  5555555555555555|05|2025|555
  5555555555555555|05|2025|555
  5555555555555555|05|2025|555
  
  # list command check multiple cards in a list, you can abbreviate it with: l
  $ schecker list my_cards.txt
```

# configuration file
your [configuration](https://github.com/J4c5/stripe-checker/blob/main/config.yaml) file should look like this:
```yaml
stripePublishKey: "pk_live..."
stripePrivateKey: "sk_live_..."
ammount: "...."
currency: "usd"
```
- note: in your file you can change keys, change the value that will be pulled in the check and change the currency
