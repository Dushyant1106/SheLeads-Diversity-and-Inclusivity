# 🌟 MahilaMitra

<div align="center">
<img width="2816" height="1536" alt="Gemini_Generated_Image_h40i2sh40i2sh40i" src="https://github.com/user-attachments/assets/1a922538-81c8-468d-8592-243f4c276a63" />

**Empowering Women Entrepreneurs Through AI-Powered Care Work Recognition**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![React](https://img.shields.io/badge/React-18.0+-61DAFB?style=for-the-badge&logo=react)](https://reactjs.org/)
[![MongoDB](https://img.shields.io/badge/MongoDB-6.0+-47A248?style=for-the-badge&logo=mongodb)](https://www.mongodb.com/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)

[Features](#-features) • [Quick Start](#-quick-start) • [Architecture](#-architecture) • [API Docs](#-api-documentation) • [Contributing](#-contributing)

</div>

---

## 📖 About

**MahilaMitra** (meaning "Women's Friend" in Hindi) is a revolutionary platform designed to recognize, value, and monetize the unpaid care work performed by women entrepreneurs in India. By combining AI-powered work verification, financial analytics, and government scheme recommendations, we're transforming invisible labor into recognized economic value.

### 🎯 Mission

To empower 10 million women entrepreneurs by 2030 through:
- 💰 **Recognition** of unpaid care work as economic value
- 🤖 **AI-powered** work verification and time tracking
- 📊 **Financial insights** for loan eligibility
- 🎓 **Access** to government schemes and training
- 🔒 **Privacy-first** approach with local data control

---

## ✨ Features

### 🎤 Voice-to-Log Work Tracker
- **Multilingual Support**: Voice input in English, Hindi, and Kannada
- **AI Verification**: Gemini 2.5 Flash validates work with image proof
- **Smart Categorization**: 11 work categories from cooking to farming
- **Automatic Time Estimation**: AI calculates realistic work hours
- **Points System**: Gamified recognition of effort

### 📈 Financial Analytics Dashboard
- **Market Value Calculator**: Converts work hours to monetary value (₹150-500/hour)
- **Loan Eligibility Scoring**: Based on work consistency and points
- **Category Breakdown**: Visual insights into work distribution
- **PDF Reports**: Professional monthly summaries with charts
- **Burnout Detection**: SMS alerts to emergency contacts via Twilio

### 🌱 Grow Module (Schemes & Training)
- **AI-Powered Recommendations**: Personalized government scheme matching
- **Real-time Search**: SerpAPI integration for latest schemes
- **Training Programs**: Skill development opportunities
- **Application Assistance**: Pre-filled forms with user data
- **24-hour Caching**: Fast, reliable scheme data

### 💬 Loan & Scheme Assistant Chatbot
- **Context-Aware AI**: Groq LLaMA 3.3 70B model
- **Personalized Advice**: References user's work logs and profile
- **Multi-topic Support**: Loans, schemes, documents, SHG, subsidies
- **Conversation History**: Maintains context across sessions
- **Instant Responses**: Sub-second AI inference

### ♿ Accessibility Features
- **8 Accessibility Options**: Bigger text, dyslexia-friendly fonts, cursor enhancement
- **Global CSS Variables**: Consistent accessibility across all pages
- **Keyboard Navigation**: Full keyboard support
- **Screen Reader Optimized**: ARIA labels and semantic HTML
- **Persistent Settings**: Saved in localStorage

### 🎨 Marketing Automation
- **AI Content Generation**: Gemini creates social media posts
- **Multi-platform Support**: Twitter, Instagram, Facebook
- **Image Generation**: Runway AI for visual content
- **Scheduling**: Plan posts in advance
- **Analytics Tracking**: Engagement metrics and insights

---

## 🚀 Quick Start

### Prerequisites

```bash
# Required
- Go 1.21+
- Node.js 18+
- MongoDB 6.0+

# API Keys (see .env.example)
- Google Gemini API
- Groq API
- Twilio (optional)
- SerpAPI (optional)
```

### Installation

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/mahilamitra.git
cd mahilamitra/sheleads-backend

# 2. Backend Setup
cp .env.example .env
# Edit .env with your API keys
go mod download
go run main.go

# 3. Frontend Setup (new terminal)
cd frontend
npm install
npm start

# 4. Access the application
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
```

### Docker Setup

```bash
# Quick start with Docker Compose
docker-compose up -d

# Access at http://localhost:3000
```

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     React Frontend (Port 3000)               │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │Dashboard │  │  Logger  │  │   Grow   │  │Marketing │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│         │              │              │              │       │
│         └──────────────┴──────────────┴──────────────┘       │
│                          │                                    │
│                    Axios API Client                          │
└──────────────────────────┼──────────────────────────────────┘
                           │
                    REST API (JSON)
                           │
┌──────────────────────────┼──────────────────────────────────┐
│                   Go Backend (Port 8080)                     │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                    Gin HTTP Router                      │ │
│  └────────────────────────────────────────────────────────┘ │
│           │              │              │              │     │
│    ┌──────────┐   ┌──────────┐  ┌──────────┐  ┌──────────┐│
│    │  Auth    │   │ WorkLog  │  │Analytics │  │  Grow    ││
│    │ Handler  │   │ Handler  │  │ Handler  │  │ Handler  ││
│    └──────────┘   └──────────┘  └──────────┘  └──────────┘│
│           │              │              │              │     │
│    ┌──────────────────────────────────────────────────────┐ │
│    │                   Services Layer                      │ │
│    │  • Gemini AI    • Groq AI     • Twilio              │ │
│    │  • SerpAPI      • Runway AI   • PDF Generator       │ │
│    └──────────────────────────────────────────────────────┘ │
│                           │                                  │
│    ┌──────────────────────────────────────────────────────┐ │
│    │              MongoDB Database                         │ │
│    │  • users  • worklogs  • business_profiles            │ │
│    │  • generated_content  • content_metrics              │ │
│    └──────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

---

## 📚 Tech Stack

### Backend
- **Framework**: Gin (Go)
- **Database**: MongoDB with official Go driver
- **AI Models**:
  - Google Gemini 2.5 Flash (work verification, content generation, scheme recommendations)
  - Groq LLaMA 3.3 70B (chatbot conversations)
  - Whisper Large V3 Turbo (voice transcription)
- **PDF Generation**: gofpdf
- **Authentication**: JWT with bcrypt
- **File Storage**: Local filesystem + AWS S3 ready

### Frontend
- **Framework**: React 18 with Hooks
- **UI Library**: shadcn/ui + Tailwind CSS
- **State Management**: React Context API
- **HTTP Client**: Axios
- **Charts**: Recharts
- **Icons**: Lucide React
- **Notifications**: Sonner
- **Routing**: React Router v6

### External Services
- **Voice Transcription**: Groq Whisper API
- **Scheme Search**: SerpAPI (Google Search)
- **SMS Alerts**: Twilio
- **Image Generation**: Runway AI (optional)

---

## 🔑 Key Endpoints

### Authentication
```http
POST   /api/v1/auth/signup          # Create new account
POST   /api/v1/auth/login           # Login with credentials
GET    /api/v1/auth/profile         # Get user profile
PUT    /api/v1/auth/profile         # Update profile
```

### Work Logging
```http
POST   /api/v1/work/logs            # Create work log with image
POST   /api/v1/work/voice-to-log    # Voice-to-log conversion
GET    /api/v1/work/logs            # Get all work logs
DELETE /api/v1/work/logs/:id        # Delete work log
```

### Analytics
```http
GET    /api/v1/analytics/summary    # Get analytics summary
GET    /api/v1/analytics/stats      # Get activity statistics
GET    /api/v1/analytics/burnout    # Check burnout status
GET    /api/v1/analytics/market-value  # Calculate market value
```

### Grow (Schemes & Training)
```http
GET    /api/v1/grow/schemes         # Get scheme recommendations
GET    /api/v1/grow/training        # Get training programs
```

### Chatbot
```http
POST   /api/v1/chatbot/query        # Send chatbot query
```

### Marketing
```http
POST   /api/v1/business/profile     # Create business profile
GET    /api/v1/business/profile     # Get business profile
POST   /api/v1/content/generate     # Generate content
GET    /api/v1/content/list         # List generated content
POST   /api/v1/metrics              # Add metrics
```

### Reports
```http
GET    /api/v1/reports/monthly/pdf  # Generate PDF report
GET    /api/v1/reports/monthly/data # Get monthly data
```

📖 **Full API Documentation**: See [API_DOCUMENTATION.md](API_DOCUMENTATION.md)

---

## 📁 Project Structure

```
sheleads-backend/
├── config/                 # Configuration management
│   └── config.go          # Environment variables
├── database/              # MongoDB connection
│   └── database.go        # DB initialization
├── handlers/              # HTTP request handlers
│   ├── auth.go           # Authentication
│   ├── worklog.go        # Work logging
│   ├── analytics.go      # Analytics & reports
│   ├── grow.go           # Schemes & training
│   ├── chatbot.go        # Chatbot queries
│   ├── content.go        # Content generation
│   └── business.go       # Business profiles
├── middleware/            # HTTP middleware
│   ├── auth.go           # JWT authentication
│   └── cors.go           # CORS configuration
├── models/                # Data models
│   ├── user.go           # User model
│   ├── worklog.go        # WorkLog model
│   ├── business.go       # Business profile
│   └── content.go        # Generated content
├── services/              # Business logic
│   ├── gemini.go         # Gemini AI service
│   ├── groq.go           # Groq API service
│   ├── twilio.go         # SMS service
│   ├── serp.go           # Scheme search
│   ├── pdf.go            # PDF generation
│   ├── burnout.go        # Burnout detection
│   └── runway.go         # Image generation
├── utils/                 # Helper functions
│   ├── jwt.go            # JWT utilities
│   └── response.go       # API responses
├── frontend/              # React application
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── pages/        # Page components
│   │   ├── contexts/     # Context providers
│   │   ├── lib/          # Utilities
│   │   └── styles/       # CSS files
│   └── public/           # Static assets
├── uploads/               # User-uploaded files
├── generated_images/      # AI-generated images
├── main.go               # Application entry point
├── go.mod                # Go dependencies
├── .env.example          # Environment template
└── README.md             # This file
```

---

## 🔐 Environment Variables

Create a `.env` file with the following variables:

```env
# Server Configuration
PORT=8080

# MongoDB Configuration
MONGO_URL=mongodb+srv://username:password@cluster.mongodb.net/
DB_NAME=mahilamitra

# JWT Secret (change in production!)
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Google Gemini AI API Key
GEMINI_API_KEY=your-gemini-api-key-here

# Groq API Key (for voice & chatbot)
GROQ_API_KEY=your-groq-api-key-here

# Twilio Configuration (optional)
TWILIO_ACCOUNT_SID=your-twilio-account-sid
TWILIO_AUTH_TOKEN=your-twilio-auth-token
TWILIO_MESSAGING_SERVICE_SID=your-messaging-service-sid

# Burnout Detection Settings
BURNOUT_HOURS_THRESHOLD=12
BURNOUT_DAYS_WINDOW=7

# SerpAPI (optional - for scheme search)
SERP_API_KEY=your-serpapi-key

# AWS S3 (optional - for file storage)
AWS_ACCESS_KEY_ID=your-aws-access-key
AWS_SECRET_ACCESS_KEY=your-aws-secret-key
AWS_S3_BUCKET_NAME=your-bucket-name
AWS_REGION=ap-south-1

# Runway AI (optional - for image generation)
RUNWAY_API_KEY=your-runway-api-key
```

---

## 🎨 Screenshots

### Dashboard
![Dashboard](docs/screenshots/dashboard.png)

### Voice-to-Log
![Voice Logger](docs/screenshots/voice-logger.png)

### Grow Module
![Schemes](docs/screenshots/schemes.png)

### Chatbot
![Chatbot](docs/screenshots/chatbot.png)

---

## 🧪 Testing

### Run Backend Tests
```bash
go test ./... -v
```

### Test API Endpoints
```bash
# Use the provided test script
chmod +x test_api.sh
./test_api.sh

# Or use curl commands from api-examples.http
```

### Test Twilio SMS
```bash
chmod +x test_twilio_sms.sh
./test_twilio_sms.sh
```

---

## 📊 Work Categories & Market Rates

| Category | Rate (₹/hour) | Description |
|----------|---------------|-------------|
| Cooking | 300 | Meal preparation, food processing |
| Cleaning | 250 | House cleaning, maintenance |
| Childcare | 350 | Child supervision, education support |
| Elderly Care | 400 | Senior care, medical assistance |
| Tutoring | 500 | Educational tutoring, homework help |
| Tailoring | 350 | Sewing, alterations, embroidery |
| Handicrafts | 400 | Craft making, artistic work |
| Farming | 200 | Agricultural work, gardening |
| Animal Husbandry | 250 | Livestock care, dairy work |
| Water Collection | 150 | Water fetching, management |
| Other | 200 | Miscellaneous care work |

---

## 🤝 Contributing

We welcome contributions! Here's how you can help:

1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/amazing-feature`
3. **Commit your changes**: `git commit -m 'Add amazing feature'`
4. **Push to the branch**: `git push origin feature/amazing-feature`
5. **Open a Pull Request**

### Development Guidelines
- Follow Go best practices and conventions
- Write tests for new features
- Update documentation
- Use meaningful commit messages
- Ensure code passes linting

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Google Gemini** for AI-powered work verification
- **Groq** for fast voice transcription and chatbot
- **Twilio** for SMS alert infrastructure
- **SerpAPI** for real-time scheme search
- **shadcn/ui** for beautiful UI components
- All the women entrepreneurs who inspired this platform

---

## 📞 Support

- **Documentation**: [Full Docs](docs/)
- **Issues**: [GitHub Issues](https://github.com/yourusername/mahilamitra/issues)
- **Email**: support@mahilamitra.com
- **Community**: [Discord Server](https://discord.gg/mahilamitra)

---

<div align="center">

**Made with ❤️ for Women Entrepreneurs in India**

⭐ Star us on GitHub if this project helped you!

</div>

