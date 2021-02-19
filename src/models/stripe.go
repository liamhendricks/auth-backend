package models

type StripeEvent string

const CheckoutComplete StripeEvent = "checkout.session.completed"
const PaymentSuccess StripeEvent = "payment_intent.succeeded"
const CustomerUpdated StripeEvent = "customer.updated"

func (s StripeEvent) String() string {
	return string(s)
}
