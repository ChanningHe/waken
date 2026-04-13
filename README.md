# Waken

![](./images/homepage.jpg)

Wake-on-LAN web application. Go backend with embedded React frontend, running in a single binary.

Supports device management via web UI and direct HTTP API calls with Bearer token authentication.

## Quick Start

```yaml
# docker-compose.yml
services:
  waken:
    image: ghcr.io/channinghe/waken:latest
    network_mode: host
    volumes:
      - ./config:/app/waken/config
    environment:
      - WOL_PORT=19527
      - WOL_AUTH_TOKEN=your-secret-token
      #- WOL_DB_PATH=/app/waken/config/wol.db
    restart: unless-stopped
```

```bash
docker compose up -d
```

Open `http://localhost:19527` to access the web UI.

`network_mode: host` is required for UDP broadcast packets to reach the local network.

## API

All endpoints except `/api/health` require `Authorization: Bearer <token>` header.

```
GET  /api/health          Health check
GET  /api/devices         List devices
POST /api/devices         Add device
PUT  /api/devices/{id}    Update device
DELETE /api/devices/{id}  Delete device
POST /api/wake/{id}       Wake device by ID
POST /api/wake            Wake by MAC address
```

Wake a device directly:

```bash
curl -X POST http://localhost:19527/api/wake \
  -H "Authorization: Bearer your-secret-token" \
  -H "Content-Type: application/json" \
  -d '{"mac": "AA:BB:CC:DD:EE:FF"}'
```

## Environment Variables

| Variable | Default | Description |
|---|---|---|
| `WOL_PORT` | `19527` | HTTP listen port |
| `WOL_AUTH_TOKEN` | empty (auth disabled) | Bearer token |
| `WOL_BROADCAST_ADDR` | `255.255.255.255` | Default broadcast address |
| `WOL_WOL_PORT` | `9` | Default WoL UDP port |
| `WOL_DB_PATH` | `/app/waken/config/wol.db` | SQLite database path |

## License

MIT
