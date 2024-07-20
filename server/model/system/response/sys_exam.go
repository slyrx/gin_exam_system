package response

type AnswerSubmitResponse struct {
	Code     int    `json:"code"`
	Response string `json:"response"`
	Message  string `json:"message"`
}
