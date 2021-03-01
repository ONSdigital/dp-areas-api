package brokenrecorder

import (
	"bytes"
	"errors"
	"net/http"
)

// BrokeRecorder is a place holder for test usage
type BrokeRecorder struct {
	Code        int
	HeaderMap   http.Header
	Body        *bytes.Buffer
	snapHeader  http.Header // snapshot of HeaderMap at first Write
	wroteHeader bool
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorder() *BrokeRecorder {
	return &BrokeRecorder{}
}

// Write implements http.ResponseWriter. The data in buf is written to
// rw.Body, if not nil.
func (rw *BrokeRecorder) Write(data []byte) (n int, err error) {
	//	rw.writeHeader(buf, "")
	if rw.Body == nil {
		b := make([]byte, 0)
		rw.Body = bytes.NewBuffer(b)
	}
	rw.Body.WriteString("This is the broken write")
	return 0, errors.New("broken write")
}

// =====

// Header implements http.ResponseWriter. It returns the response
// headers to mutate within a handler. To test the headers that were
// written after a handler completes, use the Result method and see
// the returned Response value's Header.
func (rw *BrokeRecorder) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

// WriteHeader implements http.ResponseWriter.
func (rw *BrokeRecorder) WriteHeader(code int) {
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
