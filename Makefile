#goctl api go -api app/shortener/api/shortener.api -dir app/shortener/api -style gozero
#goctl api go -api app/user/api/user.api -dir app/user/api -style gozero
#goctl rpc protoc protos/shortener/v1/shortener.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/shortener/rpc -c gen/go/shortener --style gozero
#goctl rpc protoc protos/shortener/v1/shortener.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/shortener/rpc --style=gozero  -c gen/go/shortener
#goctl rpc protoc protos/shortener/v1/shortener.proto --go_out=. --go-grpc_out=. --zrpc_out=./app/shortener/rpc --style=gozero -c ./gen/go/shortener
#

#goctl model mysql ddl -src="app/user/schema/sql/000001_users.up.sql" -dir="app/user/data/model"
#goctl model mysql ddl -src="app/shortener/schema/sql/000001_links.up.sql" -dir="app/shortener/data/model" -c
#goctl model mysql ddl -src="app/shortener/schema/sql/000002_link_access_logs.up.sql" -dir="app/shortener/data/model"
#goctl model mysql ddl -src="app/shortener/schema/sql/000003_analytics_summary_daily.up.sql" -dir="app/shortener/data/model"

#
#migrate  create -ext sql -dir . -seq analytics_summary_daily
#migrate  create -ext sql -dir schema/sql/user -seq agg_daily_summary_table
#migrate  create -ext sql -dir schema/sql/user -seq url_analytics_table
#migrate  create -ext sql -dir schema/sql/shortener -seq create_users_table
#user
#migrate  create -ext sql -dir schema/sql/user -seq create_short_urls_table
databaseURL="mysql://root:root@tcp(localhost:3306)/minify?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
# migrate -path="./app/user/schema/sql/000001_users.up.sql" -database="mysql://root:root@tcp(localhost:3306)/lucid?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai" -verbose up
migrate_up:
	migrate -path="./app/user/schema/sql/000001_users.up.sql" -database=${databaseURL} -verbose up
migrate_drop:
	migrate -path="./schema/sql/user" -database=${databaseURL} -verbose drop -f



#goctl rpc protoc protos/shortener/v1/shortener.proto --go_out=. --go-grpc_out=. --zrpc_out=./gen/go --style gozero

# goctl template init --home template