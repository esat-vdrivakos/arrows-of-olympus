package assets

import "embed"

// FS embeds all game images, normalized sprite sheets, .sprite definitions, .music scores, and .sfx cues.
//
//go:embed images/*.png sprites/*.png sprites/*.sprite audio/*.music audio/*.sfx
var FS embed.FS
