package system

import (
	"net/http"

	"github.com/slyrx/gin_exam_system/server/model/common/response"
	"github.com/slyrx/gin_exam_system/server/model/system/request"
	systemResponse "github.com/slyrx/gin_exam_system/server/model/system/response"

	"github.com/gin-gonic/gin"
	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"github.com/slyrx/gin_exam_system/server/others/utils"
)

type ExamPaperStdApi struct{}

func (h *ExamPaperStdApi) GetExamPaperPageList(c *gin.Context) {
	var req request.PageListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, systemResponse.PageListResponse{Code: 0, Message: "请求参数错误"})
		return
	}

	subjectID := int(req.SubjectID)

	var total int64
	var papers []systemMod.ExamPaper_1

	// 查询总数
	if err := global.GES_DB.Debug().Model(&systemMod.ExamPaper_1{}).
		Where("paper_type = ? AND subject_id = ?", req.PaperType, subjectID).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, systemResponse.PageListResponse{Code: 0, Message: "查询失败"})
		return
	}

	// 分页查询
	offset := (req.PageIndex - 1) * req.PageSize
	if err := global.GES_DB.Debug().Where("paper_type = ? AND subject_id = ?", req.PaperType, subjectID).
		Limit(req.PageSize).Offset(offset).
		Find(&papers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, systemResponse.PageListResponse{Code: 0, Message: "查询失败"})
		return
	}

	var responseList []systemResponse.ExamPaper
	for _, paper := range papers {
		responseList = append(responseList, systemResponse.ExamPaper{
			ID:                 paper.ID,
			Name:               paper.Name,
			QuestionCount:      paper.QuestionCount,
			Score:              paper.Score,
			CreateTime:         paper.CreateTime,
			CreateUser:         paper.CreateUser,
			SubjectID:          paper.SubjectID,
			SubjectName:        "", // 这里可以通过SubjectID查询具体的学科名称并赋值
			PaperType:          paper.PaperType,
			FrameTextContentID: paper.FrameTextContentID,
		})
	}

	response.OkWithMessageExamInterface(systemResponse.ResponseData{
		Total:            int(total),
		List:             responseList,
		PageNum:          req.PageIndex,
		PageSize:         req.PageSize,
		Size:             len(papers),
		StartRow:         offset + 1,
		EndRow:           offset + len(papers),
		Pages:            (int(total) + req.PageSize - 1) / req.PageSize,
		PrePage:          req.PageIndex - 1,
		NextPage:         req.PageIndex + 1,
		IsFirstPage:      req.PageIndex == 1,
		IsLastPage:       req.PageIndex == (int(total)+req.PageSize-1)/req.PageSize,
		HasPreviousPage:  req.PageIndex > 1,
		HasNextPage:      req.PageIndex < (int(total)+req.PageSize-1)/req.PageSize,
		NavigatePages:    8,
		NavigatePageNums: []int{req.PageIndex},
	}, c)
}

func (h *ExamPaperStdApi) GetExamPaperPageList_1(c *gin.Context) {
	var req request.PageListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, systemResponse.PageListResponse{Code: 0, Message: "请求参数错误"})
		return
	}

	subjectID := int(req.SubjectID)

	// 获取当前学生ID
	// studentID, exists := c.Get("userID")
	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, systemResponse.PageListResponse{Code: 0, Message: "未授权"})
	// 	return
	// }
	studentID := examService.GetUserInfo(utils.GetCurrentUser(c).UserName).ID // 假设这个函数已经实现

	var total int64
	var papers []systemMod.ExamPaper_1

	// 修改查询逻辑，加入 t_exam_paper_assignment 关系
	query := global.GES_DB.Debug().Table("t_exam_paper ep").
		Joins("JOIN t_exam_paper_assignment epa ON ep.id = epa.exam_paper_id").
		Where("ep.paper_type = ? AND ep.subject_id = ? AND epa.student_id = ?", req.PaperType, subjectID, studentID)

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, systemResponse.PageListResponse{Code: 0, Message: "查询失败"})
		return
	}

	// 分页查询
	offset := (req.PageIndex - 1) * req.PageSize
	if err := query.Select("ep.*").
		Limit(req.PageSize).Offset(offset).
		Find(&papers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, systemResponse.PageListResponse{Code: 0, Message: "查询失败"})
		return
	}

	var responseList []systemResponse.ExamPaper
	for _, paper := range papers {
		responseList = append(responseList, systemResponse.ExamPaper{
			ID:                 paper.ID,
			Name:               paper.Name,
			QuestionCount:      paper.QuestionCount,
			Score:              paper.Score,
			CreateTime:         paper.CreateTime,
			CreateUser:         paper.CreateUser,
			SubjectID:          paper.SubjectID,
			SubjectName:        "", // 这里可以通过SubjectID查询具体的学科名称并赋值
			PaperType:          paper.PaperType,
			FrameTextContentID: paper.FrameTextContentID,
		})
	}

	response.OkWithMessageExamInterface(systemResponse.ResponseData{
		Total:            int(total),
		List:             responseList,
		PageNum:          req.PageIndex,
		PageSize:         req.PageSize,
		Size:             len(papers),
		StartRow:         offset + 1,
		EndRow:           offset + len(papers),
		Pages:            (int(total) + req.PageSize - 1) / req.PageSize,
		PrePage:          req.PageIndex - 1,
		NextPage:         req.PageIndex + 1,
		IsFirstPage:      req.PageIndex == 1,
		IsLastPage:       req.PageIndex == (int(total)+req.PageSize-1)/req.PageSize,
		HasPreviousPage:  req.PageIndex > 1,
		HasNextPage:      req.PageIndex < (int(total)+req.PageSize-1)/req.PageSize,
		NavigatePages:    8,
		NavigatePageNums: []int{req.PageIndex},
	}, c)
}
