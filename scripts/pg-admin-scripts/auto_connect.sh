# Подождём немного, пока поднимется Postgres (можно заменить на более сложную проверку)
echo "Starting connection script..."
# sleep 10

# # Авторизируемся в pgAdmin API и получаем токен
# TOKEN=$(curl -X GET \
#   http://localhost:$PGADMIN_LISTEN_PORT/login \
#   -H 'Content-Type: application/json' \
#   -d '{"email": "$PGADMIN_DEFAULT_EMAIL", "password": "$PGADMIN_DEFAULT_PASSWORD"}' | jq -r '.token')

# # Создаём сервер PostgreSQL
# echo "Token received: $TOKEN"
# SERVER_ID=$(curl -X POST \
#   http://localhost:$PGADMIN_LISTEN_PORT/api/server \
#   -H "Authorization: Bearer $TOKEN" \
#   -H 'Content-Type: application/json' \
#   -d @server.json | jq -r '.id')

# echo "Server created with ID: $SERVER_ID"