.PHONE: run build-image run-container

run:
	go run cmd/server/*.go

run-client:
	go run cmd/client/*.go


build-image:
	docker build -t word-of-wisdom .

run-container: 
	docker run -p 8000:8000 word-of-wisdom
