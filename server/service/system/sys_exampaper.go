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
