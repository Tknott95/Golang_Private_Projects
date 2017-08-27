package routerCtrl

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/julienschmidt/httprouter"
	userCtrl "github.com/tknott95/LCW-FC/Controllers/UsersCtrl"
	httpGlobals "github.com/tknott95/LCW-FC/Globals/http"
)

func InitServer() {
	mux := httprouter.New()

	/* CORS */
	handler := cors.Default().Handler(mux)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Insert the middleware
	handler = c.Handler(handler)

	// PC LANGS
	mux.GET("/", index)

	mux.POST("/api/login", userCtrl.Login)

	/* UMBRELLA API PORTION */

	//http.Handle("/Public/", http.StripPrefix("/Public", http.FileServer(http.Dir("./Public"))))

	mux.NotFound = http.StripPrefix("/Public", http.FileServer(http.Dir("./Public")))

	// handler for serving files
	// mux.ServeFiles("/Public/*filepath", http.Dir("/var/www/Public/"))
	http.ListenAndServe(httpGlobals.PortNumber, handler)

}

func index(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	println("📝 Currently on Index page.")

}
