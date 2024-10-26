package data

import (
	"net/http"
	"strconv"
	"time"

	"github.com/classcharts-oss/server/pkg/db"
	"github.com/classcharts-oss/server/pkg/util"

	"github.com/CommunityCharts/CCModels/shared"
)

type getActivityMeta struct {
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	LastId           int    `json:"last_id"`
	StepSize         string `json:"step_size"`
	DetentionAliasUC string `json:"detention_alias_uc"`
}

func GetActivityHandler(w http.ResponseWriter, r *http.Request) {
	s := util.GetStudent(r)
	meta := getActivityMeta{
		StartDate:        time.Now().Format(time.RFC3339),
		EndDate:          time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		LastId:           1,
		StepSize:         "week",
		DetentionAliasUC: "Detention",
	}

	data := s.Activities
	response := shared.NewSuccessfulMetaResponse(data, meta)
	response.Write(w)
}

func CreateActivityHandler(w http.ResponseWriter, r *http.Request) {
	s := util.GetStudent(r)
	ran := util.RandomId()

	score, err := strconv.Atoi(r.FormValue("score"))
	if err != nil {
		panic(err)
	}

	activity := shared.Activity{
		Id:       ran,
		Type:     shared.ActivityType(r.FormValue("type")),
		Polarity: shared.ActivityPolarity(r.FormValue("polarity")),
		Reason:   r.FormValue("reason"),
		Score:    score,

		Timestamp: time.Now().Format(time.RFC3339),

		Style: shared.ActivityStyle{},

		PupilName:   s.Name,
		LessonName:  util.ToPtr(r.FormValue("lesson_name")),
		TeacherName: util.ToPtr(r.FormValue("teacher_name")),
		RoomName:    util.ToPtr(r.FormValue("room_name")),

		Note:       util.ToPtr(r.FormValue("note")),
		ParentNote: util.ToPtr(r.FormValue("parent_note")),

		DetentionDate:     nil,
		DetentionTime:     nil,
		DetentionLocation: nil,
		DetentionType:     nil,
	}

	s.Activities = append(s.Activities, activity)
	db.UpdateStudent(s)

	response := shared.NewSuccessfulResponse(activity)
	response.Write(w)
}
