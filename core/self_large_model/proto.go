package self_large_model

type ChatMessageRequest struct {
	Message string `json:"message"`
	ReChat  bool   `json:"re_chat"`
	Stream  bool   `json:"stream"`
}
type ChatMessageResp struct {
	Code            int              `json:"code"`
	Message         string           `json:"message"`
	ChatMessageData *ChatMessageData `json:"data"`
}

type ChatMessageData struct {
	ChatId  string `json:"chat_id"`
	Id      string `json:"id"`
	Operate bool   `json:"operate"`
	Content string `json:"content"`
	IsEnd   bool   `json:"is_end"`
}

type ProfileResp struct {
	Code        int          `json:"code"`
	Message     string       `json:"message"`
	ProfileData *ProfileData `json:"data"`
}

type ProfileData struct {
	Id                     string `json:"id"`
	Name                   string `json:"name"`
	Desc                   string `json:"desc"`
	Prologue               string `json:"prologue"`
	Icon                   string `json:"icon"`
	ShowSource             bool   `json:"show_source"`
	MultipleRoundsDialogue bool   `json:"multiple_rounds_dialogue"`
}

type ChatOpen struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}
