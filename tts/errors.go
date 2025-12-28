package tts

import "errors"

var (
	// ErrNoAvailableProvider is returned when no provider is available.
	ErrNoAvailableProvider = errors.New("tts: no available provider")

	// ErrVoiceNotFound is returned when a voice ID is not found.
	ErrVoiceNotFound = errors.New("tts: voice not found")

	// ErrInvalidConfig is returned when the synthesis config is invalid.
	ErrInvalidConfig = errors.New("tts: invalid configuration")

	// ErrRateLimited is returned when the provider rate limits the request.
	ErrRateLimited = errors.New("tts: rate limited")

	// ErrQuotaExceeded is returned when the provider quota is exceeded.
	ErrQuotaExceeded = errors.New("tts: quota exceeded")

	// ErrStreamClosed is returned when attempting to use a closed stream.
	ErrStreamClosed = errors.New("tts: stream closed")
)
