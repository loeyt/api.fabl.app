SHELL = /bin/bash

PROTOS = $(patsubst proto/%.proto, %, $(shell find proto -type f -name '*.proto'))
SERVICES = $(filter %_service, $(PROTOS))
GATEWAY_SERVICES = $(patsubst proto/%.yaml, %, $(shell find proto -type f -name '*_service.yaml'))
OTHER_PROTOS = $(filter-out $(SERVICES), $(PROTOS))

vpath %.proto proto
vpath %_service.yaml proto
vpath %.pb.go internal/pb
vpath %.pb.gw.go internal/pb

%.pb.go : %.proto
	protoc -I proto \
		--go_out internal/pb --go_opt paths=source_relative \
		$<

%_service.pb.go %_service_grpc.pb.go : %_service.proto
	protoc -I proto \
		--go_out internal/pb --go_opt paths=source_relative \
		--go-grpc_out internal/pb --go-grpc_opt paths=source_relative \
		$<

%_service.pb.gw.go : %_service.yaml %_service.proto
	protoc -I proto \
		--grpc-gateway_out internal/pb \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt grpc_api_configuration=$< \
		--grpc-gateway_opt paths=source_relative \
		$(filter %.proto, $^)

proto: $(patsubst %, %.pb.go, $(OTHER_PROTOS) $(SERVICES)) $(patsubst %, %.pb.gw.go, $(GATEWAY_SERVICES))
.PHONY: proto
