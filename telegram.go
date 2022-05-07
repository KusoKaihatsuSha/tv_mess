package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	UpdateID int64 `json:"update_id"`
	Message  struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"message"`
	EditedMessage struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int    `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"edited_message"`
	ChannelPost struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int    `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"channel_post"`
	EditedChannelPost struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"edited_channel_post"`
	InlineQuery struct {
		ID   string `json:"id"`
		From struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Query  string `json:"query"`
		Offset string `json:"offset"`
	} `json:"inline_query"`
	ChosenInlineResult struct {
		ResultID string `json:"result_id"`
		From     struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		InlineMessageID string `json:"inline_message_id"`
		Query           string `json:"query"`
	} `json:"chosen_inline_result"`
	CallbackQuery struct {
		ID   string `json:"id"`
		From struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Message struct {
			MessageID int64 `json:"message_id"`
			From      struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Date int `json:"date"`
			Chat struct {
				ID                          int64  `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"chat"`
			ForwardFrom struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"forward_from"`
			ForwardFromChat struct {
				ID                          int64  `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"forward_from_chat"`
			ForwardFromMessageID int    `json:"forward_from_message_id"`
			ForwardDate          int    `json:"forward_date"`
			EditDate             int    `json:"edit_date"`
			Text                 string `json:"text"`
			Entities             []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"entities"`
			CaptionEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"caption_entities"`
			Audio struct {
				FileID    string `json:"file_id"`
				Duration  int    `json:"duration"`
				Performer string `json:"performer"`
				Title     string `json:"title"`
				MimeType  string `json:"mime_type"`
				FileSize  int    `json:"file_size"`
			} `json:"audio"`
			Document struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"document"`
			Game struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Photo       []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Text         string `json:"text"`
				TextEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int64  `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"text_entities"`
				Animation struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"animation"`
			} `json:"game"`
			Photo []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Sticker struct {
				FileID string `json:"file_id"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				Emoji        string `json:"emoji"`
				SetName      string `json:"set_name"`
				MaskPosition struct {
					Point  string `json:"point"`
					XShift int    `json:"x_shift"`
					YShift int    `json:"y_shift"`
					Zoom   int    `json:"zoom"`
				} `json:"mask_position"`
				FileSize int `json:"file_size"`
			} `json:"sticker"`
			Video struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"video"`
			Voice struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"voice"`
			VideoNote struct {
				FileID   string `json:"file_id"`
				Length   int    `json:"length"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileSize int `json:"file_size"`
			} `json:"video_note"`
			Caption string `json:"caption"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				UserID      int    `json:"user_id"`
			} `json:"contact"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Venue struct {
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Title        string `json:"title"`
				Address      string `json:"address"`
				FoursquareID string `json:"foursquare_id"`
			} `json:"venue"`
			NewChatMembers []struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"new_chat_members"`
			LeftChatMember struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"left_chat_member"`
			NewChatTitle string `json:"new_chat_title"`
			NewChatPhoto []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"new_chat_photo"`
			DeleteChatPhoto       bool `json:"delete_chat_photo"`
			GroupChatCreated      bool `json:"group_chat_created"`
			SupergroupChatCreated bool `json:"supergroup_chat_created"`
			ChannelChatCreated    bool `json:"channel_chat_created"`
			MigrateToChatID       int  `json:"migrate_to_chat_id"`
			MigrateFromChatID     int  `json:"migrate_from_chat_id"`
			Invoice               struct {
				Title          string `json:"title"`
				Description    string `json:"description"`
				StartParameter string `json:"start_parameter"`
				Currency       string `json:"currency"`
				TotalAmount    int    `json:"total_amount"`
			} `json:"invoice"`
			SuccessfulPayment struct {
				Currency         string `json:"currency"`
				TotalAmount      int    `json:"total_amount"`
				InvoicePayload   string `json:"invoice_payload"`
				ShippingOptionID string `json:"shipping_option_id"`
				OrderInfo        struct {
					Name            string `json:"name"`
					PhoneNumber     string `json:"phone_number"`
					Email           string `json:"email"`
					ShippingAddress struct {
						CountryCode string `json:"country_code"`
						Stat        string `json:"stat"`
						City        string `json:"city"`
						StreetLine1 string `json:"street_line1"`
						StreetLine2 string `json:"street_line2"`
						PostCode    string `json:"post_code"`
					} `json:"shipping_address"`
				} `json:"order_info"`
				TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
				ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
			} `json:"successful_payment"`
			ForwardSignature string `json:"forward_signature"`
			AuthorSignature  string `json:"author_signature"`
			ConnectedWebsite string `json:"connected_website"`
		} `json:"message"`
		InlineMessageID string `json:"inline_message_id"`
		ChatInstance    string `json:"chat_instance"`
		Data            string `json:"data"`
		GameShortName   string `json:"game_short_name"`
	} `json:"callback_query"`
	ShippingQuery struct {
		ID   string `json:"id"`
		From struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		InvoicePayload  string `json:"invoice_payload"`
		ShippingAddress struct {
			CountryCode string `json:"country_code"`
			Stat        string `json:"stat"`
			City        string `json:"city"`
			StreetLine1 string `json:"street_line1"`
			StreetLine2 string `json:"street_line2"`
			PostCode    string `json:"post_code"`
		} `json:"shipping_address"`
	} `json:"shipping_query"`
	PreCheckoutQuery struct {
		ID   string `json:"id"`
		From struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Currency         string `json:"currency"`
		TotalAmount      int    `json:"total_amount"`
		InvoicePayload   string `json:"invoice_payload"`
		ShippingOptionID string `json:"shipping_option_id"`
		OrderInfo        struct {
			Name            string `json:"name"`
			PhoneNumber     string `json:"phone_number"`
			Email           string `json:"email"`
			ShippingAddress struct {
				CountryCode string `json:"country_code"`
				Stat        string `json:"stat"`
				City        string `json:"city"`
				StreetLine1 string `json:"street_line1"`
				StreetLine2 string `json:"street_line2"`
				PostCode    string `json:"post_code"`
			} `json:"shipping_address"`
		} `json:"order_info"`
	} `json:"pre_checkout_query"`
	Autoload bool
}

