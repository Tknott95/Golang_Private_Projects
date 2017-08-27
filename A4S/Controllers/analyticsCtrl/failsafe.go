package analytics

import "net/http"

/*  Recovers server if any coroutine panicks. This will allow for it to not crash it where it usually does */
func Recover(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				switch err {
				case ErrInvalidEmail:
					http.Error(w, ErrInvalidEmail.Error(), http.StatusUnauthorized)
				case ErrInvalidID:
					http.Error(w, ErrInvalidID.Error(), http.StatusUnauthorized)
				default:
					http.Error(w, "Unknown error, recovered from panic", http.StatusInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, req)
	})
}
