package system

import (
	"time"

	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
)

type ExamPaperStdService struct{}

var ExamPaperStdServiceApp = new(ExamPaperStdService)

func (s *ExamPaperStdService) AssignExamPaperToStudent(examPaperID uint, studentID uint, assignedBy uint) error {
	assignment := &systemMod.ExamPaperAssignment{
		ExamPaperID: examPaperID,
		StudentID:   studentID,
		AssignedBy:  assignedBy,
		CreatedAt:   time.Now(),
	}
	return global.GES_DB.Create(assignment).Error
}
