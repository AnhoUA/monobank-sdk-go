package mcc

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed mcc.json mcc-en.json
var mccFS embed.FS

// MCC represents a Merchant Category Code.
type MCC struct {
	MCC              string `json:"mcc"`
	Group            Group  `json:"group"`
	FullDescription  any    `json:"fullDescription"`  // Can be string or map[string]string
	ShortDescription any    `json:"shortDescription"` // Can be string or map[string]string
}

// Group represents a category group for MCC.
type Group struct {
	Type        string `json:"type"`
	Description any    `json:"description"` // Can be string or map[string]string
}

// MCCInfo is a normalized version of MCC data.
type MCCInfo struct {
	MCC              string
	GroupType        string
	GroupName        string
	FullDescription  string
	ShortDescription string
}

var (
	mccData   map[string]MCCInfo
	mccEnData map[string]MCCInfo
	once      sync.Once
	loadErr   error
)

func load() {
	mccData, loadErr = loadFile("mcc.json", "uk")
	if loadErr != nil {
		return
	}
	mccEnData, loadErr = loadFile("mcc-en.json", "en")
}

func loadFile(filename string, lang string) (map[string]MCCInfo, error) {
	data, err := mccFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", filename, err)
	}

	var raw []MCC
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s: %w", filename, err)
	}

	res := make(map[string]MCCInfo, len(raw))
	for _, m := range raw {
		info := MCCInfo{
			MCC:       m.MCC,
			GroupType: m.Group.Type,
		}

		info.GroupName = getString(m.Group.Description, lang)
		info.FullDescription = getString(m.FullDescription, lang)
		info.ShortDescription = getString(m.ShortDescription, lang)

		res[m.MCC] = info
	}
	return res, nil
}

func getString(v any, lang string) string {
	switch val := v.(type) {
	case string:
		return val
	case map[string]any:
		if s, ok := val[lang].(string); ok {
			return s
		}
		// Fallback to en if uk not found
		if s, ok := val["en"].(string); ok {
			return s
		}
		// Return any first available string
		for _, v := range val {
			if s, ok := v.(string); ok {
				return s
			}
		}
	}
	return ""
}

// Get returns MCC information by code.
// It uses English localization by default.
func Get(code string) (MCCInfo, bool) {
	return GetEn(code)
}

// GetUk returns MCC information by code in Ukrainian.
func GetUk(code string) (MCCInfo, bool) {
	once.Do(load)
	if loadErr != nil {
		return MCCInfo{}, false
	}
	info, ok := mccData[code]
	return info, ok
}

// GetEn returns MCC information by code in English.
func GetEn(code string) (MCCInfo, bool) {
	once.Do(load)
	if loadErr != nil {
		return MCCInfo{}, false
	}
	info, ok := mccEnData[code]
	return info, ok
}

// All returns all MCC records (English by default).
func All() ([]MCCInfo, error) {
	once.Do(load)
	if loadErr != nil {
		return nil, loadErr
	}
	res := make([]MCCInfo, 0, len(mccEnData))
	for _, v := range mccEnData {
		res = append(res, v)
	}
	return res, nil
}

// MCCService provides methods to access MCC data.
type MCCService struct{}

// Get returns MCC information by code.
func (s *MCCService) Get(code string) (MCCInfo, bool) {
	return Get(code)
}

// GetUk returns MCC information by code in Ukrainian.
func (s *MCCService) GetUk(code string) (MCCInfo, bool) {
	return GetUk(code)
}

// GetEn returns MCC information by code in English.
func (s *MCCService) GetEn(code string) (MCCInfo, bool) {
	return GetEn(code)
}

// All returns all MCC records.
func (s *MCCService) All() ([]MCCInfo, error) {
	return All()
}
