package awsSes

type FileContentType string

const (
	// ===== Documents =====
	ContentTypePDF  FileContentType = "application/pdf"
	ContentTypeTXT  FileContentType = "text/plain"
	ContentTypeDOC  FileContentType = "application/msword"
	ContentTypeDOCX FileContentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	ContentTypeXLSX FileContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentTypePPTX FileContentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"

	// ===== Images =====
	ContentTypePNG  FileContentType = "image/png"
	ContentTypeJPEG FileContentType = "image/jpeg"
	ContentTypeGIF  FileContentType = "image/gif"
	ContentTypeSVG  FileContentType = "image/svg+xml"

	// ===== Archives =====
	ContentTypeZIP FileContentType = "application/zip"
	ContentTypeRAR FileContentType = "application/vnd.rar"
	ContentType7Z  FileContentType = "application/x-7z-compressed"

	// ===== Data / Text =====
	ContentTypeJSON FileContentType = "application/json"
	ContentTypeCSV  FileContentType = "text/csv"

	// ===== Media =====
	ContentTypeMP3 FileContentType = "audio/mpeg"
	ContentTypeMP4 FileContentType = "video/mp4"

	// ===== Fallback =====
	ContentTypeOctetStream FileContentType = "application/octet-stream"
)

type EmailAttachment struct {
	Filename    string
	ContentType FileContentType
	Data        []byte
}
