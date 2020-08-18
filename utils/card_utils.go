package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

const (
	// Path for message cards
	CardPath string = "data/cards/"

	cardSpan time.Duration = time.Minute
)

// Card - message card class
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"cover_url"`
	DestURL     string `json:"destination_url"`
}

func SaveCard(title string, description string, coverURL string, destURL string) (string, error) {
	cid := getCardID(destURL)
	path := getCardPath(cid)

	// check if card already exists
	exist, _ := FileExist(path)
	if exist {
		return cid, nil
	}

	// save card
	card := Card{title, description, coverURL, destURL}
	file, err := json.MarshalIndent(card, "", "\t")
	if err != nil {
		return "", err
	}
	f, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.Write(file)
	if err != nil {
		return "", err
	}
	go RemoveWithDelay(path, cardSpan)
	return cid, nil
}

// GetCard - read card
func GetCard(cid string) (*Card, error) {
	card := Card{}
	cardBytes, err := ioutil.ReadFile(getCardPath(cid))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(cardBytes, &card)
	return &card, err
}

func getCardPath(cid string) string {
	return path.Join(Cwd, CardPath, fmt.Sprintf("temp_card_%s.json", cid))
}

func getCardID(key string) string {
	key = time.Now().Format("20060102") + key
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}
