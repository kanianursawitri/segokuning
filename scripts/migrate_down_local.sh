export SEGOKUNING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" && migrate -database ${SEGOKUNING} -path ./db/migrations down
