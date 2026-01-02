# Gateway API Documentation

Base URL: `http://localhost:3000`

All requests and responses use JSON format with `Content-Type: application/json`.

## Response Format

### Success with Data
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "createdAt": "2026-01-01T10:00:00Z"
}
```

### Success with Message
```json
{
  "message": "User registered successfully"
}
```

### Error Response
```json
{
  "error": "Invalid request payload"
}
```

## Authentication

Most endpoints require authentication via a Bearer token obtained from the login endpoint.

**Header Format:**
```
Authorization: Bearer <your_jwt_token>
```

🔓 = Public endpoint (no authentication required)  
🔒 = Protected endpoint (requires authentication)

---

## Authentication Endpoints

### 🔓 Register User

Create a new user account.

**Endpoint:** `POST /auth/register`

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Validations:**
- `username`: required, 3-30 characters
- `email`: required, valid email format
- `password`: required, minimum 8 characters

**Success Response:** `201 Created`
```json
{
  "message": "User registered successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or validation errors
- `409 Conflict`: User with this email already exists
- `500 Internal Server Error`: Server error

---

### 🔓 Login

Authenticate and receive a JWT token.

**Endpoint:** `POST /auth/login`

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Validations:**
- `email`: required
- `password`: required

**Success Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload
- `401 Unauthorized`: Invalid email or password
- `500 Internal Server Error`: Server error

---

### 🔒 Get User Profile

Get the authenticated user's profile information.

**Endpoint:** `GET /auth/user`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response:** `200 OK`
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "createdAt": "2026-01-01T10:00:00Z"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

### 🔒 Update User Profile

Update the authenticated user's profile.

**Endpoint:** `PUT /auth/user`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:** (all fields are optional)
```json
{
  "username": "new_username",
  "email": "newemail@example.com",
  "password": "newPassword123"
}
```

**Validations:**
- `username`: optional, 3-30 characters
- `email`: optional, valid email format
- `password`: optional, minimum 8 characters

**Success Response:** `200 OK`
```json
{
  "message": "User updated successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or no fields to update
- `401 Unauthorized`: Missing or invalid token
- `409 Conflict`: Email or username already in use
- `500 Internal Server Error`: Server error

---

### 🔒 Delete User Account

Delete the authenticated user's account permanently.

**Endpoint:** `DELETE /auth/user`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response:** `200 OK`
```json
{
  "message": "User deleted successfully"
}
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

## Mood Endpoints

### 🔒 Get Mood Types

Retrieve all available mood types.

**Endpoint:** `GET /mood/types`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response:** `200 OK`
```json
[
  {
    "id": 1,
    "name": "Happy",
    "description": "Feeling joyful, content, and positive about the day"
  },
  {
    "id": 2,
    "name": "Sad",
    "description": "Feeling down, melancholic, or experiencing a sense of loss"
  }
]
```

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

### 🔒 Add Mood Entry

Create a new mood entry for the authenticated user.

**Endpoint:** `POST /mood`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "moodTypeId": 1,
  "note": "Had a great day at work!",
  "date": "2026-01-02"
}
```

**Validations:**
- `moodTypeId`: required
- `note`: optional, maximum 500 characters
- `date`: required, format `YYYY-MM-DD`

**Success Response:** `201 Created`
```json
{
  "message": "Mood entry created"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or validation errors
- `401 Unauthorized`: Missing or invalid token
- `409 Conflict`: Mood entry for this date already exists
- `500 Internal Server Error`: Server error

---

### 🔒 Get Mood Entries

Retrieve mood entries for the authenticated user within a date range.

**Endpoint:** `GET /mood?from=2026-01-01&to=2026-01-31`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `from`: required, format `YYYY-MM-DD`
- `to`: required, format `YYYY-MM-DD`

**Success Response:** `200 OK`
```json
[
  {
    "id": 1,
    "userId": 1,
    "moodDate": "2026-01-01",
    "moodTypeId": 1,
    "note": "Great start to the year!",
    "createdAt": "2026-01-01T08:30:00Z"
  },
  {
    "id": 2,
    "userId": 1,
    "moodDate": "2026-01-02",
    "moodTypeId": 4,
    "note": "Feeling calm and relaxed",
    "createdAt": "2026-01-02T09:15:00Z"
  }
]
```

**Error Responses:**
- `400 Bad Request`: Invalid or missing query parameters
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

### 🔒 Get Single Mood Entry

Retrieve a specific mood entry by ID (must belong to authenticated user).

**Endpoint:** `GET /mood/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Mood entry ID

**Success Response:** `200 OK`
```json
{
  "id": 1,
  "userId": 1,
  "moodDate": "2026-01-01",
  "moodTypeId": 1,
  "note": "Great start to the year!",
  "createdAt": "2026-01-01T08:30:00Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID parameter
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Mood entry does not belong to user
- `404 Not Found`: Mood entry not found
- `500 Internal Server Error`: Server error

---

### 🔒 Get Mood Summary

Retrieve a statistical summary of mood entries for the authenticated user within a date range.

**Endpoint:** `GET /mood/summary?from=2026-01-01&to=2026-01-31`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `from`: required, format `YYYY-MM-DD`
- `to`: required, format `YYYY-MM-DD`

**Success Response:** `200 OK`
```json
[
  {
    "moodTypeId": 1,
    "count": 15,
    "percentage": 48.39
  },
  {
    "moodTypeId": 4,
    "count": 10,
    "percentage": 32.26
  },
  {
    "moodTypeId": 2,
    "count": 6,
    "percentage": 19.35
  }
]
```

**Notes:**
- Results are ordered by count (descending)
- Percentages are rounded to 2 decimal places

**Error Responses:**
- `400 Bad Request`: Invalid or missing query parameters
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error

---

### 🔒 Update Mood Entry

Update an existing mood entry (must belong to authenticated user).

**Endpoint:** `PUT /mood`

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "id": 1,
  "moodTypeId": 2,
  "note": "Updated note about my mood"
}
```

**Validations:**
- `id`: required
- `moodTypeId`: required
- `note`: required, maximum 500 characters

**Success Response:** `200 OK`
```json
{
  "message": "Mood entry updated"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request payload or validation errors
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Mood entry does not belong to user
- `404 Not Found`: Mood entry not found
- `500 Internal Server Error`: Server error

---

### 🔒 Delete Mood Entry

Delete a mood entry by ID (must belong to authenticated user).

**Endpoint:** `DELETE /mood/{id}`

**Headers:**
```
Authorization: Bearer <token>
```

**Path Parameters:**
- `id`: Mood entry ID

**Success Response:** `200 OK`
```json
{
  "message": "Mood entry deleted"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid ID parameter
- `401 Unauthorized`: Missing or invalid token
- `403 Forbidden`: Mood entry does not belong to user
- `404 Not Found`: Mood entry not found
- `500 Internal Server Error`: Server error

---

## Advice Endpoints

### 🔒 Get Advice

Get personalized advice for the authenticated user based on their mood patterns within a date range.

**Endpoint:** `GET /advice?from=2026-01-01&to=2026-01-31`

**Headers:**
```
Authorization: Bearer <token>
```

**Query Parameters:**
- `from`: required, format `YYYY-MM-DD`
- `to`: required, format `YYYY-MM-DD`

**Success Response:** `200 OK`
```json
{
  "adviceId": 42,
  "title": "Start your day with a clear goal",
  "content": "Identify one main goal for today and focus on achieving it. This gives you direction and purpose."
}
```

**How It Works:**
1. If advice already exists for the specified period, it's returned immediately
2. Otherwise, the system:
   - Analyzes your mood summary for the period
   - Selects appropriate advice based on your mood patterns
   - Saves the advice-period association
   - Returns the selected advice

**Error Responses:**
- `400 Bad Request`: Invalid or missing query parameters
- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: No advice found for the given period (no mood entries or no matching advice)
- `500 Internal Server Error`: Server error

---

## Quote Endpoints

### 🔒 Get Today's Quote

Retrieve the daily motivational quote (cached for 24 hours).

**Endpoint:** `GET /quote/today`

**Headers:**
```
Authorization: Bearer <token>
```

**Success Response:** `200 OK`
```json
{
  "quote": "The only way to do great work is to love what you do.",
  "author": "Steve Jobs",
  "attribution": "Quotes provided by https://zenquotes.io/"
}
```

**Notes:**
- Quotes are cached in Redis for 24 hours
- The same quote is returned for all users on the same day
- Quote provider: [ZenQuotes API](https://zenquotes.io/)

**Error Responses:**
- `401 Unauthorized`: Missing or invalid token
- `500 Internal Server Error`: Server error or external API failure

---

## Status Codes Summary

| Code | Description |
|------|-------------|
| `200 OK` | Request successful, data returned |
| `201 Created` | Resource created successfully |
| `400 Bad Request` | Invalid request payload or parameters |
| `401 Unauthorized` | Missing or invalid authentication token |
| `403 Forbidden` | Authenticated but not authorized to access resource |
| `404 Not Found` | Resource not found |
| `409 Conflict` | Resource conflict (e.g., duplicate entry) |
| `500 Internal Server Error` | Server or service error |

---

## Common Patterns

### Date Range Queries

Many endpoints accept `from` and `to` query parameters:
- Format: `YYYY-MM-DD`
- Both parameters are required
- `from` should be less than or equal to `to`

Example:
```
GET /mood?from=2026-01-01&to=2026-01-31
```

### Authentication Flow

1. Register: `POST /auth/register`
2. Login: `POST /auth/login` → receive token
3. Use token in subsequent requests: `Authorization: Bearer <token>`
4. Token expires after a configured period (requires re-login)

### Pagination

Currently, the API does not implement pagination. Consider limiting your date ranges for optimal performance when querying mood entries.

---

## Error Handling

All errors follow a consistent format:

```json
{
  "error": "Descriptive error message"
}
```

Common validation errors:
- Missing required fields
- Invalid data types
- Invalid date formats
- Out of range values
- Constraint violations (e.g., duplicate entries)
