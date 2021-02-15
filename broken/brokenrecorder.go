package brokenrecorder

import (
	"bytes"
	"errors"
	"net/http"
)

// BrokenRecorderS is a place holder for test usage
type BrokenRecorderS struct {
	Code        int
	HeaderMap   http.Header
	Body        *bytes.Buffer
	snapHeader  http.Header // snapshot of HeaderMap at first Write
	wroteHeader bool
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorder() *BrokenRecorderS {
	return &BrokenRecorderS{}
}

// Write implements http.ResponseWriter. The data in buf is written to
// rw.Body, if not nil.
func (rw *BrokenRecorderS) Write(data []byte) (n int, err error) {
	//	rw.writeHeader(buf, "")
	if rw.Body == nil {
		b := make([]byte, 0)
		rw.Body = bytes.NewBuffer(b)
	}
	rw.Body.WriteString("This is the broken write")
	return 0, errors.New("Broken write")
}

// =====

// Header implements http.ResponseWriter. It returns the response
// headers to mutate within a handler. To test the headers that were
// written after a handler completes, use the Result method and see
// the returned Response value's Header.
func (rw *BrokenRecorderS) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

// WriteHeader implements http.ResponseWriter.
func (rw *BrokenRecorderS) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.Code = code
	rw.wroteHeader = true
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
	rw.snapHeader = rw.HeaderMap.Clone()
}
