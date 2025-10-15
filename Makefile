goctl api go -api app/user/api/user.api -dir app/user/api -style gozero 

goctl model mysql ddl -src="schema/sql/user/001_create_users_table.sql" -dir="data/model/user"

migrate  create -ext sql -dir schema/sql/user -seq create_users_table


databaseURL="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

migrate_up:
	migrate -path="./schema/sql/user" -database=${databaseURL} -verbose up
	migrate -path="./schema/sql/user" -database="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai" -verbose up
migrate_drop:
	migrate -path="./schema/sql/user" -database=${databaseURL} -verbose drop -f