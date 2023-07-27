# edgefarm-accel-demo-frontend
Consumer Web UI for basic example that shows 

# Configuration

Set the following environment variables:

* `NATS_SERVER`, e.g. `nats://<yourIP>:4222`
* `NATS_EXPORT_SUBJECT` to `*.sensor`
* `NATS_STREAM_NAME` to `aggregate-stream` or whatever stream name is used
* `NATS_CREDS` to the path where the creds file os located e.g. `/creds/user.creds`

# Run locally

## Prerequisites

```bash	
pip install -r requirements.txt
```

## Run

```bash
bokeh serve --show serve.py
```

Automatically opens a browser window with the web UI.


## Build docker image

```bash
docker build -t example-basic-consumer:latest .
```

## Run docker image

```bash
docker run -it -p 5006:5006 -e NATS_SERVER=nats://<server> -e NATS_EXPORT_SUBJECT="*.sensor" -e NATS_STREAM_NAME=aggregate-stream example-basic-consumer:latest
```
Then open browser at http://localhost:5006/serve

