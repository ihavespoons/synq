package daemon

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
)

const debounceInterval = 2 * time.Second

// Watcher wraps fsnotify with debouncing.
type Watcher struct {
	fsw       *fsnotify.Watcher
	onChange  func()
	log       *zerolog.Logger
	mu        sync.Mutex
	timer     *time.Timer
	watching  map[string]bool
}

// NewWatcher creates a new file watcher.
func NewWatcher(onChange func(), log *zerolog.Logger) (*Watcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &Watcher{
		fsw:      fsw,
		onChange: onChange,
		log:      log,
		watching: make(map[string]bool),
	}, nil
}

// WatchPaths adds parent directories of the given paths to the watcher.
func (w *Watcher) WatchPaths(paths []string) {
	dirs := make(map[string]bool)
	for _, p := range paths {
		dir := filepath.Dir(p)
		dirs[dir] = true
	}
	for dir := range dirs {
		if w.watching[dir] {
			continue
		}
		if err := w.fsw.Add(dir); err != nil {
			w.log.Warn().Err(err).Str("dir", dir).Msg("failed to watch directory")
			continue
		}
		w.watching[dir] = true
		w.log.Debug().Str("dir", dir).Msg("watching directory")
	}
}

// Start begins listening for events in a goroutine.
func (w *Watcher) Start() {
	go func() {
		for {
			select {
			case event, ok := <-w.fsw.Events:
				if !ok {
					return
				}
				if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
					continue
				}
				w.log.Debug().Str("file", event.Name).Str("op", event.Op.String()).Msg("file changed")
				w.debounce()

			case err, ok := <-w.fsw.Errors:
				if !ok {
					return
				}
				w.log.Error().Err(err).Msg("watcher error")
			}
		}
	}()
}

func (w *Watcher) debounce() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.timer != nil {
		w.timer.Stop()
	}
	w.timer = time.AfterFunc(debounceInterval, w.onChange)
}

// Close shuts down the watcher.
func (w *Watcher) Close() error {
	return w.fsw.Close()
}
