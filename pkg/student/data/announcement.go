package data

import (
	"net/http"
	"time"

	"github.com/CommunityCharts/CCModels/school"
	"github.com/classcharts-oss/server/pkg/util"

	"github.com/CommunityCharts/CCModels/shared"
	"github.com/CommunityCharts/CCModels/student"
	"github.com/classcharts-oss/server/pkg/db"
)

func GetAnnouncementsHandler(w http.ResponseWriter, r *http.Request) {
	sch, err := db.GetSchoolByID(r.Context().Value("student").(student.DBStudentUser).SchoolId)
	if err != nil {
		panic(err)
	}

	response := shared.NewSuccessfulResponse(sch.Announcements)
	response.Write(w)
}

func CreateAnnouncementHandler(w http.ResponseWriter, r *http.Request) {
	sch, err := db.GetSchoolByID(r.Context().Value("student").(student.DBStudentUser).SchoolId)
	if err != nil {
		panic(err)
	}

	announcement := MakeAnnouncement(
		sch,
		r.FormValue("title"),
		r.FormValue("content"),
		r.FormValue("teacher"),
		shared.No,
		shared.AnnouncementAttachment{},
	)

	response := shared.NewSuccessfulResponse(announcement)
	response.Write(w)
}

func MakeAnnouncement(sch school.School,
	title, content, teacher string,
	sticky shared.YesNoBool, attachments ...shared.AnnouncementAttachment) shared.Announcement {
	randomInt := util.RandomId()

	announcement := shared.Announcement{
		Id:          randomInt,
		Title:       title,
		Description: content,

		SchoolName:  sch.Name,
		TeacherName: teacher,
		SchoolLogo:  sch.Logo,

		Sticky: sticky,
		State:  "new",

		Timestamp: time.Now().Format(time.RFC3339),

		Attachments: attachments,

		CommentVisibility: "public",

		AllowComments:    shared.Yes,
		AllowReactions:   shared.Yes,
		AllowConsent:     shared.Yes,
		PriorityPinned:   shared.No,
		RequiresConsent:  shared.No,
		CanChangeConsent: false,
	}

	sch.Announcements = append(sch.Announcements, announcement)
	db.UpdateSchool(sch)

	return announcement
}
