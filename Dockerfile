FROM python:3.10.7-slim

WORKDIR /app

COPY main.py ./main.py
COPY requirements.txt ./requirements.txt

RUN pip install -r requirements.txt

ENTRYPOINT [ "python", "-u", "main.py" ]