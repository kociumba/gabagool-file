package gabagool

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zhuangsirui/binpacker"
)

// SaveWithBitPacking opens a bit packed gabagool file
func (g *GabagoolFile) OpenWithBitPacking(path string) (*GabagoolFile, error) {
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

	// Read the packed data from the file
	packedData, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// Create an unpacker to unpack the data
	unpacker := binpacker.NewUnpacker(binary.LittleEndian, bytes.NewReader(packedData))

	// Unpack the data
	timestamp, err := unpacker.ShiftUint64()
	if err != nil {
		return nil, err
	}

	hash, err := unpacker.ShiftBytes(sha256.Size)
	if err != nil {
		return nil, err
	}

	dataTypeByte, err := unpacker.ShiftByte()
	if err != nil {
		return nil, err
	}
	dataType := DataTypes(dataTypeByte)

	length, err := unpacker.ShiftUint32()
	if err != nil {
		return nil, err
	}

	data, err := unpacker.ShiftBytes(uint64(length))
	if err != nil {
		return nil, err
	}

	if unpacker.Error() != nil {
		return nil, unpacker.Error()
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

// SaveWithBitPacking saves the gabagool file with bit packing
//
// admitedly does not seem to do anything
func (g *GabagoolFile) SaveWithBitPacking(path string, file *GabagoolFile) error {
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

	// Create a buffer to pack the data
	buf := new(bytes.Buffer)
	packer := binpacker.NewPacker(binary.LittleEndian, buf)

	// Pack the data into the buffer
	packer.PushUint64(file.Timestamp)
	packer.PushBytes(file.Hash)
	packer.PushByte(byte(file.DataType))
	packer.PushUint32(file.Length)
	packer.PushBytes(file.Data)

	if packer.Error() != nil {
		return packer.Error()
	}

	// Write the packed data to the file
	_, err = f.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
