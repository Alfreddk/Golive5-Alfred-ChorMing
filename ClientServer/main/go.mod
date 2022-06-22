module ProjectGoLive/ClientServer/main

go 1.18

require (
	github.com/gorilla/mux v1.8.0
	github.com/joho/godotenv v1.4.0
	github.com/satori/go.uuid v1.2.0
	goSch/golive5/arrds v0.0.0-00010101000000-000000000000
	goSch/golive5/database v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/kr/pretty v0.2.1 // indirect
)

replace goSch/golive5/arrds => ../../arrds

replace goSch/golive5/database => ../../database
