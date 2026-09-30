package dispatcher

import (
	"io"
	"testing"
	"time"

	"github.com/xtls/xray-core/common/buf"
)

type timeoutInterruptRecorder struct{ interrupted bool }

func (r *timeoutInterruptRecorder) ReadMultiBuffer() (buf.MultiBuffer, error) { return nil, io.EOF }
func (r *timeoutInterruptRecorder) ReadMultiBufferTimeout(time.Duration) (buf.MultiBuffer, error) {
	return nil, io.EOF
}
func (r *timeoutInterruptRecorder) Interrupt() { r.interrupted = true }

// Links handed over with DispatchLink carry readers that are not pipes, such as
// the WireGuard and TUN inbounds' connection readers; outbounds still rely on
// interrupting them to end idle sessions.
func TestCachedReaderInterruptsNonPipeReader(t *testing.T) {
	inner := &timeoutInterruptRecorder{}
	r := &cachedReader{reader: inner}
	r.Interrupt()
	if !inner.interrupted {
		t.Fatal("cachedReader did not interrupt a non-pipe reader")
	}
}
