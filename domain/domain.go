package domain

import "sync/atomic"

type Headers struct {
	H1, H2, H3, H4, H5, H6 int64
}

type Links struct {
	Internal, External, InAccessible int64
}

type Result struct {
	HtmlVersion string
	PageTitle   string
	IsLoginPage bool
	Headers
	Links
}

func (l *Links) AddInternal() {
	atomic.AddInt64(&l.Internal, 1)
}
func (l *Links) AddExternal() {
	atomic.AddInt64(&l.External, 1)
}
func (l *Links) AddInAccessible() {
	atomic.AddInt64(&l.InAccessible, 1)
}

func (h *Headers) AddH1() {
	atomic.AddInt64(&h.H1, 1)
}
func (h *Headers) AddH2() {
	atomic.AddInt64(&h.H2, 1)
}
func (h *Headers) AddH3() {
	atomic.AddInt64(&h.H3, 1)
}
func (h *Headers) AddH4() {
	atomic.AddInt64(&h.H4, 1)
}
func (h *Headers) AddH5() {
	atomic.AddInt64(&h.H5, 1)
}
func (h *Headers) AddH6() {
	atomic.AddInt64(&h.H6, 1)
}
