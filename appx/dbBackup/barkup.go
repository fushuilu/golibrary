package dbBackup

// 从 https://github.com/keighl/barkup 复制过来的，因为部分依赖无法下载

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//////////////

// Exporter is expected to export "something" to a file and return a complete `ExportResult` struct (`Path`, `MIME`, `Error`). If any error occurs during it's work, it should set the error to the result's `Error` attribute
type Exporter interface {
	Export() (*ExportResult, *Error)
}

// Error can ship a cmd output as well as the start interface. Useful for understanding why a system command (exec.Command) failed
type Error struct {
	err       error
	Cmd       string // 命令
	CmdOutput string // 命令输出结果
}

func (e Error) Error() string {
	return fmt.Sprintf("[%s]\nCmdOutput:[%s]\nErr:[%s]\n", e.Cmd, e.CmdOutput, e.err.Error())
}

func makeErr(err error, out string) *Error {
	if err != nil {
		return &Error{
			err:       err,
			CmdOutput: out,
		}
	}
	return nil
}

//////////////

// Storer takes an `ExportResult` and move it somewhere! To a cloud storage service, for instance...
type Storer interface {
	Store(result *ExportResult, directory string) *Error
}

//////////////

// ExportResult is the result of an export operation... duh
type ExportResult struct {
	// Path to exported file
	Path string
	// MIME type of the exported file (e.g. application/x-tar)
	MIME string
	// Any error that occured during `Export()`
	Error *Error
}

// To hands off an ExportResult to a `Storer` interface and invokes its Store() method. The directory argument is passed along too. If `store` is `nil`, the the method will simply move the export result to the specified directory (via the `mv` command)
func (x *ExportResult) To(directory string, store Storer) *Error {
	if !strings.HasSuffix(directory, "/") {
		directory += "/"
	}
	if store == nil {
		cmd := exec.Command("mv", x.Path, directory+x.Filename())
		out, err := cmd.Output()
		if err != nil {
			rst := makeErr(err, string(out))
			rst.Cmd = cmd.String()
			return rst
		}
		return nil
	}

	storeErr := store.Store(x, directory)
	if storeErr != nil {
		return storeErr
	}

	err := os.Remove(x.Path)
	return makeErr(err, "")
}

// Filename returns the just filename component of the `Path` attribute
func (x ExportResult) Filename() string {
	_, filename := filepath.Split(x.Path)
	return filename
}
