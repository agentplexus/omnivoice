// Package agent provides voice agent orchestration for real-time conversations.
package agent

import (
	"context"
	"io"
	"time"
)

// Config configures a voice agent.
type Config struct {
	// Name is a human-readable name for the agent.
	Name string

	// SystemPrompt is the initial system prompt for the LLM.
	SystemPrompt string

	// VoiceID is the TTS voice to use.
	VoiceID string

	// Language is the primary language (BCP-47 code).
	Language string

	// STTProvider is the speech-to-text provider name.
	STTProvider string

	// TTSProvider is the text-to-speech provider name.
	TTSProvider string

	// LLMProvider is the LLM provider name.
	LLMProvider string

	// LLMModel is the specific LLM model to use.
	LLMModel string

	// MaxTurnDuration is the maximum duration for a single turn.
	MaxTurnDuration time.Duration

	// MaxSessionDuration is the maximum total session duration.
	MaxSessionDuration time.Duration

	// InterruptionMode controls how interruptions are handled.
	InterruptionMode InterruptionMode

	// Tools defines functions the agent can call.
	Tools []Tool

	// Webhooks configures event webhooks.
	Webhooks WebhookConfig
}

// InterruptionMode controls how user interruptions are handled.
type InterruptionMode string

const (
	// InterruptImmediate stops TTS immediately when user speaks.
	InterruptImmediate InterruptionMode = "immediate"

	// InterruptAfterSentence finishes current sentence before stopping.
	InterruptAfterSentence InterruptionMode = "after_sentence"

	// InterruptDisabled ignores interruptions.
	InterruptDisabled InterruptionMode = "disabled"
)

// Tool defines a function the voice agent can call.
type Tool struct {
	// Name is the function name.
	Name string

	// Description describes what the function does.
	Description string

	// Parameters defines the function parameters (JSON Schema).
	Parameters map[string]any

	// Handler is called when the tool is invoked.
	Handler ToolHandler
}

// ToolHandler processes a tool call and returns a result.
type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

// WebhookConfig configures event webhooks.
type WebhookConfig struct {
	// OnSessionStart is called when a session begins.
	OnSessionStart string

	// OnSessionEnd is called when a session ends.
	OnSessionEnd string

	// OnTurnComplete is called after each conversation turn.
	OnTurnComplete string

	// OnToolCall is called when a tool is invoked.
	OnToolCall string
}

// Session represents an active voice conversation session.
type Session interface {
	// ID returns the unique session identifier.
	ID() string

	// Start begins the voice session.
	Start(ctx context.Context) error

	// Stop ends the voice session gracefully.
	Stop(ctx context.Context) error

	// SendAudio sends audio data to the agent.
	SendAudio(audio []byte) error

	// ReceiveAudio returns a channel for receiving agent audio.
	ReceiveAudio() <-chan []byte

	// SendText sends text input to the agent (bypass STT).
	SendText(text string) error

	// Events returns a channel for session events.
	Events() <-chan Event

	// Transcript returns the conversation transcript so far.
	Transcript() []Turn

	// Metrics returns session performance metrics.
	Metrics() Metrics
}

// Turn represents a single conversation turn.
type Turn struct {
	// Role is "user" or "agent".
	Role string

	// Text is the transcribed/generated text.
	Text string

	// Timestamp is when the turn occurred.
	Timestamp time.Time

	// DurationMs is the turn duration in milliseconds.
	DurationMs int

	// ToolCalls contains any tool calls made during this turn.
	ToolCalls []ToolCall
}

// ToolCall represents a tool invocation during conversation.
type ToolCall struct {
	// Name is the tool name.
	Name string

	// Arguments is the parsed arguments.
	Arguments map[string]any

	// Result is the tool result.
	Result string

	// Error is any error from the tool call.
	Error string

	// DurationMs is the tool execution time.
	DurationMs int
}

// Event represents a session event.
type Event struct {
	// Type is the event type.
	Type EventType

	// Timestamp is when the event occurred.
	Timestamp time.Time

	// Data contains event-specific data.
	Data any

	// Error contains any error details.
	Error error
}

// EventType identifies the type of session event.
type EventType string

const (
	// EventSessionStarted indicates the session has started.
	EventSessionStarted EventType = "session_started"

	// EventSessionEnded indicates the session has ended.
	EventSessionEnded EventType = "session_ended"

	// EventUserSpeechStart indicates the user started speaking.
	EventUserSpeechStart EventType = "user_speech_start"

	// EventUserSpeechEnd indicates the user stopped speaking.
	EventUserSpeechEnd EventType = "user_speech_end"

	// EventUserTranscript contains user speech transcription.
	EventUserTranscript EventType = "user_transcript"

	// EventAgentThinking indicates the agent is processing.
	EventAgentThinking EventType = "agent_thinking"

	// EventAgentSpeechStart indicates the agent started speaking.
	EventAgentSpeechStart EventType = "agent_speech_start"

	// EventAgentSpeechEnd indicates the agent stopped speaking.
	EventAgentSpeechEnd EventType = "agent_speech_end"

	// EventAgentTranscript contains agent response text.
	EventAgentTranscript EventType = "agent_transcript"

	// EventToolCall indicates a tool was called.
	EventToolCall EventType = "tool_call"

	// EventInterruption indicates the user interrupted.
	EventInterruption EventType = "interruption"

	// EventError indicates an error occurred.
	EventError EventType = "error"
)

// Metrics contains session performance metrics.
type Metrics struct {
	// SessionDurationMs is total session duration.
	SessionDurationMs int

	// TurnCount is the number of conversation turns.
	TurnCount int

	// UserSpeechDurationMs is total user speech time.
	UserSpeechDurationMs int

	// AgentSpeechDurationMs is total agent speech time.
	AgentSpeechDurationMs int

	// AvgSTTLatencyMs is average STT processing time.
	AvgSTTLatencyMs int

	// AvgLLMLatencyMs is average LLM processing time.
	AvgLLMLatencyMs int

	// AvgTTSLatencyMs is average TTS processing time.
	AvgTTSLatencyMs int

	// AvgTotalLatencyMs is average end-to-end latency.
	AvgTotalLatencyMs int

	// InterruptionCount is number of user interruptions.
	InterruptionCount int

	// ToolCallCount is number of tool invocations.
	ToolCallCount int

	// ErrorCount is number of errors encountered.
	ErrorCount int
}

// Provider defines the interface for voice agent providers.
type Provider interface {
	// Name returns the provider name.
	Name() string

	// CreateSession creates a new voice session.
	CreateSession(ctx context.Context, config Config) (Session, error)

	// GetSession retrieves an existing session by ID.
	GetSession(ctx context.Context, sessionID string) (Session, error)

	// ListSessions lists active sessions.
	ListSessions(ctx context.Context) ([]string, error)
}

// TransportAdapter adapts a transport to a voice agent session.
type TransportAdapter interface {
	// Connect connects the transport to a session.
	Connect(ctx context.Context, session Session) error

	// Disconnect disconnects the transport from the session.
	Disconnect(ctx context.Context) error

	// AudioIn returns a writer for incoming audio.
	AudioIn() io.Writer

	// AudioOut returns a reader for outgoing audio.
	AudioOut() io.Reader
}
