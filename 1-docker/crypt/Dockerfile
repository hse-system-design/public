FROM python:3.7

RUN pip install cryptography

COPY crypt.py /app/crypt.py

ENTRYPOINT ["python3", "/app/crypt.py"]
