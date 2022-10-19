BODY='{"title":"Bills to Pay", "descriptions":"Internet, water, light, Cable, School, loans", "completed":false}'

curl -X PATCH  -d '{"Completed" : True}' localhost:4000/v1/schools/3