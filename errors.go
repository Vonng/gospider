package gospider

import "errors"

var ErrDropItem = errors.New("drop item")

// ErrNoSuchKey occurs when access map with non-exist key
var ErrNoSuchKey = errors.New("no such key")

// ErrValueIsNotString occurs when Data.GetString() meets a non-string value
var ErrValueIsNotString = errors.New("value is not string")

var ErrDownloadFail = errors.New("download fail")

var ErrNilRequest = errors.New("nil request")

var ErrDupeRequest = errors.New("duplicate request")

var ErrStopIteration = errors.New("stop iteration")

var ErrNilProcessorList = errors.New("nil processor list for pipeline")

var ErrNilProcessor = errors.New("nil processor")

var ErrNilItem = errors.New("nil item")

var ErrParse = errors.New("parse error")

var ErrContinue = errors.New("continue")

var ErrNilResponse = errors.New("nil response")

var ErrInvalidCallback = errors.New("invalid callback")

var ErrCallbackNotFount = errors.New("callback not found")

var ErrTrashInRequestPool = errors.New("trash in request pool")

var ErrTrashInResponsePool = errors.New("trash in response pool")

var ErrTrashInItemPool = errors.New("trash in item pool")

var ErrTrashInErrorPool = errors.New("trash in error pool")

var ErrInvalidDataItem = errors.New("invalid data items")

var ErrResponseFromAnalyzer = errors.New("response from analyzer")

var ErrGenerateInvalidType = errors.New("generate invliad type")

var ErrNilParser = errors.New("nil parser")
