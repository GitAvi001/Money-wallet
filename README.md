# Money Transfer Application 💸
This application is built using a **microservices architecture** with separate services for authentication, transactions, and an API gateway for unified access. Users can register, verify their email, add funds to their wallet, withdraw money, and transfer funds to other users seamlessly.

---

## 🛠️ Tech Stack

### **Backend**
- **Language:** Go (Golang) 1.21+
- **Framework:** Gin (Web Framework)
- **Database:** PostgreSQL 12+
- **Authentication:** JWT (JSON Web Tokens)
- **Password Hashing:** bcrypt
- **Email:** SMTP (Gmail)

### **Frontend**
- **Framework:** React 18+
- **Build Tool:** Vite
- **Styling:** Tailwind CSS

### **Architecture**
```bash
#Port numbers can configured as per the requirements
```
- **API Gateway** (Port 8080) - Unified entrypoint for routeing requests to appropriate services
- **Auth Service** (Port 8081) - User authentication and authorization
- **Transaction Service** (Port 8082) - Wallet and transaction management

### **Go Packages Used**
- `github.com/gin-gonic/gin` - Web framework
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto/bcrypt` - Password hashing
- `github.com/joho/godotenv` - Environment variable loader
- `github.com/google/uuid` - UUID generation

---

## 🗄️ Database Connection Setup (Local)

### **Step 1: Install PostgreSQL**

Download and install PostgreSQL from [postgresql.org](https://www.postgresql.org/download/)

### **Step 2: Start PostgreSQL Server**

**Windows:**
```bash
# PostgreSQL should start automatically after installation
# Or use Services app to start "postgresql-x64-XX"
```

**macOS:**
```bash
brew services start postgresql
```

**Linux:**
```bash
sudo systemctl start postgresql
```

### **Step 3: Create Database**

Open PostgreSQL command line:

```bash
# Connect to PostgreSQL
psql -U postgres

# Enter your postgres password when prompted
```

Create the database:

```sql
CREATE DATABASE money_transfer_db;

-- Verify database created
\l

-- Exit
\q
```

### **Step 4: Configure Database Connection**

Update the `DATABASE_URL` in your `.env` files:

**Format:**
```
DATABASE_URL=postgres://USERNAME:PASSWORD@HOST:PORT/DATABASE?sslmode=disable
```

**Example:**
```
DATABASE_URL=postgres://postgres:<your_postgres_password>@localhost:5432/money_transfer_db?sslmode=disable
```

**Components:**
- `USERNAME`: Your PostgreSQL username (default: `postgres`)
- `PASSWORD`: Your PostgreSQL password
- `HOST`: Database host (local: `localhost`)
- `PORT`: PostgreSQL port (default: `5432`)
- `DATABASE`: Database name (`money_transfer_db`)

---

## 🚀 Start Services

Follow these steps in order to run the complete application:

### **Prerequisites**

Ensure you have installed:
- ✅ Go 1.21 or higher ([Download](https://go.dev/dl/))
- ✅ PostgreSQL 12 or higher ([Download](https://www.postgresql.org/download/))
- ✅ Node.js 16+ and npm ([Download](https://nodejs.org/))
- ✅ Git ([Download](https://git-scm.com/))

---

### **Step 1: Clone Repository & Setup Environment**

```bash
# Clone the repository
git clone https://github.com/GitAvi001/Money-wallet.git
```

```bash
# Navigate to project directory
cd money-wallet
```

---

### **Step 2: Configure Auth Service**

```bash
# Navigate to auth-service
cd auth-service

# Copy environment template
copy .env.example .env

# Open .env and configure:
# - DATABASE_URL (your PostgreSQL connection string)
# - JWT_SECRET (generate a secure random string)
# - SMTP_USERNAME and SMTP_PASSWORD (your Gmail credentials)
# - FROM_EMAIL (your Gmail address)
```

**Example `.env` for Auth Service:**
```env
PORT=8081
DATABASE_URL=postgres://postgres:<your_postgres_password>@localhost:5432/money_transfer_db?sslmode=disable
JWT_SECRET=7Kx9mP2nQ5vL8wR3yT6zA4bC1dE0fG9hJ2kM5nP8qS1tU4vW7xY0zA3bC6dE9fG //Generate any secret key
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=youremail@gmail.com
SMTP_PASSWORD=your-gmail-app-password
FROM_EMAIL=youremail@gmail.com //registers this email as the SMTP server
FRONTEND_URL=http://localhost:5173
```

**Install dependencies:**
```bash
go mod tidy
```

---

### **Step 3: Configure Transaction Service**

```bash
# Navigate to transaction-service (from project root)
cd ..\transaction-service

# Copy environment template
copy .env.example .env

# Open .env and configure:
# - DATABASE_URL (same as auth-service)
# - JWT_SECRET (MUST match auth-service)
```

**Example `.env` for Transaction Service:**
```env
PORT=8082
DATABASE_URL=postgres://postgres:<your_postgre_password>@localhost:5432/money_transfer_db?sslmode=disable
JWT_SECRET=7Kx9mP2nQ5vL8wR3yT6zA4bC1dE0fG9hJ2kM5nP8qS1tU4vW7xY0zA3bC6dE9fG
AUTH_SERVICE_URL=http://localhost:8081
```

**Install dependencies:**
```bash
go mod tidy
```

---

### **Step 4: Configure API Gateway**

```bash
# Navigate to api-gateway (from project root)
cd ..\api-gateway

