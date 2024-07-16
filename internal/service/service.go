package service

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

const (
	scheme = "http://"
	host   = "localhost:3000/"
)

var messages = make(map[string]string)

type Service struct {
	repo  Repository
	cache Cache
}

// 1. ссылка в кеше живет 10 минут после последнего запроса (задается через конфиг)
// 2. читать из кеша. если пустой, то идем в базу. если там нашли, то параллельно отдать ответ пользователю и положить запись в кеш
// 3. если нет ни в базе, ни в кеше: пишешем в базу, отдаем пользователю. (опционально) параллельно пишем в кеш (обоснавать решение)

type Cache interface {
	Get(string) (string, error)
	Set(m map[string]string) error
}

type Repository interface {
	GetByOriginalLink(link []byte) (string, error)
	Create(token string, original []byte) error
	GetByToken(token string) (string, error)
}

func NewService(repo Repository, cache Cache) *Service {
	return &Service{
		repo:  repo,
		cache: cache,
	}
}

func CreateToken(link []byte) string {
	sha256checksum := sha256.Sum256(link)
	return hex.EncodeToString(sha256checksum[:])[:5]
}

// Writer starts cache set pipeline at time.Duration
func (s *Service) Writer(t time.Duration, out chan int) {
	tick := time.NewTicker(t)
	for {
		select {
		case <-tick.C:
			reqMessages := messages
			messages = map[string]string{}
			err := s.cache.Set(reqMessages)
			if err != nil {
				out <- -1
				return
			}
		case <-out:
			tick.Stop()
			return
		}
	}
}

func (s *Service) GetShortUrl(link []byte) (string, error) {
	// Check cache
	token, err := s.cache.Get(string(link))
	if err != nil {
		return "", err
	}
	if token != "" {
		return scheme + host + token, nil
	}

	// 1. Check db
	token, err = s.repo.GetByOriginalLink(link)
	if err != nil {
		return "", err
	}
	if token != "" {
		return scheme + host + token, nil
	}

	// 2. Create short link
	token = CreateToken(link)
	err = s.repo.Create(token, link)
	if err != nil {
		return "", err
	}

	return scheme + host + token, nil
}

func (s *Service) GetLongUrl(shortLink string) (string, error) {
	//Check cache
	shortLink = shortLink[1:]
	link, err := s.cache.Get(shortLink)
	if err != nil {
		return "", err
	}
	if link != "" {
		return scheme + link, nil
	}

	// 1. Check db
	link, err = s.repo.GetByToken(shortLink)
	if err != nil {
		return "", err
	}
	if link != "" {
		messages[shortLink] = link
		return scheme + link, nil
	}

	return "", nil
}
