// Package transport provides audio transport protocols for voice agents.
package transport

import (
	"context"
	"io"
	"net"
)

// Config configures a transport connection.
type Config struct {
	// SampleRate is the audio sample rate in Hz.
	SampleRate int

	// Channels is the number of audio channels.
	Channels int

	// Encoding is the audio encoding ("pcm", "opus", "g711").
	Encoding string

	// BufferSizeMs is the audio buffer size in milliseconds.
	BufferSizeMs int
}

// Connection represents an active transport connection.
type Connection interface {
	// ID returns the connection identifier.
	ID() string

	// AudioIn returns a writer for sending audio to the remote.
	AudioIn() io.WriteCloser

	// AudioOut returns a reader for receiving audio from the remote.
	AudioOut() io.Reader

	// Events returns a channel for transport events.
	Events() <-chan Event

	// Close closes the connection.
	Close() error

	// RemoteAddr returns the remote address.
	RemoteAddr() net.Addr
}

// Event represents a transport event.
type Event struct {
	// Type is the event type.
	Type EventType

	// Data contains event-specific data.
	Data any

	// Error contains any error.
	Error error
}

// EventType identifies the type of transport event.
type EventType string

const (
	// EventConnected indicates connection established.
	EventConnected EventType = "connected"

	// EventDisconnected indicates connection closed.
	EventDisconnected EventType = "disconnected"

	// EventAudioStarted indicates audio stream started.
	EventAudioStarted EventType = "audio_started"

	// EventAudioStopped indicates audio stream stopped.
	EventAudioStopped EventType = "audio_stopped"

	// EventError indicates a transport error.
	EventError EventType = "error"

	// EventDTMF indicates DTMF tone received (telephony).
	EventDTMF EventType = "dtmf"
)

// Transport defines the interface for audio transport protocols.
type Transport interface {
	// Name returns the transport name.
	Name() string

	// Protocol returns the protocol type ("webrtc", "sip", "websocket", etc).
	Protocol() string

	// Listen starts listening for incoming connections.
	Listen(ctx context.Context, addr string) (<-chan Connection, error)

	// Connect initiates an outbound connection.
	Connect(ctx context.Context, addr string, config Config) (Connection, error)

	// Close shuts down the transport.
	Close() error
}

// WebRTCTransport provides WebRTC-based audio transport.
type WebRTCTransport interface {
	Transport

	// CreateOffer creates an SDP offer for WebRTC negotiation.
	CreateOffer(ctx context.Context) (string, error)

	// HandleAnswer processes an SDP answer.
	HandleAnswer(ctx context.Context, sdp string) error

	// AddICECandidate adds an ICE candidate.
	AddICECandidate(ctx context.Context, candidate string) error

	// OnICECandidate sets the ICE candidate callback.
	OnICECandidate(callback func(candidate string))
}

// SIPTransport provides SIP-based audio transport.
type SIPTransport interface {
	Transport

	// Register registers with a SIP server.
	Register(ctx context.Context, server, username, password string) error

	// Invite initiates a SIP INVITE.
	Invite(ctx context.Context, uri string) (Connection, error)

	// OnInvite sets the incoming INVITE handler.
	OnInvite(handler func(conn Connection, from string) bool)
}

// TelephonyTransport extends Transport with telephony features.
type TelephonyTransport interface {
	Transport

	// SendDTMF sends DTMF tones.
	SendDTMF(conn Connection, digits string) error

	// OnDTMF sets the DTMF handler.
	OnDTMF(handler func(conn Connection, digit string))

	// Transfer transfers the call to another number.
	Transfer(conn Connection, target string) error

	// Hold places the call on hold.
	Hold(conn Connection) error

	// Unhold resumes a held call.
	Unhold(conn Connection) error
}
