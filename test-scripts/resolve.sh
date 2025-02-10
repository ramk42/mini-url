#!/bin/bash

API_URL="http://localhost:8080"

if [ $# -eq 0 ] || [ -z "$1" ]; then
    echo "Error: Missing argument slug"
    exit 1
fi

TEST_SLUG=$1

# call API
response=$(curl -s -o /dev/null -w "%{http_code}\n%{redirect_url}" "${API_URL}/${TEST_SLUG}")

# Extraction des infos
status_code=$(echo "$response" | head -n1)
redirect_url=$(echo "$response" | tail -n1)

# Comparaison
if [[ "$status_code" -eq 301 || "$status_code" -eq 302 ]]; then
    if [ -n "$redirect_url" ]; then
        echo "✅ Test réussi :"
        echo "Status: $status_code"
        echo "Redirection vers: $redirect_url"
    else
        echo "❌ En-tête Location manquant"
    fi
else
    echo "❌ Test échoué:"
    echo "Code: $status_code (attendu 301/302)"
fi
