package z

import (
	"github.com/rwxrob/scan/tk"
)

var WS = I{tk.SP, tk.TAB, tk.CR, tk.LF}
var EndLine = I{tk.LF, tk.CRLF, tk.CR}
