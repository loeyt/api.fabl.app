SHELL = /bin/bash

V1_PROTOS = $(shell cd proto && find v1 -type f -name '*.proto')
V1_SERVICES = $(filter %_service.proto, $(V1_PROTOS))
V1_OTHER_PROTOS = $(filter-out $(V1_SERVICES), $(V1_PROTOS))

vpath %.proto proto

all: $(patsubst %.proto, %.pb.go, $(V1_OTHER_PROTOS) $(V1_SERVICES)) $(patsubst %.proto, %.pb.gw.go, $(V1_SERVICES)) v1.swagger.json
.PHONY: all

%.pb.go : %.proto
	protoc -I proto \
		--go_out . --go_opt paths=source_relative \
		$<

%_service.pb.go %_service_grpc.pb.go %_service.pb.gw.go : %_service.proto
	protoc -I proto \
		--go_out . --go_opt paths=source_relative \
		--go-grpc_out . --go-grpc_opt paths=source_relative \
		--grpc-gateway_out . \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		$<

v1.swagger.json: $(V1_PROTOS)
	protoc -I proto \
		--openapiv2_out . \
		--openapiv2_opt logtostderr=true \
		--openapiv2_opt allow_merge=true \
		--openapiv2_opt json_names_for_fields=false \
		--openapiv2_opt merge_file_name=v1 \
		$(filter %/api.proto, $^) $(filter-out %/api.proto, $^)
