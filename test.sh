# SEND
curl -X POST \
  http://localhost:8080/submit_data \
  -H 'Content-Type: application/json' \
  -d '{
        "question": "What is the meaning of life, the universe, and everything?",
}'

# RECEIVE
curl -X GET \
  'http://localhost:8080/get_processed_data?id=6aabdd6764ff6b2db703419bbec26281'
