# test vegeta click increment
# be careful with the rate and workers, it can crash the website redirected
echo "GET http://localhost:8080/[REPLACE-YOUR-SLUG]" | vegeta attack \
  -rate=100 \
  -workers=1000 \
  -duration=10s \
  -redirects=1 \
  -output=results.bin && \
  vegeta report -type=json results.bin > metrics.json && \
  vegeta plot results.bin > performance.html