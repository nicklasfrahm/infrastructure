package namecheap

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Reference: https://www.namecheap.com/support/api/global-parameters/
var (
	ErrApiUserMissing            = NewGlobalError(1010101, "parameter missing: APIUser")
	ErrUnsupportedAuth           = NewGlobalError(1030408, "unsupported authentication type")
	ErrCommandMissing            = NewGlobalError(1010104, "parameter missing: Command")
	ErrApiKeyMissing0            = NewGlobalError(1010102, "parameter missing: APIKey")
	ErrApiKeyMissing1            = NewGlobalError(1011102, "parameter missing: APIKey")
	ErrClientIpMissing0          = NewGlobalError(1010105, "parameter missing: ClientIP")
	ErrClientIpMissing1          = NewGlobalError(1011105, "parameter missing: ClientIP")
	ErrApiUserUnknownValidation  = NewGlobalError(1050900, "unknown validation error: APIUser")
	ErrRequestIpInvalid          = NewGlobalError(1011150, "parameter invalid: RequestIP")
	ErrRequestIpDisabledOrLocked = NewGlobalError(1017150, "parameter disabled or locked: RequestIP")
	ErrClientIpDisabledOrLocked  = NewGlobalError(1017105, "parameter disabled or locked: ClientIP")
	ErrApiUserDisabledOrLocked   = NewGlobalError(1017101, "parameter disabled or locked: APIUser")
	ErrTooManyDeclinedPayments   = NewGlobalError(1017410, "too many declined payments")
	ErrTooManyLoginAttempts      = NewGlobalError(1017411, "too many login attempts")
	ErrUserNameNotAvailable      = NewGlobalError(1019103, "parameter not available: UserName")
	ErrUserNameNotAuthorized     = NewGlobalError(1016103, "parameter no authorized: UserName")
	ErrUserNameDisabledOrLocked  = NewGlobalError(1017103, "parameter disabled or locked: UserName")
)

type GlobalError struct {
	Code        int
	Description string
}

func (e *GlobalError) Error() string {
	return fmt.Sprintf("namecheap: %d - %s", e.Code, e.Description)
}

func NewGlobalError(code int, description string) *GlobalError {
	return &GlobalError{
		Code:        code,
		Description: description,
	}
}

// For more information on how to write a provider, consider the following code:
// - https://github.com/pulumi/pulumi-gcp/blob/master/sdk/go/gcp/provider.go
type Provider struct {
	pulumi.ProviderResourceState

	ApiUser  pulumi.StringPtrOutput `pulumi:"apiUser"`
	ApiKey   pulumi.StringPtrOutput `pulumi:"apiKey"`
	UserName pulumi.StringPtrOutput `pulumi:"userName"`
	ClientIp pulumi.StringPtrOutput `pulumi:"clientIp"`
}

// TODO: Import provider args.

func NewCustomDNS() {
}
