while true ; do
curl -sL -w "%{http_code} \n" -i "https://api-pointsbet-prod.azure-api.net/api/v2/events/nextup?limit=10" -o /dev/null
sleep 5
done