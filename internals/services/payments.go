package services

import (
	"errors"
	"fmt"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
	"gitlab.com/saas-template1/go-monolith/config"
	"log"
)

type PaymentsService struct {
	config config.PaymentsConfig
}

func NewPaymentService(conf config.PaymentsConfig) (*PaymentsService, error) {
	stripe.Key = conf.StripeKey

	return &PaymentsService{
		config: conf,
	}, nil
}

func (p *PaymentsService) Checkout(product string, billingCycle string) (string, error) {
	priceId := p.config.PlansIDs[product]
	if priceId == "" {
		return "", errors.New("current product is missing")
	}
	//successURL := "https://86ab-176-102-59-24.eu.ngrok.io/api/v1/success?session_id={CHECKOUT_SESSION_ID}"
	//cancelURL := "https://example.com/canceled.html"
	params := &stripe.CheckoutSessionParams{
		SuccessURL: &p.config.SuccessUrl,
		CancelURL:  &p.config.CancelUrl,
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceId),
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, _ := session.New(params)
	return s.URL, nil
}

func (p *PaymentsService) HandleWebhook(payload []byte, signatureHeader string) error {
	event, err := webhook.ConstructEvent(payload, signatureHeader, "whsec_FTJaRR4E9tvMVGJ6eHRv2Penh0BxjvsJ")
	if err != nil {
		log.Printf("webhook.ConstructEvent: %v", err)
		return err
	}

	switch event.Type {
	case "checkout.session.completed":
		fmt.Println(">>>>>> checkout.session.completed")
		// Payment is successful and the subscription is created.
		// You should provision the subscription and save the customer ID to your database.
	case "invoice.paid":
		fmt.Println(">>>>>> invoice.paid")
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
	case "invoice.payment_failed":
		fmt.Println(">>>>>> invoice.payment_failed")
		// The payment failed or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
	default:
		// unhandled event type
	}

	return nil
}

func (p *PaymentsService) Success(sessionId string) string {
	s, _ := session.Get(
		sessionId,
		nil,
	)

	fmt.Println(">>>>> Success")

	return string(s.PaymentStatus)
}
