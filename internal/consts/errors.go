package consts

// http code 400
const (
	ErrEventNotFoundInDB               = "event not found"
	ErrCountryNotFoundInDB             = "country not found"
	ErrRegionNotFoundInDB              = "region not found"
	ErrRegionNotInCountry              = "region not in country"
	ErrCountryHasNoRegions             = "country has no regions"
	ErrNominationNotFoundInDB          = "nomination not found"
	ErrEmailAlreadyInUse               = "email already in use"
	ErrAtoi                            = "string to int error"
	ErrTimeParse                       = "string to time error"
	ErrIncorrectPasswordOrEmail        = "incorrect password or email"
	ErrNotFoundInDB                    = "not found"
	ErrShortPassword                   = "please input password, at least 8 symbols"
	ErrPasswordResetTokenInvalid       = "password reset token invalid"
	ErrPasswordResetTokenExpired       = "password reset token expired"
	ErrUserWithEmailNotFound           = "user with this email not found"
	ErrApplicationAlreadySubmitted     = "application already submitted"
	ErrEventRoleAlreadyAssigned        = "event role already assigned for this user"
	ErrCountryAlreadyAssigned          = "country is already assigned to the event"
	ErrRegionAlreadyAssigned           = "region is already assigned to the event"
	ErrRegionCountryNotLinkedToEvent   = "region cannot be assigned because its country is not linked to the event"
	ErrRegionRequired                  = "region is required for this country"
	ErrTranslationForEventAlreadyExist = "translation for the event already exist"
)

// http code 401
const (
	ErrTokenExpired     = "token expired"
	ErrNotStandardToken = "token claims are not of type *StandardClaims"
)

// http code 403
const (
	ErrUserIsNotActive          = "user is not active. please check your email"
	ErrAccessDenied             = "access denied"
	ErrDoesNotMatchAgeCategory  = "does not match the age category"
	ErrSelfAssignmentNotAllowed = "self assignment not allowed"
	ErrEventNotAccessible       = "event is not accessible to the user"
)

// ErrActivationLinkUnavailable have http code 503
const (
	ErrActivationLinkUnavailable = "activation link is currently unavailable"
)
