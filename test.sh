#!/bin/bash

# Variables
SERVER_URL="http://localhost:8080"
ROUTE="/posts"
DATA="{\"key\":\"value\"}"  # Example data to send in the POST request

# Send POST request and capture response
RESPONSE=$(curl -s -X POST -H "Content-Type: application/json" -d "$DATA" "$SERVER_URL$ROUTE")

# Check HTTP status code
HTTP_STATUS=$(echo "$RESPONSE" | head -n 1 | cut -d$' ' -f2)

# Check response body or perform additional checks as needed
# For example, you can check if a specific key exists in the response JSON
# You might need to install jq for JSON parsing: sudo apt-get install jq
# For example, check if "success" key exists and its value is true
# SUCCESS=$(echo "$RESPONSE" | jq -r '.success')

# Check HTTP status code
if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "POST request successful: HTTP status $HTTP_STATUS"
    # Additional checks or actions if needed
else
    echo "POST request failed: HTTP status $HTTP_STATUS"
    # Additional error handling if needed
fi

