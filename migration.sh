

source .env
cd sql/schema
goose postgres $POSTGRES_URL $1
