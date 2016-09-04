

## CLI

    aws dynamodb list-tables --endpoint-url http://localhost:8000

    aws dynamodb query --endpoint-url http://localhost:8000 --table-name patients --key-conditions file://patients-keys.json
