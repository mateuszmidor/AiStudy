install:
	cd embedding-localhost && virtualenv venv && source venv/bin/activate && pip install -U sentence-transformers

run:
	docker run --rm --name=qdrant-db -d -p 6333:6333 -p 6334:6334 qdrant/qdrant
	source ./embedding-localhost/venv/bin/activate && go run . || true
	docker stop qdrant-db

dashboard:
	firefox http://127.0.0.1:6333/dashboard
