// Package converter 提供時間戳轉換功能
package converter

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimestampFormat 定義支援的時間格式
type TimestampFormat int

const (
	UnixSeconds TimestampFormat = iota
	UnixMilliseconds
	UnixMicroseconds
	UnixNanoseconds
	RFC3339
	RFC3339Nano
	DateTime
	DateOnly
	TimeOnly
)

// Converter 時間戳轉換器
type Converter struct {
	Location *time.Location
}

// NewConverter 建立新的轉換器
func NewConverter(timezone string) (*Converter, error) {
	if timezone == "" {
		timezone = "Local"
	}
	
	var loc *time.Location
	var err error
	
	if timezone == "Local" {
		loc = time.Local
	} else {
		loc, err = time.LoadLocation(timezone)
		if err != nil {
			return nil, fmt.Errorf("無法載入時區 %s: %v", timezone, err)
		}
	}
	
	return &Converter{Location: loc}, nil
}

// GetLocalTimezone 取得本機時區名稱
func GetLocalTimezone() string {
	// 取得本機時區名稱
	name, _ := time.Now().Zone()
	if name == "" {
		// 如果無法取得名稱，嘗試取得完整的時區資訊
		zone := time.Now().Location().String()
		if zone == "Local" {
			// 嘗試從系統環境變數取得
			if tz := os.Getenv("TZ"); tz != "" {
				return tz
			}
			return "Local (系統時區)"
		}
		return zone
	}
	return name
}

// ParseTimeOffset 解析相對時間偏移
func ParseTimeOffset(offset string) (time.Duration, error) {
	if offset == "" {
		return 0, nil
	}
	
	// 移除前後空格
	offset = strings.TrimSpace(offset)
	
	// 檢查符號
	var sign int = 1
	if strings.HasPrefix(offset, "+") {
		offset = offset[1:]
	} else if strings.HasPrefix(offset, "-") {
		sign = -1
		offset = offset[1:]
	}
	
	// 使用正則表達式解析數字和單位
	re := regexp.MustCompile(`^(\d+)([dwMyhms])$`)
	matches := re.FindStringSubmatch(offset)
	if len(matches) != 3 {
		return 0, fmt.Errorf("無效的時間偏移格式: %s (支援格式如: 1d, 2w, 3M, 4y, 5h, 6m, 7s)", offset)
	}
	
	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("無效的數字: %s", matches[1])
	}
	
	unit := matches[2]
	var duration time.Duration
	
	switch unit {
	case "s": // 秒
		duration = time.Duration(num) * time.Second
	case "m": // 分鐘
		duration = time.Duration(num) * time.Minute
	case "h": // 小時
		duration = time.Duration(num) * time.Hour
	case "d": // 天
		duration = time.Duration(num) * 24 * time.Hour
	case "w": // 週
		duration = time.Duration(num) * 7 * 24 * time.Hour
	case "M": // 月 (近似為30天)
		duration = time.Duration(num) * 30 * 24 * time.Hour
	case "y": // 年 (近似為365天)
		duration = time.Duration(num) * 365 * 24 * time.Hour
	default:
		return 0, fmt.Errorf("不支援的時間單位: %s", unit)
	}
	
	return time.Duration(sign) * duration, nil
}

// AddTimeOffset 為時間添加偏移量
func (c *Converter) AddTimeOffset(baseTime time.Time, offset string) (time.Time, error) {
	if offset == "" {
		return baseTime, nil
	}
	
	// 移除前後空格
	offset = strings.TrimSpace(offset)
	
	// 檢查符號
	var sign int = 1
	if strings.HasPrefix(offset, "+") {
		offset = offset[1:]
	} else if strings.HasPrefix(offset, "-") {
		sign = -1
		offset = offset[1:]
	}
	
	// 使用正則表達式解析數字和單位
	re := regexp.MustCompile(`^(\d+)([dwMyhms])$`)
	matches := re.FindStringSubmatch(offset)
	if len(matches) != 3 {
		return baseTime, fmt.Errorf("無效的時間偏移格式: %s (支援格式如: 1d, 2w, 3M, 4y, 5h, 6m, 7s)", offset)
	}
	
	num, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return baseTime, fmt.Errorf("無效的數字: %s", matches[1])
	}
	
	unit := matches[2]
	var result time.Time
	
	switch unit {
	case "s": // 秒
		result = baseTime.Add(time.Duration(int64(sign) * num) * time.Second)
	case "m": // 分鐘
		result = baseTime.Add(time.Duration(int64(sign) * num) * time.Minute)
	case "h": // 小時
		result = baseTime.Add(time.Duration(int64(sign) * num) * time.Hour)
	case "d": // 天
		result = baseTime.AddDate(0, 0, int(int64(sign) * num))
	case "w": // 週
		result = baseTime.AddDate(0, 0, int(int64(sign) * num * 7))
	case "M": // 月
		result = baseTime.AddDate(0, int(int64(sign) * num), 0)
	case "y": // 年
		result = baseTime.AddDate(int(int64(sign) * num), 0, 0)
	default:
		return baseTime, fmt.Errorf("不支援的時間單位: %s", unit)
	}
	
	return result.In(c.Location), nil
}

