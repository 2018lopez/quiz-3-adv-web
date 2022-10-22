BODY='{"title":"School", "description":"Do Senior seminar essay", "completed":false}'

Post

Update 
curl -X PATCH  -d '{"Completed" : True}' localhost:4000/v1/schools/3

Delete

get by id 


get by all by sorting, pagination

curl "localhost:4000/v1/todos?title=errands"
curl "localhost:4000/v1/todos?completed=false"