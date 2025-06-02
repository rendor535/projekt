package repositories

import (
	"backend/internal/domain/models"
	"database/sql"
	"fmt"
	"strings"
)

type DietPlanRepository struct {
	DB *sql.DB
}

func NewDietPlanRepository(db *sql.DB) *DietPlanRepository {
	return &DietPlanRepository{DB: db}
}

// ===================== TRAINER =====================

func (r *DietPlanRepository) GetTrainerDietPlans(trainerID int, clientID *int, isActive *bool) ([]models.DietPlanWithClient, error) {
	args := []interface{}{trainerID}
	conditions := []string{"dp.trainer_id = $1"}

	if clientID != nil {
		args = append(args, *clientID)
		conditions = append(conditions, fmt.Sprintf("dp.client_id = $%d", len(args)))
	}
	if isActive != nil {
		args = append(args, *isActive)
		conditions = append(conditions, fmt.Sprintf("dp.is_active = $%d", len(args)))
	}

	query := `
	SELECT 
		dp.id, dp.client_id, dp.trainer_id, dp.name, dp.description, dp.calories,
		COALESCE(dp.meals, '') as meals, dp.start_date, dp.end_date, dp.is_active, dp.created_at,
		u.username as client_username,
		COUNT(m.id) as meals_count
	FROM diet_plans dp
	JOIN users u ON dp.client_id = u.id
	LEFT JOIN meals m ON dp.id = m.diet_plan_id
	WHERE ` + strings.Join(conditions, " AND ") + `
	GROUP BY dp.id, u.username
	ORDER BY dp.created_at DESC
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.DietPlanWithClient
	for rows.Next() {
		var plan models.DietPlanWithClient
		err := rows.Scan(
			&plan.ID, &plan.ClientID, &plan.TrainerID, &plan.Name, &plan.Description,
			&plan.Calories, &plan.Meals, &plan.StartDate, &plan.EndDate, &plan.IsActive,
			&plan.CreatedAt, &plan.ClientUsername,
			&plan.MealsCount,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// ===================== TRAINEE =====================

func (r *DietPlanRepository) GetTraineeDietPlans(clientID int, isActive *bool) ([]models.DietPlanWithTrainer, error) {
	args := []interface{}{clientID}
	conditions := []string{"dp.client_id = $1"}

	if isActive != nil {
		args = append(args, *isActive)
		conditions = append(conditions, fmt.Sprintf("dp.is_active = $%d", len(args)))
	}

	query := `
	SELECT 
		dp.id, dp.client_id, dp.trainer_id, dp.name, dp.description, dp.calories,
		COALESCE(dp.meals, '') as meals, dp.start_date, dp.end_date, dp.is_active, dp.created_at,
		u.username as trainer_username,
		COALESCE(tp.first_name, '') || ' ' || COALESCE(tp.last_name, '') as trainer_full_name,
		COUNT(m.id) as meals_count
	FROM diet_plans dp
	JOIN users u ON dp.trainer_id = u.id
	LEFT JOIN trainer_profiles tp ON u.id = tp.user_id
	LEFT JOIN meals m ON dp.id = m.diet_plan_id
	WHERE ` + strings.Join(conditions, " AND ") + `
	GROUP BY dp.id, u.username, tp.first_name, tp.last_name
	ORDER BY dp.created_at DESC
	`

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.DietPlanWithTrainer
	for rows.Next() {
		var plan models.DietPlanWithTrainer
		err := rows.Scan(
			&plan.ID, &plan.ClientID, &plan.TrainerID, &plan.Name, &plan.Description,
			&plan.Calories, &plan.Meals, &plan.StartDate, &plan.EndDate, &plan.IsActive,
			&plan.CreatedAt, &plan.TrainerUsername, &plan.TrainerFullName, &plan.MealsCount,
		)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

// ===================== CREATE PLAN =====================

func (r *DietPlanRepository) CreateWithMeals(req *models.CreateDietPlanRequest, trainerID int) (*models.DietPlan, error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var plan models.DietPlan
	query := `
		INSERT INTO diet_plans (client_id, trainer_id, name, description, start_date, end_date, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, false, NOW()) 
		RETURNING id, created_at
	`
	fmt.Println(">>>> DEBUG: trainerID =", trainerID)

	err = tx.QueryRow(query, req.ClientID, trainerID, req.Name, req.Description, req.StartDate, req.EndDate).
		Scan(&plan.ID, &plan.CreatedAt)
	if err != nil {
		return nil, err
	}

	plan.ClientID = req.ClientID
	plan.TrainerID = trainerID
	plan.Name = req.Name
	plan.Description = req.Description
	plan.StartDate = req.StartDate
	plan.EndDate = req.EndDate
	plan.IsActive = false

	totalCalories := 0
	for _, mealReq := range req.Meals {
		var mealID int
		err = tx.QueryRow(`INSERT INTO meals (diet_plan_id, meal_order, name, calories, proteins, fats, carbs)
			VALUES ($1, $2, $3, 0, 0, 0, 0) RETURNING id`,
			plan.ID, mealReq.Order, mealReq.Name).Scan(&mealID)
		if err != nil {
			return nil, err
		}

		mealCalories := 0
		mealProteins := 0
		mealFats := 0
		mealCarbs := 0

		for _, productReq := range mealReq.Products {
			var product models.FoodProduct
			err = tx.QueryRow(`SELECT id, name, calories, protein, fat, carbs FROM food_products WHERE id = $1`,
				productReq.ProductID).Scan(&product.ID, &product.Name, &product.Calories, &product.Protein, &product.Fat, &product.Carbs)
			if err != nil {
				return nil, fmt.Errorf("product not found: %v", err)
			}

			factor := productReq.AmountGrams / 100.0
			mealCalories += int(product.Calories * factor)
			mealProteins += int(product.Protein * factor)
			mealFats += int(product.Fat * factor)
			mealCarbs += int(product.Carbs * factor)

			_, err = tx.Exec(`INSERT INTO meal_products (meal_id, product_id, amount_grams) VALUES ($1, $2, $3)`,
				mealID, productReq.ProductID, productReq.AmountGrams)
			if err != nil {
				return nil, err
			}
		}

		_, err = tx.Exec(`UPDATE meals SET calories = $1, proteins = $2, fats = $3, carbs = $4 WHERE id = $5`,
			mealCalories, mealProteins, mealFats, mealCarbs, mealID)
		if err != nil {
			return nil, err
		}

		totalCalories += mealCalories
	}

	_, err = tx.Exec("UPDATE diet_plans SET calories = $1 WHERE id = $2", totalCalories, plan.ID)
	if err != nil {
		return nil, err
	}
	plan.Calories = totalCalories

	return &plan, tx.Commit()
}

// ===================== OTHERS =====================

func (r *DietPlanRepository) GetDietPlanMeals(dietPlanID int) ([]models.MealWithProducts, error) {
	query := `
		SELECT id, diet_plan_id, meal_order, name, calories, proteins, fats, carbs
		FROM meals 
		WHERE diet_plan_id = $1
		ORDER BY meal_order
	`

	rows, err := r.DB.Query(query, dietPlanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []models.MealWithProducts
	for rows.Next() {
		var meal models.MealWithProducts
		if err := rows.Scan(
			&meal.ID, &meal.DietPlanID, &meal.MealOrder, &meal.Name,
			&meal.Calories, &meal.Proteins, &meal.Fats, &meal.Carbs); err != nil {
			return nil, err
		}

		productRows, err := r.DB.Query(`
			SELECT mp.id, mp.meal_id, mp.product_id, mp.amount_grams,
			       fp.name, fp.calories, fp.protein, fp.fat, fp.carbs
			FROM meal_products mp
			JOIN food_products fp ON mp.product_id = fp.id
			WHERE mp.meal_id = $1`, meal.ID)
		if err != nil {
			return nil, err
		}

		var products []models.MealProductDetail
		for productRows.Next() {
			var product models.MealProductDetail
			if err := productRows.Scan(&product.ID, &product.MealID, &product.ProductID, &product.AmountGrams,
				&product.ProductName, &product.CaloriesPer100g, &product.ProteinPer100g,
				&product.FatPer100g, &product.CarbsPer100g); err != nil {
				productRows.Close()
				return nil, err
			}

			factor := product.AmountGrams / 100.0
			product.TotalCalories = product.CaloriesPer100g * factor
			product.TotalProtein = product.ProteinPer100g * factor
			product.TotalFat = product.FatPer100g * factor
			product.TotalCarbs = product.CarbsPer100g * factor

			products = append(products, product)
		}
		productRows.Close()

		meal.Products = products
		meals = append(meals, meal)
	}
	return meals, nil
}

func (r *DietPlanRepository) Update(planID int, req *models.UpdateDietPlanRequest, userID int, userRole string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Uaktualnij dane planu
	var query string
	if userRole == "trainer" {
		query = "UPDATE diet_plans SET name = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5 AND trainer_id = $6"
	} else {
		query = "UPDATE diet_plans SET name = $1, description = $2, start_date = $3, end_date = $4 WHERE id = $5 AND client_id = $6"
	}

	_, err = tx.Exec(query, req.Name, req.Description, req.StartDate, req.EndDate, planID, userID)
	if err != nil {
		return err
	}

	// Usuń stare posiłki (kaskadowo usunie też meal_products)
	_, err = tx.Exec("DELETE FROM meals WHERE diet_plan_id = $1", planID)
	if err != nil {
		return err
	}

	// Dodaj nowe posiłki
	totalCalories := 0
	for _, mealReq := range req.Meals {
		var mealID int
		err = tx.QueryRow(`
			INSERT INTO meals (diet_plan_id, meal_order, name, calories, proteins, fats, carbs)
			VALUES ($1, $2, $3, 0, 0, 0, 0) RETURNING id
		`, planID, mealReq.Order, mealReq.Name).Scan(&mealID)
		if err != nil {
			return err
		}

		mealCalories := 0
		mealProteins := 0
		mealFats := 0
		mealCarbs := 0

		for _, productReq := range mealReq.Products {
			var product models.FoodProduct
			err = tx.QueryRow("SELECT calories, protein, fat, carbs FROM food_products WHERE id = $1", productReq.ProductID).
				Scan(&product.Calories, &product.Protein, &product.Fat, &product.Carbs)
			if err != nil {
				return fmt.Errorf("product not found: %v", err)
			}

			factor := productReq.AmountGrams / 100.0
			mealCalories += int(product.Calories * factor)
			mealProteins += int(product.Protein * factor)
			mealFats += int(product.Fat * factor)
			mealCarbs += int(product.Carbs * factor)

			_, err = tx.Exec(`
				INSERT INTO meal_products (meal_id, product_id, amount_grams)
				VALUES ($1, $2, $3)
			`, mealID, productReq.ProductID, productReq.AmountGrams)
			if err != nil {
				return err
			}
		}

		_, err = tx.Exec(`
			UPDATE meals SET calories = $1, proteins = $2, fats = $3, carbs = $4
			WHERE id = $5
		`, mealCalories, mealProteins, mealFats, mealCarbs, mealID)
		if err != nil {
			return err
		}

		totalCalories += mealCalories
	}

	// Uaktualnij kaloryczność całego planu
	_, err = tx.Exec("UPDATE diet_plans SET calories = $1 WHERE id = $2", totalCalories, planID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *DietPlanRepository) Delete(planID int, trainerID int) error {
	_, err := r.DB.Exec("DELETE FROM diet_plans WHERE id = $1 AND trainer_id = $2", planID, trainerID)
	return err
}

func (r *DietPlanRepository) Create(plan *models.DietPlan) error {
	query := `
		INSERT INTO diet_plans (client_id, trainer_id, name, description, calories, meals, start_date, end_date, is_active, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW()) RETURNING id, created_at
	`
	return r.DB.QueryRow(query, plan.ClientID, plan.TrainerID, plan.Name, plan.Description,
		plan.Calories, plan.Meals, plan.StartDate, plan.EndDate, plan.IsActive).Scan(&plan.ID, &plan.CreatedAt)
}

func (r *DietPlanRepository) GetByClientID(clientID int) ([]models.DietPlanWithTrainer, error) {
	query := `
		SELECT dp.id, dp.client_id, dp.trainer_id, dp.name, dp.description, dp.calories,
		      COALESCE(dp.meals, '') as meals, dp.start_date, dp.end_date, dp.is_active, dp.created_at, u.username
		FROM diet_plans dp
		JOIN users u ON dp.trainer_id = u.id
		WHERE dp.client_id = $1
		ORDER BY dp.created_at DESC
	`
	rows, err := r.DB.Query(query, clientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []models.DietPlanWithTrainer
	for rows.Next() {
		var plan models.DietPlanWithTrainer
		err := rows.Scan(&plan.ID, &plan.ClientID, &plan.TrainerID, &plan.Name, &plan.Description,
			&plan.Calories, &plan.Meals, &plan.StartDate, &plan.EndDate, &plan.IsActive, &plan.CreatedAt, &plan.TrainerUsername)
		if err != nil {
			return nil, err
		}
		plans = append(plans, plan)
	}
	return plans, nil
}

func (r *DietPlanRepository) GetActivePlan(clientID int) (*models.DietPlanWithTrainer, error) {
	query := `
		SELECT dp.id, dp.client_id, dp.trainer_id, dp.name, dp.description, dp.calories,
		      COALESCE(dp.meals, '') as meals, dp.start_date, dp.end_date, dp.is_active, dp.created_at, u.username
		FROM diet_plans dp
		JOIN users u ON dp.trainer_id = u.id
		WHERE dp.client_id = $1 AND dp.is_active = true
		ORDER BY dp.created_at DESC
		LIMIT 1
	`
	var plan models.DietPlanWithTrainer
	err := r.DB.QueryRow(query, clientID).Scan(&plan.ID, &plan.ClientID, &plan.TrainerID, &plan.Name,
		&plan.Description, &plan.Calories, &plan.Meals, &plan.StartDate, &plan.EndDate, &plan.IsActive, &plan.CreatedAt, &plan.TrainerUsername)
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *DietPlanRepository) ActivatePlan(planID int, clientID int) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`UPDATE diet_plans SET is_active = false WHERE client_id = $1`, clientID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`UPDATE diet_plans SET is_active = true WHERE id = $1 AND client_id = $2`, planID, clientID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