// DetectFormat 自動偵測輸入的時間格式
func (c *Converter) DetectFormat(input string) (TimestampFormat, error) {
	input = strings.TrimSpace(input)
	
	// 檢查是否為純數字 (Unix timestamp)
	if matched, _ := regexp.MatchString(`^\d+$`, input); matched {
		num, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("無法解析數字: %v", err)
		}
		
		// 根據數字長度判斷時間戳格式
		switch len(input) {
		case 10:
			return UnixSeconds, nil
		case 13:
			return UnixMilliseconds, nil
		case 16:
			return UnixMicroseconds, nil
		case 19:
			return UnixNanoseconds, nil
		default:
			// 嘗試作為秒級時間戳
			if num > 0 && num < 4102444800 { // 2100年前
				return UnixSeconds, nil
			}
			return 0, fmt.Errorf("無法識別的數字格式: %s", input)
		}
	}
	
	// 檢查 RFC3339 格式
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}`, input); matched {
		if strings.Contains(input, ".") {
			return RFC3339Nano, nil
		}
		return RFC3339, nil
	}
	
	// 檢查日期時間格式
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, input); matched {
		return DateTime, nil
	}
	
	// 檢查日期格式
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, input); matched {
		return DateOnly, nil
	}
	
	// 檢查時間格式
	if matched, _ := regexp.MatchString(`^\d{2}:\d{2}:\d{2}$`, input); matched {
		return TimeOnly, nil
	}
	
	return 0, fmt.Errorf("無法識別的時間格式: %s", input)
}

// Parse 解析輸入的時間字串
func (c *Converter) Parse(input string, format TimestampFormat) (time.Time, error) {
	input = strings.TrimSpace(input)
	
	switch format {
	case UnixSeconds:
		timestamp, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 Unix 秒級時間戳: %v", err)
		}
		return time.Unix(timestamp, 0).In(c.Location), nil
		
	case UnixMilliseconds:
		timestamp, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 Unix 毫秒級時間戳: %v", err)
		}
		return time.Unix(timestamp/1000, (timestamp%1000)*1000000).In(c.Location), nil
		
	case UnixMicroseconds:
		timestamp, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 Unix 微秒級時間戳: %v", err)
		}
		return time.Unix(timestamp/1000000, (timestamp%1000000)*1000).In(c.Location), nil
		
	case UnixNanoseconds:
		timestamp, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 Unix 納秒級時間戳: %v", err)
		}
		return time.Unix(timestamp/1000000000, timestamp%1000000000).In(c.Location), nil
		
	case RFC3339:
		t, err := time.Parse(time.RFC3339, input)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 RFC3339 格式: %v", err)
		}
		return t.In(c.Location), nil
		
	case RFC3339Nano:
		t, err := time.Parse(time.RFC3339Nano, input)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析 RFC3339Nano 格式: %v", err)
		}
		return t.In(c.Location), nil
		
	case DateTime:
		t, err := time.ParseInLocation("2006-01-02 15:04:05", input, c.Location)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析日期時間格式: %v", err)
		}
		return t, nil
		
	case DateOnly:
		t, err := time.ParseInLocation("2006-01-02", input, c.Location)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析日期格式: %v", err)
		}
		return t, nil
		
	case TimeOnly:
		today := time.Now().In(c.Location).Format("2006-01-02")
		fullTime := today + " " + input
		t, err := time.ParseInLocation("2006-01-02 15:04:05", fullTime, c.Location)
		if err != nil {
			return time.Time{}, fmt.Errorf("無法解析時間格式: %v", err)
		}
		return t, nil
		
	default:
		return time.Time{}, fmt.Errorf("不支援的格式")
	}
}

// ConvertResult 轉換結果
type ConvertResult struct {
	Original        string `json:"original"`
	DetectedFormat  string `json:"detected_format"`
	UnixSeconds     int64  `json:"unix_seconds"`
	UnixMillis      int64  `json:"unix_milliseconds"`
	UnixMicros      int64  `json:"unix_microseconds"`
	UnixNanos       int64  `json:"unix_nanoseconds"`
	RFC3339         string `json:"rfc3339"`
	RFC3339Nano     string `json:"rfc3339_nano"`
	DateTime        string `json:"datetime"`
	DateOnly        string `json:"date_only"`
	TimeOnly        string `json:"time_only"`
	Weekday         string `json:"weekday"`
	Timezone        string `json:"timezone"`
}

// Convert 轉換時間到所有格式
func (c *Converter) Convert(input string, inputFormat *TimestampFormat) (*ConvertResult, error) {
	var format TimestampFormat
	var err error
	
	if inputFormat != nil {
		format = *inputFormat
	} else {
		format, err = c.DetectFormat(input)
		if err != nil {
			return nil, err
		}
	}
	
	t, err := c.Parse(input, format)
	if err != nil {
		return nil, err
	}
	
	result := &ConvertResult{
		Original:        input,
		DetectedFormat:  c.formatName(format),
		UnixSeconds:     t.Unix(),
		UnixMillis:      t.UnixMilli(),
		UnixMicros:      t.UnixMicro(),
		UnixNanos:       t.UnixNano(),
		RFC3339:         t.Format(time.RFC3339),
		RFC3339Nano:     t.Format(time.RFC3339Nano),
		DateTime:        t.Format("2006-01-02 15:04:05"),
		DateOnly:        t.Format("2006-01-02"),
		TimeOnly:        t.Format("15:04:05"),
		Weekday:         c.weekdayName(t.Weekday()),
		Timezone:        c.getTimezoneInfo(t),
	}
	
	return result, nil
}

// formatName 返回格式名稱
func (c *Converter) formatName(format TimestampFormat) string {
	switch format {
	case UnixSeconds:
		return "Unix 秒級時間戳"
	case UnixMilliseconds:
		return "Unix 毫秒級時間戳"
	case UnixMicroseconds:
		return "Unix 微秒級時間戳"
	case UnixNanoseconds:
		return "Unix 納秒級時間戳"
	case RFC3339:
		return "RFC3339 格式"
	case RFC3339Nano:
		return "RFC3339Nano 格式"
	case DateTime:
		return "日期時間格式"
	case DateOnly:
		return "日期格式"
	case TimeOnly:
		return "時間格式"
	default:
		return "未知格式"
	}
}

// weekdayName 返回中文星期名稱
func (c *Converter) weekdayName(weekday time.Weekday) string {
	switch weekday {
	case time.Sunday:
		return "星期日"
	case time.Monday:
		return "星期一"
	case time.Tuesday:
		return "星期二"
	case time.Wednesday:
		return "星期三"
	case time.Thursday:
		return "星期四"
	case time.Friday:
		return "星期五"
	case time.Saturday:
		return "星期六"
	default:
		return "未知"
	}
}

// getTimezoneInfo 取得時區資訊
func (c *Converter) getTimezoneInfo(t time.Time) string {
	zone, offset := t.Zone()
	location := t.Location().String()
	
	// 計算時區偏移量
	offsetHours := offset / 3600
	offsetMinutes := (offset % 3600) / 60
	
	var offsetStr string
	if offset >= 0 {
		offsetStr = fmt.Sprintf("+%02d:%02d", offsetHours, offsetMinutes)
	} else {
		offsetStr = fmt.Sprintf("-%02d:%02d", -offsetHours, -offsetMinutes)
	}
	
	if location == "Local" {
		if zone == "" {
			return fmt.Sprintf("Local (UTC%s)", offsetStr)
		}
		return fmt.Sprintf("Local (%s, UTC%s)", zone, offsetStr)
	}
	
	if zone == "" {
		return fmt.Sprintf("%s (UTC%s)", location, offsetStr)
	}
	
	return fmt.Sprintf("%s (%s, UTC%s)", location, zone, offsetStr)
}
