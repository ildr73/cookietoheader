package cookietoheader

import (
	"context"
	"net/http"
)

// Config - конфигурация плагина.
type Config struct {
	CookieName  string `json:"cookieName,omitempty"`  // Имя cookie
	QueryParam  string `json:"queryParam,omitempty"`  // Имя параметра query string
	HeaderName  string `json:"headerName,omitempty"`  // Имя заголовка
	CookiePath  string `json:"cookiePath,omitempty"`  // Путь cookie (по умолчанию "/")
	CookieMaxAge int    `json:"cookieMaxAge,omitempty"` // Время жизни cookie в секундах
}

// CreateConfig - создаёт конфигурацию по умолчанию.
func CreateConfig() *Config {
	return &Config{
		CookiePath:  "/",
		CookieMaxAge: 3600, // 1 час
	}
}

// CookieToHeader - структура плагина.
type CookieToHeader struct {
	next        http.Handler
	cookieName  string
	queryParam  string
	headerName  string
	cookiePath  string
	cookieMaxAge int
	name        string
}

// New - создаёт экземпляр плагина.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &CookieToHeader{
		next:        next,
		cookieName:  config.CookieName,
		queryParam:  config.QueryParam,
		headerName:  config.HeaderName,
		cookiePath:  config.CookiePath,
		cookieMaxAge: config.CookieMaxAge,
		name:        name,
	}, nil
}

// ServeHTTP - основная логика плагина.
func (p *CookieToHeader) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var value string

	// Проверяем наличие cookie
	cookie, err := req.Cookie(p.cookieName)
	if err == nil {
		// Если cookie найден, извлекаем его значение
		value = cookie.Value
	} else {
		// Если cookie не найден, проверяем параметр в query string
		query := req.URL.Query()
		value = query.Get(p.queryParam)

		if value != "" {
			// Если параметр найден, создаем cookie
			http.SetCookie(rw, &http.Cookie{
				Name:     p.cookieName,
				Value:    value,
				Path:     p.cookiePath,
				MaxAge:   p.cookieMaxAge,
				HttpOnly: true,
			})
		}
	}

	if value != "" {
		// Добавляем значение в заголовок
		req.Header.Set(p.headerName, value)
	}

	// Передаем запрос дальше
	p.next.ServeHTTP(rw, req)
}
