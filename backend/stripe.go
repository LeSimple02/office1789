package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
	checkoutsession "github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/subscription"
	"github.com/stripe/stripe-go/v81/webhook"
)

// CreateCheckoutSessionRequest structure pour la requête
type CreateCheckoutSessionRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	PlanID   int    `json:"plan_id"` // 1=Standard, 2=Professional, 3=Enterprise
}

// getStripePriceID récupère le Price ID Stripe pour un plan donné
func getStripePriceID(planID int) string {
	priceIDs := map[int]string{
		1: os.Getenv("STRIPE_PRICE_STANDARD"),      // 5€/mois
		2: os.Getenv("STRIPE_PRICE_PROFESSIONAL"),  // 12€/mois
		3: os.Getenv("STRIPE_PRICE_ENTERPRISE"),    // 49€/mois
	}
	priceID := priceIDs[planID]
	fmt.Printf("🔍 DEBUG: Price ID récupéré pour plan %d: '%s'\n", planID, priceID)
	return priceID
}

// UpdateStripeSubscription met à jour un abonnement Stripe existant (upgrade/downgrade avec prorata)
func UpdateStripeSubscription(userID int, username string, newPlanID int, stripeCustomerID string, stripeSubscriptionID string) error {
	priceID := getStripePriceID(newPlanID)
	if priceID == "" {
		return fmt.Errorf("invalid price ID for plan %d", newPlanID)
	}

	fmt.Printf("🔄 Updating Stripe subscription for user %s (ID: %d)\n", username, userID)
	fmt.Printf("   Customer: %s, Subscription: %s, New Plan: %d\n", stripeCustomerID, stripeSubscriptionID, newPlanID)

	// Récupérer l'abonnement actuel
	sub, err := subscription.Get(stripeSubscriptionID, nil)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %v", err)
	}

	// Mettre à jour l'abonnement avec le nouveau prix
	// Stripe applique automatiquement le prorata
	updateParams := &stripe.SubscriptionParams{
		CancelAtPeriodEnd: stripe.Bool(false),
		ProrationBehavior: stripe.String("create_prorations"), // Active le prorata automatique
		Items: []*stripe.SubscriptionItemsParams{
			{
				ID:    stripe.String(sub.Items.Data[0].ID),
				Price: stripe.String(priceID),
			},
		},
	}

	updatedSub, err := subscription.Update(stripeSubscriptionID, updateParams)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %v", err)
	}

	fmt.Printf("✅ Subscription updated successfully: %s\n", updatedSub.ID)
	fmt.Printf("   New price: %s, Status: %s\n", priceID, updatedSub.Status)

	return nil
}

