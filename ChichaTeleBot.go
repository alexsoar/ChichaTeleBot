package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
  "time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var (
	tempDir     = "/tmp"
	tempFileExt = ".ogg"
)

var lastRequestTime time.Time
var cachedText string
var cacheDuration = time.Hour

func main() {

	// Получаем текущее значение переменной PATH
	currentPath := os.Getenv("PATH")

	// Добавляем новый путь к текущему значению переменной PATH
	newPath := "/root/.local/bin/"
	newPathValue := fmt.Sprintf("%s:%s", newPath, currentPath)

	// Устанавливаем новое значение переменной PATH
	err := os.Setenv("PATH", newPathValue)
	if err != nil {
		fmt.Println("Ошибка при установке переменной PATH:", err)
		return
	}

	godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not found in environment or .env file")
	}

	debug := os.Getenv("DEBUG")

	Model := os.Getenv("MODEL")

	if Model == "" {
		Model = "medium"
	} else if Model == "small" ||  Model == "Small"  ||  Model == "SMALL" {
		Model = "small"
	} else if Model == "medium" ||  Model == "Medium"  ||  Model == "MEDIUM" {
		Model = "medium"
	} else if Model == "large" ||  Model == "large"  ||  Model == "large" {
		Model = "large"
	} else {
		Model = "medium"
	}
	
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	if debug == "" {
		bot.Debug = false
	} else if debug == "false" ||  debug == "False"  ||  debug == "FALSE" {
		bot.Debug = false
	} else if debug == "true" ||  debug == "True"  ||  debug == "TRUE" {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	var wg sync.WaitGroup

	for update := range updates {
		if update.Message == nil {
			continue
		}

		switch {
		case update.Message.Voice != nil:
			wg.Add(1)
			go handleVoiceMessage(bot, update.Message, &wg, Model)
		default:
			// Handle other message types or commands
		}
	}

	wg.Wait()
}

func handleVoiceMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, wg *sync.WaitGroup, Model string) {
	defer wg.Done()

	// Отправляем typing action
	typingMsg := tgbotapi.NewChatAction(message.Chat.ID, tgbotapi.ChatTyping)
	_, err := bot.Send(typingMsg)
	if err != nil {
		log.Printf("Error sending typing action: %v", err)
		return
	}

	fileID := message.Voice.FileID

	// Download the voice file
	voiceFilePath, err := downloadVoiceFile(bot, fileID)
	if err != nil {
		log.Printf("Error downloading voice file: %v", err)
		return
	}

	transcribedVoiceMessage, err := transcribeWithWhisper(voiceFilePath, Model)
	if err != nil {
		log.Printf("Error transcribing voice file with whisper: %v", err)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, transcribedVoiceMessage)
	_, err = bot.Send(reply)
	if err != nil {
		log.Printf("Error sending reply: %v", err)
	}

}

// Функция для изменения расширения файла
func changeFileExtension(voiceFilePath, newExtension string) string {
	fileName := filepath.Base(voiceFilePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	return filepath.Join(filepath.Dir(voiceFilePath), fileNameWithoutExt+newExtension)
}

func transcribeWithWhisper(audioFilePath string, Model string) (string, error) {
	// Изменяем расширение файла на .txt
	textFilePath := changeFileExtension(audioFilePath, ".txt")

	defer os.Remove(audioFilePath)
	defer os.Remove(textFilePath)

	// Выполняем транскрибацию с использованием Whisper
	cmd := exec.Command(
		"whisper",
		audioFilePath,
		"--model", Model,
		"--task", "transcribe",
		"--output_format", "txt",
		"--max_line_width", "0",
		"--highlight_words", "False",
		"--max_line_count", "0",
		"--max_words_per_line", "0",
		"--word_timestamps", "False",
		"--output_dir", tempDir,
	)
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("transcribing the voice file with Whisper: %v", err)
	}

	// Читаем текст из файла
	text, err := ioutil.ReadFile(textFilePath)
	if err != nil {
		return string(text), err
	}

	// Удаляем лишние переводы строк и пробелы
	formattedText := strings.ReplaceAll(string(text), "\n", " ")
	formattedText = strings.Join(strings.Fields(formattedText), " ")

	return updateFormattedText(formattedText), nil
}

// Function to download the voice file
func downloadVoiceFile(bot *tgbotapi.BotAPI, fileID string) (string, error) {
	voiceFile, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", fmt.Errorf("error getting voice file: %v", err)
	}

	response, err := http.Get("https://api.telegram.org/file/bot" + bot.Token + "/" + voiceFile.FilePath)
	if err != nil {
		return "", fmt.Errorf("error downloading voice file: %v", err)
	}
	defer response.Body.Close()

	outFile, err := os.Create(filepath.Join(tempDir, voiceFile.FileID+tempFileExt))
	if err != nil {
		return "", fmt.Errorf("error creating output file: %v", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		return "", fmt.Errorf("error copying response body to file: %v", err)
	}

	return outFile.Name(), nil
}

func updateFormattedText(currentText string) string {
	// Проверка времени последнего запроса
	if time.Since(lastRequestTime) < cacheDuration {
		// Использование кэшированного текста
		return currentText + "\n" + cachedText
	}

	// URL для запроса текста
	url := "https://raw.githubusercontent.com/matveynator/ChichaTeleBot/main/COPYRIGHT.md"

	// Создание HTTP-клиента с тайм-аутом 2 секунды
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Попытка получения текста с веб-сайта
	response, err := client.Get(url)
	if err != nil {
		// В случае ошибки (например, тайм-аута), использование текста "@ChichaTeleBot"
		copyrightText := "@ChichaTeleBot"
		return currentText + "\n" + copyrightText
	}
	defer response.Body.Close()

	// Чтение текста с веб-сайта
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// В случае ошибки чтения, использование текста "@ChichaTeleBot"
		copyrightText := "@ChichaTeleBot"
		return currentText + "\n" + copyrightText
	}

	// Использование полученного текста
	copyrightText := string(body)

	// Обновление кэшированных данных
	cachedText = copyrightText
	lastRequestTime = time.Now()

	return currentText + "\n" + cachedText
}

