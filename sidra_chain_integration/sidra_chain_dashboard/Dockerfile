# Dockerfile
FROM python:3.13.3-slim

WORKDIR /app

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY . .

CMD ["python", "sidra_chain_dashboard.py"]
