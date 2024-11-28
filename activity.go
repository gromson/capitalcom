package capitalcom

type ActivitySource string

const (
	ActivityCloseOut ActivitySource = "CLOSE_OUT"
	ActivityDealer   ActivitySource = "DEALER"
	ActivitySl       ActivitySource = "SL"
	ActivitySystem   ActivitySource = "SYSTEM"
	ActivityTp       ActivitySource = "TP"
	ActivityUser     ActivitySource = "USER"
)

type ActivityStatus string

const (
	ActivityStatusAccepted     ActivityStatus = "ACCEPTED"
	ActivityStatusCreated      ActivityStatus = "CREATED"
	ActivityStatusExecuted     ActivityStatus = "EXECUTED"
	ActivityStatusExpired      ActivityStatus = "EXPIRED"
	ActivityStatusRejected     ActivityStatus = "REJECTED"
	ActivityStatusModified     ActivityStatus = "MODIFIED"
	ActivityStatusModifyReject ActivityStatus = "MODIFY_REJECT"
	ActivityStatusCancelled    ActivityStatus = "CANCELLED"
	ActivityStatusCancelReject ActivityStatus = "CANCEL_REJECT"
	ActivityStatusUnknown      ActivityStatus = "UNKNOWN"
)

type ActivityType string

const (
	ActivityTypePosition         ActivityType = "POSITION"
	ActivityTypeWorkingOrder     ActivityType = "WORKING_ORDER"
	ActivityTypeEditStopAndLimit ActivityType = "EDIT_STOP_AND_LIMIT"
	ActivityTypeSwap             ActivityType = "SWAP"
	ActivityTypeSystem           ActivityType = "SYSTEM"
)
