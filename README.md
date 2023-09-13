# CodeFlix Encoder

Microservice responsible for encoding videos on mp4 format and convert them to MPEG-DASH. The service will encode and upload the video in parallel/concurrently.

The service will be developed using hexagonal architecture/ports and adapters.

## Flow

1. Receive a message from RabbitMQ with the video to be encoded
2. Download the video from the storage service (Google Cloud Storage)
3. Break the video into smaller parts (fragments)
4. Convert to MPEG-DASH
5. Upload the video to the storage service (Google Cloud Storage)
6. Send a notification to the queue that the video has been encoded or failed.
7. In case of failure, the message will be rejected and sent to the dead letter exchange.

The workflow for message:

![image](https://github.com/Twsouza/codeflix-encoder/assets/8239709/9e080133-b830-41ca-ac56-967a2b02e85a)

In case of error:

![image](https://github.com/Twsouza/codeflix-encoder/assets/8239709/9b491967-7cf1-46c1-b9dd-47c7ee210d92)

## Input message format

It must be sent as JSON with the following format:

```json
{
  "resource_id": "my-resource-id-can-be-a-uuid-type",
  "file_path": "video.mp4"
}
```

Where:

- resource_id: is the ID of the resource that will be encoded. It is a string type and refers to the ID on the source.
- file_path: is the path of the video file in the storage service, string type.

## Output message format (notification message)

The service will send a message once the encoding is finished or failed. The message will be sent to the queue on `.env`.

In case of success, the following message will be sent:

```json
{
  "id":"uuid",
  "output_bucket_path":"destination_bucket",
  "status":"COMPLETED",
  "video":{
    "encoded_video_folder":"folder_where_it_was_saved",
    "resource_id":"uuid",
    "file_path":"video.mp4"
  },
  "Error":"",
  "created_at":"2023-08-29T00:00:00.000000-03:00",
  "updated_at":"2023-08-29T00:00:00.000000-03:00"
}
```

For errored messages, the following message will be sent:

```json
{
  "message": {
    "resource_id": "uuid",
    "file_path": "video.mp4"
  },
  "error":"error message"
}
```

Additionally, the service will send the original message to the dead letter exchange, set the `RABBITMQ_DLX` on `.env`.
