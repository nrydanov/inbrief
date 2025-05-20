proto:
    protoc --go_out=. \
           --twirp_out=. \
           proto/*.proto

gen-docs:
    twirp-openapi-gen \
        -in proto/fetch.proto \
        -out docs/openapi.json
