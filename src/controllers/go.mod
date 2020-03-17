module gitlab.com/digiresilience/link/quepasa/controllers

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dlclark/regexp2 v1.1.6 // indirect
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/jwtauth v4.0.3+incompatible
	github.com/gorilla/websocket v1.4.1
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.2.0
	github.com/nbutton23/zxcvbn-go v0.0.0-20180912185939-ae427f1e4c1d
	github.com/prometheus/client_golang v1.1.0
	github.com/trustelem/zxcvbn v1.0.1
	gitlab.com/digiresilience/link/quepasa/models v0.0.0
)

replace gitlab.com/digiresilience/link/quepasa/controllers => ./

replace gitlab.com/digiresilience/link/quepasa/models => ../models

go 1.13
