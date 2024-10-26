package user

import (
	"net/http"

	"github.com/CommunityCharts/CCModels/shared"
	"github.com/CommunityCharts/CCModels/student"
	"github.com/classcharts-oss/server/pkg/db"
)

type userResponseMeta struct {
	Version   *string `json:"version,omitempty"`
	SessionId *string `json:"session_id,omitempty"`
}

const version = "20.0.1"

func StudentUserHandler(w http.ResponseWriter, r *http.Request) {
	var sessionId *string
	s := r.Context().Value("student").(student.DBStudentUser)

	if r.FormValue("include_data") != "false" {
		sessionId = db.GetStudentJWTForLogin(s)
	}

	version := version

	meta := userResponseMeta{
		Version:   &version,
		SessionId: sessionId,
	}

	data := shared.Object{
		"user": s.StudentUser,
	}

	response := shared.NewSuccessfulMetaResponse(data, meta)
	response.Write(w)
}
