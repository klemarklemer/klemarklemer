# 03 - Officers & Policies API Reference

This document covers supporting endpoints for managing Claims Officers, capacity allocations, and Policy contract verification.

---

## 1. Claims Officers API

### 1.1 List Claims Officers (`GET /v1/officer`)
Retrieves all claims officers with their current live workload count, motor skill ratings, and availability status.

- **Query Parameters:**
  - `page` *(int, default: 1)*
  - `limit` *(int, default: 10)*

```bash
curl -X GET "http://localhost:8000/v1/officer"
```

#### Sample Response:
```json
{
  "code": 200,
  "message": "Success",
  "data": [
    {
      "id": 1,
      "name": "Alex Rivera",
      "email": "alex.rivera@klemarklemer.com",
      "role": "Senior Motor Claims Officer",
      "current_workload": 4,
      "motor_skill_rating": 4.80,
      "is_available": true
    },
    {
      "id": 2,
      "name": "David Chen",
      "email": "david.chen@klemarklemer.com",
      "role": "Lead Claims Specialist",
      "current_workload": 8,
      "motor_skill_rating": 4.90,
      "is_available": true
    },
    {
      "id": 3,
      "name": "Elena Rostova",
      "email": "elena.rostova@klemarklemer.com",
      "role": "Junior Claims Officer",
      "current_workload": 2,
      "motor_skill_rating": 4.20,
      "is_available": true
    }
  ]
}
```

---

## 2. Policies API

### 2.1 Get Policy Detail (`GET /v1/policy/:id`)
Fetches policyholder coverage details, max coverage limits, and deductible rules.

```bash
curl -X GET "http://localhost:8000/v1/policy/1"
```

#### Sample Response:
```json
{
  "code": 200,
  "message": "Success",
  "data": {
    "id": 1,
    "policy_number": "POL-MOTOR-2026-8819",
    "policy_holder_name": "Jordan Hayes",
    "vehicle_plate": "B 1284 UQ",
    "vehicle_model": "2024 Hyundai Ioniq 5 EV",
    "coverage_type": "COMPREHENSIVE",
    "max_coverage_amount": 45000.00,
    "deductible_amount": 500.00,
    "effective_date": "2026-01-01T00:00:00Z",
    "expiry_date": "2027-01-01T00:00:00Z",
    "status": "ACTIVE"
  }
}
```
