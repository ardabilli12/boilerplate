package util

import "strings"

func MaskCreditCardNumber(cardNumber string) string {
	if len(cardNumber) < 10 {
		return cardNumber
	}

	maskedLength := len(cardNumber) - 10
	maskedCard := cardNumber[:6] + strings.Repeat("*", maskedLength) + cardNumber[len(cardNumber)-4:]
	return maskedCard
}
