package gospider

import "errors"

// Constant that used as item key for special use
var (
	KeyDefault = "_default"
	KeyData    = "_data"
	KeyBody    = "_body"
)

/**************************************************************
* errors: Scheduler
**************************************************************/
var ErrDupeRequest = errors.New("duplicate request")

/**************************************************************
* errors: Downloader
**************************************************************/
var ErrNilRequest = errors.New("nil request")

/**************************************************************
* errors: Analyzer
**************************************************************/
var ErrParse = errors.New("parse error")

/**************************************************************
* errors: Pipeline
**************************************************************/
var ErrDropItem = errors.New("drop item")
var ErrNilProcessor = errors.New("nil processor")

/**************************************************************
* errors: Nil Entity
**************************************************************/
var ErrNilResponse = errors.New("nil response")
var ErrNilParser = errors.New("nil parser")
var ErrNilItem = errors.New("nil item")

var ErrValueIsNotString = errors.New("value is not string")

var ErrDownloadFail = errors.New("download fail")

var ErrStopIteration = errors.New("stop iteration")

var ErrNilProcessorList = errors.New("nil processor list for pipeline")

var ErrContinue = errors.New("continue")

var ErrInvalidCallback = errors.New("invalid callback")

var ErrCallbackNotFount = errors.New("callback not found")

var ErrTrashInRequestPool = errors.New("trash in request pool")

var ErrTrashInResponsePool = errors.New("trash in response pool")

var ErrTrashInItemPool = errors.New("trash in item pool")

var ErrTrashInErrorPool = errors.New("trash in error pool")

var ErrInvalidDataItem = errors.New("invalid data items")

var ErrResponseFromAnalyzer = errors.New("response from analyzer")

var ErrGenerateInvalidType = errors.New("generate invliad type")
