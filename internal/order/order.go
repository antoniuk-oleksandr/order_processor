package order

// Order represents a customer order with user and payment information.
type Order struct {
	// ID is the unique identifier for the order.
	ID int
	// UserID is the identifier of the user placing the order.
	UserID int
	// Amount is the order amount in the smallest currency unit (e.g., cents).
	Amount int
}
