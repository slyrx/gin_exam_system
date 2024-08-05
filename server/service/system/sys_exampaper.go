package system

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/model/system/request"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ExamPaperService struct{}

var ExamPaperServiceApp = new(ExamPaperService)

func (s *ExamPaperService) CreateExamPaper(userId int, req request.CreateExamPaperRequest) (int, error) {
	var paperID int
	err := global.GES_DB.Debug().Transaction(func(tx *gorm.DB) error {
		// 1. 创建试卷记录
		examPaper := systemMod.ExamPaper_1{
			SubjectID:  req.SubjectID,
			PaperType:  req.PaperType,
			Name:       req.Name,
			CreateUser: userId, // 设置为当前用户ID
			CreateTime: time.Now(),
			GradeLevel: req.Level, // 如果需要，设置年级等级
			Deleted:    []byte{0},
		}

		// 解析 SuggestTime
		suggestTime, err := s.parseSuggestTime(req.SuggestTime)
		if err != nil {
			return err
		}

		examPaper.SuggestTime = suggestTime

		questionCount := 0
		totalScore := 0
		for _, titleItem := range req.TitleItems {
			questionCount += len(titleItem.QuestionItems)
			for _, question := range titleItem.QuestionItems {
				score, _ := strconv.Atoi(question.Score)
				totalScore += score
			}
		}
		examPaper.QuestionCount = questionCount
		examPaper.Score = totalScore

		if len(req.LimitDateTime) == 2 {
			startTime, err := time.Parse(time.RFC3339, req.LimitDateTime[0])
			if err == nil {
				examPaper.LimitStartTime = &startTime
			}
			endTime, err := time.Parse(time.RFC3339, req.LimitDateTime[1])
			if err == nil {
				examPaper.LimitEndTime = &endTime
			}
		}

		// 2. 保存试卷结构
		now := time.Now()
		frameTextContent := systemMod.ExamPaperTextContent1{
			Content:    s.generateFrameTextContent(req),
			CreateTime: &now,
		}
		if err := tx.Create(&frameTextContent).Error; err != nil {
			return err
		}

		examPaper.FrameTextContentID = frameTextContent.ID
		global.GES_LOG.Info("exam4", zap.Any("frameTextContent.ID", frameTextContent.ID))
		// 3. 更新试卷记录
		// if err := tx.Save(&examPaper).Error; err != nil {
		// 	return err
		// }

		if req.ID == nil || *req.ID == 0 {
			// ID 为空，创建新记录
			if err := tx.Omit("ID").Create(&examPaper).Error; err != nil {
				return err
			}
		} else {
			// ID 不为空，更新现有记录
			examPaper.ID = *req.ID
			if err := tx.Save(&examPaper).Error; err != nil {
				return err
			}
		}

		// 4. 创建题目到 t_question 表
		// for _, titleItem := range req.TitleItems {
		// 	for _, questionItem := range titleItem.QuestionItems {
		// 		question := systemMod.Question_1{
		// 			QuestionType: questionItem.QuestionType,
		// 			SubjectID:    questionItem.SubjectID,
		// 			GradeLevel:   questionItem.GradeLevel,
		// 			Difficult:    questionItem.Difficult,
		// 			CreateUser:   userId, // 需要从上下文或请求中获取
		// 			Status:       1,      // 假设 1 表示正常状态
		// 			CreateTime:   time.Now(),
		// 		}

		// 		// 转换 Score 为整数
		// 		score, err := strconv.Atoi(questionItem.Score)
		// 		if err != nil {
		// 			log.Printf("Error converting score to integer for question %d: %v", questionItem.ID, err)
		// 			score = 0 // 或者使用其他默认值
		// 		}
		// 		question.Score = score

		// 		// 创建 TextContent 结构
		// 		textContent := request.TextContent{
		// 			TitleContent:        questionItem.Title,
		// 			Analyze:             questionItem.Analyze,
		// 			QuestionItemObjects: make([]request.QuestionItemObject, len(questionItem.Items)),
		// 			Correct:             questionItem.Correct,
		// 		}

		// 		// 填充 QuestionItemObjects
		// 		for i, item := range questionItem.Items {
		// 			textContent.QuestionItemObjects[i] = request.QuestionItemObject{
		// 				Prefix:   item.Prefix,
		// 				Content:  item.Content,
		// 				Score:    item.Score,
		// 				ItemUUID: item.ItemUUID,
		// 			}
		// 		}

		// 		// 将 TextContent 转换为 JSON
		// 		contentJSON, err := json.Marshal(textContent)
		// 		if err != nil {
		// 			return err
		// 		}

		// 		// 创建 t_text_content 记录
		// 		textContentRecord := systemMod.ExamPaperTextContent1{
		// 			Content: string(contentJSON),
		// 		}
		// 		if err := tx.Create(&textContentRecord).Error; err != nil {
		// 			return err
		// 		}

		// 		question.InfoTextContentID = textContentRecord.ID

		// 		// 创建题目
		// 		if err := tx.Create(&question).Error; err != nil {
		// 			return err
		// 		}
		// 	}
		// }

		paperID = examPaper.ID
		return nil
	})

	return paperID, err
}

