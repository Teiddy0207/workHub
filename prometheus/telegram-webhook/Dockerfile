FROM python:3.11-slim

WORKDIR /app

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY telegram-webhook.py .

EXPOSE 5000

CMD ["python", "telegram-webhook.py"]

