package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"gopass/internal/models"
)

type Storage struct {
	passwords []models.Password
	notes     []models.Note
	key       []byte
	mu        sync.RWMutex
}

func NewStorage(pin string) *Storage {
	// Use PIN to derive encryption key
	key := sha256.Sum256([]byte(pin))
	return &Storage{
		passwords: make([]models.Password, 0),
		notes:     make([]models.Note, 0),
		key:       key[:],
	}
}

func (s *Storage) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (s *Storage) decrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func (s *Storage) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data := models.ExportData{
		Passwords: s.passwords,
		Notes:     s.notes,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encrypted, err := s.encrypt(jsonData)
	if err != nil {
		return err
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dataPath := filepath.Join(configDir, "gopass", "data.enc")
	return os.WriteFile(dataPath, encrypted, 0600)
}

func (s *Storage) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	dataPath := filepath.Join(configDir, "gopass", "data.enc")
	encrypted, err := os.ReadFile(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No data file yet
		}
		return err
	}

	decrypted, err := s.decrypt(encrypted)
	if err != nil {
		return err
	}

	var data models.ExportData
	if err := json.Unmarshal(decrypted, &data); err != nil {
		return err
	}

	s.passwords = data.Passwords
	s.notes = data.Notes
	return nil
}

// Password operations
func (s *Storage) AddPassword(p models.Password) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.passwords = append(s.passwords, p)
	return s.Save()
}

func (s *Storage) UpdatePassword(p models.Password) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.passwords {
		if existing.ID == p.ID {
			s.passwords[i] = p
			return s.Save()
		}
	}
	return errors.New("password not found")
}

func (s *Storage) DeletePassword(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.passwords {
		if p.ID == id {
			s.passwords = append(s.passwords[:i], s.passwords[i+1:]...)
			return s.Save()
		}
	}
	return errors.New("password not found")
}

func (s *Storage) GetPasswords() []models.Password {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Password{}, s.passwords...)
}

// Note operations
func (s *Storage) AddNote(n models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notes = append(s.notes, n)
	return s.Save()
}

func (s *Storage) UpdateNote(n models.Note) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, existing := range s.notes {
		if existing.ID == n.ID {
			s.notes[i] = n
			return s.Save()
		}
	}
	return errors.New("note not found")
}

func (s *Storage) DeleteNote(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, n := range s.notes {
		if n.ID == id {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return s.Save()
		}
	}
	return errors.New("note not found")
}

func (s *Storage) GetNotes() []models.Note {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]models.Note{}, s.notes...)
}

// Search operations
func (s *Storage) Search(query string) models.SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result models.SearchResult

	// Search passwords
	for _, p := range s.passwords {
		if contains(p.Name, query) || contains(p.Username, query) || contains(p.Note, query) {
			result.Passwords = append(result.Passwords, p)
		}
	}

	// Search notes
	for _, n := range s.notes {
		if contains(n.Title, query) || contains(n.Content, query) {
			result.Notes = append(result.Notes, n)
		}
	}

	return result
}

func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) > 0 && (s == substr || contains(s, substr))
}

// Export/Import operations
func (s *Storage) Export() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := models.ExportData{
		Passwords: s.passwords,
		Notes:     s.notes,
	}

	return data.ToJSON()
}

func (s *Storage) Import(data []byte) error {
	var importData models.ExportData
	if err := importData.FromJSON(data); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Merge imported data with existing data
	s.passwords = append(s.passwords, importData.Passwords...)
	s.notes = append(s.notes, importData.Notes...)

	return s.Save()
}
