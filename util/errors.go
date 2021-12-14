package util

import werr "github.com/newclarity/rep-go-offstage/wraperr"

var (
	ErrorApplyingTemplate     = werr.New("unable to apply template from %s")
	ErrorMatchingAndCapturing = werr.New("unable to match '%s' and capture all values specified in regexp '%s'")
)
