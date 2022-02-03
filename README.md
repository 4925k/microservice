#microservice

Folder structure required.

service
    certs
        <server certificate>
        <server key>
    config
        config.yaml
    logging 
        microservice.log
    bin
        <executable>

NOTE: 
-> logging folder should exist
-> the config.yaml file should have the required fields.
    --- 
    address: ":8080"
    certificate: ./certs/server.crt
    key: ./certs/server.key
-> the path to certificate and key is relative to the executable
