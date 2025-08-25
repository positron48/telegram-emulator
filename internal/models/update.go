package models

import (
	"time"
)

// Update представляет обновление от Telegram Bot API
type Update struct {
	UpdateID           int64               `json:"update_id"`
	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member,omitempty"`
	ChatMember         *ChatMemberUpdated  `json:"chat_member,omitempty"`
	ChatJoinRequest    *ChatJoinRequest    `json:"chat_join_request,omitempty"`
	Timestamp          time.Time           `json:"timestamp"`
}

// CallbackQuery представляет callback query от inline кнопок
type CallbackQuery struct {
	ID              string   `json:"id"`
	From            User     `json:"from"`
	Message         *Message `json:"message,omitempty"`
	InlineMessageID string   `json:"inline_message_id,omitempty"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data,omitempty"`
	GameShortName   string   `json:"game_short_name,omitempty"`
}

// InlineQuery представляет inline запрос
type InlineQuery struct {
	ID       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType string    `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

// ChosenInlineResult представляет выбранный inline результат
type ChosenInlineResult struct {
	ResultID        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageID string    `json:"inline_message_id,omitempty"`
	Query           string    `json:"query"`
}

// ShippingQuery представляет запрос на доставку
type ShippingQuery struct {
	ID              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

// PreCheckoutQuery представляет предварительный запрос на оплату
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             User       `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// Poll представляет опрос
type Poll struct {
	ID                    string          `json:"id"`
	Question              string          `json:"question"`
	Options               []PollOption    `json:"options"`
	TotalVoterCount       int             `json:"total_voter_count"`
	IsClosed              bool            `json:"is_closed"`
	IsAnonymous           bool            `json:"is_anonymous"`
	Type                  string          `json:"type"`
	AllowsMultipleAnswers bool            `json:"allows_multiple_answers"`
	CorrectOptionID       int             `json:"correct_option_id,omitempty"`
	Explanation           string          `json:"explanation,omitempty"`
	ExplanationEntities   []MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            int             `json:"open_period,omitempty"`
	CloseDate             int64           `json:"close_date,omitempty"`
}

// PollOption представляет опцию опроса
type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

// PollAnswer представляет ответ на опрос
type PollAnswer struct {
	PollID    string `json:"poll_id"`
	User      User   `json:"user"`
	OptionIDs []int  `json:"option_ids"`
}

// ChatMemberUpdated представляет обновление участника чата
type ChatMemberUpdated struct {
	Chat          Chat               `json:"chat"`
	From          User               `json:"from"`
	Date          int64              `json:"date"`
	OldChatMember TelegramChatMember `json:"old_chat_member"`
	NewChatMember TelegramChatMember `json:"new_chat_member"`
	InviteLink    *ChatInviteLink    `json:"invite_link,omitempty"`
}

// ChatJoinRequest представляет запрос на присоединение к чату
type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	From       User            `json:"from"`
	UserChatID int64           `json:"user_chat_id"`
	Date       int64           `json:"date"`
	Bio        string          `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

// TelegramChatMember представляет участника чата в формате Telegram Bot API
type TelegramChatMember struct {
	User                User   `json:"user"`
	Status              string `json:"status"`
	CustomTitle         string `json:"custom_title,omitempty"`
	IsAnonymous         bool   `json:"is_anonymous,omitempty"`
	CanBeEdited         bool   `json:"can_be_edited,omitempty"`
	CanManageChat       bool   `json:"can_manage_chat,omitempty"`
	CanDeleteMessages   bool   `json:"can_delete_messages,omitempty"`
	CanManageVideoChats bool   `json:"can_manage_video_chats,omitempty"`
	CanRestrictMembers  bool   `json:"can_restrict_members,omitempty"`
	CanPromoteMembers   bool   `json:"can_promote_members,omitempty"`
	CanChangeInfo       bool   `json:"can_change_info,omitempty"`
	CanInviteUsers      bool   `json:"can_invite_users,omitempty"`
	CanPostMessages     bool   `json:"can_post_messages,omitempty"`
	CanEditMessages     bool   `json:"can_edit_messages,omitempty"`
	CanPinMessages      bool   `json:"can_pin_messages,omitempty"`
	CanPostStories      bool   `json:"can_post_stories,omitempty"`
	CanEditStories      bool   `json:"can_edit_stories,omitempty"`
	CanDeleteStories    bool   `json:"can_delete_stories,omitempty"`
	UntilDate           int64  `json:"until_date,omitempty"`
}

// ChatInviteLink представляет ссылку-приглашение в чат
type ChatInviteLink struct {
	InviteLink              string `json:"invite_link"`
	Creator                 User   `json:"creator"`
	CreatesJoinRequest      bool   `json:"creates_join_request"`
	IsPrimary               bool   `json:"is_primary"`
	IsRevoked               bool   `json:"is_revoked"`
	Name                    string `json:"name,omitempty"`
	ExpireDate              int64  `json:"expire_date,omitempty"`
	MemberLimit             int    `json:"member_limit,omitempty"`
	MemberCount             int    `json:"member_count,omitempty"`
	PendingJoinRequestCount int    `json:"pending_join_request_count,omitempty"`
}

// Location представляет географическое местоположение
type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           int     `json:"live_period,omitempty"`
	Heading              int     `json:"heading,omitempty"`
	ProximityAlertRadius int     `json:"proximity_alert_radius,omitempty"`
}

// ShippingAddress представляет адрес доставки
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

// OrderInfo представляет информацию о заказе
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// MessageEntity представляет сущность в сообщении
type MessageEntity struct {
	Type          string `json:"type"`
	Offset        int    `json:"offset"`
	Length        int    `json:"length"`
	URL           string `json:"url,omitempty"`
	User          *User  `json:"user,omitempty"`
	Language      string `json:"language,omitempty"`
	CustomEmojiID string `json:"custom_emoji_id,omitempty"`
}

// TelegramMessage представляет сообщение в формате Telegram Bot API
type TelegramMessage struct {
	MessageID                    int64            `json:"message_id"`
	From                         TelegramUser     `json:"from"`
	SenderChat                   *TelegramChat    `json:"sender_chat,omitempty"`
	Date                         int64            `json:"date"`
	Chat                         TelegramChat     `json:"chat"`
	ForwardFrom                  *TelegramUser    `json:"forward_from,omitempty"`
	ForwardFromChat              *TelegramChat    `json:"forward_from_chat,omitempty"`
	ForwardFromMessageID         int64            `json:"forward_from_message_id,omitempty"`
	ForwardSignature             string           `json:"forward_signature,omitempty"`
	ForwardSenderName            string           `json:"forward_sender_name,omitempty"`
	ForwardDate                  int64            `json:"forward_date,omitempty"`
	IsAutomaticForward           bool             `json:"is_automatic_forward,omitempty"`
	ReplyToMessage               *TelegramMessage `json:"reply_to_message,omitempty"`
	ViaBot                       *TelegramUser    `json:"via_bot,omitempty"`
	EditDate                     int64            `json:"edit_date,omitempty"`
	HasProtectedContent          bool             `json:"has_protected_content,omitempty"`
	MediaGroupID                 string           `json:"media_group_id,omitempty"`
	AuthorSignature              string           `json:"author_signature,omitempty"`
	Text                         string           `json:"text,omitempty"`
	Entities                     []MessageEntity  `json:"entities,omitempty"`
	Animation                    interface{}      `json:"animation,omitempty"`
	Audio                        interface{}      `json:"audio,omitempty"`
	Document                     interface{}      `json:"document,omitempty"`
	Photo                        []interface{}    `json:"photo,omitempty"`
	Sticker                      interface{}      `json:"sticker,omitempty"`
	Story                        interface{}      `json:"story,omitempty"`
	Video                        interface{}      `json:"video,omitempty"`
	VideoNote                    interface{}      `json:"video_note,omitempty"`
	Voice                        interface{}      `json:"voice,omitempty"`
	Caption                      string           `json:"caption,omitempty"`
	CaptionEntities              []MessageEntity  `json:"caption_entities,omitempty"`
	HasMediaSpoiler              bool             `json:"has_media_spoiler,omitempty"`
	Contact                      interface{}      `json:"contact,omitempty"`
	Dice                         interface{}      `json:"dice,omitempty"`
	Game                         interface{}      `json:"game,omitempty"`
	Poll                         *Poll            `json:"poll,omitempty"`
	Venue                        interface{}      `json:"venue,omitempty"`
	Location                     *Location        `json:"location,omitempty"`
	NewChatMembers               []TelegramUser   `json:"new_chat_members,omitempty"`
	LeftChatMember               *TelegramUser    `json:"left_chat_member,omitempty"`
	NewChatTitle                 string           `json:"new_chat_title,omitempty"`
	NewChatPhoto                 []interface{}    `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto              bool             `json:"delete_chat_photo,omitempty"`
	GroupChatCreated             bool             `json:"group_chat_created,omitempty"`
	SupergroupChatCreated        bool             `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated           bool             `json:"channel_chat_created,omitempty"`
	MessageAutoDeleteTime        int              `json:"message_auto_delete_time,omitempty"`
	MigrateToChatID              int64            `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID            int64            `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                *TelegramMessage `json:"pinned_message,omitempty"`
	Invoice                      interface{}      `json:"invoice,omitempty"`
	SuccessfulPayment            interface{}      `json:"successful_payment,omitempty"`
	UserShared                   interface{}      `json:"user_shared,omitempty"`
	ChatShared                   interface{}      `json:"chat_shared,omitempty"`
	ConnectedWebsite             string           `json:"connected_website,omitempty"`
	WriteAccessAllowed           interface{}      `json:"write_access_allowed,omitempty"`
	PassportData                 interface{}      `json:"passport_data,omitempty"`
	ProximityAlertTriggered      interface{}      `json:"proximity_alert_triggered,omitempty"`
	ForumTopicCreated            interface{}      `json:"forum_topic_created,omitempty"`
	ForumTopicEdited             interface{}      `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed             interface{}      `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened           interface{}      `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden      interface{}      `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden    interface{}      `json:"general_forum_topic_unhidden,omitempty"`
	GiveawayCreated              interface{}      `json:"giveaway_created,omitempty"`
	Giveaway                     interface{}      `json:"giveaway,omitempty"`
	GiveawayWinners              interface{}      `json:"giveaway_winners,omitempty"`
	GiveawayCompleted            interface{}      `json:"giveaway_completed,omitempty"`
	VideoChatScheduled           interface{}      `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted             interface{}      `json:"video_chat_started,omitempty"`
	VideoChatEnded               interface{}      `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited interface{}      `json:"video_chat_participants_invited,omitempty"`
	WebAppData                   interface{}      `json:"web_app_data,omitempty"`
	ReplyMarkup                  interface{}      `json:"reply_markup,omitempty"`
}

// TelegramUser представляет пользователя в формате Telegram Bot API
type TelegramUser struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
}

// TelegramChat представляет чат в формате Telegram Bot API
type TelegramChat struct {
	ID                                 int64            `json:"id"`
	Type                               string           `json:"type"`
	Title                              string           `json:"title,omitempty"`
	Username                           string           `json:"username,omitempty"`
	FirstName                          string           `json:"first_name,omitempty"`
	LastName                           string           `json:"last_name,omitempty"`
	IsForum                            bool             `json:"is_forum,omitempty"`
	ActiveUsernames                    []string         `json:"active_usernames,omitempty"`
	EmojiStatusCustomEmojiID           string           `json:"emoji_status_custom_emoji_id,omitempty"`
	EmojiStatusExpirationDate          int64            `json:"emoji_status_expiration_date,omitempty"`
	Bio                                string           `json:"bio,omitempty"`
	HasPrivateForwards                 bool             `json:"has_private_forwards,omitempty"`
	HasRestrictedVoiceAndVideoMessages bool             `json:"has_restricted_voice_and_video_messages,omitempty"`
	JoinToSendMessages                 bool             `json:"join_to_send_messages,omitempty"`
	JoinByRequest                      bool             `json:"join_by_request,omitempty"`
	Description                        string           `json:"description,omitempty"`
	InviteLink                         string           `json:"invite_link,omitempty"`
	PinnedMessage                      *TelegramMessage `json:"pinned_message,omitempty"`
	Permissions                        *ChatPermissions `json:"permissions,omitempty"`
	SlowModeDelay                      int              `json:"slow_mode_delay,omitempty"`
	MessageAutoDeleteTime              int              `json:"message_auto_delete_time,omitempty"`
	HasAggressiveAntiSpamEnabled       bool             `json:"has_aggressive_anti_spam_enabled,omitempty"`
	HasHiddenMembers                   bool             `json:"has_hidden_members,omitempty"`
	HasProtectedContent                bool             `json:"has_protected_content,omitempty"`
	StickerSetName                     string           `json:"sticker_set_name,omitempty"`
	CanSetStickerSet                   bool             `json:"can_set_sticker_set,omitempty"`
	LinkedChatID                       int64            `json:"linked_chat_id,omitempty"`
	Location                           *ChatLocation    `json:"location,omitempty"`
}

// ChatPermissions представляет разрешения чата
type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`
	CanSendPolls          bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"`
	CanChangeInfo         bool `json:"can_change_info,omitempty"`
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       bool `json:"can_manage_topics,omitempty"`
}

