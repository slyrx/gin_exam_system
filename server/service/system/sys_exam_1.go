package system

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	systemMod "github.com/slyrx/gin_exam_system/server/model/system"
	"github.com/slyrx/gin_exam_system/server/model/system/request"
	"github.com/slyrx/gin_exam_system/server/others/global"
	"gorm.io/gorm"
)

func (s *ExamService) SubmitAnswer1(req request.AnswerSubmitRequest, userID int) (int, error) {
	var totalScore int
	err := global.GES_DB.Transaction(func(tx *gorm.DB) error {
		// 创建试卷答案记录
		examPaperAnswer := systemMod.ExamPaperAnswer_1{
			ID:         req.ID,
			DoTime:     req.DoTime,
			CreateTime: time.Now(),
			Status:     1, // 待判分状态
			CreateUser: userID,
		}
		if err := tx.Create(&examPaperAnswer).Error; err != nil {
			return err
		}

		correctCount := 0

		// 处理每个答案项
		for _, item := range req.AnswerItems {
			var question systemMod.Question_1
			if err := tx.First(&question, item.QuestionID).Error; err != nil {
				return err
			}

			isCorrect := item.Content == question.Correct
			score := 0
			if isCorrect {
				score = question.Score
				correctCount++
			}
			totalScore += score

			// 创建题目答案记录
			answer := systemMod.ExamPaperQuestionCustomerAnswer_1{
				QuestionID:        item.QuestionID,
				ExamPaperID:       req.ID,
				ExamPaperAnswerID: examPaperAnswer.ID,
				Answer:            &item.Content,
				DoRight:           isCorrect,
				CustomerScore:     score,
				QuestionScore:     question.Score,
				CreateUser:        userID,
				CreateTime:        time.Now(),
				ItemOrder:         item.ItemOrder,
			}
			if err := tx.Create(&answer).Error; err != nil {
				return err
			}

			// 如果答错，更新错题统计
			if !isCorrect {
				if err := s.updateErrorCount(tx, item.QuestionID); err != nil {
					return err
				}
			}
		}

		// 更新试卷答案记录
		examPaperAnswer.UserScore = totalScore
		examPaperAnswer.QuestionCorrect = correctCount
		examPaperAnswer.QuestionCount = len(req.AnswerItems)
		examPaperAnswer.Status = 2 // 完成状态
		if err := tx.Save(&examPaperAnswer).Error; err != nil {
			return err
		}

		return nil
	})

	return totalScore, err
}

func (s *ExamService) updateErrorCount2(tx *gorm.DB, questionID int) error {
	var errorCount systemMod.QuestionErrorCount
	result := tx.FirstOrCreate(&errorCount, systemMod.QuestionErrorCount{QuestionID: questionID})
	if result.Error != nil {
		return result.Error
	}
	errorCount.ErrorCount++
	errorCount.LastErrorTime = time.Now()
	return tx.Save(&errorCount).Error
}

