package server

import (
	"net/http"

	"github.com/rutmir/services/web-rest/server/controllers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"index", "GET", "/", func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "/doc", 301) }},
	Route{"alive", "GET", "/alive", controllers.Alive},
	Route{"auth", "POST", "/oauth/token", controllers.Authentication},
	Route{"newAccount", "POST", "/account", controllers.CreateAccount},
	//Route{"GetCurrentData", "GET", "/api/{format}/{symbol}/current", controllers.Current, },
	//Route{"GetAllData", "GET", "/api/{format}/{symbol}/all", controllers.All, },
	//Route{"GetAdvisor", "GET", "/api/{format}/{symbol}/advisor", controllers.Advisor, },
	//Route{"RefreshData", "GET", "/api/refresh", controllers.Refresh, },
	//Route{"CleanData", "GET", "/api/clean", controllers.ClearDB, },
}