// CreateCheckoutSession crée une session de paiement Stripe
func CreateCheckoutSession(c *gin.Context) {
	var req CreateCheckoutSessionRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
		})
		return
	}

	// Valider la session utilisateur
	sess, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid session",
		})
		return
	}

	// Vérifier que le plan est valide (pas de paiement pour Free)
	if req.PlanID < 1 || req.PlanID > 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid plan. Only paid plans (1-3) require payment.",
		})
		return
	}

	// Vérifier si l'utilisateur a déjà un abonnement Stripe actif
	var stripeCustomerID, stripeSubscriptionID string
	err := db.QueryRow("SELECT stripe_customer_id, stripe_subscription_id FROM Users WHERE user_id=$1", sess.UserID).Scan(&stripeCustomerID, &stripeSubscriptionID)
	
	// Si l'utilisateur a déjà un abonnement Stripe actif, on le met à jour au lieu d'en créer un nouveau
	if err == nil && stripeCustomerID != "" && stripeSubscriptionID != "" {
		fmt.Printf("🔄 User %s already has a Stripe subscription, updating instead of creating new one\n", req.Username)
		
		// Mettre à jour l'abonnement Stripe existant avec prorata
		err := UpdateStripeSubscription(sess.UserID, req.Username, req.PlanID, stripeCustomerID, stripeSubscriptionID)
		if err != nil {
			fmt.Printf("❌ Error updating Stripe subscription: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update subscription: " + err.Error(),
			})
			return
		}
		
		// Mettre à jour le plan dans la base de données
		_, err = db.Exec("UPDATE Users SET nboffer=$1 WHERE user_id=$2", req.PlanID, sess.UserID)
		if err != nil {
			fmt.Printf("❌ Error updating plan in database: %v\n", err)
		}
		
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"url":     fmt.Sprintf("%s/account/subscription-success", os.Getenv("FRONTEND_URL")),
			"message": "Subscription updated successfully with prorated billing",
		})
		return
	}

	// Récupérer le Price ID Stripe
	priceID := getStripePriceID(req.PlanID)
	fmt.Printf("🔍 DEBUG: Test mode check - priceID='%s'\n", priceID)
	if priceID == "" || priceID == "price_test_standard" || priceID == "price_test_professional" || priceID == "price_test_enterprise" {
		// Mode test : utiliser l'endpoint de changement direct sans Stripe
		fmt.Printf("⚠️  Stripe not configured, using direct subscription change for user %s, plan %d\n", req.Username, req.PlanID)
		
		// Mettre à jour directement l'abonnement
		_, err := db.Exec("UPDATE Users SET nboffer=$1 WHERE user_id=$2", req.PlanID, sess.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update subscription",
			})
			return
		}
		
		// Rediriger vers la page de succès
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"url":     fmt.Sprintf("%s/account/subscription-success", os.Getenv("FRONTEND_URL")),
			"message": "Subscription updated (test mode - Stripe not configured)",
		})
		return
	}

	// Créer la session Stripe Checkout
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(fmt.Sprintf("%s/account/subscription-success?session_id={CHECKOUT_SESSION_ID}", 
			os.Getenv("FRONTEND_URL"))),
		CancelURL: stripe.String(fmt.Sprintf("%s/account/change", 
			os.Getenv("FRONTEND_URL"))),
		ClientReferenceID: stripe.String(fmt.Sprintf("%d:%s:%d", sess.UserID, req.Username, req.PlanID)),
		CustomerEmail: stripe.String(fmt.Sprintf("%s@office1789.com", req.Username)),
	}

	s, err := checkoutsession.New(params)
	if err != nil {
		fmt.Printf("Error creating Stripe session: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create payment session",
		})
		return
	}

	fmt.Printf("Stripe checkout session created for user %s (ID: %d), Plan: %d, Session: %s\n",
		req.Username, sess.UserID, req.PlanID, s.ID)

	c.JSON(http.StatusOK, gin.H{
		"success":    true,
		"session_id": s.ID,
		"url":        s.URL,
	})
}

