package repositories

import (
	"database/sql"
	"fmt"
	"time"
)

type TrainingPlanRepository interface {
	GetTrainerCalendarEvents(trainerID int) ([]map[string]interface{}, error)
	GetTraineeCalendarEvents(traineeID int) ([]map[string]interface{}, error)
	GetTrainerPlans(trainerID int) ([]map[string]interface{}, error)
	GetTraineePlans(traineeID int) ([]map[string]interface{}, error)
	GetExercises() ([]map[string]interface{}, error)
	GetTrainerTrainees(trainerID int) ([]map[string]interface{}, error)
	CreateTrainingPlan(planData map[string]interface{}) (int, error)
	UpdateTrainingPlan(planID int, planData map[string]interface{}) error // 🆕 NOWA METODA
	DeleteTrainingPlan(planID, userID int) error
	GetPlanDetails(planID, userID int) (map[string]interface{}, error)
}

type trainingPlanRepository struct {
	db *sql.DB
}

func NewTrainingPlanRepository(db *sql.DB) TrainingPlanRepository {
	return &trainingPlanRepository{db: db}
}

func (r *trainingPlanRepository) GetTrainerCalendarEvents(trainerID int) ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetTrainerCalendarEvents for trainer %d\n", trainerID)

	query := `
		SELECT 
			tp.id,
			tp.name,
			tp.description,
			tp.start_date,
			tp.end_date,
			u.username as client_name,
			u.id as client_id,
			COUNT(DISTINCT tpd.id) as exercises_count
		FROM training_plans tp
		JOIN users u ON tp.client_id = u.id
		LEFT JOIN training_plan_days tpd ON tp.id = tpd.training_plan_id
		WHERE tp.trainer_id = $1
		GROUP BY tp.id, tp.name, tp.description, tp.start_date, tp.end_date, u.username, u.id
		ORDER BY tp.start_date
	`

	rows, err := r.db.Query(query, trainerID)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var id, clientID, exercisesCount int
		var name, description, clientName sql.NullString
		var startDate, endDate sql.NullTime

		err := rows.Scan(&id, &name, &description, &startDate, &endDate, &clientName, &clientID, &exercisesCount)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		// Konwertuj plan na wydarzenie kalendarza
		if startDate.Valid {
			event := map[string]interface{}{
				"id":    fmt.Sprintf("plan-%d", id),
				"title": name.String,
				"start": startDate.Time.Format("2006-01-02"),
				"end":   startDate.Time.AddDate(0, 0, 1).Format("2006-01-02"), // Domyślnie jeden dzień
				"extendedProps": map[string]interface{}{
					"planId":         id,
					"clientName":     clientName.String,
					"clientId":       clientID,
					"exercisesCount": exercisesCount,
					"description":    description.String,
				},
			}

			if endDate.Valid {
				event["end"] = endDate.Time.Format("2006-01-02")
			}

			events = append(events, event)
		}
	}

	fmt.Printf("✅ Repository: Found %d events for trainer %d\n", len(events), trainerID)
	return events, nil
}

func (r *trainingPlanRepository) GetTraineeCalendarEvents(traineeID int) ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetTraineeCalendarEvents for trainee %d\n", traineeID)

	query := `
		SELECT 
			tp.id,
			tp.name,
			tp.description,
			tp.start_date,
			tp.end_date,
			u.username as trainer_name,
			COUNT(DISTINCT tpd.id) as exercises_count
		FROM training_plans tp
		JOIN users u ON tp.trainer_id = u.id
		LEFT JOIN training_plan_days tpd ON tp.id = tpd.training_plan_id
		WHERE tp.client_id = $1
		GROUP BY tp.id, tp.name, tp.description, tp.start_date, tp.end_date, u.username
		ORDER BY tp.start_date
	`

	rows, err := r.db.Query(query, traineeID)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var events []map[string]interface{}
	for rows.Next() {
		var id, exercisesCount int
		var name, description, trainerName sql.NullString
		var startDate, endDate sql.NullTime

		err := rows.Scan(&id, &name, &description, &startDate, &endDate, &trainerName, &exercisesCount)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		if startDate.Valid {
			event := map[string]interface{}{
				"id":    fmt.Sprintf("plan-%d", id),
				"title": name.String,
				"start": startDate.Time.Format("2006-01-02"),
				"end":   startDate.Time.AddDate(0, 0, 1).Format("2006-01-02"),
				"extendedProps": map[string]interface{}{
					"planId":         id,
					"trainerName":    trainerName.String,
					"exercisesCount": exercisesCount,
					"description":    description.String,
				},
			}

			if endDate.Valid {
				event["end"] = endDate.Time.Format("2006-01-02")
			}

			events = append(events, event)
		}
	}

	fmt.Printf("✅ Repository: Found %d events for trainee %d\n", len(events), traineeID)
	return events, nil
}

