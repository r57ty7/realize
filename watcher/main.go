package watcher

import (
	"github.com/tockins/realize/settings"
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

// Watcher interface used by polling/fsnotify watching
type watcher interface {
	Add(path string) error
}

// Polling watcher
type pollWatcher struct {
	paths map[string]bool
}

// Log struct
type logWriter struct{}

// Blueprint struct contains a projects list
type Blueprint struct {
	*settings.Settings `yaml:"-"`
	Projects           []Project   `yaml:"projects,omitempty" json:"projects,omitempty"`
	Sync               chan string `yaml:"-"`
}

type tools struct {
	Fmt, Test, Generate, Vet tool
}

type tool struct {
	status  bool
	cmd     string
	options []string
	name    string
}

// Cmds go supported
type Cmds struct {
	Vet      bool `yaml:"vet,omitempty" json:"vet,omitempty"`
	Fmt      Cmd  `yaml:"fmt,omitempty" json:"fmt,omitempty"`
	Test     Cmd  `yaml:"test,omitempty" json:"test,omitempty"`
	Generate Cmd  `yaml:"generate,omitempty" json:"generate,omitempty"`
	Install  Cmd  `yaml:"install" json:"install"`
	Build    Cmd  `yaml:"build,omitempty" json:"build,omitempty"`
	Run      bool `yaml:"run,omitempty" json:"run,omitempty"`
}

// Cmd buildmode options
type Cmd struct {
	Status bool     `yaml:"status,omitempty" json:"status,omitempty"`
	Args   []string `yaml:"args,omitempty" json:"args,omitempty"`
}

// Watcher struct defines the livereload's logic
type Watcher struct {
	Paths   []string  `yaml:"paths" json:"paths"`
	Exts    []string  `yaml:"extensions" json:"extensions"`
	Ignore  []string  `yaml:"ignored_paths,omitempty" json:"ignored_paths,omitempty"`
	Preview bool      `yaml:"preview,omitempty" json:"preview,omitempty"`
	Scripts []Command `yaml:"scripts,omitempty" json:"scripts,omitempty"`
}

// Command options
type Command struct {
	Type    string `yaml:"type" json:"type"`
	Command string `yaml:"command" json:"command"`
	Path    string `yaml:"path,omitempty" json:"path,omitempty"`
	Global  bool   `yaml:"global,omitempty" json:"global,omitempty"`
	Output  bool   `yaml:"output,omitempty" json:"output,omitempty"`
}

// Buffer define an array buffer for each log files
type Buffer struct {
	StdOut []BufferOut `json:"stdOut"`
	StdLog []BufferOut `json:"stdLog"`
	StdErr []BufferOut `json:"stdErr"`
}

// BufferOut is used for exchange information between "realize cli" and "web realize"
type BufferOut struct {
	Time   time.Time `json:"time"`
	Text   string    `json:"text"`
	Path   string    `json:"path"`
	Type   string    `json:"type"`
	Stream string    `json:"stream"`
	Errors []string  `json:"errors"`
}

// Initialize the application
func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}
