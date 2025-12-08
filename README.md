->Architecture/Approach

1.service/libhttp

Gin routes

Request parsing & validation

JSON responses

Admin authentication middleware

2.service/handler

Business logic (issue, return, fine calculation)

Validations (copies, availability, duplicates)

Keeps controllers clean and thin

3.service/repository

MySQL queries using sqlx

CRUD operations separated by entity

SQL stored separately in queries.go

4.service/models

Structs representing Books, Members, Issues

JSON & DB field mappings

5.Auto Database Schema

Schema is created on server startup

No migrations needed for initial setup
--------------------------

1. Start MySQL

Create database:

CREATE DATABASE library_db;

2. Configure Environment Variables (optional)

Defaults exist — override if needed:

DB_USER=root
DB_PASS=yourpassword
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=database_name

ADMIN_PASSWORD=admin123
ADMIN_TOKEN=admintoken

3. Run the server
go run main.go


Server runs at:

http://localhost:8080


Database tables are auto-created at startup.

Admin Authentication
Step 1 — Login
POST /admin/login


Request:

{ "password": "admin123" }


Response:

{ "token": "admintoken" }

Step 2 — Use token for all admin routes
X-Admin-Token: admintoken

 API Endpoints
 Public Endpoints
Books
Method	Route	Description
GET	/books	List all books
GET	/books/:id	Get book by ID
GET	/books/search?q=keyword	Search by title/author
Members
Method	Route	Description
GET	/members/:id	Get member details
 Admin Endpoints (Require Token)
Books
Method	Route	Description
POST	/admin/books	Create book
PUT	/admin/books/:id	Update book
DELETE	/admin/books/:id	Delete book
GET	/admin/books	Admin list books
Members
Method	Route	Description
POST	/admin/members	Create member
PUT	/admin/members/:id	Update member
DELETE	/admin/members/:id	Delete member
GET	/admin/members	List all
Issues
Method	Route	Description
POST	/admin/issues	Issue a book
POST	/admin/issues/:id/return	Return a book (auto fine calc)
GET	/admin/issues/member/:member_id	All issues for member

 Sample API Requests
Create Book
POST /admin/books
{
  "title": "Go Programming Language",
  "author": "Alan Donovan",
  "copies": 3
}

Issue a Book
POST /admin/issues
{
  "book_id": 1,
  "member_id": 2,
  "due_days": 7
}

Return Book
POST /admin/issues/10/return


Response example:

{ "fine": 4 }

NOTE: I used postman api tool for test the all payloads for responses.