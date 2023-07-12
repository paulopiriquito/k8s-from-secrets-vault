FROM python:3.10.7-slim

WORKDIR /app

COPY ./src .

RUN pip install -r requirements.txt

ENTRYPOINT [ "python", "-u", "main.py" ]