package models

import (
	"time"
)

type TrainingPlan struct {
	ID          int        `json:"id" db:"id"`
	ClientID    int        `json:"client_id" db:"client_id"`
	TrainerID   int        `json:"trainer_id" db:"trainer_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description" db:"description"`
	StartDate   time.Time  `json:"start_date" db:"start_date"`
	EndDate     *time.Time `json:"end_date" db:"end_date"`
}

type TrainingDay struct {
	ID        int    `json:"id" db:"id"`
	PlanID    int    `json:"plan_id" db:"plan_id"`
	DayOfWeek string `json:"day_of_week" db:"day_of_week"`
}

type TrainingDayExercise struct {
	ID            int     `json:"id" db:"id"`
	TrainingDayID int     `json:"training_day_id" db:"training_day_id"`
	ExerciseID    int     `json:"exercise_id" db:"exercise_id"`
	Sets          int     `json:"sets" db:"sets"`
	Reps          int     `json:"reps" db:"reps"`
	Weight        int     `json:"weight" db:"weight"`
	Instructions  *string `json:"instructions" db:"instructions"`
}

type ExerciseLibrary struct {
	ID        int     `json:"id" db:"id"`
	Name      string  `json:"name" db:"name"`
	Muscles   string  `json:"muscles" db:"muscles"`
	Equipment *string `json:"equipment" db:"equipment"`
	Level     string  `json:"level" db:"level"`
}

// DTOs for API responses
type TrainingPlanWithDetails struct {
	TrainingPlan
	ClientName     string  `json:"client_name"`
	TrainerName    string  `json:"trainer_name"`
	DayOfWeek      *string `json:"day_of_week"`
	ExercisesCount int     `json:"exercises_count"`
}

type TrainingDayWithExercises struct {
	TrainingDay
	Exercises []TrainingDayExerciseWithDetails `json:"exercises"`
}

type TrainingDayExerciseWithDetails struct {
	TrainingDayExercise
	ExerciseName string `json:"exercise_name"`
}

type CreateTrainingPlanRequest struct {
	ClientID     int                        `json:"client_id" binding:"required"`
	Name         string                     `json:"name" binding:"required"`
	Description  *string                    `json:"description"`
	StartDate    time.Time                  `json:"start_date" binding:"required"`
	EndDate      *time.Time                 `json:"end_date"`
	TrainingDays []CreateTrainingDayRequest `json:"training_days" binding:"required"`
}

type CreateTrainingDayRequest struct {
	DayOfWeek string                  `json:"day_of_week" binding:"required"`
	Exercises []CreateExerciseRequest `json:"exercises" binding:"required"`
}

type CreateExerciseRequest struct {
	ExerciseID   int     `json:"exercise_id" binding:"required"`
	Sets         int     `json:"sets" binding:"required,min=1"`
	Reps         int     `json:"reps" binding:"required,min=1"`
	Weight       int     `json:"weight" binding:"required,min=0"`
	Instructions *string `json:"instructions"`
}

type CalendarEvent struct {
	ID            string             `json:"id"`
	Title         string             `json:"title"`
	Start         time.Time          `json:"start"`
	End           time.Time          `json:"end"`
	ExtendedProps CalendarEventProps `json:"extendedProps"`
}

type CalendarEventProps struct {
	PlanID         int    `json:"planId"`
	ClientName     string `json:"clientName,omitempty"`
	ClientID       int    `json:"clientId,omitempty"`
	TrainerName    string `json:"trainerName,omitempty"`
	ExercisesCount int    `json:"exercisesCount"`
	Description    string `json:"description,omitempty"`
}
