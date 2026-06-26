package main

import (
	"github.com/thelastvideostore/internal/models"
	"github.com/thelastvideostore/internal/store"
)

func seedSnackBar(s *store.Store) {
	items := []models.SnackBarItem{
		{ID: "snack-popcorn", Name: "Popcorn", Description: "Freshly popped, buttery classic.", Price: 3.99, Category: "snack", Stock: 50, Emoji: "🍿"},
		{ID: "snack-nachos", Name: "Nachos", Description: "Crispy tortilla chips with warm cheese sauce.", Price: 5.99, Category: "snack", Stock: 30, Emoji: "🧀"},
		{ID: "snack-hotdog", Name: "Hot Dog", Description: "All-beef frank on a toasted bun with condiments.", Price: 4.99, Category: "snack", Stock: 25, Emoji: "🌭"},
		{ID: "snack-pizza", Name: "Pizza Slice", Description: "Pepperoni pizza slice, hot and ready.", Price: 4.49, Category: "snack", Stock: 20, Emoji: "🍕"},
		{ID: "snack-pretzel", Name: "Soft Pretzel", Description: "Warm salted pretzel with cheese dip.", Price: 3.99, Category: "snack", Stock: 30, Emoji: "🥨"},
		{ID: "snack-fries", Name: "French Fries", Description: "Crispy golden fries with ketchup.", Price: 3.49, Category: "snack", Stock: 40, Emoji: "🍟"},
		{ID: "snack-burger", Name: "Cheeseburger", Description: "Quarter-pound patty with lettuce, tomato, and cheese.", Price: 6.99, Category: "snack", Stock: 15, Emoji: "🍔"},
		{ID: "snack-candy", Name: "Candy Assortment", Description: "Mixed box of theater candies — Sour Patch, M&Ms, Skittles.", Price: 2.99, Category: "candy", Stock: 60, Emoji: "🍬"},
		{ID: "snack-chocolate", Name: "Chocolate Bar", Description: "Large milk chocolate bar. Classic concession companion.", Price: 2.49, Category: "candy", Stock: 45, Emoji: "🍫"},
		{ID: "snack-icecream", Name: "Ice Cream Cup", Description: "Vanilla soft serve with your choice of topping.", Price: 3.49, Category: "candy", Stock: 35, Emoji: "🍦"},
		{ID: "snack-soda", Name: "Fountain Soda", Description: "Large 32oz soda — Coke, Sprite, Fanta, or Dr Pepper.", Price: 2.99, Category: "drink", Stock: 80, Emoji: "🥤"},
		{ID: "snack-water", Name: "Bottled Water", Description: "Pure spring water, 500ml.", Price: 1.49, Category: "drink", Stock: 60, Emoji: "💧"},
		{ID: "snack-slushie", Name: "Slushie", Description: "Ice-cold slushie — cherry, blue raspberry, or cola.", Price: 3.99, Category: "drink", Stock: 40, Emoji: "🧊"},
		{ID: "snack-coffee", Name: "Coffee", Description: "Fresh brewed hot coffee. Regular or decaf.", Price: 2.49, Category: "drink", Stock: 50, Emoji: "☕"},
		{ID: "snack-milkshake", Name: "Milkshake", Description: "Thick milkshake — chocolate, vanilla, or strawberry.", Price: 4.99, Category: "drink", Stock: 25, Emoji: "🥛"},
	}
	for i := range items {
		s.CreateSnackBarItem(&items[i])
	}
}

