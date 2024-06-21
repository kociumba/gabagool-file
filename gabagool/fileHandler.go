package gabagool

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	kio "github.com/flanglet/kanzi-go/v2/io"
)

var (
	openFile  = (*os.File)(nil)
	err       = error(nil)
	fileCache = FileCache{}
	fileMutex = &sync.RWMutex{}
)

type DataTypes int

const (
	Bytes DataTypes = iota // data stored is an array of bytes
	Text                   // data stored is a string
	Image                  // data stored is an image, as an array of pixels that are compatible with the go image package
)

// Cache of the last file opened,
// usefull for saving on syscalls
//
// # This is an internal component and should not be used directly
type FileCache struct {
	Path     string
	FileData GabagoolFile
}

// Structure of the Gabagool™️ file, also houses all of it's methods
type GabagoolFile struct {
	Header    uint64    // header for veryfying the file
	Timestamp uint64    // timestamp when the file was created
	Hash      []byte    // sha256 hash of the data in the file
	DataType  DataTypes // type of the data stored in the file
	Length    uint32    // length of the data section
	Data      []byte    // data stored
}

// The same as calling CreateFile() followed by Save()
func (g *GabagoolFile) CreateAndSave(path string, dataType DataTypes, data []byte) error {
	fileMutex.Lock()
	defer fileMutex.Unlock()

	if fileCache.Path == path {
		return fmt.Errorf("file already exists in cache, get it with Open()")
	}

	if dataType != Bytes && dataType != Text && dataType != Image {
		return fmt.Errorf("unsupported data type")
	}

	file, err := g.CreateFile(dataType, data)
	if err != nil {
		return err
	}

	err = file.Save(path, file)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new file with the data provided and returns it for use
func (g *GabagoolFile) CreateFile(dataType DataTypes, data []byte) (*GabagoolFile, error) {
	header := uint64(0x67616261676F6F6C)
	length := uint32(len(data))
	timestamp := uint64(time.Now().Unix())
	hash := sha256.Sum256(data)

	g.Header = header
	g.Timestamp = timestamp
	g.Hash = hash[:]
	g.Length = length
	g.DataType = dataType
	g.Data = data

	if length != uint32(len(data)) {
		return nil, fmt.Errorf("error calculating length of the data section")
	}

	return g, nil
}

// Same as Open() with extra checks to make sure the file is valid
func (g *GabagoolFile) ParseFile(path string) (*GabagoolFile, error) {
	// Open the file using the Open function
	f, err := g.Open(path)
	if err != nil {
		return nil, err
	}

	// Ensure the opened file is not nil
	if f == nil {
		return nil, fmt.Errorf("file is nil")
	}

	// check if the data type is supported
	if f.DataType != Bytes && f.DataType != Text && f.DataType != Image {
		return nil, fmt.Errorf("unsupported data type")
	}

	// Verify the header
	if f.Header != 0x67616261676F6F6C {
		return nil, fmt.Errorf("invalid header")
	}

	// The rest of the fields have already been populated by the Open function
	return f, nil
}

// Opens the files from the provided path and returns it
func (g *GabagoolFile) Open(path string) (*GabagoolFile, error) {
	fileMutex.RLock()

	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		dir, _ := filepath.Split(path)
		os.Mkdir(dir, 0755)
		if err != nil {
			return nil, err
		}
		path = absPath
	}

	if fileCache.Path == path {
		return &fileCache.FileData, nil
	}

	// if not managed manually creats a deadlock
	fileMutex.RUnlock()

	f, err := os.Open(path + ".gabagool")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read the header
	header := make([]byte, 8)
	if _, err := f.Read(header); err != nil {
		return nil, err
	}

	// Read the timestamp
	timestampBytes := make([]byte, 8)
	if _, err := f.Read(timestampBytes); err != nil {
		return nil, err
	}
	timestamp := binary.LittleEndian.Uint64(timestampBytes)

	// Read the hash
	hash := make([]byte, sha256.Size)
	if _, err := f.Read(hash); err != nil {
		return nil, err
	}

	// Read the data type
	dataTypeBytes := make([]byte, 1)
	if _, err := f.Read(dataTypeBytes); err != nil {
		return nil, err
	}
	dataType := DataTypes(dataTypeBytes[0])

	// Read the length
	lengthBytes := make([]byte, 4)
	if _, err := f.Read(lengthBytes); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint32(lengthBytes)

	// Read the data
	data := make([]byte, length)
	if _, err := f.Read(data); err != nil {
		return nil, err
	}

	decoded, err := Decompress(string(data))
	data = []byte(decoded)
	if err != nil {
		return nil, err
	}

	// Populate the GabagoolFile structure
	g.Header = binary.BigEndian.Uint64(header)
	g.Timestamp = timestamp
	g.Hash = hash
	g.DataType = dataType
	g.Length = length
	g.Data = data

	// Cache the file data
	fileMutex.Lock()
	defer fileMutex.Unlock()
	fileCache.Path = path
	fileCache.FileData = *g

	return g, nil
}

// Saves the file to the provided path
func (g *GabagoolFile) Save(path string, file *GabagoolFile) error {
	if file == nil {
		return fmt.Errorf("file is nil")
	}

	// eval relative paths like "./test"
	if !filepath.IsAbs(path) {
		absPath, err := filepath.Abs(path)
		dir, _ := filepath.Split(path)
		os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
		path = absPath
	}

	f, err := os.Create(path + ".gabagool")
	if err != nil {
		return err
	}
	defer f.Close()

	header := make([]byte, 8)
	binary.BigEndian.PutUint64(header, file.Header)
	if _, err := f.Write(header); err != nil {
		return err
	}

	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, file.Timestamp)
	if _, err := f.Write(timestamp); err != nil {
		return err
	}

	if _, err := f.Write(file.Hash); err != nil {
		return err
	}

	dataType := make([]byte, 1)
	dataType[0] = byte(file.DataType)
	if _, err := f.Write(dataType); err != nil {
		return err
	}

	length := make([]byte, 4)
	binary.LittleEndian.PutUint32(length, file.Length)
	if _, err := f.Write(length); err != nil {
		return err
	}

	// Encode data with huffman for lossless compression
	data, err := Compress(string(file.Data))
	if err != nil {
		return err
	}

	log.Info(data) // for testing

	if _, err := f.Write([]byte(data)); err != nil {
		return err
	}

	return nil
}

func Compress(data string) (string, error) {
	var buf bytes.Buffer
	pr, pw := io.Pipe()

	go func() {
		defer pw.Close()
		w, err := kio.NewWriter(pw, "RLT+TEXT", "HUFFMAN", 1024, 4, false, 0, false)
		if err != nil {
			pw.CloseWithError(err)
			return
		}
		defer w.Close()

		_, err = w.Write([]byte(data))
		if err != nil {
			pw.CloseWithError(err)
			return
		}
	}()

	_, err := io.Copy(&buf, pr)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Decompress(data string) (string, error) {
	r, err := kio.NewReader(io.NopCloser(strings.NewReader(data)), 4)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 1024)
	var decoded strings.Builder
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
		if n == 0 {
			break
		}
		decoded.Write(buf[:n])
	}

	return decoded.String(), nil
}