// ChatLocation представляет местоположение чата
type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

// ToTelegramMessage конвертирует внутреннее сообщение в формат Telegram Bot API
func (m *Message) ToTelegramMessage() TelegramMessage {
	// Все ID уже int64, конвертация не нужна
	telegramMessage := TelegramMessage{
		MessageID: m.ID,
		From: TelegramUser{
			ID:        m.FromID,
			IsBot:     m.From.IsBot,
			FirstName: m.From.FirstName,
			LastName:  m.From.LastName,
			Username:  m.From.Username,
		},
		Chat: TelegramChat{
			ID:       m.ChatID,
			Type:     "private", // TODO: получить тип чата
			Title:    m.From.GetFullName(),
			Username: m.From.Username,
		},
		Date: m.Timestamp.Unix(),
		Text: m.Text,
	}

	// Добавляем сущности, если они есть
	if entities := m.GetEntities(); len(entities) > 0 {
		telegramMessage.Entities = entities
	}

	// Добавляем клавиатуру, если она есть
	if replyMarkup := m.GetReplyMarkup(); replyMarkup != nil {
		telegramMessage.ReplyMarkup = replyMarkup
	}

	return telegramMessage
}

// FromTelegramMessage конвертирует сообщение из формата Telegram Bot API во внутренний формат
func FromTelegramMessage(tgMsg TelegramMessage, chatID int64) *Message {
	return &Message{
		ID:        tgMsg.MessageID, // Все ID уже int64
		ChatID:    chatID,
		FromID:    tgMsg.From.ID,
		Text:      tgMsg.Text,
		Type:      MessageTypeText,
		Status:    MessageStatusSent,
		Timestamp: time.Unix(tgMsg.Date, 0),
		CreatedAt: time.Now(),
	}
}

