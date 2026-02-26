# Deployment Guide

## Environment Configuration

Aplikasi ini menggunakan environment variables untuk konfigurasi. Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env.example .env
```

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `HTTP_PORT` | `8080` | Port untuk REST API |
| `GRPC_PORT` | `50051` | Port untuk gRPC server |
| `DB_DRIVER` | `sqlite` | Database driver |
| `DB_PATH` | `test.db` | Path ke database file |
| `JWT_SECRET` | - | Secret key untuk JWT (WAJIB di production) |
| `JWT_EXPIRY_HOURS` | `24` | Durasi token dalam jam |
| `ENV` | `development` | Environment: development/production |

## Docker Deployment

### Build dan Run dengan Docker Compose

```bash
# Build dan start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

### Build Manual

```bash
# Build image
docker build -t user-crud-api .

# Run container
docker run -d \
  -p 8080:8080 \
  -p 50051:50051 \
  -e JWT_SECRET=your-secret-key \
  -e ENV=production \
  -v $(pwd)/data:/data \
  --name user-crud-api \
  user-crud-api
```

## Production Checklist

- [ ] Set `JWT_SECRET` ke nilai yang aman dan unik
- [ ] Set `ENV=production`
- [ ] Gunakan database yang persistent (bukan in-memory)
- [ ] Setup HTTPS/TLS untuk REST API
- [ ] Setup TLS untuk gRPC server
- [ ] Configure proper logging
- [ ] Setup monitoring dan alerting
- [ ] Backup database secara berkala
- [ ] Review dan update security headers
- [ ] Rate limiting untuk API endpoints

## Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status": "ok"}
```
