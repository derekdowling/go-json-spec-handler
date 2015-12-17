package jshapi

// Relationship helps define the relationship between two resources
type Relationship string

const (
	// ToOne signifies a one to one relationship
	ToOne Relationship = "One-To-One"
	// ToMany signifies a one to many relationship
	ToMany Relationship = "One-To-Many"
)
