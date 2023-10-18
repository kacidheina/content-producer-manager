Content producer and consumer manager service

# Description of the service:

This is a service which contains two REST APIs for handling content and metadata about the content.
The sender API will listen to port :8002 and the consumer API will listen to port :8001.
The sender API will process the request and publish a rabbitmq message to the rabbitmq 'content' exchange.
The consumer API will set up a queue, bind it to the exchange and consume the messages from the queue.To handle the
messages, a go routine is set up in order to process the messages and then store them in the postgres database.
In order to query for the content, you should fire a request to the consumer API and it will return a list with contents
for the specific sender.

# How to run the service

run `docker-compose up --build` and the service will start and listen on port :8001 and the consumer API will listen on port :8002

# API Endpoints

## Sender API

    Endpoint: `/api/send`
    Method: POST
    Description: Publish endpoint in the rabbitmq.
    Body Request: 
      ```json {   
        "sender_id": 1,
        "file_type": "application/msword",
        "receiver_id": 3,
        "is_payable": true
        } ```
    Response:
      Status Code: 200 (OK)
    example: http://localhost:8002/api/send

## Consumer API

    Endpoint: `/api/consumer/{sender_id}`
    Method: GET
    Description: Query a list of content for a specific sender_id.
    Parameters:
        sender_id

    Response:
    Status Code: 200 (OK)
    Body:
        [
            {
                "id": 1,
                "sender_id": 1,
                "receiver_id": 4,
                "file": "SGVsbG8gV29ybGQh",
                "file_type": "application/json",
                "is_payable": false,
                "is_paid": false,
                "created_at": "2023-10-17T23:23:28.754341Z"
              },
              {
                "id": 2,
                "sender_id": 1,
                "receiver_id": 4,
                "file": "SGVsbG8gV29ybGQh",
                "file_type": "application/json",
                "is_payable": true,
                "is_paid": true,
                "created_at": "2023-10-17T23:25:24.521795Z"
     }] 

Things I wished I would have added to the service:

* Authentication and a middleware to check if the user is authenticated
* Not using postgresDB for storage of the files but other options such as S3 bucket
* Metrics to monitor the service
* More tests