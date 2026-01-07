package middlewares

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/AhmedHossam777/go-mongo/internal/handlers"
)

type client struct {
	lastSeen time.Time
	count    int
}

var (
	clients = make(map[string]*client)
	mu      sync.Mutex
)

func init() {
	// Periodic cleanup of the clients map to prevent memory leaks
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// RateLimitMiddleware limits the number of requests a user can make per minute
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{lastSeen: time.Now(), count: 0}
			}

			clients[ip].lastSeen = time.Now()
			clients[ip].count++

			if clients[ip].count > 60 {
				mu.Unlock()
				handlers.RespondWithError(
					w, http.StatusTooManyRequests,
					"Too many requests. Please try again later.",
				)
				return
			}

			go func(ip string) {
				time.Sleep(time.Minute)
				mu.Lock()
				if c, ok := clients[ip]; ok {
					c.count--
				}
				mu.Unlock()
			}(ip)

			mu.Unlock()
			next.ServeHTTP(w, r)
		},
	)
}
