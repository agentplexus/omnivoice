# Whisper MLX Provider

Whisper provider for OmniVoice speech-to-text, running locally on Apple Silicon via MLX.

## Features

- Local STT transcription (no cloud API required)
- Word-level timestamps
- Per-segment confidence scores
- Automatic language detection or explicit language selection
- gRPC over Unix Domain Socket

## Requirements

- Apple Silicon Mac (M1/M2/M3/M4)
- macOS 13+
- Python 3.11+ (ARM64)

MLX ships arm64-only wheels, so the server must run under a native arm64
Python. Build the virtual environment with `arch -arm64` (important if your
shell is running under Rosetta).

## Usage

### 1. Start the Server

```bash
cd server

# Create an ARM64 virtual environment
arch -arm64 python3 -m venv .venv

# Install dependencies and generate the proto stubs
arch -arm64 .venv/bin/pip install -r requirements.txt
./generate_proto.sh

# Start the server (default model: large-v3-turbo)
arch -arm64 .venv/bin/python whisper_server.py

# Or select a model explicitly
arch -arm64 .venv/bin/python whisper_server.py --model large-v3-turbo
```

The model downloads from Hugging Face on first use (~1.6GB for
`large-v3-turbo`) and is cached for subsequent runs. The server listens on
`unix:///tmp/omnivoice-whisper.sock` by default.

### 2. Use from Go

```go
import (
    "github.com/plexusone/omnivoice"
    _ "github.com/plexusone/omnivoice-core/providers/whisper-mlx"

    "github.com/plexusone/omnivoice-core/stt"
)

func main() {
    provider, _ := omnivoice.GetSTTProvider("whisper-mlx")

    // Transcribe raw audio bytes
    result, _ := provider.Transcribe(ctx, audioBytes, stt.TranscriptionConfig{
        Language:             "en-US",
        EnableWordTimestamps: true,
    })

    // Or transcribe a file directly
    result, _ = provider.TranscribeFile(ctx, "narration.wav", stt.TranscriptionConfig{})
}
```

Whisper expects an ISO 639-1 language code (e.g. `en`), but the Go provider
normalizes BCP-47 locale tags for you — `en-US` and `zh-Hans` become `en` and
`zh` before the request is sent, so either form works. Leave `Language` empty
to let Whisper auto-detect.

## Architecture

```
Go Application                     Python Server
┌───────────────────┐  gRPC/UDS   ┌───────────────────┐
│ whisper.Provider  │◄───────────►│ whisper_server.py │
│ (gRPC client)     │             │ (MLX + Whisper)   │
└───────────────────┘             └───────────────────┘
        unix:///tmp/omnivoice-whisper.sock
```

## Capability Interfaces

| Interface | Description |
|-----------|-------------|
| `stt.Provider` | Transcribe audio bytes, files, and URLs |
| Health check | Provider and model status |
| Model management | Load / unload / list Whisper models |
| Runtime info | Device type, memory usage, framework version |
