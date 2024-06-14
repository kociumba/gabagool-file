package main

import (
	"fmt"

	"github.com/charmbracelet/log"
	"golang.org/x/sys/windows/registry"
)

func getDefaultTextEditor() string {
	// Open the registry key for .txt file association
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.txt\UserChoice`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("Failed to open registry key: %v", err)
	}
	defer key.Close()

	log.Info(key)

	// Read the ProgId value
	progId, _, err := key.GetStringValue("ProgId")
	if err != nil {
		log.Fatalf("Failed to read ProgId value: %v", err)
	}

	log.Info(progId)

	// Open the registry key for the ProgId
	progIdKeyPath := fmt.Sprintf(`%s\shell\open\command`, progId)
	progIdKey, err := registry.OpenKey(registry.CLASSES_ROOT, progIdKeyPath, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("Failed to open ProgId registry key: %v", err)
	}
	defer progIdKey.Close()

	log.Info(progIdKey)

	// Read the command value
	command, _, err := progIdKey.GetStringValue("")
	if err != nil {
		log.Fatalf("Failed to read command value: %v", err)
	}

	log.Info(command)

	return command
}

func getDefaultImageViewer() string {
	// Open the registry key for .png file association
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.png\UserChoice`, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("Failed to open registry key: %v", err)
	}
	defer key.Close()

	// Read the ProgId value
	progId, _, err := key.GetStringValue("ProgId")
	if err != nil {
		log.Fatalf("Failed to read ProgId value: %v", err)
	}

	// Open the registry key for the ProgId
	progIdKeyPath := fmt.Sprintf(`%s\shell\open\command`, progId)
	progIdKey, err := registry.OpenKey(registry.CLASSES_ROOT, progIdKeyPath, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("Failed to open ProgId registry key: %v", err)
	}
	defer progIdKey.Close()

	// Read the command value
	command, _, err := progIdKey.GetStringValue("")
	if err != nil {
		log.Fatalf("Failed to read command value: %v", err)
	}

	return command
}

func getDefaultBinaryEditor() string {
	// Open the registry key for .bin file association
	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.bin\UserChoice`, registry.QUERY_VALUE)
	if err != nil {
		// If .bin association is not found, try .hex
		key, err = registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Explorer\FileExts\.hex\UserChoice`, registry.QUERY_VALUE)
		if err != nil {
			log.Fatalf("Failed to open registry key: %v", err)
		}
	}
	defer key.Close()

	// Read the ProgId value
	progId, _, err := key.GetStringValue("ProgId")
	if err != nil {
		log.Fatalf("Failed to read ProgId value: %v", err)
	}

	// Open the registry key for the ProgId
	progIdKeyPath := fmt.Sprintf(`%s\shell\open\command`, progId)
	progIdKey, err := registry.OpenKey(registry.CLASSES_ROOT, progIdKeyPath, registry.QUERY_VALUE)
	if err != nil {
		log.Fatalf("Failed to open ProgId registry key: %v", err)
	}
	defer progIdKey.Close()

	// Read the command value
	command, _, err := progIdKey.GetStringValue("")
	if err != nil {
		log.Fatalf("Failed to read command value: %v", err)
	}

	return command
}
