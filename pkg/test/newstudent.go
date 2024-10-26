package test

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/CommunityCharts/CCModels/shared"
	"github.com/CommunityCharts/CCModels/student"
	"github.com/classcharts-oss/server/pkg/db"
	"github.com/classcharts-oss/server/pkg/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateStudentHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("schoolId"))
	if err != nil {
		response := shared.NewErrorfulResponse("Invalid school ID, or none provided.")
		response.Write(w)
		return
	}

	school, err := db.GetSchoolByID(i)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			response := shared.NewErrorfulResponse("School not found")
			response.Write(w)
			return
		} else {
			panic(err)
		}
	}

	st := db.CreateStudent(student.DBStudentUser{
		StudentUser: student.NewUser(db.GetNextID(), r.FormValue("name"), r.FormValue("avatarUrl")),
		SchoolId:    school.Id,

		DOB:  util.ToPtr("2010-01-01"),
		Code: r.FormValue("code"),
	})

	response := shared.NewSuccessfulResponse(st)
	response.Write(w)
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.URL.Query().Get("studentId"))
	if err != nil {
		response := shared.NewErrorfulResponse("Invalid student ID, or none provided.")
		response.Write(w)
		return
	}

	st := db.GetStudentByID(i)

	response := shared.NewSuccessfulResponse(st)
	response.Write(w)
}