func (s *ExamService) SubmitAnswer(req request.AnswerSubmitRequest, userID int) (totalScore, doTime int, paperName string, err error) {
	err = global.GES_DB.Debug().Transaction(func(tx *gorm.DB) error {
		// 1. 查找试卷模板
		var examPaper systemMod.ExamPaper_1
		if err := tx.First(&examPaper, req.ID).Error; err != nil {
			return err
		}

		// 2. 获取试卷内容
		var textContent systemMod.ExamPaperTextContent1
		if err := tx.First(&textContent, examPaper.FrameTextContentID).Error; err != nil {
			return err
		}

		// 3. 解析试卷内容
		var paperStructWrapper systemMod.PaperStructWrapper
		if err := json.Unmarshal([]byte(textContent.Content), &paperStructWrapper); err != nil {
			return err
		}

		if len(paperStructWrapper) == 0 {
			return errors.New("no paper structure found")
		}

		paperStruct := paperStructWrapper[0]

		// 4. 获取题目ID列表
		questionIDs := make([]int, len(paperStruct.QuestionItems))
		for i, item := range paperStruct.QuestionItems {
			questionIDs[i] = item.ID
		}

		// 5. 查询题目内容
		var questions []systemMod.Question_1
		if err := tx.Where("id IN ?", questionIDs).Find(&questions).Error; err != nil {
			return err
		}

		// 7-1. 创建或更新试卷答案记录
		examPaperAnswer := systemMod.ExamPaperAnswer_1{
			ExamPaperID: examPaper.ID,
			PaperName:   examPaper.Name,
			PaperType:   examPaper.PaperType,
			SubjectID:   examPaper.SubjectID,
			PaperScore:  examPaper.Score,
			DoTime:      req.DoTime,
			Status:      1, // 假设1表示进行中状态
			CreateUser:  userID,
			CreateTime:  time.Now(),
		}
		if err := tx.Create(&examPaperAnswer).Error; err != nil {
			return err
		}

		paperName = examPaperAnswer.PaperName
		doTime = examPaperAnswer.DoTime

		// 6. 判卷和统计
		correctCount := 0
		systemScore := 0
		for _, answerItem := range req.AnswerItems {
			var question systemMod.Question_1
			for _, q := range questions {
				if q.ID == answerItem.QuestionID {
					question = q
					break
				}
			}

			if question.ID == 0 {
				return errors.New(fmt.Sprintf("question not found: %d", answerItem.QuestionID))
			}

			isCorrect := answerItem.Content == question.Correct
			score := 0
			if isCorrect {
				score = question.Score
				correctCount++
			}
			systemScore += score

			// 创建题目答案记录
			answer := systemMod.ExamPaperQuestionCustomerAnswer_1{
				QuestionID:        answerItem.QuestionID,
				ExamPaperID:       req.ID,
				ExamPaperAnswerID: examPaperAnswer.ID,
				QuestionType:      question.QuestionType,
				SubjectID:         question.SubjectID,
				DoRight:           isCorrect,
				CustomerScore:     score,
				QuestionScore:     question.Score,
				CreateUser:        userID,
				CreateTime:        time.Now(),
				ItemOrder:         answerItem.ItemOrder,
			}

			// 仅在 answerItem.Content 不为空时设置 Answer 字段
			if answerItem.Content != "" {
				answer.Answer = &answerItem.Content
			}

			if err := tx.Create(&answer).Error; err != nil {
				return err
			}

			// 如果答错，更新错题统计
			if !isCorrect {
				if err := s.updateErrorCount(tx, answerItem.QuestionID); err != nil {
					return err
				}

				var wrongBook systemMod.UserWrongBook
				result := tx.Where(systemMod.UserWrongBook{UserID: userID, QuestionID: answerItem.QuestionID}).
					FirstOrCreate(&wrongBook)

				if result.Error != nil {
					return fmt.Errorf("failed to find or create wrong book entry: %w", result.Error)
				}

				// 如果是新创建的记录，result.RowsAffected 会等于 1
				if result.RowsAffected == 0 {
					// 记录已存在，增加错误计数
					wrongBook.ErrorCount++
				} else {
					// 新记录，设置初始值
					wrongBook.ExamPaperID = req.ID
					wrongBook.SubjectID = question.SubjectID
					wrongBook.CreateTime = time.Now()
				}

				wrongBook.UpdateTime = time.Now()

				if err := tx.Save(&wrongBook).Error; err != nil {
					return fmt.Errorf("failed to update wrong book: %w", err)
				}

			}
		}

		// 7. 创建或更新试卷答案记录
		examPaperAnswer.SystemScore = systemScore
		examPaperAnswer.UserScore = systemScore // 假设用户得分等于系统得分
		examPaperAnswer.QuestionCorrect = correctCount
		examPaperAnswer.QuestionCount = len(req.AnswerItems)
		examPaperAnswer.Status = 2 // 假设2表示完成状态

		if err := tx.Save(&examPaperAnswer).Error; err != nil {
			return err
		}

		totalScore = systemScore
		return nil
	})

	return
}

func (s *ExamService) updateErrorCount(tx *gorm.DB, questionID int) error {
	var errorCount systemMod.QuestionErrorCount
	result := tx.FirstOrCreate(&errorCount, systemMod.QuestionErrorCount{QuestionID: questionID})
	if result.Error != nil {
		return result.Error
	}

	// 如果是新创建的记录，设置 LastErrorTime 的值
	if result.RowsAffected == 1 {
		errorCount.LastErrorTime = time.Now()
	}

	errorCount.ErrorCount++
	errorCount.LastErrorTime = time.Now()
	return tx.Save(&errorCount).Error
}
