docker run --name ad-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=ad_db \
  -p 5433:5432 \
  -d postgres:15

# создаём и запускаем контейнер PosgreSQL и сразу создаем БД  
# В теории один первый раз