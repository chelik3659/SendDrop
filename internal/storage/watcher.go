package storage

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	Events  chan struct{}
	done    chan struct{}
}

func NewWatcher(baseDir string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	
	if err := watcher.Add(baseDir); err != nil {
		return nil, err
	}
	
	w := &Watcher{
		watcher: watcher,
		Events:  make(chan struct{}),
		done:    make(chan struct{}),
	}
	
	go w.watch()
	return w, nil
}

func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			// Игнорируем временные файлы
			if filepath.Ext(event.Name) == ".tmp" {
				continue
			}
			select {
			case w.Events <- struct{}{}:
			default:
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			// Логируем ошибку
			_ = err
		case <-w.done:
			return
		}
	}
}

func (w *Watcher) Close() error {
	close(w.done)
	return w.watcher.Close()
}