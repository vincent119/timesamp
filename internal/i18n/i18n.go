// Package i18n 提供國際化支援
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*/*.json
var localeFS embed.FS

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
	currentLang string
)

// SupportedLanguages 支援的語言列表
var SupportedLanguages = []string{
	"en",    // English
	"zh-TW", // Traditional Chinese
	"zh-CN", // Simplified Chinese
	"ja",    // Japanese
}

// Init 初始化 i18n
func Init() error {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// 載入所有語言檔案
	for _, lang := range SupportedLanguages {
		fileName := fmt.Sprintf("locales/%s/messages.json", lang)
		data, err := localeFS.ReadFile(fileName)
		if err != nil {
			continue
		}

		// 使用包含語言標籤的假檔案名稱，讓 go-i18n 能正確識別語言
		// 例如: messages.zh-TW.json
		fakeFileName := fmt.Sprintf("messages.%s.json", lang)
		_, err = bundle.ParseMessageFileBytes(data, fakeFileName)
		if err != nil {
			continue
		}
	}

	// 設定預設語言
	SetLanguage(DetectLanguage())
	return nil
}

// DetectLanguage 偵測系統語言
func DetectLanguage() string {
	// 1. 檢查環境變數 TIMESTAMP_LANG
	if lang := os.Getenv("TIMESTAMP_LANG"); lang != "" {
		return normalizeLanguage(lang)
	}

	// 2. 檢查標準語言環境變數
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if lang := os.Getenv(env); lang != "" {
			return normalizeLanguage(lang)
		}
	}

	// 3. 預設使用英文
	return "en"
}

// normalizeLanguage 正規化語言標籤
func normalizeLanguage(lang string) string {
	// 解析語言標籤 (如: zh_TW.UTF-8 -> zh-TW)
	lang = strings.Split(lang, ".")[0]
	lang = strings.Replace(lang, "_", "-", -1)

	// 檢查是否支援
	for _, supported := range SupportedLanguages {
		if strings.HasPrefix(lang, supported) || strings.HasPrefix(supported, lang) {
			return supported
		}
	}

	// 特殊處理中文
	if strings.HasPrefix(lang, "zh") {
		if strings.Contains(lang, "TW") || strings.Contains(lang, "HK") || strings.Contains(lang, "MO") {
			return "zh-TW"
		}
		return "zh-CN"
	}

	return "en"
}

// SetLanguage 設定語言
func SetLanguage(lang string) {
	// 驗證語言是否支援
	supported := false
	for _, l := range SupportedLanguages {
		if l == lang {
			supported = true
			break
		}
	}

	if !supported {
		lang = "en" // 回退到英文
	}

	currentLang = lang
	localizer = i18n.NewLocalizer(bundle, lang)
}

// T 翻譯函數
func T(messageID string, templateData ...map[string]interface{}) string {
	if localizer == nil {
		return messageID // 如果未初始化，返回原始 ID
	}

	config := &i18n.LocalizeConfig{
		MessageID: messageID,
	}

	if len(templateData) > 0 && templateData[0] != nil {
		config.TemplateData = templateData[0]
	}

	translated, err := localizer.Localize(config)
	if err != nil {
		return messageID // 如果翻譯失敗，返回原始 ID
	}

	return translated
}

// Tf 帶格式化參數的翻譯函數
func Tf(messageID string, templateData map[string]interface{}) string {
	return T(messageID, templateData)
}

// GetCurrentLanguage 取得當前語言
func GetCurrentLanguage() string {
	if currentLang == "" {
		return "en"
	}
	return currentLang
}

// ListSupportedLanguages 列出支援的語言
func ListSupportedLanguages() []string {
	return SupportedLanguages
}
