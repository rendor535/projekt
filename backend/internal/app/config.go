package app

import (
	"backend/internal/domain/repositories"
	"backend/internal/handlers"
	"backend/internal/middleware"
	"database/sql"
	"github.com/gin-gonic/gin"
)

func SetupMiddleware(router *gin.Engine) {
	router.Use(middleware.CORSMiddleware())
}

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Backend is up and running"})
	})

	// Auth routes
	auth := handlers.NewAuthHandler(db)
	router.POST("/register", auth.Register)
	router.POST("/login", auth.Login)
	router.GET("/me", middleware.AuthMiddleware(), auth.Me)
	// User status handler
	userHandler := handlers.NewUserHandler(db)
	router.POST("/me/online", middleware.AuthMiddleware(), userHandler.SetOnline)
	router.POST("/me/offline", middleware.AuthMiddleware(), userHandler.SetOffline)
	// Profile routes
	profileRepo := repositories.NewProfileRepository(db)
	profileHandler := handlers.NewProfileHandler(profileRepo, db) // Dodaj db jako drugi parametr
	authGroup := router.Group("/profile")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.GET("", profileHandler.GetProfile)
		authGroup.POST("", profileHandler.SaveProfile)
		authGroup.GET("/check", profileHandler.CheckProfile)
		authGroup.GET("/trainee/:traineeId", profileHandler.GetTraineeProfile) // Dodaj nowy endpoint
	}
	trainerProfileRepo := repositories.NewTrainerProfileRepository(db)
	trainerProfileHandler := handlers.NewTrainerProfileHandler(trainerProfileRepo)

	trainerGroup := router.Group("/trainer-profile")
	trainerGroup.Use(middleware.AuthMiddleware())
	{
		trainerGroup.GET("", trainerProfileHandler.GetProfile)
		trainerGroup.POST("", trainerProfileHandler.SaveProfile)
		trainerGroup.GET("/check", trainerProfileHandler.CheckProfile)
		trainerGroup.GET("/:id", trainerProfileHandler.GetTrainerByID) // ← DODAJ TĘ LINIĘ
	}

	// Training plan routes
	trainingPlanRepo := repositories.NewTrainingPlanRepository(db)
	trainingPlanHandler := handlers.NewTrainingPlanHandler(trainingPlanRepo)

	trainingGroup := router.Group("/training-plans")
	trainingGroup.Use(middleware.AuthMiddleware())
	{
		// Pobieranie planów
		trainingGroup.GET("/trainer", trainingPlanHandler.GetTrainerPlans)    // GET /training-plans/trainer
		trainingGroup.GET("/trainee", trainingPlanHandler.GetTraineePlans)    // GET /training-plans/trainee
		trainingGroup.GET("/calendar", trainingPlanHandler.GetCalendarEvents) // GET /training-plans/calendar

		// CRUD operacje
		trainingGroup.POST("", trainingPlanHandler.CreatePlan)                // POST /training-plans
		trainingGroup.DELETE("/:id", trainingPlanHandler.DeletePlan)          // DELETE /training-plans/:id
		trainingGroup.GET("/:id/details", trainingPlanHandler.GetPlanDetails) // GET /training-plans/:id/details
		trainingGroup.PUT("/:id", trainingPlanHandler.UpdatePlan)
		// Pomocnicze endpointy
		trainingGroup.GET("/exercises", trainingPlanHandler.GetExercises) // GET /training-plans/exercises
		trainingGroup.GET("/trainees", trainingPlanHandler.GetTrainees)   // GET /training-plans/trainees
	}

	//// Workout routes
	//workoutRepo := repositories.NewWorkoutRepository(db)
	//workoutHandler := handlers.NewWorkoutHandler(workoutRepo)
	//router.Group("/plans/:planID/workouts").
	//	Use(middleware.AuthMiddleware()).
	//	POST("", workoutHandler.Create).
	//	GET("", workoutHandler.List)
	//
	//// Exercise routes
	//exerciseRepo := repositories.NewExerciseRepository(db)
	//exerciseHandler := handlers.NewExerciseHandler(exerciseRepo)
	//router.Group("/workouts/:workoutID/exercises").
	//	Use(middleware.AuthMiddleware()).
	//	POST("", exerciseHandler.Create).
	//	GET("", exerciseHandler.List)

	// Calculator route
	calculator := handlers.NewCalculatorHandler(db)
	router.POST("/calculate", middleware.AuthMiddleware(), calculator.Calculate)

	// Product routes
	productHandler := handlers.NewProductHandler(db)
	router.GET("/products", middleware.AuthMiddleware(), productHandler.List)

	// Message routes
	messageRepo := repositories.NewMessageRepository(db)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	messageGroup := router.Group("/messages")
	messageGroup.Use(middleware.AuthMiddleware())
	{
		messageGroup.POST("", messageHandler.SendMessage)
		messageGroup.GET("/conversation/:userID", messageHandler.GetConversation)
		messageGroup.GET("/unread", messageHandler.GetUnreadCount)
		messageGroup.PUT("/:messageID/read", messageHandler.MarkAsRead)
		messageGroup.GET("/partners", messageHandler.GetConversationPartners)
		messageGroup.PUT("/conversation/:senderID/read", messageHandler.MarkConversationAsRead)

	}

	// Trainer-Trainee relationship routes
	trainerTraineeRepo := repositories.NewTrainerTraineeRepository(db)
	trainerTraineeHandler := handlers.NewTrainerTraineeHandler(trainerTraineeRepo)

	router.POST("/trainees/assign", middleware.AuthMiddleware(), trainerTraineeHandler.AssignTrainee)
	router.GET("/trainees", middleware.AuthMiddleware(), trainerTraineeHandler.GetTrainees)
	router.DELETE("/trainees/:traineeID", middleware.AuthMiddleware(), trainerTraineeHandler.RemoveTrainee)
	router.GET("/my-trainer", middleware.AuthMiddleware(), trainerTraineeHandler.GetMyTrainer)

	// NEW: Trainer Request routes for TraineeDashboard functionality
	trainerRequestRepo := repositories.NewTrainerRequestRepository(db)
	trainerRequestHandler := handlers.NewTrainerRequestHandler(trainerRequestRepo)
	trainerRequestGroup := router.Group("/trainer-requests")
	trainerRequestGroup.Use(middleware.AuthMiddleware())
	{
		trainerRequestGroup.POST("", trainerRequestHandler.SendRequest)
		trainerRequestGroup.GET("/my", trainerRequestHandler.GetMyRequests)
		trainerRequestGroup.DELETE("/:id", trainerRequestHandler.CancelRequest)
		trainerRequestGroup.GET("/my-trainer", trainerRequestHandler.GetMyTrainer)

		// New endpoints for trainers
		trainerRequestGroup.GET("/for-me", trainerRequestHandler.GetRequestsForMe)
		trainerRequestGroup.PUT("/:id/status", trainerRequestHandler.UpdateRequestStatus)
	}

	// Trainer search route (publicly accessible for registered users)
	router.GET("/trainers/search", middleware.AuthMiddleware(), trainerRequestHandler.SearchTrainers)

	// NEW: Progress Report routes
	progressReportRepo := repositories.NewProgressReportRepository(db)
	progressReportHandler := handlers.NewProgressReportHandler(progressReportRepo)

	reportsGroup := router.Group("/reports")
	reportsGroup.Use(middleware.AuthMiddleware())
	{
		reportsGroup.POST("/", progressReportHandler.Create)
		reportsGroup.GET("/my", progressReportHandler.GetMyReports)
		reportsGroup.GET("/my/stats", progressReportHandler.GetMyStats)
		reportsGroup.GET("/my/chart", progressReportHandler.GetChartData)
		reportsGroup.GET("/latest-weight", progressReportHandler.GetLatestWeight)

		// Endpointy dla trenera
		reportsGroup.GET("/trainer/stats", progressReportHandler.GetTrainerStats)
		reportsGroup.GET("/trainer/trainee/:traineeId", progressReportHandler.GetTraineeReports)
		reportsGroup.GET("/trainer/trainee/:traineeId/chart", progressReportHandler.GetTraineeChartData)
	}

	// NEW: User Media routes (for shared videos/files)
	mediaRepo := repositories.NewUserMediaRepository(db)
	mediaHandler := handlers.NewUserMediaHandler(mediaRepo, db)

	mediaGroup := router.Group("/media")
	mediaGroup.Use(middleware.AuthMiddleware())
	{
		mediaGroup.POST("/upload", mediaHandler.Upload)
		mediaGroup.GET("/shared", mediaHandler.GetSharedMedia)
		mediaGroup.GET("/uploaded", mediaHandler.GetUploadedMedia)
		mediaGroup.GET("/users", mediaHandler.GetUsers)
		mediaGroup.DELETE("/:id", mediaHandler.DeleteMedia)
		mediaGroup.PATCH("/:id/read", mediaHandler.MarkAsRead)
		mediaGroup.GET("/exercise-videos", mediaHandler.GetExerciseVideos)             // Pobierz bibliotekę filmów
		mediaGroup.POST("/assign-exercise-video", mediaHandler.AddExerciseVideoToUser) // Przypisz film z biblioteki
	}

	// NEW: Diet Plan routes
	dietPlanRepo := repositories.NewDietPlanRepository(db)
	dietPlanHandler := handlers.NewDietPlanHandler(dietPlanRepo)

	dietGroup := router.Group("/diet-plans")
	dietGroup.Use(middleware.AuthMiddleware())
	{
		// Trainer endpoints
		dietGroup.POST("", dietPlanHandler.CreateDietPlan)             // POST /diet-plans
		dietGroup.GET("/trainer", dietPlanHandler.GetTrainerDietPlans) // GET /diet-plans/trainer?client_id=&is_active=
		dietGroup.PUT("/:id", dietPlanHandler.UpdateDietPlan)          // PUT /diet-plans/:id
		dietGroup.DELETE("/:id", dietPlanHandler.DeleteDietPlan)       // DELETE /diet-plans/:id

		// Trainee endpoints
		dietGroup.GET("/trainee", dietPlanHandler.GetTraineeDietPlans) // GET /diet-plans/trainee?is_active=

		// Shared endpoints
		dietGroup.GET("/:id/details", dietPlanHandler.GetDietPlanDetails) // GET /diet-plans/:id/details

		// Legacy endpoints (zachowane dla kompatybilności)
		dietGroup.POST("/legacy", dietPlanHandler.Create)
		dietGroup.GET("/my", dietPlanHandler.GetMyDietPlans)
		dietGroup.GET("/active", dietPlanHandler.GetActiveDietPlan)
		dietGroup.PUT("/:id/activate", dietPlanHandler.ActivateDietPlan)
	}
	notificationHandler := handlers.NewNotificationHandler(db)
	notifGroup := router.Group("/")
	notifGroup.Use(middleware.AuthMiddleware())
	{
		// Endpointy dla powiadomień
		notifGroup.GET("/reports/trainer/trainee/:traineeId/unread", notificationHandler.GetUnreadReportsCount)
		notifGroup.GET("/media/trainee/:traineeId/unread", notificationHandler.GetUnreadMediaCount)
		notifGroup.GET("/diet-plans/trainee/:traineeId/pending-feedback", notificationHandler.GetPendingDietFeedbackCount)
		notifGroup.GET("/reports/trainer/pending", notificationHandler.GetTrainerPendingReportsCount)
	}

}
