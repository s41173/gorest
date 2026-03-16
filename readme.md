# Go REST API Demo

A lightweight **Go REST API** project built for demonstration purposes.

---

## 🚀 Features

- RESTful API using **Go** and **GORM**
- **MySQL (Railway Cloud)** for persistent data
- **Redis (Railway Cloud)** for caching
- **JWT** for Token
- Whatsapp Notification for OTP Request
- Environment-based configuration via `.env`
- Quick deployment to cloud (Railway demo)
- TLS connection to cloud database (demo-friendly)

---

## 🏗 Architecture / Tech Stack

Client (Browser / Postman)
│
▼
Go REST API (GORM)
│
▼
┌───────────────┐
│ MySQL DB │ ← Aiven Cloud
└───────────────┘
│
▼
┌───────────────┐
│ Redis Cache │ ← Redis Cloud
└───────────────┘


- **Language:** Go  
- **ORM:** GORM  
- **Database:** MySQL (Aiven Cloud)  
- **Cache:** Redis Cloud  
- **Deployment:** Railway (or local)  
- **Environment Variables:** `.env`  

---

## 📦 Collections / Endpoints

### Auth
| Method | Endpoint      | Description           |
|--------|---------------|---------------------|
| POST   | /login        | Login               |       
| GET    | /decode       | Decode Token        |
| GET    | /logout       | Logout              |
| POST   | /forgot       | Forgot Password     |
| POST   | /otp          | Request OTP         |
| POST   | /verify       | Verify User         |

### Product
| Method | Endpoint      | Description           |
|--------|---------------|---------------------|
| GET    | /product       | Get all products    |
| GET    | /product/:id   | Get product by ID   |
| POST   | /product       | Add new product     |

### Users (example)
| Method | Endpoint   | Description          |
|--------|------------|----------------------|
| POST   | /user      | Register user        |
| PUT    | /user      | Update user          |
| PUT    | /password  | Change user password |
| POST   | /image     | Upload user image    |

---

## ⚙ Usage

1. Create a `.env` file in project root with these variables:
