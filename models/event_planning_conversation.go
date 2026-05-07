package models

import (
	"time"

	"hotbox-adm-backend/cli"
	"hotbox-adm-backend/form"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// 活动策划gpt对话
type EventPlanningConversationModel struct {
	ID             uint           `gorm:"primarykey;type:int(0)" json:"id" redis:"id"`
	InputCondition datatypes.JSON `gorm:"input_condition" json:"input_condition"`
	Content        string         `gorm:"content" json:"content"`
	UserId         int            `gorm:"user_id" json:"user_id"`
	Prompt         string         `gorm:"prompt" json:"prompt"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (q *EventPlanningConversationModel) TableName() string {
	return "event_planning_conversation"
}

type EventPlanningConversation struct {
	Ctx *gin.Context
}

func (q EventPlanningConversation) GetEventPlanningConversationList(req form.GetEventPlanningConversationListReq, userId int, order []string) (list *[]EventPlanningConversationModel, total int64, err error) {
	query := cli.HotDogGormDB.WithContext(q.Ctx).Model(&EventPlanningConversationModel{})

	if userId > 0 {
		query.Where("user_id", userId)
	}
	err = query.Count(&total).Error
	if err != nil {
		logrus.Errorf("GetEventPlanningConversationList query count error: %v", err)
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, nil
	}
	if req.PageNumber != 0 && req.PageSize != 0 {
		query.Offset((int(req.PageNumber) - 1) * req.PageSize).Limit(req.PageSize)
	}
	for _, v := range order {
		query = query.Order(v)
	}
	err = query.Scan(&list).Error
	if err != nil {
		logrus.Errorf("GetEventPlanningConversationList error: %v", err)
		return nil, 0, err
	}
	return
}

func (q EventPlanningConversation) CreateEventPlanningConversation(userId int, body form.CreateEventPlanningConversationListReq) (dm EventPlanningConversationModel, err error) {
	dm = EventPlanningConversationModel{
		Content:        body.Content,
		InputCondition: datatypes.JSON(body.InputCondition),
		UserId:         userId,
		Prompt:         body.Prompt,
	}
	err = cli.HotDogGormDB.WithContext(q.Ctx).Save(&dm).Error
	return dm, err
}
