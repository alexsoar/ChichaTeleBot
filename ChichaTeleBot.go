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

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

var (
	tempDir     = "/tmp"
	tempFileExt = ".ogg"
)

func main() {
	godotenv.Load()

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not found in environment or .env file")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	debug := os.Getenv("DEBUG")

	if debug == "" {
		bot.Debug = false
	} else if debug == "false" ||  debug == "False"  ||  debug == "FALSE" {
		bot.Debug = false
	} else if debug == "true" ||  debug == "True"  ||  debug == "TRUE" {
		bot.Debug = true
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
			go handleVoiceMessage(bot, update.Message, &wg)
		default:
			// Handle other message types or commands
		}
	}

	wg.Wait()
}

func handleVoiceMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, wg *sync.WaitGroup) {
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

	transcribedVoiceMessage, err := transcribeWithWhisper(voiceFilePath)
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

func transcribeWithWhisper(audioFilePath string) (string, error) {
	// Изменяем расширение файла на .txt
	textFilePath := changeFileExtension(audioFilePath, ".txt")

	defer os.Remove(audioFilePath)
	defer os.Remove(textFilePath)

	// Выполняем транскрибацию с использованием Whisper
	cmd := exec.Command(
		"whisper",
		audioFilePath,
		"--model", "medium",
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

	return formattedText, nil
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
