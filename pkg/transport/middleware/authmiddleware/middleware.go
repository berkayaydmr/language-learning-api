package authmiddleware

import (
	"net/http"

	customerr "github.com/berkayaydmr/language-learning-api/pkg/error"
	"github.com/berkayaydmr/language-learning-api/pkg/transport/middleware"
	"github.com/berkayaydmr/language-learning-api/pkg/utils"
)

type authMiddleware struct {
	apiKey string
	next   http.Handler
}

func NewAuthMiddleware(apiKey string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return &authMiddleware{
			apiKey: apiKey,
			next:   next,
		}
	}
}

func (m *authMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Api-Key") != m.apiKey {
		utils.RespondWithError(w, customerr.ErrAuthorizationFailed)
		return
	}

	m.next.ServeHTTP(w, r)
}
