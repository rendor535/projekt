package models

import (
	"time"
)

// Wykorzystuje istniejącą tabelę diet_plans
type DietPlan struct {
	ID          int        `json:"id"`
	ClientID    int        `json:"client_id"`
	TrainerID   int        `json:"trainer_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Calories    int        `json:"calories"`
	Meals       string     `json:"meals"` // TEXT field w bazie
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}
type DietPlanWithClient struct {
	DietPlan
	ClientUsername string `json:"client_username"`
	ClientFullName string `json:"client_full_name"`
	MealsCount     int    `json:"meals_count"`
}

type DietPlanWithTrainer struct {
	DietPlan
	TrainerUsername string `json:"trainer_username"`
	TrainerFullName string `json:"trainer_full_name"`
	MealsCount      int    `json:"meals_count"`
}

// Wykorzystuje istniejącą tabelę meals
type Meal struct {
	ID         int    `json:"id"`
	DietPlanID int    `json:"diet_plan_id"`
	MealOrder  int    `json:"meal_order"`
	Name       string `json:"name"`
	Calories   int    `json:"calories"`
	Proteins   int    `json:"proteins"`
	Fats       int    `json:"fats"`
	Carbs      int    `json:"carbs"`
}

// Wykorzystuje istniejącą tabelę meal_products
type MealProduct struct {
	ID          int     `json:"id"`
	MealID      int     `json:"meal_id"`
	ProductID   int     `json:"product_id"`
	AmountGrams float64 `json:"amount_grams"`
}

// Rozszerzony widok posiłku z produktami
type MealWithProducts struct {
	Meal
	Products []MealProductDetail `json:"products"`
}

type MealProductDetail struct {
	MealProduct
	ProductName     string  `json:"product_name"`
	CaloriesPer100g float64 `json:"calories_per_100g"`
	ProteinPer100g  float64 `json:"protein_per_100g"`
	FatPer100g      float64 `json:"fat_per_100g"`
	CarbsPer100g    float64 `json:"carbs_per_100g"`
	TotalCalories   float64 `json:"total_calories"`
	TotalProtein    float64 `json:"total_protein"`
	TotalFat        float64 `json:"total_fat"`
	TotalCarbs      float64 `json:"total_carbs"`
}

// Request DTOs
type CreateDietPlanRequest struct {
	ClientID    int                 `json:"client_id" binding:"required"`
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	StartDate   *time.Time          `json:"start_date"`
	EndDate     *time.Time          `json:"end_date"`
	Meals       []CreateMealRequest `json:"meals" binding:"required"`
}

type CreateMealRequest struct {
	Name     string                     `json:"name" binding:"required"`
	Order    int                        `json:"order" binding:"required"`
	Products []CreateMealProductRequest `json:"products" binding:"required"`
}

type CreateMealProductRequest struct {
	ProductID   int     `json:"product_id" binding:"required"`
	AmountGrams float64 `json:"amount_grams" binding:"required"`
}

type UpdateDietPlanRequest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	StartDate   *time.Time          `json:"start_date"`
	EndDate     *time.Time          `json:"end_date"`
	Meals       []CreateMealRequest `json:"meals"`
}
