package stt

import "errors"

var (
	// ErrNoAvailableProvider is returned when no provider is available.
	ErrNoAvailableProvider = errors.New("stt: no available provider")

	// ErrStreamingNotSupported is returned when streaming is not supported.
	ErrStreamingNotSupported = errors.New("stt: streaming not supported by any provider")

	// ErrInvalidAudio is returned when the audio data is invalid.
	ErrInvalidAudio = errors.New("stt: invalid audio data")

	// ErrInvalidConfig is returned when the transcription config is invalid.
	ErrInvalidConfig = errors.New("stt: invalid configuration")

	// ErrAudioTooLong is returned when audio exceeds provider limits.
	ErrAudioTooLong = errors.New("stt: audio too long")

	// ErrAudioTooShort is returned when audio is too short to transcribe.
	ErrAudioTooShort = errors.New("stt: audio too short")

	// ErrRateLimited is returned when the provider rate limits the request.
	ErrRateLimited = errors.New("stt: rate limited")

	// ErrQuotaExceeded is returned when the provider quota is exceeded.
	ErrQuotaExceeded = errors.New("stt: quota exceeded")

	// ErrUnsupportedLanguage is returned when the language is not supported.
	ErrUnsupportedLanguage = errors.New("stt: unsupported language")

	// ErrUnsupportedFormat is returned when the audio format is not supported.
	ErrUnsupportedFormat = errors.New("stt: unsupported audio format")

	// ErrStreamClosed is returned when attempting to use a closed stream.
	ErrStreamClosed = errors.New("stt: stream closed")
)
