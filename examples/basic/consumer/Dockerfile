FROM python:3.10-slim

WORKDIR /app

RUN apt-get update && apt-get install -y gcc

COPY requirements.txt .
RUN pip install --trusted-host pypi.python.org -r requirements.txt

COPY . .

EXPOSE 5006

CMD ["bokeh", "serve", "--allow-websocket-origin=*", "serve.py"]
