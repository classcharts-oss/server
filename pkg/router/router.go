package router

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/classcharts-oss/server/pkg/frontend"
	data2 "github.com/classcharts-oss/server/pkg/student/data"
	"github.com/classcharts-oss/server/pkg/student/user"
	"github.com/classcharts-oss/server/pkg/test"

	"github.com/gorilla/mux"
)

func CreateMuxRouter() *mux.Router {
	// StrictSlash(true) is used to make sure that the router will automatically redirect requests with a trailing slash to the equivalent URL without a trailing slash.
	// Honestly, this is more of a preference thing. I don't think it's necessary to have this, but it's good to have it anyway. Some dev out there might forget that this doesn't use trailing slashes, and end up spending 2 hours debugging why their code isn't working.
	// This is just a safety net to save some unobservant people :)
	router := mux.NewRouter().StrictSlash(true)
	router.Use(ErrorHandler)
	router.Use(RequestHandler)

	CreateStudentRoutes(router.PathPrefix("/apiv2student").Subrouter(), true)
	CreateStudentV1Routes(router.PathPrefix("/student").Subrouter())

	CreateParentRoutes(router.PathPrefix("/apiv2parent").Subrouter())
	//CreateParentReportAbsenceRoutes(router.PathPrefix("/apiv2parentreportabsence").Subrouter())

	CreateTestRouter(router.PathPrefix("/test").Subrouter())

	CreateFrontendRoutes(router)

	return router
}

func CreateStudentRoutes(v2student *mux.Router, includeExtras bool) *mux.Router {
	restrictedv2Student := v2student.PathPrefix("").Subrouter()
	restrictedv2Student.Use(AuthHandler)

	if includeExtras {
		v2student.HandleFunc("/hasdob", user.HasDOBHandler).Methods(http.MethodPost)
		v2student.HandleFunc("/login", user.LoginHandler).Methods(http.MethodPost)
		restrictedv2Student.HandleFunc("/ping", user.StudentUserHandler).Methods(http.MethodPost)
		restrictedv2Student.HandleFunc("/getcode", user.GetCodeHandler).Methods(http.MethodPost)
		v2student.HandleFunc("/logout", user.LogoutHandler).Methods(http.MethodPost)
	}

	restrictedv2Student.HandleFunc("/behaviour/{studentId}", data2.GetBehaviourHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/activity/{studentId}", data2.GetActivityHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/announcements/{studentId}", data2.GetAnnouncementsHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/addconcern", data2.AddConcernHandler).Methods(http.MethodPost)

	restrictedv2Student.HandleFunc("/getacademicreports", data2.ListAcademicReportsHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getacademicreport/{id}", data2.GetAcademicReportHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/getpupilreportcards", data2.ListOnReportCardsHandler).Methods(http.MethodPost)
	restrictedv2Student.HandleFunc("/getpupilreportcard/{id}", data2.GetOnReportCardHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getpupilreportcardsummarycomment/{id}", data2.GetOnReportCardSummaryCommentHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getpupilreportcardtarget/{id}", data2.GetOnReportCardTargetHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/timetable/{studentId}", data2.TimetableHandler).Methods(http.MethodGet)

	return v2student
}

func CreateParentRoutes(v2parent *mux.Router) *mux.Router {
	//v2parent.HandleFunc("/login", user2.LoginHandler).Methods(http.MethodPost)
	//v2parent.HandleFunc("/ping", user2.ParentUserHandler).Methods(http.MethodPost)
	//v2parent.HandleFunc("/logout", user2.LogoutHandler).Methods(http.MethodPost)
	//
	//v2parent.HandleFunc("/pupils", user2.GetPupilsHandler).Methods(http.MethodGet)
	//v2parent.HandleFunc("/announcements", data.GetAnnouncementsHandler).Methods(http.MethodGet)

	CreateStudentRoutes(v2parent, false) // Creates all the /apiv2parent/behaviour, /apiv2parent/activity, etc. routes

	return v2parent
}

//func CreateParentReportAbsenceRoutes(v2parentreportabs *mux.Router) *mux.Router {
//	v2parentreportabs.HandleFunc("/getreportedabsences/{studentId}", parentData.ListReportedAbsencesHandler).Methods(http.MethodGet)
//	return v2parentreportabs
//}

func CreateStudentV1Routes(v1student *mux.Router) *mux.Router {
	v1student.HandleFunc("/checkpupilcode/{code}", user.CheckPupilCodeHandler).Methods(http.MethodPost)

	return v1student
}

func CreateTestRouter(router *mux.Router) *mux.Router {
	restrictedTest := router.PathPrefix("").Subrouter()
	restrictedTest.Use(AuthHandler)

	router.HandleFunc("/newstudent", test.CreateStudentHandler).Methods(http.MethodPost)
	router.HandleFunc("/getstudent", test.GetStudentHandler).Methods(http.MethodGet)

	router.HandleFunc("/newschool", test.CreateSchoolHandler).Methods(http.MethodPost)

	restrictedTest.HandleFunc("/newannouncement", data2.CreateAnnouncementHandler).Methods(http.MethodPost)

	restrictedTest.HandleFunc("/newactivity", data2.CreateActivityHandler).Methods(http.MethodPost)

	return router
}

func CreateFrontendRoutes(router *mux.Router) *mux.Router {
	router.Handle("/", templ.Handler(frontend.Hello("NERD")))

	return router
}