# Copy environment template
copy .env.example .env
```

**Example `.env` for API Gateway:**
```env
PORT=8080
AUTH_SERVICE_URL=http://localhost:8081
TRANSACTION_SERVICE_URL=http://localhost:8082
```

**Install dependencies:**
```bash
go mod tidy
```

---

### **Step 5: Setup Frontend**

```bash
# Navigate to frontend (from project root)
cd ..\frontend

# Install dependencies
npm install
```

---

### **Step 6: Start All Services**

Open **4 separate terminal windows** and run each service:

#### **Terminal 1: Start Auth Service**
```bash
cd auth-service
go run main.go
```
**Output:**
```
Database connected successfully
Database migrations completed successfully
Auth Service starting on port 8081
```

#### **Terminal 2: Start Transaction Service**
```bash
cd transaction-service
go run main.go
```
**Output:**
```
Database connected successfully
Database migrations completed successfully
Transaction Service starting on port 8082
```

#### **Terminal 3: Start API Gateway**
```bash
cd api-gateway
go run main.go
```
**Output:**
```
API Gateway starting on port 8080
Auth Service URL: http://localhost:8081
Transaction Service URL: http://localhost:8082
```

#### **Terminal 4: Start Frontend**
```bash
cd frontend
npm run dev
```
**Output:**
```
VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
```

---

### **Step 7: Access the Application**

Open your browser and navigate to:

🌐 **http://localhost:5173**

---

## 🔌 Service URLs

| Service | URL | Purpose |
|---------|-----|---------|
| **Frontend** | http://localhost:5173 | User interface |
| **API Gateway** | http://localhost:8080 | Unified API access |
| **Auth Service** | http://localhost:8081 | Authentication (direct) |
| **Transaction Service** | http://localhost:8082 | Transactions (direct) |

---

## 📡 API Endpoints

### **Auth Service** (via `/api/auth`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/auth/register` | Register new user | ❌ |
| POST | `/api/auth/login` | User login | ❌ |
| GET | `/api/auth/verify-email` | Verify email | ❌ |
| POST | `/api/auth/send-verification` | Resend verification | ❌ |
| GET | `/api/auth/me` | Get current user | ✅ |
| GET | `/api/auth/users` | Get all users | ✅ |

### **Transaction Service** (via `/api`)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/wallet` | Get user wallet | ✅ |
| POST | `/api/wallet/add` | Add funds to wallet | ✅ |
| POST | `/api/wallet/withdraw` | Withdraw funds | ✅ |
| POST | `/api/transactions/transfer` | Transfer money | ✅ |
| GET | `/api/transactions` | Get transaction history | ✅ |
| GET | `/api/transactions/:id` | Get specific transaction | ✅ |

---

## ✨ Features

- ✅ **Microservices Architecture** - Scalable and maintainable
- ✅ **API Gateway** - Unified access point with CORS support
- ✅ **User Registration** - Create account with email verification
- ✅ **Email Verification** - Secure email confirmation via SMTP
- ✅ **JWT Authentication** - Stateless token-based authentication
- ✅ **Wallet Management** - Add and withdraw funds
- ✅ **Money Transfers** - Send money to other verified users
- ✅ **Transaction History** - View all past transactions
- ✅ **Auto Migrations** - Database tables created automatically
- ✅ **Password Security** - bcrypt hashing with salt

---

## 🐛 Troubleshooting

### **1. Database Connection Failed**

**Error:** `Failed to connect to database`

**Solution:**
- Verify PostgreSQL is running
- Check `DATABASE_URL` in `.env` files
- Ensure database `money_transfer_db` exists
- Verify username and password are correct

```bash
# Test connection
psql -U postgres -d money_transfer_db
```

---

### **2. JWT Token Invalid**

**Error:** `Invalid or expired token`

**Solution:**
- Ensure `JWT_SECRET` is **identical** in both:
  - `auth-service/.env`
  - `transaction-service/.env`
- Re-login to get a new token

---

### **3. Import Errors**

**Error:** `no required module provides package`

**Solution:**
```bash
# Run in service directory
go mod tidy
```

---

### **4. SMTP Email Not Sending**

**Error:** `Failed to send email`

**Solution:**
- Use Gmail App Password (not regular password)
- Generate at: https://myaccount.google.com/apppasswords
- Enable 2-Step Verification first
- Update `SMTP_PASSWORD` in `auth-service/.env`

---

### **5. Frontend Not Starting**

**Error:** `npm run dev` fails

**Solution:**
```bash
# Delete node_modules and reinstall
rm -rf node_modules
rm package-lock.json
npm install
npm run dev
```

---

### **6. CORS Errors**

**Error:** `CORS policy: No 'Access-Control-Allow-Origin' header`

**Solution:**
- Ensure all services are running
- Check `FRONTEND_URL` in `auth-service/.env` matches frontend URL
- Verify API Gateway is routing correctly

---

### **7. Port Already in Use**

**Error:** `bind: address already in use`

**Solution:**
```bash
# Find and kill process using the port (example for port 8080)
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

---

## 🔒 Security Notes

⚠️ **Important for Production:**

1. **Change JWT Secret**
   - Generate a long, random string
   - Use different secrets for different environments

2. **Environment Files**
   - Never commit `.env` files to Git
   - Use environment-specific configurations

3. **Database**
   - Enable SSL/TLS: Change `sslmode=disable` to `sslmode=require`
   - Use strong passwords
   - Restrict database access

4. **SMTP Credentials**
   - Use Gmail App Passwords
   - Never use plain passwords
   - Rotate credentials regularly

5. **Additional Security**
   - Implement rate limiting
   - Add request validation
   - Enable HTTPS
   - Use secure headers
   - Implement logging and monitoring

