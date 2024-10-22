package router

import (
	"net/http"

	studentData "github.com/daydreme/classcharts-server-mock/pkg/student/data"
	studentUser "github.com/daydreme/classcharts-server-mock/pkg/student/user"
	"github.com/daydreme/classcharts-server-mock/pkg/test"

	parentUser "github.com/daydreme/classcharts-server-mock/pkg/parent/user"

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

	return router
}

func CreateStudentRoutes(v2student *mux.Router, includeExtras bool) *mux.Router {
	restrictedv2Student := v2student.PathPrefix("").Subrouter()
	restrictedv2Student.Use(AuthHandler)

	if includeExtras {
		v2student.HandleFunc("/hasdob", studentUser.HasDOBHandler).Methods(http.MethodPost)
		v2student.HandleFunc("/login", studentUser.LoginHandler).Methods(http.MethodPost)
		restrictedv2Student.HandleFunc("/ping", studentUser.StudentUserHandler).Methods(http.MethodPost)
		restrictedv2Student.HandleFunc("/getcode", studentUser.GetCodeHandler).Methods(http.MethodPost)
		v2student.HandleFunc("/logout", studentUser.LogoutHandler).Methods(http.MethodPost)
	}

	restrictedv2Student.HandleFunc("/behaviour/{studentId}", studentData.GetBehaviourHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/activity/{studentId}", studentData.GetActivityHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/announcements/{studentId}", studentData.GetAnnouncementHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/addconcern", studentData.AddConcernHandler).Methods(http.MethodPost)

	restrictedv2Student.HandleFunc("/getacademicreports", studentData.ListAcademicReportsHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getacademicreport/{id}", studentData.GetAcademicReportHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/getpupilreportcards", studentData.ListOnReportCardsHandler).Methods(http.MethodPost)
	restrictedv2Student.HandleFunc("/getpupilreportcard/{id}", studentData.GetOnReportCardHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getpupilreportcardsummarycomment/{id}", studentData.GetOnReportCardSummaryCommentHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/getpupilreportcardtarget/{id}", studentData.GetOnReportCardTargetHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/timetable/{studentId}", studentData.TimetableHandler).Methods(http.MethodGet)

	restrictedv2Student.HandleFunc("/rewards/{studentId}", studentData.GetRewardHandler).Methods(http.MethodGet)
	restrictedv2Student.HandleFunc("/purchase/{itemId}", studentData.GetPurchaseHandler).Methods(http.MethodPost)

	return v2student
}

func CreateParentRoutes(v2parent *mux.Router) *mux.Router {
	v2parent.HandleFunc("/login", parentUser.LoginHandler).Methods(http.MethodPost)
	v2parent.HandleFunc("/ping", parentUser.ParentUserHandler).Methods(http.MethodPost)
	v2parent.HandleFunc("/logout", parentUser.LogoutHandler).Methods(http.MethodPost)

	v2parent.HandleFunc("/pupils", parentUser.GetPupilsHandler).Methods(http.MethodGet)
	v2parent.HandleFunc("/announcements", studentData.GetAnnouncementHandler).Methods(http.MethodGet)

	CreateStudentRoutes(v2parent, false) // Creates all the /apiv2parent/behaviour, /apiv2parent/activity, etc. routes

	return v2parent
}

//func CreateParentReportAbsenceRoutes(v2parentreportabs *mux.Router) *mux.Router {
//	v2parentreportabs.HandleFunc("/getreportedabsences/{studentId}", parentData.ListReportedAbsencesHandler).Methods(http.MethodGet)
//	return v2parentreportabs
//}

func CreateStudentV1Routes(v1student *mux.Router) *mux.Router {
	v1student.HandleFunc("/checkpupilcode/{code}", studentUser.CheckPupilCodeHandler).Methods(http.MethodPost)

	return v1student
}

func CreateTestRouter(router *mux.Router) *mux.Router {
	router.HandleFunc("/newstudent", test.CreateStudentHandler).Methods(http.MethodPost)
	router.HandleFunc("/getstudent", test.GetStudentHandler).Methods(http.MethodGet)

	router.HandleFunc("/newschool", test.CreateSchoolHandler).Methods(http.MethodPost)

	return router
}
