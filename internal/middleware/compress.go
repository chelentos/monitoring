package middleware

import (
	"compress/gzip"
	"mime"
	"net/http"
	"strconv"
	"strings"
)

type ResponseWriterWithCompress struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
	headerSet  bool
}

func (rw *ResponseWriterWithCompress) WriteHeader(status int) {
	if rw.headerSet {
		return
	}

	rw.headerSet = true

	header := rw.Header()
	contentType, _, _ := mime.ParseMediaType(header.Get("Content-Type"))
	if contentType == "application/json" || contentType == "text/html" {
		header.Set("Content-Encoding", "gzip")
		rw.gzipWriter = gzip.NewWriter(rw.ResponseWriter)
	}

	rw.ResponseWriter.WriteHeader(status)
}

func (rw *ResponseWriterWithCompress) Write(data []byte) (int, error) {
	if !rw.headerSet {
		rw.WriteHeader(http.StatusOK)
	}

	if rw.gzipWriter != nil {
		return rw.gzipWriter.Write(data)
	}

	return rw.ResponseWriter.Write(data)
}

func WithCompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !acceptsGzip(r.Header) {
			next.ServeHTTP(w, r)
			return
		}

		rwc := ResponseWriterWithCompress{
			ResponseWriter: w,
		}

		next.ServeHTTP(&rwc, r)

		if rwc.gzipWriter != nil {
			defer rwc.gzipWriter.Close()
		}
	})
}

func encodingList(header http.Header, key string) []string {
	var encodings []string
	for _, line := range header.Values(key) {
		for _, part := range strings.Split(line, ",") {
			if value := strings.TrimSpace(part); value != "" {
				encodings = append(encodings, value)
			}
		}
	}
	return encodings
}

func acceptsGzip(header http.Header) bool {
	wildcardAllowed := false
	for _, value := range encodingList(header, "Accept-Encoding") {
		name, parameter, hasParameter := strings.Cut(value, ";")
		allowed := true
		if hasParameter {
			key, weight, _ := strings.Cut(parameter, "=")
			q, err := strconv.ParseFloat(strings.TrimSpace(weight), 64)
			allowed = strings.EqualFold(strings.TrimSpace(key), "q") && err == nil && q > 0 && q <= 1
		}
		switch strings.ToLower(strings.TrimSpace(name)) {
		case "gzip":
			return allowed
		case "*":
			wildcardAllowed = allowed
		}
	}
	return wildcardAllowed
}

func WithDecompress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encodings := encodingList(r.Header, "Content-Encoding")
		if len(encodings) == 0 {
			next.ServeHTTP(w, r)
			return
		}
		if len(encodings) != 1 || !strings.EqualFold(encodings[0], "gzip") {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}

		rdr, err := gzip.NewReader(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		defer rdr.Close()

		r.Body = rdr
		next.ServeHTTP(w, r)
	})
}
