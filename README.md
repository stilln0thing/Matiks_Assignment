# 🏆 Matiks Leaderboard

A scalable real-time leaderboard system supporting 10,000+ users with instant rank updates and search.

## 🚀 Live Demo

- **Frontend**: https://matiks-ranklist.vercel.app/
- **Backend API**: https://matiks-leaderboard-backend-gdj4.onrender.com

## 📁 Project Structure

```
├── backend/          # Go backend API
│   ├── cmd/server/   # Entry point
│   └── internal/     # Business logic
└── frontend/         # React Native (Expo) app
    └── src/
        ├── api/      # API client
        ├── components/
        └── screens/
```

## ⚡ Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | React Native (Expo), TypeScript |
| Backend | Go, Gin Framework |
| Database | PostgreSQL |
| Cache | Redis (Sorted Sets) |

## 🎯 Features

- **Real-time Leaderboard** - Paginated with infinite scroll
- **Instant Search** - Debounced username search with live ranks
- **Tie-aware Ranking** - Accurate rankings using Redis sorted sets
- **Auto Score Updates** - Background simulator updates ratings every second
- **Async DB Writes** - Batched writes for high throughput

## 🛠️ Local Development

### Backend

```bash
cd backend

# Start PostgreSQL & Redis
docker-compose up -d

# Run server
go run cmd/server/main.go
```

### Frontend

```bash
cd frontend
npm install
npm run start
# Press 'w' for web, 'a' for Android, 'i' for iOS
```

## 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/leaderboard?limit=50&offset=0` | Paginated leaderboard |
| GET | `/api/search?q=player` | Search users |
| GET | `/api/user/:id/rank` | Get user's rank |
| POST | `/api/rating` | Update rating |
| GET | `/health` | Health check |

## 🌐 Deployment

- **Backend**: Render.com (with PostgreSQL + Redis)
- **Frontend**: Vercel
