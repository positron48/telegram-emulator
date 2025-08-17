package models

// ReplyKeyboardMarkup представляет обычную клавиатуру
type ReplyKeyboardMarkup struct {
	Keyboard        [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool               `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard bool               `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder string       `json:"input_field_placeholder,omitempty"`
	Selective       bool               `json:"selective,omitempty"`
}

// ReplyKeyboardRemove представляет удаление клавиатуры
type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
	Selective      bool `json:"selective,omitempty"`
}

// InlineKeyboardMarkup представляет inline клавиатуру
type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

// KeyboardButton представляет кнопку обычной клавиатуры
type KeyboardButton struct {
	Text            string `json:"text"`
	RequestContact  bool   `json:"request_contact,omitempty"`
	RequestLocation bool   `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo             `json:"web_app,omitempty"`
}

// InlineKeyboardButton представляет кнопку inline клавиатуры
type InlineKeyboardButton struct {
	Text                         string `json:"text"`
	URL                          string `json:"url,omitempty"`
	CallbackData                 string `json:"callback_data,omitempty"`
	WebApp                       *WebAppInfo `json:"web_app,omitempty"`
	LoginURL                     *LoginURL   `json:"login_url,omitempty"`
	SwitchInlineQuery            string      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat string      `json:"switch_inline_query_current_chat,omitempty"`
	CallbackGame                 interface{} `json:"callback_game,omitempty"`
	Pay                          bool        `json:"pay,omitempty"`
}

// KeyboardButtonPollType представляет тип опроса для кнопки клавиатуры
type KeyboardButtonPollType struct {
	Type string `json:"type,omitempty"` // "quiz" или "regular"
}

// WebAppInfo представляет информацию о веб-приложении
type WebAppInfo struct {
	URL string `json:"url"`
}

// LoginURL представляет URL для авторизации
type LoginURL struct {
	URL                string `json:"url"`
	ForwardText        string `json:"forward_text,omitempty"`
	BotUsername        string `json:"bot_username,omitempty"`
	RequestWriteAccess bool   `json:"request_write_access,omitempty"`
}

// ReplyMarkup представляет общий тип разметки ответа
type ReplyMarkup interface {
	// Интерфейс для всех типов разметки
}

// IsReplyKeyboardMarkup проверяет, является ли разметка обычной клавиатурой
func (r *ReplyKeyboardMarkup) IsReplyKeyboardMarkup() bool {
	return true
}

// IsReplyKeyboardRemove проверяет, является ли разметка удалением клавиатуры
func (r *ReplyKeyboardRemove) IsReplyKeyboardRemove() bool {
	return true
}

// IsInlineKeyboardMarkup проверяет, является ли разметка inline клавиатурой
func (r *InlineKeyboardMarkup) IsInlineKeyboardMarkup() bool {
	return true
}
