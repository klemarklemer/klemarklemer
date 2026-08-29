# 04 - Error Handling & Response Envelopes

## 1. Standard Response Envelope

The Taskmaster Claims Operations API wraps all HTTP responses in a standard Candi JSON envelope:

```json
{
  "code": 200,
  "message": "Success",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "totalRecords": 1,
    "totalPages": 1
  }
}
```

---

## 2. HTTP Status Code Matrix

| Status Code | Reason Phrase | Typical Trigger Condition |
|---|---|---|
| `200 OK` | Success | Successful resource retrieval, update, or approval execution |
| `201 Created` | Created | Successful claim submission |
| `400 Bad Request` | Validation Error | Missing required JSON body fields, malformed types, or invalid stage transition attempt |
| `401 Unauthorized` | Unauthorized | Missing or expired JWT Bearer token |
| `403 Forbidden` | Forbidden | Autonomous system token attempting human approval endpoint (`/approval`) |
| `404 Not Found` | Not Found | Target claim ID, policy ID, or officer ID does not exist in the database |
| `422 Unprocessable Entity` | Precondition Failed | Document upload rejected due to unrecognized file format or corrupted metadata |
| `500 Internal Server Error` | Server Error | Database transaction timeout or infrastructure failure |

---

## 3. Error Payload Example

```json
{
  "code": 400,
  "message": "Failed to submit human approval: claim is in stage DOCUMENT_VERIFICATION, human approval requires stage DECISION",
  "error": "Bad Request"
}
```
