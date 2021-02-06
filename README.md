# Notifications Microservice

## Make a config.json inside configuration folder using following syntax:

```javascript
{
    "server": {
        "port": "8080",
        "domain": "0.0.0.0"
    },
    "database": {
        "user": "<your_postgres_user_name>",
        "password": "<your_postgres_user_password>",
        "dbname": "<your_database_name>",
        "sslmode": "<ssl_mode>"
    },
    "token": {
        "secret_key": "<token_key>",
        "header_prefix": "Bearer",
        "expiry_time": {
            "validation_token": 2
        }
    },
    "email_notification": {
        "email": "<email_id>",
        "password": "<password>",
        "smtp_host": "smtp.gmail.com",
        "smtp_port": "587"
    },
    "password_hash": "randomString"
}
```
