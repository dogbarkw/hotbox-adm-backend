package form

type GetEventPlanningConversationListReq struct {
	PagingReq
}

type CreateEventPlanningConversationListReq struct {
	InputCondition string `json:"input_condition" binding:"required"`
	Content        string `gorm:"content" json:"content" binding:"required"`
	Prompt         string `json:"prompt" `
}
