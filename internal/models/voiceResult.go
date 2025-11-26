package models

type VoiceResult struct {
	QID    string `json:"qid"`
	Result struct {
		Texts []struct {
			Text           string  `json:"text"`
			Confidence     float64 `json:"confidence"`
			PunctuatedText string  `json:"punctuated_text"`
		} `json:"texts"`
	} `json:"result"`
}
