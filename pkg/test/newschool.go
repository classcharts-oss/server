package test

import (
	"net/http"
	"strconv"

	"github.com/CommunityCharts/CCModels/school"
	"github.com/CommunityCharts/CCModels/shared"
	"github.com/classcharts-oss/server/pkg/db"
)

func CreateSchoolHandler(w http.ResponseWriter, r *http.Request) {
	i, err := strconv.Atoi(r.FormValue("schoolId"))
	if err != nil {
		response := shared.NewErrorfulResponse("Invalid school ID, or none provided.")
		response.Write(w)
		return
	}

	db.CreateSchool(school.NewSchool(i, r.FormValue("name"), r.FormValue("logo")))

	response := shared.NewSuccessfulResponse(shared.Object{})
	response.Write(w)
}
