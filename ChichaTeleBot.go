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

// Define global variables
var (
	tempDir     = "/tmp"
	tempFileExt = ".ogg"
)

// Define global variables for caching
var lastRequestTime time.Time
var cachedText string
var cacheDuration = time.Hour

func main() {
	// Get the current value of the PATH variable
	currentPath := os.Getenv("PATH")

	// Add a new path to the current PATH value
	newPath := "/venv/bin"
	newPathValue := fmt.Sprintf("%s:%s", newPath, currentPath)

	// Set the new value for the PATH variable
	err := os.Setenv("PATH", newPathValue)
	if err != nil {
		fmt.Println("Error setting PATH variable:", err)
		return
	}

	// Load environment variables from a .env file
	godotenv.Load()

	// Get the Telegram bot token from environment variables
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN not found in environment or .env file")
	}

	// Get the DEBUG and MODEL environment variables with default values
	debug := os.Getenv("DEBUG")
	Model := os.Getenv("MODEL")
	if Model == "" {
		Model = "medium"
	} else if strings.EqualFold(Model, "small") {
		Model = "small"
	} else if strings.EqualFold(Model, "medium") {
		Model = "medium"
	} else if strings.EqualFold(Model, "large") {
		Model = "large"
	} else {
		Model = "medium"
	}

	// Create a new Telegram bot
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	// Set bot debugging based on the DEBUG environment variable
	if strings.EqualFold(debug, "true") {
		bot.Debug = true
	} else {
		bot.Debug = false
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Move "medium" model to memory partition to speed up starting of the bot transcription routines:
	err = pushModelToMemory()
	if err != nil {
	   log.Println(err)
	}

	// Set up updates channel and wait group for handling voice messages concurrently
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	var wg sync.WaitGroup

	// Process incoming updates
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

	// Wait for all voice message processing to finish
	wg.Wait()
}

// Function to handle incoming voice messages
func handleVoiceMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, wg *sync.WaitGroup, Model string) {
	defer wg.Done()

	// Send a typing action
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

	// Transcribe the voice message using Whisper
	transcribedVoiceMessage, err := transcribeWithWhisper(voiceFilePath, Model)
	if err != nil {
		log.Printf("Error transcribing voice file with whisper: %v", err)
		return
	}

	// Send the transcribed text as a reply
	reply := tgbotapi.NewMessage(message.Chat.ID, transcribedVoiceMessage)
	_, err = bot.Send(reply)
	if err != nil {
		log.Printf("Error sending reply: %v", err)
	}
}

// Function to change the file extension
func changeFileExtension(voiceFilePath, newExtension string) string {
	fileName := filepath.Base(voiceFilePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	return filepath.Join(filepath.Dir(voiceFilePath), fileNameWithoutExt+newExtension)
}

// Function to transcribe a voice file using Whisper
func transcribeWithWhisper(audioFilePath string, Model string) (string, error) {
	// Change the file extension to .txt
	textFilePath := changeFileExtension(audioFilePath, ".txt")

	defer os.Remove(audioFilePath)
	defer os.Remove(textFilePath)

	// Execute transcription using Whisper
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

	// Read text from the file
	text, err := ioutil.ReadFile(textFilePath)
	if err != nil {
		return string(text), err
	}

	// Remove extra line breaks and spaces
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

// Function to update formatted text with additional information
func updateFormattedText(currentText string) string {
	// Check the time of the last request
	if time.Since(lastRequestTime) < cacheDuration {
		// Use cached text
		return currentText + "\n" + cachedText
	}

	// URL for obtaining additional text
	url := "https://raw.githubusercontent.com/matveynator/ChichaTeleBot/main/COPYRIGHT.md"

	// Create an HTTP client with a 5-second timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Attempt to get text from the website
	response, err := client.Get(url)
	if err != nil {
		// In case of an error (e.g., timeout), use default text "@ChichaTeleBot"
		copyrightText := "@ChichaTeleBot"
		return currentText + "\n" + copyrightText
	}
	defer response.Body.Close()

	// Read text from the website
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// In case of a reading error, use default text "@ChichaTeleBot"
		copyrightText := "@ChichaTeleBot"
		return currentText + "\n" + copyrightText
	}

	// Use the obtained text
	copyrightText := string(body)

	// Update cached data
	cachedText = copyrightText
	lastRequestTime = time.Now()

	return currentText + "\n" + cachedText
}

func pushModelToMemory() error {
	cmd := exec.Command("rsync", "-avP", "/root/models/*", "/root/.cache/whisper/")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error executing rsync: %v", err)
	}
	return nil
}

