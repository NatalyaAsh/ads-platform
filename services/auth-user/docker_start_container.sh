docker run --name test-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=auth_user_db \
  -p 5432:5432 \
  -d postgres:15