func (r *trainingPlanRepository) GetTrainerPlans(trainerID int) ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetTrainerPlans for trainer %d\n", trainerID)

	query := `
		SELECT 
			tp.id,
			tp.name,
			tp.description,
			tp.start_date,
			tp.end_date,
			u.username as client_name,
			COUNT(DISTINCT tpd.id) as exercises_count
		FROM training_plans tp
		JOIN users u ON tp.client_id = u.id
		LEFT JOIN training_plan_days tpd ON tp.id = tpd.training_plan_id
		WHERE tp.trainer_id = $1
		GROUP BY tp.id, tp.name, tp.description, tp.start_date, tp.end_date, u.username
	`

	rows, err := r.db.Query(query, trainerID)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var plans []map[string]interface{}
	for rows.Next() {
		var id, exercisesCount int
		var name, description, clientName sql.NullString
		var startDate, endDate sql.NullTime

		err := rows.Scan(&id, &name, &description, &startDate, &endDate, &clientName, &exercisesCount)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		plan := map[string]interface{}{
			"id":              id,
			"name":            name.String,
			"description":     description.String,
			"client_name":     clientName.String,
			"exercises_count": exercisesCount,
		}

		if startDate.Valid {
			plan["start_date"] = startDate.Time.Format("2006-01-02")
		}
		if endDate.Valid {
			plan["end_date"] = endDate.Time.Format("2006-01-02")
		}

		plans = append(plans, plan)
	}

	fmt.Printf("✅ Repository: Found %d plans for trainer %d\n", len(plans), trainerID)
	return plans, nil
}

func (r *trainingPlanRepository) GetTraineePlans(traineeID int) ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetTraineePlans for trainee %d\n", traineeID)

	query := `
		SELECT 
			tp.id,
			tp.name,
			tp.description,
			tp.start_date,
			tp.end_date,
			u.username as trainer_name,
			COUNT(DISTINCT tpd.id) as exercises_count
		FROM training_plans tp
		JOIN users u ON tp.trainer_id = u.id
		LEFT JOIN training_plan_days tpd ON tp.id = tpd.training_plan_id
		WHERE tp.client_id = $1
		GROUP BY tp.id, tp.name, tp.description, tp.start_date, tp.end_date, u.username
	`

	rows, err := r.db.Query(query, traineeID)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var plans []map[string]interface{}
	for rows.Next() {
		var id, exercisesCount int
		var name, description, trainerName sql.NullString
		var startDate, endDate sql.NullTime

		err := rows.Scan(&id, &name, &description, &startDate, &endDate, &trainerName, &exercisesCount)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		plan := map[string]interface{}{
			"id":              id,
			"name":            name.String,
			"description":     description.String,
			"trainer_name":    trainerName.String,
			"exercises_count": exercisesCount,
		}

		if startDate.Valid {
			plan["start_date"] = startDate.Time.Format("2006-01-02")
		}
		if endDate.Valid {
			plan["end_date"] = endDate.Time.Format("2006-01-02")
		}

		plans = append(plans, plan)
	}

	fmt.Printf("✅ Repository: Found %d plans for trainee %d\n", len(plans), traineeID)
	return plans, nil
}

