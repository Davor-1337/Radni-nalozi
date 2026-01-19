🛠️ Work Order Management System

A full-stack web application designed to manage work orders, users, materials, and financial documentation in an efficient and centralized way.
The system supports administrators, technicians, and clients, each with clearly defined roles and responsibilities.

✨ Key Features
👨‍💼 Administrator Panel
Work Order Management

Create, view, update, and archive work orders

Assign work orders to technicians

Monitor work order lifecycle and status

User Management

Approve or reject registration requests (clients and technicians)

Manage user roles and access privileges

Material & Inventory Management

Track material usage per work order

Monitor stock availability in real time

Financial Module

Generate invoices

Export detailed PDF reports

Notification System

Real-time notifications for new work order requests

Alerts for pending user registration approvals

👤 Client Portal
Work Orders

Submit new work order requests

View all personal work orders

Track status (Open / Accepted / Completed)

Documentation

View issued invoices

Download PDF reports

Notifications

Receive updates about work order status changes

🧑‍🔧 Technician Dashboard
Work Process

View assigned work orders

Log time spent on tasks

Record used materials

Complete work orders

Manage and track materials

Communication

Notifications about new assignments

Alerts for urgent work orders

🧱 Technology Stack
Backend

Go (1.x+)

REST API

Microsoft SQL Server

Frontend

Angular CLI 19.1.8

TypeScript

HTML5 / SCSS

Database

Microsoft SQL Server

⚙️ Prerequisites

Before running the application, make sure the following tools are installed:

Go
 (version 1.x or newer)

Node.js

Angular CLI

Microsoft SQL Server

🚀 Installation & Setup
Backend (Go)

Clone the repository:

git clone https://github.com/Davor-1337/Radni-nalozi
cd radni-nalozi/backend


Install dependencies:

go mod tidy


Configure the database (see Database Configuration below)

Run the backend server:

go run main.go

Frontend (Angular)

Navigate to the frontend directory:

cd radni-nalozi/frontend


Install dependencies:

npm install


Start the Angular application:

ng serve


Open the application in your browser:
👉 http://localhost:4200

🗄️ Database Setup & Configuration
1. Database Initialization
Option A: Using SQL Server Management Studio (SSMS)

Open SQL Server Management Studio

Create a new query window

Execute the script:

backend/database/baza.sql


Confirm that the database RadniNaloziDB has been successfully created

2. Environment Configuration

Create a .env file in the backend directory with the following content:

DB_CONNECTION=sqlserver://<username>:<password>@localhost:1433?database=RadniNaloziDB


Replace:

<username> with your SQL Server username

<password> with your SQL Server password

3. Verification

Start the backend service

Ensure the application successfully connects to the database

Verify that all tables and initial data are loaded correctly

🔐 Test Accounts

For testing purposes, the following accounts are available:

Administrator

Username: davor

Password: davor

Technician

Username: Petar00

Password: Petar00

Client

Username: SwiftP

Password: SwiftP

📌 Notes

This project was developed as part of an academic thesis and demonstrates a real-world business workflow

Emphasis is placed on role-based access, clean architecture, and scalable backend design using Go




















   
