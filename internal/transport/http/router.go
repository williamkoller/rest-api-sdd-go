package http

import (
	"github.com/gin-gonic/gin"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/handler"
	"github.com/williamkoller/rest-api-sdd-go/internal/transport/http/middleware"
)

type Handlers struct {
	Auth         *handler.AuthHandler
	Health       *handler.HealthHandler
	School       *handler.SchoolHandler
	Unit         *handler.UnitHandler
	Classroom    *handler.ClassroomHandler
	Class        *handler.ClassHandler
	Student      *handler.StudentHandler
	Attendance   *handler.AttendanceHandler
	Grade        *handler.GradeHandler
	Invoice      *handler.InvoiceHandler
	Ticket       *handler.TicketHandler
	Agenda       *handler.AgendaHandler
	Calendar     *handler.CalendarHandler
	Feed         *handler.FeedHandler
	Referral     *handler.ReferralHandler
	Menu         *handler.MenuHandler
	Curriculum   *handler.CurriculumHandler
	Reenrollment *handler.ReenrollmentHandler
	Waitlist     *handler.WaitlistHandler
}

func NewRouter(h *Handlers, jwtSecret string) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger(), middleware.Recovery())

	r.GET("/health", h.Health.Health)

	auth := r.Group("/api/v1/auth")
	{
		auth.POST("/login", h.Auth.Login)
		auth.POST("/refresh", h.Auth.Refresh)
		auth.POST("/logout", h.Auth.Logout)
	}

	// Waitlist public registration
	r.POST("/api/v1/units/:unit_id/waitlist", h.Waitlist.Register)

	api := r.Group("/api/v1")

	// Schools context — all school-scoped sub-resources share :school_id
	schools := api.Group("/schools")
	{
		schools.GET("", h.School.List)
		schools.POST("", h.School.Create)
		schools.GET("/:school_id", h.School.Get)
		schools.PUT("/:school_id", h.School.Update)
		schools.GET("/:school_id/units", h.Unit.List)
		schools.POST("/:school_id/units", h.Unit.Create)
		schools.GET("/:school_id/calendar", h.Calendar.List)
		schools.POST("/:school_id/calendar", h.Calendar.Create)
		schools.GET("/:school_id/financial/delinquency", h.Invoice.Delinquency)
		schools.GET("/:school_id/tickets/report", h.Ticket.Report)
		schools.GET("/:school_id/feed", h.Feed.List)
		schools.GET("/:school_id/referrals", h.Referral.ListReferrals)
	}

	// Units context — standalone unit operations share :unit_id
	units := api.Group("/units")
	{
		units.GET("/:unit_id", h.Unit.Get)
		units.PUT("/:unit_id", h.Unit.Update)
		units.GET("/:unit_id/classrooms", h.Classroom.List)
		units.POST("/:unit_id/classrooms", h.Classroom.Create)
		units.GET("/:unit_id/classes", h.Class.List)
		units.POST("/:unit_id/classes", h.Class.Create)
		units.GET("/:unit_id/waitlist", h.Waitlist.List)
		units.GET("/:unit_id/menu", h.Menu.Get)
		units.POST("/:unit_id/menu", h.Menu.Create)
	}

	// Classrooms context
	classrooms := api.Group("/classrooms")
	{
		classrooms.PUT("/:classroom_id", h.Classroom.Update)
	}

	// Classes context — all class-scoped sub-resources share :class_id
	classes := api.Group("/classes")
	{
		classes.GET("/:class_id", h.Class.Get)
		classes.PUT("/:class_id", h.Class.Update)
		classes.GET("/:class_id/students", h.Student.List)
		classes.POST("/:class_id/students", h.Student.Enroll)
		classes.DELETE("/:class_id/students/:student_id", h.Student.Unenroll)
		classes.POST("/:class_id/attendance", h.Attendance.BatchRecord)
		classes.POST("/:class_id/grades", h.Grade.BatchUpsert)
		classes.GET("/:class_id/agenda", h.Agenda.List)
		classes.POST("/:class_id/agenda", h.Agenda.Create)
		classes.GET("/:class_id/curriculum", h.Curriculum.List)
		classes.POST("/:class_id/curriculum", h.Curriculum.BatchCreate)
	}

	// Students context — all student-scoped sub-resources share :student_id
	students := api.Group("/students")
	{
		students.GET("/:student_id", h.Student.Get)
		students.PUT("/:student_id", h.Student.Update)
		students.GET("/:student_id/attendance", h.Attendance.GetByStudent)
		students.GET("/:student_id/grades", h.Grade.GetByStudent)
		students.GET("/:student_id/invoices", h.Invoice.ListByStudent)
	}

	// Invoices context
	invoices := api.Group("/invoices")
	{
		invoices.POST("/generate", h.Invoice.Generate)
		invoices.GET("/:invoice_id", h.Invoice.Get)
		invoices.POST("/:invoice_id/pay", h.Invoice.Pay)
		invoices.GET("/:invoice_id/receipt", h.Invoice.Receipt)
	}

	// Re-enrollment context
	reenrollment := api.Group("/reenrollment/campaigns")
	{
		reenrollment.POST("", h.Reenrollment.CreateCampaign)
		reenrollment.GET("/:campaign_id/dashboard", h.Reenrollment.Dashboard)
		reenrollment.POST("/:campaign_id/respond", h.Reenrollment.Respond)
	}

	// Waitlist management (authenticated)
	waitlist := api.Group("/waitlist")
	{
		waitlist.PUT("/:waitlist_id/status", h.Waitlist.UpdateStatus)
	}

	// Agenda context
	agenda := api.Group("/agenda")
	{
		agenda.PUT("/:agenda_id", h.Agenda.Update)
		agenda.DELETE("/:agenda_id", h.Agenda.Delete)
	}

	// Calendar context
	calendar := api.Group("/calendar")
	{
		calendar.PUT("/:calendar_id", h.Calendar.Update)
		calendar.DELETE("/:calendar_id", h.Calendar.Delete)
	}

	// Tickets context
	tickets := api.Group("/tickets")
	{
		tickets.POST("", h.Ticket.Create)
		tickets.GET("", h.Ticket.List)
		tickets.GET("/:ticket_id", h.Ticket.Get)
		tickets.POST("/:ticket_id/reply", h.Ticket.Reply)
		tickets.PUT("/:ticket_id/status", h.Ticket.UpdateStatus)
	}

	// Feed context
	feed := api.Group("/feed")
	{
		feed.POST("", h.Feed.Create)
		feed.DELETE("/:feed_id", h.Feed.Delete)
	}

	// Menu context
	menu := api.Group("/menu")
	{
		menu.PUT("/:menu_id", h.Menu.Update)
	}

	// Curriculum context
	curriculum := api.Group("/curriculum")
	{
		curriculum.PUT("/:curriculum_id", h.Curriculum.Update)
		curriculum.DELETE("/:curriculum_id", h.Curriculum.Delete)
	}

	// Me context — authenticated user own resources
	me := api.Group("/me")
	{
		me.GET("/referral-link", h.Referral.GetMyLink)
	}

	return r
}