// ToTelegramCallbackQuery конвертирует внутренний CallbackQuery в формат Telegram Bot API
func (cq *CallbackQuery) ToTelegramCallbackQuery() map[string]interface{} {
	// Все ID уже int64, конвертация не нужна
	telegramCallbackQuery := map[string]interface{}{
		"id": cq.ID,
		"from": map[string]interface{}{
			"id":         cq.From.ID, // Уже int64
			"is_bot":     cq.From.IsBot,
			"first_name": cq.From.FirstName,
			"last_name":  cq.From.LastName,
			"username":   cq.From.Username,
		},
		"chat_instance": cq.ChatInstance,
	}

	// Добавляем данные, если они есть
	if cq.Data != "" {
		telegramCallbackQuery["data"] = cq.Data
	}

	// Добавляем сообщение, если оно есть
	if cq.Message != nil {
		telegramCallbackQuery["message"] = cq.Message.ToTelegramMessage()
	}

	// Добавляем inline_message_id, если он есть
	if cq.InlineMessageID != "" {
		telegramCallbackQuery["inline_message_id"] = cq.InlineMessageID
	}

	// Добавляем game_short_name, если он есть
	if cq.GameShortName != "" {
		telegramCallbackQuery["game_short_name"] = cq.GameShortName
	}

	return telegramCallbackQuery
}
