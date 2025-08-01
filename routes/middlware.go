package routes

import (
	"net/http"
	// ""
)

func (h *Handler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("sessionID")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if ok := h.InMemorySession.Authed(cookie.Value); !ok {
			w.WriteHeader(http.StatusUnauthorized)
			// m := models.PostListModel{}
			// PrintPosts()
			// m.Posts = posts
			h.Tmpl.ExecuteTemplate(w, "login.html", "")
			return
		}

		next.ServeHTTP(w, r)
	})

}
