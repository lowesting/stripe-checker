# stripe-checker
Credit card checker using stripe payment gateway.

# building
to build it is very simple, you need [git](https://git-scm.com/) and [golang](https://go.dev/) installed to start
```bash
  # clone repo
  $ git clone https://github.com/41337/stripe-checker.git
  
  # enter the folder
  $ cd stripe-checker
  
  # do the build
  $ go build -o schecker cli.go
```
after that you will have the binary called schecker (stripe-checker) now just [use it](https://github.com/41337/stripe-checker/edit/main/README.md#how-to-use)

# how to use
to check only one card you can use
```bash
  # command once arrives only one card, you can abbreviate with: o, 0
  $ ./schecker once "5555555555555555|05|2025|555"
```
#### Once Mode Print
![](https://github.com/41337/stripe-checker/blob/assets/once.jpg)

to check card list
```bash
  # your card list should look like this:
  5555555555555555|05|2025|555
  5555555555555555|05|2025|555
  5555555555555555|05|2025|555
  
  # list command check multiple cards in a list, you can abbreviate it with: l
  $ ./schecker list my_cards.txt
```
#### List Mode Print
![](https://github.com/41337/stripe-checker/blob/assets/list.jpg)

# configuration file
your [configuration](https://github.com/41337/stripe-checker/blob/main/config.yaml) file should look like this:
```yaml
stripePublishKey: "pk_live..."
stripePrivateKey: "sk_live_..."
ammount: "...."
currency: "usd"
```
- note: in your file you can change keys, change the value that will be pulled in the check and change the currency

if you have more than one configuration file you can pass it using flags
```bash
  # this way you can use several configuration files just passing their name in the config path flag
  $ ./schecker —config-path myconfig.yaml list my_cards.txt
```
