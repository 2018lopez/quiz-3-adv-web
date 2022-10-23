BODY='{"title":"Birds", "description":"Toucans, Dove, Peacock", "completed":false}'

Post

curl -i -d "$BODY" localhost:4000/v1/todos


Update 
curl -X PATCH -d '{"Completed" : true}' "localhost:4000/v1/todos/5"

Delete
curl -X DELETE  localhost:4000/v1/schools/1

get by id 

curl localhost:4000/v1/todos/4
curl localhost:4000/v1/todos/1


get by all by sorting, pagination

curl "localhost:4000/v1/todos?title=errands"
curl "localhost:4000/v1/todos?sort=title"
curl "localhost:4000/v1/tods?page=1&page_size=2"
