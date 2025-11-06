#!/usr/bin/env python3
"""
Telegram webhook service để nhận alerts từ Alertmanager và gửi vào Telegram
"""
import os
import json
import requests
from flask import Flask, request, jsonify
from datetime import datetime

app = Flask(__name__)

# Telegram Bot configuration
TELEGRAM_BOT_TOKEN = os.getenv('TELEGRAM_BOT_TOKEN', '')
TELEGRAM_CHAT_ID = os.getenv('TELEGRAM_CHAT_ID', '')

def format_alert_message(alert_data):
    """Format alert data thành message cho Telegram"""
    status = alert_data.get('status', 'unknown')
    group_labels = alert_data.get('groupLabels', {})
    common_labels = alert_data.get('commonLabels', {})
    alerts = alert_data.get('alerts', [])
    
    # Emoji theo status
    if status == 'firing':
        emoji = '🔴'
        status_text = 'FIRING'
    else:
        emoji = '✅'
        status_text = 'RESOLVED'
    
    # Header
    alertname = group_labels.get('alertname', 'Unknown Alert')
    severity = common_labels.get('severity', 'unknown')
    service = common_labels.get('service', 'unknown')
    
    message = f"{emoji} *{status_text}*\n\n"
    message += f"*Alert:* {alertname}\n"
    message += f"*Severity:* {severity}\n"
    message += f"*Service:* {service}\n\n"
    
    # Alerts details
    for alert in alerts:
        annotations = alert.get('annotations', {})
        labels = alert.get('labels', {})
        
        if annotations.get('description'):
            message += f"*Description:* {annotations['description']}\n"
        if annotations.get('summary'):
            message += f"*Summary:* {annotations['summary']}\n"
        
        starts_at = alert.get('startsAt', '')
        if starts_at:
            try:
                dt = datetime.fromisoformat(starts_at.replace('Z', '+00:00'))
                message += f"*Started:* {dt.strftime('%Y-%m-%d %H:%M:%S')}\n"
            except:
                message += f"*Started:* {starts_at}\n"
        
        message += "\n"
    
    # Common labels
    if common_labels:
        message += "*Labels:*\n"
        for key, value in sorted(common_labels.items()):
            if key not in ['alertname', 'severity', 'service']:
                message += f"  • {key}: {value}\n"
    
    return message

def send_to_telegram(message):
    """Gửi message vào Telegram"""
    if not TELEGRAM_BOT_TOKEN or not TELEGRAM_CHAT_ID:
        print("ERROR: Telegram bot token or chat ID not configured")
        return False
    
    url = f"https://api.telegram.org/bot{TELEGRAM_BOT_TOKEN}/sendMessage"
    payload = {
        'chat_id': TELEGRAM_CHAT_ID,
        'text': message,
        'parse_mode': 'Markdown'
    }
    
    try:
        response = requests.post(url, json=payload, timeout=10)
        response.raise_for_status()
        return True
    except Exception as e:
        print(f"ERROR sending to Telegram: {e}")
        return False

@app.route('/webhook', methods=['POST'])
def webhook():
    """Endpoint để nhận alerts từ Alertmanager"""
    try:
        data = request.get_json()
        
        if not data:
            return jsonify({'error': 'No data received'}), 400
        
        # Format và gửi message
        message = format_alert_message(data)
        success = send_to_telegram(message)
        
        if success:
            return jsonify({'status': 'success'}), 200
        else:
            return jsonify({'status': 'error', 'message': 'Failed to send to Telegram'}), 500
            
    except Exception as e:
        print(f"ERROR processing webhook: {e}")
        return jsonify({'error': str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    """Health check endpoint"""
    return jsonify({'status': 'healthy'}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000, debug=False)

