FROM python:3.7

RUN pip install Flask

COPY server.py /app/server.py

WORKDIR /app
ENTRYPOINT ["python3", "server.py"]
