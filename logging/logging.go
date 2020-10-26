package logging

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-logr/logr"
)

type BackupSummary struct {
	MessageType         string  `json:"message_type"`
	FilesNew            int     `json:"files_new"`
	FilesChanged        int     `json:"files_changed"`
	FilesUnmodified     int     `json:"files_unmodified"`
	DirsNew             int     `json:"dirs_new"`
	DirsChanged         int     `json:"dirs_changed"`
	DirsUnmodified      int     `json:"dirs_unmodified"`
	DataBlobs           int     `json:"data_blobs"`
	TreeBlobs           int     `json:"tree_blobs"`
	DataAdded           int64   `json:"data_added"`
	TotalFilesProcessed int     `json:"total_files_processed"`
	TotalBytesProcessed int     `json:"total_bytes_processed"`
	TotalDuration       float64 `json:"total_duration"`
	SnapshotID          string  `json:"snapshot_id"`
}

type BackupEnvelope struct {
	MessageType string `json:"message_type,omitempty"`
	BackupStatus
	BackupSummary
	BackupError
}

type BackupStatus struct {
	PercentDone  float64  `json:"percent_done"`
	TotalFiles   int      `json:"total_files"`
	FilesDone    int      `json:"files_done"`
	TotalBytes   int      `json:"total_bytes"`
	BytesDone    int      `json:"bytes_done"`
	CurrentFiles []string `json:"current_files"`
	ErrorCount   int      `json:"error_count"`
}

// SummaryFunc takes the summed up status of the backup and will process this further like
// logging, metrics and webhooks.
type SummaryFunc func(summary BackupSummary, errorCount int, folder string, startTimestamp, endTimestamp int64)

type BackupOutputParser struct {
	log         logr.Logger
	errorCount  int
	lineCounter int
	summaryfunc SummaryFunc
	folder      string
}

type BackupError struct {
	Error struct {
		Op   string `json:"Op"`
		Path string `json:"Path"`
		Err  int    `json:"Err"`
	} `json:"error"`
	During string `json:"during"`
	Item   string `json:"item"`
}

type outFunc func(string) error

// New creates a writer which directly writes to the given logger function.
func New(out outFunc) io.Writer {
	return &writer{out}
}

// NewInfoWriter creates a writer which directly writes to the given logger using info level.
// It ensures that each line is handled seperately. This avoids mangled lines when parsing
// JSON outputs.
func NewInfoWriter(l logr.InfoLogger) io.Writer {
	return New((&LogInfoPrinter{l}).out)
}

// NewInfoWriter creates a writer which directly writes to the given logger using error level.
// It ensures that each line is handled seperately. This avoids mangled lines when parsing
// JSON outputs.
func NewErrorWriter(l logr.Logger) io.Writer {
	return New((&LogErrPrinter{l}).out)
}

type writer struct {
	out outFunc
}

func (w writer) Write(p []byte) (int, error) {

	scanner := bufio.NewScanner(bytes.NewReader(p))

	for scanner.Scan() {
		err := w.out(scanner.Text())
		if err != nil {
			return len(p), err
		}
	}

	return len(p), nil
}

type LogInfoPrinter struct {
	log logr.InfoLogger
}

func (l *LogInfoPrinter) out(s string) error {
	l.log.Info(s)
	return nil
}

type LogErrPrinter struct {
	Log logr.Logger
}

func (l *LogErrPrinter) out(s string) error {
	l.Log.Error(fmt.Errorf("error during command"), s)
	return nil
}

func NewBackupOutputParser(logger logr.Logger, folderName string, summaryfunc SummaryFunc) io.Writer {
	bop := &BackupOutputParser{
		log:         logger,
		folder:      folderName,
		summaryfunc: summaryfunc,
	}
	return New(bop.out)
}

func (b *BackupOutputParser) out(s string) error {
	envelope := &BackupEnvelope{}

	err := json.Unmarshal([]byte(s), envelope)
	if err != nil {
		b.log.Error(err, "can't decode restic json output", "string", s)
		return err
	}

	switch envelope.MessageType {
	case "error":
		b.errorCount++
		b.log.Error(fmt.Errorf("error occurred during backup"), envelope.Item+" during "+envelope.During+" "+envelope.Error.Op)
	case "status":
		// Restic does the json output with 60hz, which is a bit much...
		if b.lineCounter%60 == 0 {
			percent := envelope.PercentDone * 100
			b.log.Info("progress of backup", "percentage", fmt.Sprintf("%.2f%%", percent))
		}
		b.lineCounter++
	case "summary":
		b.log.Info("backup finished", "new files", envelope.FilesNew, "changed files", envelope.FilesChanged, "errors", b.errorCount)
		b.log.Info("stats", "time", envelope.TotalDuration, "bytes added", envelope.DataAdded, "bytes processed", envelope.TotalBytesProcessed)
		b.summaryfunc(envelope.BackupSummary, b.errorCount, b.folder, 1, time.Now().Unix())
	}
	return nil
}