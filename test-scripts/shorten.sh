#!/bin/bash

API_URL="http://localhost:8080/shorten"
TEST_URL="https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html"

# First request
response1=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{\"long_url\":\"$TEST_URL\"}")

slug1=$(echo "$response1" | jq -r '.short_url')

# Second request with same URL
response2=$(curl -s -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d "{\"long_url\":\"$TEST_URL\"}")

slug2=$(echo "$response2" | jq -r '.short_url')

# Verification
if [ "$slug1" == "$slug2" ]; then
  echo "✅ Test passed: Slugs are identical ($slug1)"
else
  echo "❌ Test failed:"
  echo "First slug: $slug1"
  echo "Second slug: $slug2"
fi