func (s *ExamPaperService) generateFrameTextContent(req request.CreateExamPaperRequest) string {
	var frameContent []struct {
		Name          string `json:"name"`
		QuestionItems []struct {
			ID        int `json:"id"`
			ItemOrder int `json:"itemOrder"`
		} `json:"questionItems"`
	}

	for _, titleItem := range req.TitleItems {
		item := struct {
			Name          string `json:"name"`
			QuestionItems []struct {
				ID        int `json:"id"`
				ItemOrder int `json:"itemOrder"`
			} `json:"questionItems"`
		}{
			Name: titleItem.Name,
		}

		for i, questionItem := range titleItem.QuestionItems {
			item.QuestionItems = append(item.QuestionItems, struct {
				ID        int `json:"id"`
				ItemOrder int `json:"itemOrder"`
			}{
				ID:        questionItem.ID,
				ItemOrder: i + 1, // 假设题目顺序从1开始
			})
		}

		frameContent = append(frameContent, item)
	}

	jsonContent, err := json.Marshal(frameContent)
	if err != nil {
		// 处理错误，可能需要记录日志或返回一个错误
		global.GES_LOG.Info("exam", zap.Any("generateFrameTextContent", "generateFrameTextContent"))
		return "[]"
	}

	return string(jsonContent)
}

func (e *ExamPaperService) parseSuggestTime(suggestTime interface{}) (int, error) {
	switch v := suggestTime.(type) {
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case int:
		return v, nil
	default:
		return 0, errors.New("无效的 SuggestTime 类型")
	}
}

