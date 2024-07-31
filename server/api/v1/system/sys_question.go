package system

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system/request"
	"github.com/slyrx/gin_exam_system/server/others/global"

	// "github.com/slyrx/gin_exam_system/server/others/utils"
	"go.uber.org/zap"
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
	// count, err := questionService.GetCountQuestionsBySubject(query.SubjectID)
	// if err != nil {
	// 	response.FailWithMessage(err.Error(), c)
	// 	return
	// }
	path := c.Request.URL.Path
	// waitUse, _ := utils.GetClaims(c)
	obj := strings.TrimPrefix(path, global.GES_CONFIG.System.RouterPrefix)
	// 获取请求方法
	act := c.Request.Method
	// 获取用户的角色
	// sub := strconv.Itoa(int(waitUse.AuthorityId))
	global.GES_LOG.Info("CasbinHandler", zap.Any("obj", obj), zap.Any("act", act))
	fmt.Println("ccc0", obj, act)

	// 查询题目
	questions, totalCount, err := questionService.GetQuestionsBySubject(query.SubjectID, query.PageSize, query.PageIndex)
	global.GES_LOG.Info("questions", zap.Any("questions", questions))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 查询题目内容
	total := totalCount // 总记录数

	respQuestions := questionService.MapSourceToTargetQuestions(questions)
	global.GES_LOG.Info("exam8", zap.Any("total", total))
	pageResult := questionService.MapToPageQuestionResult(questions, respQuestions, int(total), query.PageIndex, query.PageSize)
	// 返回结果
	response.OkWithDetailedExam(pageResult, c)
}

func (e *QuestionApi) CreateUserQuestionTable(c *gin.Context) {
	err := questionService.CreateUserQuestionTable()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建用户做题表成功", c)
}
