package errors

import "github.com/joomcode/errorx"

var FileSystemErrors = errorx.NewNamespace("file_system")
var CannotOpenFileError = FileSystemErrors.NewType("cannot_open_file")
var CannotWriteFileError = FileSystemErrors.NewType("cannot_write_file")
var MimeTypeError = FileSystemErrors.NewType("mime_type_error")
var CannotReadFileError = FileSystemErrors.NewType("cannot_read_file")
var FileSizeError = FileSystemErrors.NewType("file_size_error")