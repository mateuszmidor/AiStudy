start:
	docker run --rm --name=qdrant-container -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
	firefox http://127.0.0.1:6333/dashboard
	#  -v ${PWD}/qdrant:/qdrant/storage \

run:
	go run .

stop:
	docker stop qdrant-container

