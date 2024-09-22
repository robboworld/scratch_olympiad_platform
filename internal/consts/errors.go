package consts

// http code 400
const (
	ErrCountryNotFoundInDB      = "country not found"
	ErrNominationNotFoundInDB   = "nomination not found"
	ErrEmailAlreadyInUse        = "email already in use"
	ErrAtoi                     = "string to int error"
	ErrTimeParse                = "string to time error"
	ErrIncorrectPasswordOrEmail = "incorrect password or email"
	ErrNotFoundInDB             = "not found"
	ErrShortPassword            = "please input password, at least 8 symbols"
	ErrPasswordResetLinkInvalid = "password reset link invalid"
	ErrPasswordResetLinkExpired = "password reset link expired"
)

// http code 401
const (
	ErrTokenExpired     = "token expired"
	ErrNotStandardToken = "token claims are not of type *StandardClaims"
)

// http code 403
const (
	ErrUserIsNotActive         = "user is not active. please check your email"
	ErrProjectPageIsBanned     = "the projectPage is banned. no access"
	ErrAccessDenied            = "access denied"
	ErrDoesNotMatchAgeCategory = "does not match the age category"
)

// ErrActivationLinkUnavailable have http code 503
const (
	ErrActivationLinkUnavailable = "activation link is currently unavailable"
)