func (e *ExamPaperService) CreateErrorQuestionPaper(subjectID int, gradeLevel int, createUserID int) (int, error) {
	var paperID int
	err := global.GES_DB.Debug().Transaction(func(tx *gorm.DB) error {
		// 1. 获取错题
		var errorQuestions []struct {
			QuestionID int
			ErrorCount int
		}
		if err := tx.Table("t_question_error_count").
			Select("question_id, error_count").
			Where("error_count > ?", 5).
			Order("error_count DESC").
			Limit(20).
			Scan(&errorQuestions).Error; err != nil {
			return err
		}

		// 2. 创建试卷
		examPaper := systemMod.ExamPaper_1{
			Name:        "错题重组试卷",
			SubjectID:   subjectID,
			PaperType:   1, // 假设 6 表示错题试卷
			GradeLevel:  gradeLevel,
			CreateUser:  createUserID,
			CreateTime:  time.Now(),
			SuggestTime: 10,
			Deleted:     []byte{0},
		}
		if err := tx.Create(&examPaper).Error; err != nil {
			return err
		}

		// 3. 获取题目详情并生成 JSON

		var titleItems []systemMod.TitleItem
		var totalScore int
		itemOrder := 1

		for _, eq := range errorQuestions {
			var question systemMod.Question_1
			if err := tx.First(&question, eq.QuestionID).Error; err != nil {
				return err
			}

			var textContent systemMod.ExamPaperTextContent1
			if err := tx.First(&textContent, question.InfoTextContentID).Error; err != nil {
				return err
			}

			questionItem := systemMod.QuestionItem_1{
				ID:        question.ID,
				ItemOrder: itemOrder,
			}

			// 假设所有题目属于同一个 TitleItem
			if len(titleItems) == 0 {
				titleItems = append(titleItems, systemMod.TitleItem{
					Name:          "错题重组",
					QuestionItems: []systemMod.QuestionItem_1{questionItem},
				})
			} else {
				titleItems[0].QuestionItems = append(titleItems[0].QuestionItems, questionItem)
			}

			totalScore += question.Score
			itemOrder++
		}

		// 4. 存储试卷内容
		now := time.Now()
		questionJson, err := json.Marshal(titleItems)
		if err != nil {
			return err
		}

		frameTextContent := systemMod.ExamPaperTextContent1{
			Content:    string(questionJson),
			CreateTime: &now,
		}
		if err := tx.Create(&frameTextContent).Error; err != nil {
			return err
		}

		// 5. 更新试卷信息
		examPaper.Score = totalScore
		examPaper.QuestionCount = itemOrder - 1
		examPaper.FrameTextContentID = frameTextContent.ID
		if err := tx.Save(&examPaper).Error; err != nil {
			return err
		}
		paperID = examPaper.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return paperID, nil
}

func (e *ExamPaperService) CreateErrorQuestionPaperByUser(subjectID int, gradeLevel int, createUserID int, examUserID int) (int, error) {
	var paperID int
	err := global.GES_DB.Debug().Transaction(func(tx *gorm.DB) error {
		// 1. 获取错题
		var errorQuestions []struct {
			QuestionID int
			ErrorCount int
		}
		if err := tx.Table("t_user_wrong_book").
			Select("question_id, error_count").
			Where("error_count > ? AND user_id = ?", 5, examUserID).
			Order("error_count DESC").
			Limit(20).
			Scan(&errorQuestions).Error; err != nil {
			return err
		}

		// 2. 创建试卷
		examPaper := systemMod.ExamPaper_1{
			Name:        "错题重组试卷",
			SubjectID:   subjectID,
			PaperType:   1, // 假设 6 表示错题试卷
			GradeLevel:  gradeLevel,
			CreateUser:  createUserID,
			CreateTime:  time.Now(),
			SuggestTime: 10,
			Deleted:     []byte{0},
		}
		if err := tx.Create(&examPaper).Error; err != nil {
			return err
		}

		// 3. 获取题目详情并生成 JSON

		var titleItems []systemMod.TitleItem
		var totalScore int
		itemOrder := 1

		for _, eq := range errorQuestions {
			var question systemMod.Question_1
			if err := tx.First(&question, eq.QuestionID).Error; err != nil {
				return err
			}

			var textContent systemMod.ExamPaperTextContent1
			if err := tx.First(&textContent, question.InfoTextContentID).Error; err != nil {
				return err
			}

			questionItem := systemMod.QuestionItem_1{
				ID:        question.ID,
				ItemOrder: itemOrder,
			}

			// 假设所有题目属于同一个 TitleItem
			if len(titleItems) == 0 {
				titleItems = append(titleItems, systemMod.TitleItem{
					Name:          "错题重组",
					QuestionItems: []systemMod.QuestionItem_1{questionItem},
				})
			} else {
				titleItems[0].QuestionItems = append(titleItems[0].QuestionItems, questionItem)
			}

			totalScore += question.Score
			itemOrder++
		}

		// 4. 存储试卷内容
		now := time.Now()
		questionJson, err := json.Marshal(titleItems)
		if err != nil {
			return err
		}

		frameTextContent := systemMod.ExamPaperTextContent1{
			Content:    string(questionJson),
			CreateTime: &now,
		}
		if err := tx.Create(&frameTextContent).Error; err != nil {
			return err
		}

		// 5. 更新试卷信息
		examPaper.Score = totalScore
		examPaper.QuestionCount = itemOrder - 1
		examPaper.FrameTextContentID = frameTextContent.ID
		if err := tx.Save(&examPaper).Error; err != nil {
			return err
		}
		paperID = examPaper.ID
		return nil
	})

	if err != nil {
		return 0, err
	}

	return paperID, nil
}

func (e *ExamPaperService) AssignPaperVisibility(req request.AssignPaperVisibilityRequest, teacherID int) error {
	// 开始事务
	tx := global.GES_DB.Debug().Begin()

	for _, userID := range req.UserIDs {
		visibility := systemMod.PaperVisibility{
			PaperID:   int(req.PaperID),
			UserID:    int(userID),
			CreatedBy: teacherID,
			CreatedAt: time.Now(),
		}

		if err := tx.Create(&visibility).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (e *ExamPaperService) AssignExamPaperToStudent(examPaperID uint, studentID uint, assignedBy uint) error {
	assignment := &systemMod.ExamPaperAssignment{
		ExamPaperID: examPaperID,
		StudentID:   studentID,
		AssignedBy:  assignedBy,
		CreatedAt:   time.Now(),
	}
	return global.GES_DB.Create(assignment).Error
}
