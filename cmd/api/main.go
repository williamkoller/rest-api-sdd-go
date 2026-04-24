package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/williamkoller/rest-api-sdd-go/config"
	"github.com/williamkoller/rest-api-sdd-go/internal/application/usecase"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/cache"
	memorycache "github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/cache/memory"
	rediscache "github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/cache/redis"
	"github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/database"
	infrarepo "github.com/williamkoller/rest-api-sdd-go/internal/infrastructure/repository"
	transporthttp "github.com/williamkoller/rest-api-sdd-go/internal/transport/http"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/handler"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	var cacheImpl cache.Cache
	if cfg.Cache.Driver == "redis" {
		cacheImpl = rediscache.New(cfg)
	} else {
		cacheImpl = memorycache.New()
	}

	// Repositories
	schoolRepo := infrarepo.NewSchoolRepository(db)
	unitRepo := infrarepo.NewUnitRepository(db)
	classroomRepo := infrarepo.NewClassroomRepository(db)
	classRepo := infrarepo.NewClassRepository(db)
	userRepo := infrarepo.NewUserRepository(db)
	studentRepo := infrarepo.NewStudentRepository(db)
	enrollmentRepo := infrarepo.NewEnrollmentRepository(db)
	attendanceRepo := infrarepo.NewAttendanceRepository(db)
	gradeRepo := infrarepo.NewGradeRepository(db)
	invoiceRepo := infrarepo.NewInvoiceRepository(db)
	reenrollRepo := infrarepo.NewReenrollmentRepository(db)
	waitlistRepo := infrarepo.NewWaitlistRepository(db)
	agendaRepo := infrarepo.NewAgendaRepository(db)
	calendarRepo := infrarepo.NewCalendarRepository(db)
	ticketRepo := infrarepo.NewTicketRepository(db)
	feedRepo := infrarepo.NewFeedRepository(db)
	referralRepo := infrarepo.NewReferralRepository(db)
	menuRepo := infrarepo.NewMenuRepository(db)
	curriculumRepo := infrarepo.NewCurriculumRepository(db)

	// Use Cases
	baseURL := fmt.Sprintf("http://localhost:%s", cfg.App.Port)

	authUC := usecase.NewAuthUseCase(userRepo, cacheImpl, cfg)
	schoolUC := usecase.NewSchoolUseCase(schoolRepo)
	unitUC := usecase.NewUnitUseCase(unitRepo)
	classroomUC := usecase.NewClassroomUseCase(classroomRepo)
	classUC := usecase.NewClassUseCase(classRepo)
	studentUC := usecase.NewStudentUseCase(studentRepo, enrollmentRepo)
	attendanceUC := usecase.NewAttendanceUseCase(attendanceRepo, userRepo, studentRepo)
	gradeUC := usecase.NewGradeUseCase(gradeRepo, userRepo, studentRepo)
	invoiceUC := usecase.NewInvoiceUseCase(invoiceRepo, studentRepo)
	reenrollUC := usecase.NewReenrollmentUseCase(reenrollRepo, invoiceRepo)
	waitlistUC := usecase.NewWaitlistUseCase(waitlistRepo)
	agendaUC := usecase.NewAgendaUseCase(agendaRepo, userRepo)
	calendarUC := usecase.NewCalendarUseCase(calendarRepo)
	ticketUC := usecase.NewTicketUseCase(ticketRepo)
	feedUC := usecase.NewFeedUseCase(feedRepo)
	referralUC := usecase.NewReferralUseCase(referralRepo, baseURL)
	menuUC := usecase.NewMenuUseCase(menuRepo)
	curriculumUC := usecase.NewCurriculumUseCase(curriculumRepo)

	// Handlers
	handlers := &transporthttp.Handlers{
		Auth:         handler.NewAuthHandler(authUC),
		Health:       handler.NewHealthHandler(db, cacheImpl),
		School:       handler.NewSchoolHandler(schoolUC),
		Unit:         handler.NewUnitHandler(unitUC),
		Classroom:    handler.NewClassroomHandler(classroomUC),
		Class:        handler.NewClassHandler(classUC),
		Student:      handler.NewStudentHandler(studentUC),
		Attendance:   handler.NewAttendanceHandler(attendanceUC),
		Grade:        handler.NewGradeHandler(gradeUC),
		Invoice:      handler.NewInvoiceHandler(invoiceUC),
		Reenrollment: handler.NewReenrollmentHandler(reenrollUC),
		Waitlist:     handler.NewWaitlistHandler(waitlistUC),
		Agenda:       handler.NewAgendaHandler(agendaUC),
		Calendar:     handler.NewCalendarHandler(calendarUC),
		Ticket:       handler.NewTicketHandler(ticketUC),
		Feed:         handler.NewFeedHandler(feedUC),
		Referral:     handler.NewReferralHandler(referralUC),
		Menu:         handler.NewMenuHandler(menuUC),
		Curriculum:   handler.NewCurriculumHandler(curriculumUC),
	}

	router := transporthttp.NewRouter(handlers, cfg.JWT.SecretKey)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.App.Port, "env", cfg.App.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forced shutdown", "error", err)
	}
	slog.Info("server stopped")
}