// StripeWebhook gère les événements webhook de Stripe
func StripeWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	
	payload, err := c.GetRawData()
	if err != nil {
		fmt.Printf("Error reading webhook body: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	// Vérifier la signature du webhook
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	signatureHeader := c.GetHeader("Stripe-Signature")
	
	event, err := webhook.ConstructEvent(payload, signatureHeader, endpointSecret)
	if err != nil {
		fmt.Printf("Webhook signature verification failed: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		return
	}

	// Traiter les différents types d'événements
	switch event.Type {
	case "checkout.session.completed":
		// Paiement réussi
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			fmt.Printf("Error parsing webhook JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		// Extraire les informations depuis ClientReferenceID
		// Format: "userID:username:planID"
		var userID int
		var username string
		var planID int
		fmt.Sscanf(checkoutSession.ClientReferenceID, "%d:%[^:]:%d", &userID, &username, &planID)

		// Récupérer les IDs Stripe depuis la session
		stripeCustomerID := checkoutSession.Customer.ID
		stripeSubscriptionID := checkoutSession.Subscription.ID

		// Déterminer le account_type selon le plan
		accountType := "personal"
		if planID >= 2 { // Professional (2) ou Enterprise (3)
			accountType = "organization_owner"
		}

		// Mettre à jour l'abonnement de l'utilisateur avec les IDs Stripe
		_, err = db.Exec(
			"UPDATE users SET nboffer=$1, stripe_customer_id=$2, stripe_subscription_id=$3, account_type=$4 WHERE user_id=$5",
			planID, stripeCustomerID, stripeSubscriptionID, accountType, userID,
		)
		if err != nil {
			fmt.Printf("❌ Error updating subscription for user %d: %v\n", userID, err)
		} else {
			fmt.Printf("✅ Subscription activated: User %s (ID: %d) upgraded to plan %d\n", username, userID, planID)
			fmt.Printf("   Account type: %s\n", accountType)
			fmt.Printf("   Stripe Customer: %s, Subscription: %s, Session: %s\n",
				stripeCustomerID, stripeSubscriptionID, checkoutSession.ID)

			// Si plan Professional ou Enterprise, créer automatiquement l'organisation
			if planID >= 2 {
				var orgID int
				maxMembers := 3 // Professional = 3 membres
				if planID == 3 {
					maxMembers = 20 // Enterprise = 20 membres
				}

				// Vérifier si l'organisation existe déjà
				err = db.QueryRow("SELECT organization_id FROM organizations WHERE owner_user_id=$1", userID).Scan(&orgID)
				if err != nil {
					// Créer l'organisation
					orgName := username + " Organization"
					err = db.QueryRow(
						"INSERT INTO organizations (organization_name, owner_user_id, max_members, created_at) VALUES ($1, $2, $3, NOW()) RETURNING organization_id",
						orgName, userID, maxMembers,
					).Scan(&orgID)

					if err != nil {
						fmt.Printf("⚠️  Warning: Could not create organization for user %d: %v\n", userID, err)
					} else {
						// Lier l'utilisateur à son organisation
						_, err = db.Exec("UPDATE users SET organization_id=$1 WHERE user_id=$2", orgID, userID)
						if err != nil {
							fmt.Printf("⚠️  Warning: Could not link user to organization: %v\n", err)
						} else {
							fmt.Printf("✅ Organization created: %s (ID: %d, max_members: %d)\n", orgName, orgID, maxMembers)
						}
					}
				} else {
					// Mettre à jour max_members si le plan change
					_, err = db.Exec("UPDATE organizations SET max_members=$1 WHERE organization_id=$2", maxMembers, orgID)
					if err == nil {
						fmt.Printf("✅ Organization updated: max_members set to %d\n", maxMembers)
					}
				}
			}
		}

	case "customer.subscription.updated":
		// Abonnement mis à jour (renouvellement, changement de plan, etc.)
		fmt.Printf("Subscription updated event received\n")

	case "customer.subscription.deleted":
		// Abonnement annulé - rétrograder l'utilisateur
		var subscription stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &subscription)
		if err != nil {
			fmt.Printf("Error parsing webhook JSON: %v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event data"})
			return
		}

		// Trouver l'utilisateur par subscription_id
		var userID int
		var username string
		err = db.QueryRow(
			"SELECT user_id, username FROM users WHERE stripe_subscription_id=$1",
			subscription.ID,
		).Scan(&userID, &username)

		if err != nil {
			fmt.Printf("⚠️  Warning: Could not find user for subscription %s: %v\n", subscription.ID, err)
		} else {
			// Rétrograder vers le plan Free (0) et type personal
			_, err = db.Exec(
				"UPDATE users SET nboffer=0, account_type='personal', stripe_subscription_id=NULL WHERE user_id=$1",
				userID,
			)
			if err != nil {
				fmt.Printf("❌ Error downgrading user %d: %v\n", userID, err)
			} else {
				fmt.Printf("✅ Subscription cancelled: User %s (ID: %d) downgraded to Free plan\n", username, userID)
			}
		}

	case "invoice.payment_failed":
		// Échec de paiement
		fmt.Printf("Payment failed event received\n")
		// TODO: Envoyer notification à l'utilisateur

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"received": true})
}

// InitStripe initialise la clé API Stripe
func InitStripe() {
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if stripeKey == "" {
		fmt.Println("⚠️  WARNING: STRIPE_SECRET_KEY not set in environment")
	} else {
		stripe.Key = stripeKey
		fmt.Println("✅ Stripe initialized successfully")
	}
}
