Unknown Request

If the system does not understand:

Guest
 ↓
AI
 ↓
UNKNOWN
 ↓
Clarification

Example:

"Saya ingin memastikan saya membantu dengan tepat. Apakah Anda ingin bertanya tentang fasilitas hotel, melakukan service request, atau melaporkan masalah kamar?"



---


# 4. `docs/DATABASE.md`


```markdown
# 🗄️ Database Specification


Database engine:


```text
MySQL

1. Guests

Stores guest session information.

CREATE TABLE guests (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    session_id VARCHAR(100) NOT NULL UNIQUE,
    room_number VARCHAR(20),
    language VARCHAR(10) DEFAULT 'id',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
2. Conversations

Stores AI conversation history.

CREATE TABLE conversations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    session_id VARCHAR(100) NOT NULL,
    role ENUM('user', 'assistant', 'system') NOT NULL,
    message TEXT NOT NULL,
    intent VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
3. Tickets

Main operational ticket table.

CREATE TABLE tickets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ticket_number VARCHAR(30) NOT NULL UNIQUE,
    room_number VARCHAR(20) NOT NULL,
    department ENUM('HOUSEKEEPING', 'ENGINEERING', 'FRONT_OFFICE') NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    priority ENUM('LOW', 'MEDIUM', 'HIGH') DEFAULT 'MEDIUM',
    status ENUM('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CANCELLED') DEFAULT 'OPEN',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    resolved_at TIMESTAMP NULL
);
4. Hotel Information

Stores verified hotel knowledge.

CREATE TABLE hotel_information (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    category VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

Examples:

BREAKFAST
WIFI
POOL
GYM
CHECKIN
CHECKOUT
RESTAURANT
ROOM
POLICY
5. Recommendations

Stores verified recommendation data.

CREATE TABLE recommendations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100) NOT NULL,
    description TEXT,
    price_min INT,
    price_max INT,
    distance_km DECIMAL(5,2),
    address TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
6. Staff

Stores staff accounts.

CREATE TABLE staff (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    department ENUM('HOUSEKEEPING', 'ENGINEERING', 'FRONT_OFFICE', 'ADMIN') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
7. Ticket Assignment

Optional table for staff assignment.

CREATE TABLE ticket_assignments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    ticket_id BIGINT NOT NULL,
    staff_id BIGINT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (ticket_id) REFERENCES tickets(id),
    FOREIGN KEY (staff_id) REFERENCES staff(id)
);
8. Relationships
Guest
 │
 └── Conversation
        │
        └── AI Interaction


Guest
 │
 └── Ticket
        │
        └── Ticket Assignment
                │
                └── Staff


Hotel Information
        │
        └── AI Knowledge


Recommendations
        │
        └── AI Recommendation
9. Database Rules
Room number must never be invented.
Ticket status must use predefined values.
Department must use predefined values.
Priority must use predefined values.
AI cannot execute raw SQL.
Repository layer is responsible for database interaction.
Service layer contains business logic.
Handler layer only handles HTTP concerns.
10. Recommended Go Architecture
Handler
   ↓
Service
   ↓
Repository
   ↓
MySQL

Example:

ChatHandler
    ↓
AIService
    ↓
TicketService
    ↓
TicketRepository
    ↓
MySQL


---


# 5. `docs/API.md`


```markdown
# 🌐 API Specification


Base URL:


```text
/api

All API responses use JSON.