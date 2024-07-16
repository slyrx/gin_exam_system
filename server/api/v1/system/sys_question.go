package system

import (
	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system/request"
)

type QuestionApi struct{}

func (e *QuestionApi) GetPageInfo(c *gin.Context) {
	var query systemMod.QuestionQuery
	err := c.ShouldBindJSON(&query)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 查询题目
	count, err := questionService.GetCountQuestionsBySubject(query.SubjectID)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 查询题目
	questions, err := questionService.GetQuestionsBySubject(query.SubjectID, query.PageSize, query.PageIndex)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 查询题目内容
	total := count // 总记录数

	respQuestions := questionService.MapSourceToTargetQuestions(questions)
	pageResult := questionService.MapToPageQuestionResult(questions, respQuestions, int(total), query.PageIndex, query.PageSize)
	// 返回结果
	response.OkWithDetailedExam(pageResult, c)
}
