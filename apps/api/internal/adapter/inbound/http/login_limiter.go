package http

import (
	"sync"
	"time"
)

// loginLimiter est un verrou anti-brute-force en mémoire, par clé (IP).
// Après maxAttempts échecs consécutifs, la clé est verrouillée pendant lockDuration.
type loginLimiter struct {
	mu           sync.Mutex
	attempts     map[string]*attemptInfo
	maxAttempts  int
	lockDuration time.Duration
	now          func() time.Time
}

type attemptInfo struct {
	failures    int
	lockedUntil time.Time
}

func newLoginLimiter(maxAttempts int, lockDuration time.Duration, now func() time.Time) *loginLimiter {
	return &loginLimiter{
		attempts:     make(map[string]*attemptInfo),
		maxAttempts:  maxAttempts,
		lockDuration: lockDuration,
		now:          now,
	}
}

// allowed indique si une tentative est autorisée pour cette clé.
func (l *loginLimiter) allowed(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.attempts[key]
	if a == nil || a.lockedUntil.IsZero() {
		return true
	}
	return !l.now().Before(a.lockedUntil)
}

// recordFailure enregistre un échec ; verrouille au-delà du seuil.
func (l *loginLimiter) recordFailure(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	a := l.attempts[key]
	if a == nil {
		a = &attemptInfo{}
		l.attempts[key] = a
	}
	a.failures++
	if a.failures >= l.maxAttempts {
		a.lockedUntil = l.now().Add(l.lockDuration)
		a.failures = 0
	}
}

// reset efface le compteur d'une clé (après un succès).
func (l *loginLimiter) reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}
