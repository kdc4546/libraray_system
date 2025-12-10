
Library Management System:

->Managing books (copies, availability)
->Managing members
->Issuing and returning books

implemented with
->gin-framework for HTTP routing
->sqlx + MySQL for database

Architecture
in main.go file:
->Connects to MySQL
->Ensures DB schema exists
->Initializes Gin router
->Registers all routes
->Starts the HTTP server

HTTP / (service/libhttp):
Contains Gin handlers for all endpoints (books, members, issues, admin login).
->Parsing path params, query params, and JSON bodies
->Validating requests
->Mapping HTTP status codes and JSON responses
->calls to the service/handler layer.

Service  (service/handler):
-> List / create / update / delete books and members
-> Issue a book
-> Return a book and add fines

Repository (service/repository):
have interfaces:
BookRepo
MemberRepo
IssueRepo

Database (repository/db):
MySQL connection and schema creation
EnsureSchema creates tables books, members, issues with foreign keys and timestamps

----------------------------------------
I used postman api tool to check the payloads
 login:
 POST /admin/login
-> admin login with username/password(implemented : admin/123)

 Books:
GET /admin/books - list all books
POST /admin/books - create a book
{
  "title": "book A",
  "author": "author_a",
  "copies": 5
}

PUT /admin/books/:id - update book
{
  "title": "book B" - Updated",
  "author": "author_b",
  "copies": 10
}

DELETE /admin/books/:id - delete book

 Members:
GET /admin/members - list all members
POST /admin/members - create a member
{
  "name": "p_1",
  "email": "p_1@gmail.com",
  "roll_no": "1"
}

PUT /admin/members/:id - update member
{
  "name": "p_2",
  "email": "p_2@gmail.com"
}

DELETE /admin/members/:id - delete member

 Issues:
POST /admin/issues - Issue a book to a member.
{
  "book_id": 1,
  "member_id": 1,
  "due_days": 7
}

POST /admin/issues/:id/return - Return a book and compute any fine.
{
  "message": "book returned",
  "fine_paid": 20
}
GET /admin/issues/member/:member_id - List all issues for a member.

--------------------------------------------
I used atomic conditional UPDATE queries in the database
(
UPDATE books SET 
available = available - 1 
WHERE id = ? 
  available > 0
)
 which have race-free borrow and return operations. Borrow and return logic now runs inside database transactions, that no two users can borrow the same last copy.