type UpdateReturn struct {
	Result      []Result `json:"result"`
	ErrorCode   int      `json:"error_code"`
	Ok          bool     `json:"ok"`
	Description string   `json:"description"`
}

type DeleteMessageReturn struct {
	Result      bool   `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type SendMessageReturn struct {
	Result struct {
		MessageID int64 `json:"message_id"`
		From      struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int64  `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int64  `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int64  `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Dice struct {
			Emoji string `json:"emoji"`
			Value int    `json:"value"`
		} `json:"dice"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int64  `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type InlineReturn struct {
	Result      bool   `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type Button struct {
	InlineKeyboard [][]ButtonOne `json:"inline_keyboard"`
}

type ButtonOne struct {
	Text         string `json:"text"`
	Url          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type Keyboard struct {
	CustomKeyboard [][]KeyboardLabel `json:"keyboard"`
	Resize         bool              `json:"resize_keyboard"`
	OneTime        bool              `json:"one_time_keyboard"`
}

type RemoveKeyboard struct {
	Remove bool `json:"remove_keyboard"`
}

type KeyboardLabel struct {
	Text string `json:"text"`
}

type PayloadMesageSend struct {
	ChatID                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageID      int    `json:"reply_to_message_id"`
}

type MinimalMessage struct {
	MessageId       int64  `json:"message_id"`
	ChatID          int64  `json:"chat_id"`
	Text            string `json:"text"`
	ParseMode       string `json:"parse_mode"`
	ReturnMessageId int64
}

type Message struct {
	DisableWebPagePreview bool        `json:"disable_web_page_preview"`
	DisableNotification   bool        `json:"disable_notification"`
	ReplyToMessageID      int         `json:"reply_to_message_id"`
	ReplyMarkup           interface{} `json:"reply_markup"`
	Caption               string      `json:"caption"`
	CallbackQueryId       string      `json:"callback_query_id"`
	ShowAlert             bool        `json:"show_alert"`
	Url                   string      `json:"url"`
	CacheTime             int         `json:"cache_time"`
	Emoji                 string      `json:"emoji"`
	MinimalMessage
	MessageExt
	TypeSrc
}

type MessageExt struct {
	User            string
	MessageIdStr    string
	ChatIDStr       string
	Command         string
	MessageIdClient int64
	MessageIdServer int64
	TextEdit        string
	FileID          string
	Cbmid           int64
	DelBefore       bool
	DelAfter        bool
	DelAfterDelay   time.Duration
	Random          int
	LanguageCode    string
	Src             string
	Folder          string
	UUID            string
}

type ItemDeleteMessage struct {
	Del string
}

type TypeSrc struct {
	Picture  bool
	Video    bool
	Document bool
	Audio    bool
}

type DocumentMessage struct {
	Src   string
	Title string
	Check bool
}

type Updates struct {
	Offset         int64    `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type Buttons struct {
	Lines []*ButtonsLines
}

type ButtonsLines struct {
	Parent *Buttons
	Childs []*ButtonsElements
}

type ButtonsElements struct {
	Parent *ButtonsLines
	Name   string
}

type Commands struct {
	Items       []*Command
	Db          *DataBase
	B           *Buttons
	CurrentLine *ButtonsLines
}

type StringMarkers struct {
	Positive     string
	Neutral      string
	Negative     string
	PositiveFlag string
	NeutralFlag  string
	NegativeFlag string
}

type Command struct {
	StringMarkers
	Points        map[string]string
	Messages      []*Message
	B             bool
	Parent        *Commands
	Name          string
	HumanRead     string
	TrHumanRead   map[string]string
	DataRead      string
	CurrentFlag   bool
	Current       string
	Binary        bool
	Tid           int64
	Mid           int64
	NextFlag      bool
	Next          string
	ReturnButtons bool
	Buttons       *Keyboard
	Json          url.Values
	Text          string
	F             func(Message, *botUser)
	IsCommand     bool
}

// New() *Buttons
// dumpy function
func (o *Buttons) New() *Buttons {
	return o
}

// NewMessage(int64, string, mid int64, cbmid int64, command string) *Message
// init new message
func NewMessage(chatId int64, text string, mid int64, cbmid int64, command string) *Message {
	o := new(Message)
	o.ChatID = chatId
	o.MessageId = mid
	o.Text = text
	o.Command = command
	return o
}

// SendMessageWrapperTask(*Tasker, Thing, *Message)
// send text with tasker
func (o Message) SendMessageWrapperTask(T *Tasker, task Thing, message *Message) {
	o.extensionMessaging(T, sendParam, false, o.sendMessage)
}

// SendRandomWrapperTask(*Tasker, Thing, *Message)
// send random with tasker
func (o Message) SendRandomWrapperTask(T *Tasker, task Thing, message *Message) {
	o.extensionMessaging(T, diceParam, false, o.sendRandom)
}

// SendDocumentWrapperTask(*Tasker, Thing, *Message)
// send document with tasker
func (o Message) SendDocumentWrapperTask(T *Tasker, task Thing, message *Message) {
	o.sendDocument(T, task.Input.(DocumentMessage))
}

// extensionMessaging(*Tasker, string, bool, func(*Tasker, string) int64)
// wrapper for text message. Timeout erase, erase prev message and typing inform.
func (o *Message) extensionMessaging(T *Tasker, source string, isFile bool, fn func(*Tasker, string) int64) {
	if isFile {
		go o.sendTyping("upload_document")
	} else {
		go o.sendTyping("typing")
	}
	if o.DelBefore {
		go o.deleteMessage()
	}
	if o.DelAfter {
		o.Text = timerLabel + o.Text
	}
	answerMessage := o
	answerMessage.MessageId = fn(T, source)
	o.ReturnMessageId = answerMessage.MessageId
	defer func() {
		if o.DelAfter {
			<-time.After(o.DelAfterDelay)
			answerMessage.deleteMessage()
		}
	}()
}

// sendMessage(*Tasker, source string) int64
// send text message
func (o *Message) sendMessage(T *Tasker, source string) int64 {
	r := new(SendMessageReturn)
	telegramQuery(source, o, r, false, "")
	return r.Result.MessageID
}

// sendRandom(*Tasker, string) int64
// send randomize to chat.
func (o *Message) sendRandom(T *Tasker, source string) int64 {
	o.Emoji = "🎲"
	r := new(SendMessageReturn)
	telegramQuery(source, o, r, false, "")
	o.AddCtx(T, "random", r.Result.Dice.Value)
	return r.Result.MessageID
}

// telegramQuery[T any](string, any, T, bool, string) error
// Wrapper for query. Send raw text or file use multipart
func telegramQuery[T any](url string, body any, ret T, multipartFlag bool, boundary string) error {
	queryTo := new(Query)
	queryTo.Host = telegramUrl + api + url
	queryTo.Parameters = make(map[string]string)
	queryTo.Type = http.MethodPost
	if multipartFlag {
		queryTo.Parameters["Content-Type"] = "multipart/form-data; boundary=" + boundary
		queryTo.Parameters["Connection"] = "keep-alive"
		queryTo.Data = body.(*bytes.Buffer)
	} else {
		queryTo.Parameters["Accept"] = "application/json"
		queryTo.Parameters["Content-Type"] = "application/json"
		tmpp, err := json.Marshal(body)
		if err != nil {
			toLog(err)
		}
		queryTo.Data = bytes.NewReader(tmpp)
	}
	return json.Unmarshal(queryTo.Query(), ret)
}

// DeleteMessageWrapperTask(*Tasker, Thing, *Message)
// delete message from chat
func (o *Message) DeleteMessageWrapperTask(T *Tasker, task Thing, message *Message) {
	o.deleteMessage()
}

// PROBABLY NOT USING
func (o *Message) editMessage(text string) {
	o.Text = text
	r := new(DeleteMessageReturn)
	telegramQuery(editParam, o, r, false, "")
}

// deleteMessage()
// delete message from chat
func (o *Message) deleteMessage() {
	r := new(DeleteMessageReturn)
	telegramQuery(delParam, o, r, false, "")
}

// deleteMessageById(int64)
// PROBABLY NOT USING
func (o *Message) deleteMessageById(id int64) {
	o.MessageId = id
	r := new(DeleteMessageReturn)
	telegramQuery(delParam, o, r, false, "")
}

// sendDocument(*Tasker, DocumentMessage)
// send file to chat
func (o *Message) sendDocument(T *Tasker, src DocumentMessage) {
	check := func(end ...string) bool {
		for _, v := range end {
			if strings.HasSuffix(strings.ToLower(src.Src), v) {
				return true
			}
		}
		return false
	}
	if !src.Check {
		switch {
		case check("json"):
			if !debug {
				os.Remove(src.Src) //delete photo if not send
			}
		}
		o.AddCtx(T, paramGood+src.Src, true)
		return
	}
	if src.Src == "" {
		o.Text = infoLabel + src.Title
		T.Add(nil, o.SendMessageWrapperTask, o)
		o.AddCtx(T, paramGood+src.Src, true)
		return
	}

	go o.sendTyping("upload_document")
	if o.DelBefore {
		go o.deleteMessage()
	}
	func() {
		typeFile := "document"
		switch {
		case check("jpg", "png", "bmp"):
			typeFile = "photo"
		case check("avi", "mp4"):
			typeFile = "video" // or "animation"
		case check("mp3", "ogg", "wav"):
			typeFile = "audio"
		}
		file, err := os.Open(src.Src)
		if err != nil {
			toLog(err)
			return
		}
		data := &bytes.Buffer{}
		writer := multipart.NewWriter(data)
		if o.FileID != "" {
			writer.WriteField("file_id", o.FileID)
		}
		writer.WriteField("chat_id", sprintf("%d", o.ChatID))
		writer.WriteField("caption", src.Title)
		writer.WriteField("disable_notification", strconv.FormatBool(o.DisableNotification))
		base := filepath.Base(file.Name())
		part, err := writer.CreateFormFile(typeFile, base)
		if err != nil {
			toLog(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			toLog(err)
		}
		writer.Close() //unblock file
		file.Close()
		ret := new(SendMessageReturn)
		telegramQuery("/send"+strings.ToTitle(string(typeFile[0]))+typeFile[1:], data, ret, true, writer.Boundary())
		fileId := ""
		switch {
		case check("jpg", "png", "bmp"):
			for _, v := range ret.Result.Photo {
				fileId = v.FileID
				break
			}
		case check("avi", "mp4"):
			fileId = ret.Result.Video.FileID
		case check("mp3", "ogg", "wav"):
			fileId = ret.Result.Audio.FileID
		default:
			fileId = ret.Result.Document.FileID
		}
		o.FileID = fileId
		o.AddCtx(T, paramGood+src.Src, true)
		switch {
		case check("json"):
			if !debug {
				os.Remove(src.Src) //delete photo if not send
			}
		case check(mp4):
			jpgDel := strings.TrimSuffix(src.Src, mp4) + jpg
			if !debug {
				os.Remove(src.Src) //if not send not delete
				os.Remove(jpgDel)
			}
		case check(mp3):
			mp4Del := strings.TrimSuffix(src.Src, mp3) + mp4
			jpgDel := strings.TrimSuffix(src.Src, mp3) + jpg
			if !debug {
				os.Remove(src.Src) //if send mp4 may be delete
				os.Remove(mp4Del)
				os.Remove(jpgDel)
			}
		}
	}()
}

// sendTyping(string)
// send typing action to chat
func (o Message) sendTyping(text string) {
	r := new(DeleteMessageReturn) // rename?
	telegramQuery("/sendChatAction?chat_id="+sprintf("%d", o.ChatID)+"&action="+text, nil, r, false, "")
}

// buttomsMap(map[string]string) [][]ButtonOne
// set buttons in message
func buttomsMap(val map[string]string) [][]ButtonOne {
	maxFill := 44
    i := 0
	for k := range val {
		i += len([]rune(k))
	}
	lenElements := i / maxFill
	if i%maxFill != 0 {
		lenElements++
	}
	lenElements = len(val) / lenElements
	keys := make([]string, 0, len(val))
	for key := range val {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return val[keys[i]] > val[keys[j]]
	})
	but := [][]ButtonOne{}
	sbut := []ButtonOne{}
	for i, k := range keys {
		if i%lenElements == 0 {
			but = append(but, sbut)
			sbut = []ButtonOne{}
		}
		sbut = append(sbut, ButtonOne{Text: k, CallbackData: val[k]})
	}
	but = append(but, sbut)
	return but
}

// removeKeyboard() RemoveKeyboard
// remove keyboard
func removeKeyboard() RemoveKeyboard {
	but := RemoveKeyboard{Remove: true}
	return but
}

// NewLine() *ButtonsLines
// new line
func (o *Buttons) NewLine() *ButtonsLines {
	newLine := ButtonsLines{Parent: o}
	o.Lines = append(o.Lines, &newLine)
	return &newLine
}

// Return() Keyboard
// get keyboard
func (o *Buttons) Return() Keyboard {
	buttons := [][]KeyboardLabel{}
	for _, val0 := range o.Lines {
		line := []KeyboardLabel{}
		for _, val := range val0.Childs {
			line = append(line, KeyboardLabel{Text: val.Name})
		}
		buttons = append(buttons, line)
	}
	return Keyboard{CustomKeyboard: buttons, Resize: true, OneTime: false}
}

// Add(string) *ButtonsElements
// new button
func (o *ButtonsLines) Add(name string) *ButtonsElements {
	newLine := ButtonsElements{Name: name, Parent: o}
	o.Childs = append(o.Childs, &newLine)
	return &newLine
}

// Return() Keyboard
// get keyboard
func (o *ButtonsLines) Return() Keyboard {
	return o.Parent.Return()
}

// Add(string) *ButtonsElements
// new button
func (o *ButtonsElements) Add(name string) *ButtonsElements {
	newLine := ButtonsElements{Name: name, Parent: o.Parent}
	o.Parent.Childs = append(o.Parent.Childs, &newLine)
	return &newLine
}

// NewLine() *ButtonsLines
// new line
func (o *ButtonsElements) NewLine() *ButtonsLines {
	return o.Parent.Parent.NewLine()
}

// Return() Keyboard
// get keyboard
func (o *ButtonsElements) Return() Keyboard {
	return o.Parent.Parent.Return()
}

// setKb([][]string) Keyboard
// set new keyboard on bottom
func setKb(elements [][]string) Keyboard {
	buttons := [][]KeyboardLabel{}
	for _, val0 := range elements {
		line := []KeyboardLabel{}
		for _, val := range val0 {
			line = append(line, KeyboardLabel{Text: val})
		}
		buttons = append(buttons, line)
	}
	return Keyboard{CustomKeyboard: buttons, Resize: true, OneTime: false}
}

// New() *Updates
// new return struct
func (o *Updates) New() *Updates {
	o.Timeout = 0 //1
	o.Limit = 1
	o.Offset = 1
	o.AllowedUpdates = append(o.AllowedUpdates, "message")
	o.AllowedUpdates = append(o.AllowedUpdates, "callback_query")
	o.AllowedUpdates = append(o.AllowedUpdates, "inline_query")
	return o
}

// SetLast(last int64) *Updates
// set num id of return message
func (o *Updates) SetLast(last int64) *Updates {
	o.Offset = last + 1
	return o
}

// GetLast() int64
// get last num id
func (o *Updates) GetLast() int64 {
	return o.Offset
}

// AddMessage(msg *Message) *Command
// add message
func (o *Command) AddMessage(msg *Message) *Command {
	o.Messages = append(o.Messages, msg)
	return o
}

// AddButtons() *Command
// create buttons
func (o *Command) AddButtons() *Command {
	o.Parent.CurrentLine.Add(o.HumanRead)
	return o
}

// AddNewLine() *Command
// commands go on new line
func (o *Command) AddNewLine() *Command {
	o.Parent.B.NewLine()
	return o
}

// AddButtons() *Commands
// create buttons
func (o *Commands) AddButtons() *Commands {
	for _, v := range o.Items {
		if v.B {
			o.CurrentLine.Add(v.HumanRead)
		}
	}
	return o
}

// DbConnect(*DataBase) *Commands
// pair with database
func (o *Commands) DbConnect(db *DataBase) *Commands {
	o.Db = db
	return o
}

// AddNewLine() *Commands
// commands go on new line
func (o *Commands) AddNewLine() *Commands {
	o.CurrentLine = o.B.NewLine()
	return o
}

// AddSender(*botUser) *Command
// dumpy function
func (o *Command) AddSender(user *botUser) *Command {
	return o
}

// GetAllUsers() map[string]string
// get users from database
func (o *Commands) GetAllUsers() map[string]string {
	return o.Db.FindCreate("users").PrintAll()
}

// GetAllPoints() *Command
// dumpy function
func (o *Command) GetAllPoints() *Command {
	return o
}

// AddMessage(*Message) *Commands
// dumpy function
func (o *Commands) AddMessage(msg *Message) *Commands {
	return o
}

// Add(string, string, bool, bool, func(Message, *botUser)) *Commands
// add command
func (o *Commands) Add(name, humanRead string, binary, toBut bool, f func(Message, *botUser)) *Commands {
	c := new(Command)
	c.F = f
	c.Parent = o
	c.B = toBut
	c.Name = "/" + name
	c.HumanRead = humanRead
	c.Binary = binary
	o.Items = append(o.Items, c)
	if toBut {
		c.AddButtons()
	}
	return o
}

// Find(string) *Command
// find command
func (o *Commands) Find(text string) *Command {
	for _, v := range o.Items {
		if strings.HasPrefix(text, v.Name) || strings.HasPrefix(text, v.HumanRead) {
			v.Text = text
			v.DataRead = strings.TrimPrefix(text, v.Name)
			v.IsCommand = true
			return v
		}
	}
	return o.Find("/help")
}

// Add(string, string, bool) *Command
// set name command
func (o *Command) Add(name, humanRead string, binary bool) *Command {
	o.Name = "/" + name
	o.HumanRead = humanRead
	return o
}

// sendFiles(*Tasker, Thing, string, *Message) bool
// sendfiles or split and send
func sendFiles(T *Tasker, task Thing, format string, message *Message) bool {
	v := task.Input.(*JsonPls)
    select {
	case <-T.Branch.Context.Done():
        os.RemoveAll(v.URLSaved+format)
		return false
	default:
	}
	var splitFiles []string
	user := GetCtx[*botUser](T, userParam, message)
	splitMp4 := true
	if format == mp4 {
		splitMp4 = !sBool(user.getParameter(paramParam, mp3))
	}
	if fileSize(v.URLSaved+format) >= limitFileTelegram && format != jpg && splitMp4 {
		SplitMp(T, task, 0, limitFileTelegram, format, message)
		splitFiles = searchFiles(v.URLSaved+"__", v.UUID, format)
		for k, val := range splitFiles {
			param := DocumentMessage{}
			param.Src = val
			param.Check = sBool(user.getParameter(paramParam, format))
			param.Title = strconv.Itoa(k+1) + ") " + v.Artist + " [" + v.Song + "]"
			message.sendDocument(T, param)
			<-time.After(1 * time.Second)
		}
		return true
	} else {
		param := DocumentMessage{}
		param.Src = v.URLSaved + format
		if format == mp4 {
			param.Check = !sBool(user.getParameter(paramParam, mp3))
		} else {
			param.Check = sBool(user.getParameter(paramParam, format))
		}
		param.Title = v.Artist + " [" + v.Song + "]"
		T.Add(param, message.SendDocumentWrapperTask, message)
		return false
	}
}
