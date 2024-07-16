package system

import (
	"encoding/json"
	"strconv"

	"github.com/slyrx/gin_exam_system/server/model/common/response"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
)

type QuestionService struct{}

var QuestionServiceApp = new(QuestionService)

func (questionService *QuestionService) GetCountQuestionsBySubject(subjectID int) (int64, error) {
	var count int64
	err := global.GES_DB.Debug().Model(&systemMod.Question{}).
		Where("deleted = ? AND subject_id = ?", false, subjectID).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (questionService *QuestionService) GetQuestionsBySubject(subjectID int, limit int, pageIndex int) ([]systemMod.Question, error) {
	var questions []systemMod.Question
	offset := (pageIndex - 1) * limit
	err := global.GES_DB.Debug().
		Where("deleted = ? AND subject_id = ?", false, subjectID).
		Order("err_count desc").
		Limit(limit).
		Offset(offset).
		Find(&questions).
		Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (questionService *QuestionService) MapToPageQuestionResult(questions []systemMod.Question, respQuestions []response.Question, total int, pageNum int, pageSize int) response.PageQuestionResult {
	result := response.PageQuestionResult{
		Total:    total,
		List:     respQuestions,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	// 计算其他分页相关字段
	result.Pages = (total + pageSize - 1) / pageSize
	result.Size = len(questions)
	result.StartRow = (pageNum-1)*pageSize + 1
	result.EndRow = result.StartRow + result.Size - 1

	// 设置导航页码
	result.NavigatePages = 8 // 假设我们想显示8个导航页码
	result.NavigatePageNums = QuestionServiceApp.calculateNavigatePageNums(pageNum, result.Pages, result.NavigatePages)
	result.NavigateFirstPage = result.NavigatePageNums[0]
	result.NavigateLastPage = result.NavigatePageNums[len(result.NavigatePageNums)-1]

	// 设置页面状态
	result.IsFirstPage = pageNum == 1
	result.IsLastPage = pageNum == result.Pages
	result.HasPreviousPage = pageNum > 1
	result.HasNextPage = pageNum < result.Pages

	// 设置上一页和下一页
	if result.HasPreviousPage {
		result.PrePage = pageNum - 1
	}
	if result.HasNextPage {
		result.NextPage = pageNum + 1
	}

	return result
}

// 计算导航页码
func (questionService *QuestionService) calculateNavigatePageNums(pageNum, totalPages, navigatePages int) []int {
	var navigatePageNums []int
	startNum := pageNum - navigatePages/2
	endNum := pageNum + navigatePages/2

	if startNum < 1 {
		startNum = 1
		endNum = navigatePages
	}

	if endNum > totalPages {
		endNum = totalPages
		startNum = totalPages - navigatePages + 1
		if startNum < 1 {
			startNum = 1
		}
	}

	for i := startNum; i <= endNum; i++ {
		navigatePageNums = append(navigatePageNums, i)
	}

	return navigatePageNums
}

func (questionService *QuestionService) MapSourceToTargetQuestions(srcQuestions []systemMod.Question) []response.Question {
	targetQuestions := make([]response.Question, len(srcQuestions))
	for i, srcQuestion := range srcQuestions {
		targetQuestions[i] = QuestionServiceApp.MapSourceToTargetQuestion(srcQuestion)
	}
	return targetQuestions
}

func (questionService *QuestionService) MapSourceToTargetQuestion(src systemMod.Question) response.Question {
	global.GES_LOG.Info("exam", zap.Any("src questionIds", src))
	dst := response.Question{
		ID:           src.ID,
		QuestionType: src.QuestionType,
		SubjectID:    src.SubjectID,
		CreateTime:   src.CreateTime,
		Correct:      src.Correct,
	}

	if src.InfoTextContentID != 0 {
		dst.TextContentID = &src.InfoTextContentID
	}

	createUser, err := strconv.Atoi(src.CreateUser)
	if err == nil {
		dst.CreateUser = createUser
	}

	dst.Score = strconv.Itoa(src.Score)

	status, err := strconv.Atoi(src.Status)
	if err == nil {
		dst.Status = status
	}

	difficult, err := strconv.Atoi(src.Difficult)
	if err == nil {
		dst.Difficult = difficult
	}

	dst.AnalyzeTextContentID = nil
	var frameTextContent string
	frameTextContent, err = QuestionServiceApp.selectTextContentByID(src.InfoTextContentID)
	var questionObject systemMod.QuestionObject
	if err != nil {
		dst.ShortTitle = ""
	} else {
		json.Unmarshal([]byte(frameTextContent), &questionObject)
		dst.ShortTitle = questionObject.TitleContent
		global.GES_LOG.Info("ShortTitle", zap.Any("ShortTitle", frameTextContent))
	}

	return dst
}

func (questionService *QuestionService) selectTextContentByID(id int) (s string, err error) {
	examPaperTextContent := &systemMod.ExamPaperTextContent{}
	err = global.GES_DB.Debug().Where("id = ?", id).First(examPaperTextContent).Error
	s = examPaperTextContent.Content
	return
}
