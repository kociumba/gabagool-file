package gabagool

import (
	"bytes"
	"crypto/sha256"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {
	data := []byte("test data")
	expectedHeader := uint64(0x67616261676F6F6C)
	expectedLength := uint32(len(data))
	expectedHash := sha256.Sum256(data)
	expectedDataType := Text

	g := &GabagoolFile{}
	file, err := g.CreateFile(expectedDataType, data)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if file.Header != expectedHeader {
		t.Errorf("expected header %v, got %v", expectedHeader, file.Header)
	}

	if file.Length != expectedLength {
		t.Errorf("expected length %v, got %v", expectedLength, file.Length)
	}

	if !bytes.Equal(file.Hash, expectedHash[:]) {
		t.Errorf("expected hash %x, got %x", expectedHash, file.Hash)
	}

	if file.DataType != expectedDataType {
		t.Errorf("expected data type %v, got %v", expectedDataType, file.DataType)
	}

	if !bytes.Equal(file.Data, data) {
		t.Errorf("expected data %v, got %v", data, file.Data)
	}

	// Check if the timestamp is within a reasonable range (e.g., ±5 seconds)
	now := time.Now().Unix()
	if file.Timestamp < uint64(now-5) || file.Timestamp > uint64(now+5) {
		t.Errorf("expected timestamp within range, got %v", file.Timestamp)
	}
}