func (r *trainingPlanRepository) GetExercises() ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetExercises\n")

	query := `
		SELECT id, name, muscles, equipment, level
		FROM exercise_library
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var exercises []map[string]interface{}
	for rows.Next() {
		var id int
		var name, muscles, equipment, level sql.NullString

		err := rows.Scan(&id, &name, &muscles, &equipment, &level)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		exercise := map[string]interface{}{
			"id":        id,
			"name":      name.String,
			"muscles":   muscles.String,
			"equipment": equipment.String,
			"level":     level.String,
		}

		exercises = append(exercises, exercise)
	}

	fmt.Printf("✅ Repository: Found %d exercises\n", len(exercises))
	return exercises, nil
}

func (r *trainingPlanRepository) GetTrainerTrainees(trainerID int) ([]map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetTrainerTrainees for trainer %d\n", trainerID)

	query := `
		SELECT 
			u.id,
			u.username,
			u.role,
			u.experience,
			u.is_online,
			u.created_at,
			u.last_seen,
			p.age,
			p.height,
			p.weight,
			p.gender,
			p.goal,
			p.activity_level
		FROM trainer_trainee ttr
		JOIN users u ON ttr.trainee_id = u.id
		LEFT JOIN profiles p ON u.id = p.user_id
		WHERE ttr.trainer_id = $1
		ORDER BY u.username
	`

	rows, err := r.db.Query(query, trainerID)
	if err != nil {
		fmt.Printf("❌ Repository error: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var trainees []map[string]interface{}
	for rows.Next() {
		var id int
		var experience sql.NullInt64
		var username, role sql.NullString
		var isOnline bool
		var createdAt time.Time
		var lastSeen sql.NullTime
		var age, height sql.NullInt64
		var weight sql.NullFloat64
		var gender, goal sql.NullString
		var activityLevel sql.NullFloat64

		err := rows.Scan(&id, &username, &role, &experience, &isOnline, &createdAt, &lastSeen,
			&age, &height, &weight, &gender, &goal, &activityLevel)
		if err != nil {
			fmt.Printf("❌ Repository scan error: %v\n", err)
			continue
		}

		trainee := map[string]interface{}{
			"id":         id,
			"username":   username.String,
			"role":       role.String,
			"experience": experience,
			"is_online":  isOnline,
			"created_at": createdAt.Format("2006-01-02 15:04:05"),
		}

		if lastSeen.Valid {
			trainee["last_seen"] = lastSeen.Time.Format("2006-01-02 15:04:05")
		}

		// Dodaj profil jeśli istnieje
		if age.Valid || height.Valid || weight.Valid || gender.Valid || goal.Valid || activityLevel.Valid {
			profile := map[string]interface{}{
				"user_id": id,
			}

			if age.Valid {
				profile["age"] = age.Int64
			}
			if height.Valid {
				profile["height"] = height.Int64
			}
			if weight.Valid {
				profile["weight"] = weight.Float64
			}
			if gender.Valid {
				profile["gender"] = gender.String
			}
			if goal.Valid {
				profile["goal"] = goal.String
			}
			if activityLevel.Valid {
				profile["activity_level"] = activityLevel.Float64
			}

			trainee["profile"] = profile
		}

		trainees = append(trainees, trainee)
	}

	fmt.Printf("✅ Repository: Found %d trainees for trainer %d\n", len(trainees), trainerID)
	return trainees, nil
}

func (r *trainingPlanRepository) CreateTrainingPlan(planData map[string]interface{}) (int, error) {
	fmt.Printf("🔍 Repository: CreateTrainingPlan with data: %+v\n", planData)

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Wstaw plan treningowy
	var planID int
	err = tx.QueryRow(`
		INSERT INTO training_plans (client_id, trainer_id, name, description, start_date, end_date, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING id
	`, planData["client_id"], planData["trainer_id"], planData["name"],
		planData["description"], planData["start_date"], planData["end_date"]).Scan(&planID)

	if err != nil {
		fmt.Printf("❌ Error creating training plan: %v\n", err)
		return 0, err
	}

	// Wstaw dni treningowe
	if trainingDays, ok := planData["training_days"].([]interface{}); ok {
		for _, dayInterface := range trainingDays {
			if day, ok := dayInterface.(map[string]interface{}); ok {
				var dayID int
				err = tx.QueryRow(`
					INSERT INTO training_plan_days (training_plan_id, day_of_week)
					VALUES ($1, $2)
					RETURNING id
				`, planID, day["day_of_week"]).Scan(&dayID)

				if err != nil {
					fmt.Printf("❌ Error creating training day: %v\n", err)
					return 0, err
				}

				// Wstaw ćwiczenia
				if exercises, ok := day["exercises"].([]interface{}); ok {
					for _, exerciseInterface := range exercises {
						if exercise, ok := exerciseInterface.(map[string]interface{}); ok {
							_, err = tx.Exec(`
								INSERT INTO training_day_exercises (training_day_id, exercise_id, sets, reps, weight, instructions)
								VALUES ($1, $2, $3, $4, $5, $6)
							`, dayID, exercise["exercise_id"], exercise["sets"], exercise["reps"], exercise["weight"], exercise["instructions"])

							if err != nil {
								fmt.Printf("❌ Error creating exercise: %v\n", err)
								return 0, err
							}
						}
					}
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("❌ Error committing transaction: %v\n", err)
		return 0, err
	}

	fmt.Printf("✅ Repository: Created training plan with ID: %d\n", planID)
	return planID, nil
}

// 🆕 NOWA METODA - UpdateTrainingPlan
func (r *trainingPlanRepository) UpdateTrainingPlan(planID int, planData map[string]interface{}) error {
	fmt.Printf("🔍 Repository: UpdateTrainingPlan - PlanID: %d, Data: %+v\n", planID, planData)

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Aktualizuj podstawowe informacje o planie
	_, err = tx.Exec(`
		UPDATE training_plans 
		SET name = $1, description = $2, start_date = $3, end_date = $4, client_id = $5
		WHERE id = $6
	`, planData["name"], planData["description"], planData["start_date"],
		planData["end_date"], planData["client_id"], planID)

	if err != nil {
		fmt.Printf("❌ Error updating training plan: %v\n", err)
		return err
	}

	// Usuń istniejące dni i ćwiczenia (cascade delete)
	_, err = tx.Exec("DELETE FROM training_plan_days WHERE training_plan_id = $1", planID)
	if err != nil {
		fmt.Printf("❌ Error deleting old training days: %v\n", err)
		return err
	}

	// Dodaj nowe dni treningowe
	if trainingDays, ok := planData["training_days"].([]interface{}); ok {
		for _, dayInterface := range trainingDays {
			if day, ok := dayInterface.(map[string]interface{}); ok {
				var dayID int
				err = tx.QueryRow(`
					INSERT INTO training_plan_days (training_plan_id, day_of_week)
					VALUES ($1, $2)
					RETURNING id
				`, planID, day["day_of_week"]).Scan(&dayID)

				if err != nil {
					fmt.Printf("❌ Error creating training day: %v\n", err)
					return err
				}

				// Dodaj ćwiczenia dla tego dnia
				if exercises, ok := day["exercises"].([]interface{}); ok {
					for _, exerciseInterface := range exercises {
						if exercise, ok := exerciseInterface.(map[string]interface{}); ok {
							_, err = tx.Exec(`
								INSERT INTO training_day_exercises (training_day_id, exercise_id, sets, reps, weight, instructions)
								VALUES ($1, $2, $3, $4, $5, $6)
							`, dayID, exercise["exercise_id"], exercise["sets"], exercise["reps"], exercise["weight"], exercise["instructions"])

							if err != nil {
								fmt.Printf("❌ Error creating exercise: %v\n", err)
								return err
							}
						}
					}
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Printf("❌ Error committing update transaction: %v\n", err)
		return err
	}

	fmt.Printf("✅ Repository: Updated training plan with ID: %d\n", planID)
	return nil
}

func (r *trainingPlanRepository) DeleteTrainingPlan(planID, userID int) error {
	fmt.Printf("🔍 Repository: DeleteTrainingPlan - PlanID: %d, UserID: %d\n", planID, userID)

	// Sprawdź czy plan należy do trenera
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM training_plans WHERE id = $1 AND trainer_id = $2", planID, userID).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("plan nie znaleziony lub nie masz uprawnień")
	}

	// Usuń plan (cascade usunie powiązane dane)
	_, err = r.db.Exec("DELETE FROM training_plans WHERE id = $1", planID)
	if err != nil {
		fmt.Printf("❌ Error deleting training plan: %v\n", err)
		return err
	}

	fmt.Printf("✅ Repository: Deleted training plan %d\n", planID)
	return nil
}

// 🔧 ULEPSZONA METODA - GetPlanDetails z client_id i exercise_id
func (r *trainingPlanRepository) GetPlanDetails(planID, userID int) (map[string]interface{}, error) {
	fmt.Printf("🔍 Repository: GetPlanDetails - PlanID: %d, UserID: %d\n", planID, userID)

	var count int
	err := r.db.QueryRow(`
		SELECT COUNT(*) FROM training_plans 
		WHERE id = $1 AND (trainer_id = $2 OR client_id = $2)
	`, planID, userID).Scan(&count)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, fmt.Errorf("plan nie znaleziony lub brak uprawnień")
	}

	// Pobierz szczegóły planu z dodatkowymi polami
	query := `
		SELECT 
			tp.id,
			tp.client_id,     
			tp.trainer_id,    
			tp.name,
			tp.description,
			tp.start_date,
			tp.end_date,
			tpd.id as day_id,
			tpd.day_of_week,
			tpe.id as exercise_id,
			el.id as exercise_library_id,  
			el.name as exercise_name,
			tpe.sets,
			tpe.reps,
			tpe.weight,
			tpe.instructions
		FROM training_plans tp
		LEFT JOIN training_plan_days tpd ON tp.id = tpd.training_plan_id
		LEFT JOIN training_day_exercises tpe ON tpe.training_day_id = tpd.id
		LEFT JOIN exercise_library el ON tpe.exercise_id = el.id
		WHERE tp.id = $1
		ORDER BY tpd.id, tpe.id
	`

	rows, err := r.db.Query(query, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := map[string]interface{}{
		"days": []map[string]interface{}{},
	}

	daysMap := make(map[int]map[string]interface{})

	for rows.Next() {
		var planID, clientID, trainerID, dayID sql.NullInt64
		var exerciseID, exerciseLibraryID sql.NullInt64
		var name, description, dayOfWeek, exerciseName, instructions sql.NullString
		var startDate, endDate sql.NullTime
		var sets, reps sql.NullInt64
		var weight sql.NullFloat64

		err := rows.Scan(&planID, &clientID, &trainerID, &name, &description, &startDate, &endDate,
			&dayID, &dayOfWeek, &exerciseID, &exerciseLibraryID, &exerciseName, &sets, &reps, &weight, &instructions)
		if err != nil {
			continue
		}

		// Ustaw podstawowe informacje o planie
		if result["id"] == nil {
			result["id"] = planID.Int64
			result["client_id"] = clientID.Int64
			result["trainer_id"] = trainerID.Int64
			result["name"] = name.String
			result["description"] = description.String
			if startDate.Valid {
				result["start_date"] = startDate.Time.Format("2006-01-02")
			}
			if endDate.Valid {
				result["end_date"] = endDate.Time.Format("2006-01-02")
			}
		}

		// Dodaj dni i ćwiczenia
		if dayID.Valid {
			if _, exists := daysMap[int(dayID.Int64)]; !exists {
				daysMap[int(dayID.Int64)] = map[string]interface{}{
					"id":          dayID.Int64,
					"day_of_week": dayOfWeek.String,
					"exercises":   []map[string]interface{}{},
				}
			}

			if exerciseID.Valid {
				exercise := map[string]interface{}{
					"id":            exerciseID.Int64,
					"exercise_id":   exerciseLibraryID.Int64, // ← POPRAWIONE dla frontendu
					"exercise_name": exerciseName.String,
					"sets":          sets.Int64,
					"reps":          reps.Int64,
					"weight":        weight.Float64,
					"instructions":  instructions.String,
				}

				day := daysMap[int(dayID.Int64)]
				exercises := day["exercises"].([]map[string]interface{})
				day["exercises"] = append(exercises, exercise)
			}
		}
	}

	// Konwertuj mapę dni na slice
	days := []map[string]interface{}{}
	for _, day := range daysMap {
		days = append(days, day)
	}
	result["days"] = days

	fmt.Printf("✅ Repository: Got plan details for plan %d\n", planID)
	return result, nil
}
