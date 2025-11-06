# Hướng dẫn setup Telegram Alerts

## Bước 1: Tạo Telegram Bot

1. Mở Telegram và tìm **@BotFather**
2. Gửi lệnh `/newbot`
3. Đặt tên cho bot (ví dụ: `WorkHub Alert Bot`)
4. Đặt username cho bot (phải kết thúc bằng `bot`, ví dụ: `workhub_alert_bot`)
5. BotFather sẽ trả về **Bot Token** - lưu lại token này

Ví dụ token: `123456789:ABCdefGHIjklMNOpqrsTUVwxyz`

## Bước 2: Lấy Chat ID

Có 2 cách:

### Cách 1: Chat với bot
1. Tìm bot vừa tạo và bắt đầu chat
2. Gửi bất kỳ message nào cho bot
3. Truy cập: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
4. Tìm `"chat":{"id":123456789}` - đây là Chat ID của bạn

### Cách 2: Dùng @userinfobot
1. Tìm **@userinfobot** trên Telegram
2. Bắt đầu chat và gửi `/start`
3. Bot sẽ trả về Chat ID của bạn

## Bước 3: Cấu hình Environment Variables

Tạo file `.env` trong thư mục root của project:

```bash
TELEGRAM_BOT_TOKEN=123456789:ABCdefGHIjklMNOpqrsTUVwxyz
TELEGRAM_CHAT_ID=123456789
```

Hoặc export trực tiếp trong terminal:

```bash
export TELEGRAM_BOT_TOKEN="123456789:ABCdefGHIjklMNOpqrsTUVwxyz"
export TELEGRAM_CHAT_ID="123456789"
```

## Bước 4: Khởi động services

```bash
docker compose up -d --build
```

## Bước 5: Test

1. Kiểm tra webhook service:
```bash
curl http://localhost:5000/health
```

2. Test gửi alert thủ công:
```bash
curl -X POST http://localhost:5000/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "status": "firing",
    "groupLabels": {"alertname": "TestAlert"},
    "commonLabels": {"severity": "warning", "service": "test"},
    "alerts": [{
      "annotations": {
        "description": "This is a test alert",
        "summary": "Test alert summary"
      },
      "startsAt": "2024-01-01T12:00:00Z"
    }]
  }'
```

3. Kiểm tra logs:
```bash
docker logs workhub_telegram_webhook
```

## Troubleshooting

### Bot không nhận được messages
- Kiểm tra Bot Token và Chat ID đã đúng chưa
- Đảm bảo đã bắt đầu chat với bot trước
- Kiểm tra logs: `docker logs workhub_telegram_webhook`

### Alerts không được gửi
- Kiểm tra Alertmanager có kết nối được với webhook:
  ```bash
  docker logs workhub_alertmanager
  ```
- Kiểm tra Prometheus có alerts không:
  - Truy cập: `http://localhost:9090/alerts`

### Lỗi "Bad Request" từ Telegram API
- Kiểm tra Bot Token format
- Đảm bảo bot không bị block
- Kiểm tra message format (Markdown syntax)

## Cấu hình nâng cao

### Gửi vào nhiều chat
Sửa `telegram-webhook.py` để hỗ trợ nhiều chat IDs:

```python
TELEGRAM_CHAT_IDS = os.getenv('TELEGRAM_CHAT_IDS', '').split(',')
```

### Thêm filters
Chỉ gửi critical alerts vào Telegram, warning vào email, etc.

Sửa `alertmanager.yml` để route khác nhau.

