package gallery

import (
	"bytes"
	"github.com/egoist/mygo/plugins/terminal"
	"io"
)

// echoConnection is a byte-stream fixture, not a shell or input system.
// MyGo owns keyboard encoding, selection, rendering and scrolling.
type echoConnection struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func newEchoConnection() *echoConnection {
	r, w := io.Pipe()
	return &echoConnection{r, w}
}
func (c *echoConnection) Read(p []byte) (int, error) { return c.reader.Read(p) }
func (c *echoConnection) Write(p []byte) (int, error) {
	data := bytes.ReplaceAll(p, []byte("\r"), []byte("\r\n"))
	data = bytes.ReplaceAll(data, []byte{127}, []byte("\b \b"))
	_, err := c.writer.Write(data)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
func (c *echoConnection) Close() error {
	_ = c.reader.Close()
	return c.writer.Close()
}

func newPreviewTerminal() (*terminal.Terminal, error) {
	conn := newEchoConnection()
	term, err := terminal.New(terminal.Options{Conn: conn, Scrollback: 1 << 20})
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	term.Feed([]byte(terminalSample))
	return term, nil
}

const terminalSample = "\x1b[1mMyGo Terminal\x1b[0m\r\n" +
	"\x1b[31mRed  \x1b[32mGreen  \x1b[33mYellow  \x1b[34mBlue  \x1b[35mMagenta  \x1b[36mCyan\x1b[0m\r\n" +
	"\x1b[1mBold\x1b[0m  \x1b[3mItalic\x1b[0m  \x1b[4mUnderline\x1b[0m  \x1b[7mInverse\x1b[0m\r\n" +
	"\x1b[38;2;255;128;64mTrue color\x1b[0m  中文 日本語 🌍\r\n\r\n"
