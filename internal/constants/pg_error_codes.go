package constants

// PostgreSQL common error codes
const (
	// Class 23 — Integrity Constraint Violation
	IntegrityConstraintViolation = "23000"
	RestrictViolation            = "23001"
	NotNullViolation             = "23502"
	ForeignKeyViolation          = "23503"
	UniqueViolation              = "23505"
	CheckViolation               = "23514"
	ExclusionViolation           = "23P01"

	// Class 22 — Data Exception
	StringDataRightTruncationDataException = "22001"
	NumericValueOutOfRange                 = "22003"
	InvalidDatetimeFormat                  = "22007"

	// Class 40 — Transaction Rollback
	TransactionRollback                     = "40000"
	TransactionIntegrityConstraintViolation = "40002"
	SerializationFailure                    = "40001"
	StatementCompletionUnknown              = "40003"
	DeadlockDetected                        = "40P01"
)

// case "22001": // String Data Right Truncation
