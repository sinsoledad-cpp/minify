goctl api go -api app/user/api/user.api -dir app/user/api -style gozero 
goctl api go -api app/shortener/api/shortener.api -dir app/shortener/api -style gozero 
goctl rpc protoc protos/user/v1/user.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/user/rpc -c gen/go/user --style gozero 
goctl rpc protoc protos/user/v1/user.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/user/rpc --style=gozero  -c gen/go/user
goctl rpc protoc protos/user/v1/user.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/user/rpc --style=gozero -c ./gen/go/user

goctl model mysql ddl -src="schema/sql/user/001_create_users_table.sql" -dir="data/model/user"
goctl model mysql ddl -src="schema/sql/shortener/000001_create_short_urls_table.up.sql" -dir="data/model/shortener"
goctl model mysql ddl -src="schema/sql/shortener/000002_url_analytics_table.up.sql" -dir="data/model/shortener"
goctl model mysql ddl -src="schema/sql/shortener/000003_agg_daily_summary_table.up.sql" -dir="data/model/shortener"
goctl model mysql ddl -src="schema/sql/shortener/000001_create_url_analytics_table.up.sql" -dir="data/model/shortener"


migrate  create -ext sql -dir schema/sql/shortener -seq agg_daily_summary_table
migrate  create -ext sql -dir schema/sql/shortener -seq url_analytics_table
migrate  create -ext sql -dir schema/sql/user -seq create_users_table
shortener
migrate  create -ext sql -dir schema/sql/shortener -seq create_short_urls_table
databaseURL="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

migrate_up:
	migrate -path="./schema/sql/user" -database=${databaseURL} -verbose up
	migrate -path="./schema/sql/user" -database="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai" -verbose up
	migrate -path="./schema/sql/shortener" -database="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai" -verbose up
migrate_drop:
	migrate -path="./schema/sql/user" -database=${databaseURL} -verbose drop -f



goctl rpc protoc protos/user/v1/user.proto --go_out=. --go-grpc_out=. --zrpc_out=./gen/go --style gozero 