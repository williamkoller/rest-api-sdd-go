# API Contracts: School Management System

**Feature**: `001-school-management-system`
**Date**: 2026-04-23
**Base URL**: `/api/v1`
**Auth**: Bearer JWT in `Authorization` header (except `/auth/*` and `GET /health`)

---

## Response Envelope

All responses use a consistent envelope:

```json
{
  "data": { ... } | [ ... ] | null,
  "error": null | { "code": "ERROR_CODE", "message": "Human-readable message" },
  "meta": { "page": 1, "per_page": 20, "total": 150 }
}
```

`meta` is only present on paginated list responses.

---

## Authentication

### POST /auth/login

**Auth**: None

**Request**:
```json
{ "email": "user@school.com", "password": "s3cr3t" }
```

**Response 200**:
```json
{
  "data": {
    "access_token": "<JWT>",
    "refresh_token": "<token>",
    "expires_in": 900,
    "user": { "id": "uuid", "name": "Name", "role": "guardian", "school_id": "uuid" }
  }
}
```

**Errors**: `401 INVALID_CREDENTIALS`, `403 ACCOUNT_INACTIVE`

---

### POST /auth/refresh

**Auth**: None (refresh token in body)

**Request**:
```json
{ "refresh_token": "<token>" }
```

**Response 200**: Same as `/auth/login` response.

**Errors**: `401 TOKEN_EXPIRED`, `401 TOKEN_INVALID`

---

### POST /auth/logout

**Auth**: Bearer JWT

**Request**: Empty body. Invalidates the refresh token from the JWT claims.

**Response 204**: No content.

---

## Health

### GET /health

**Auth**: None

**Response 200**:
```json
{ "data": { "status": "ok", "version": "1.0.0", "db": "ok", "cache": "ok" } }
```

---

## Schools

> Role required: `super_admin`

### GET /schools

**Query**: `?page=1&per_page=20&active=true`

**Response 200**: Paginated list of Schools.

### POST /schools

**Request**:
```json
{ "name": "Escola X", "cnpj": "12345678000199", "email": "contato@x.com", "phone": "11999990000" }
```

**Response 201**: Created School object.

### GET /schools/:id

**Response 200**: School object.

### PUT /schools/:id

**Request**: Partial School fields (name, email, phone, active).

**Response 200**: Updated School object.

---

## Units

> Role required: `school_admin` or `super_admin`

### GET /schools/:school_id/units

**Query**: `?active=true`

**Response 200**: List of Units for the school.

### POST /schools/:school_id/units

**Request**:
```json
{
  "name": "Unidade Centro",
  "address": "Rua das Flores, 100",
  "city": "São Paulo",
  "state": "SP",
  "zip_code": "01234-000",
  "phone": "11988880000"
}
```

**Response 201**: Created Unit object.

### GET /units/:id

**Response 200**: Unit object.

### PUT /units/:id

**Request**: Partial Unit fields.

**Response 200**: Updated Unit object.

---

## Classrooms

> Role: `school_admin`, `unit_staff`

### GET /units/:unit_id/classrooms

**Response 200**: List of Classrooms.

### POST /units/:unit_id/classrooms

**Request**: `{ "code": "A101", "capacity": 35 }`

**Response 201**: Created Classroom.

### PUT /classrooms/:id

**Request**: `{ "code": "A101", "capacity": 40, "active": true }`

**Response 200**: Updated Classroom.

---

## Classes (Turmas)

> Role: `school_admin`, `unit_staff`

### GET /units/:unit_id/classes

**Query**: `?academic_year=2026&active=true`

**Response 200**: List of Classes with teacher count and enrollment count.

### POST /units/:unit_id/classes

**Request**:
```json
{
  "name": "1º Ano A",
  "grade_level": "1",
  "shift": "morning",
  "academic_year": 2026,
  "classroom_id": "uuid"
}
```

**Response 201**: Created Class.

### GET /classes/:id

**Response 200**: Class object with enrolled student count.

### PUT /classes/:id

**Request**: Partial Class fields.

**Response 200**: Updated Class.

---

## Students

> Role: `school_admin`, `unit_staff` (write); `teacher` (read own classes); `guardian` (read own children)

### GET /classes/:class_id/students

**Response 200**: List of Students enrolled in the class.

### POST /classes/:class_id/students

Enrolls an existing student in this class, or creates and enrolls.

**Request**:
```json
{
  "student_id": "uuid",
  "enrolled_at": "2026-02-01"
}
```

Or create + enroll:
```json
{
  "name": "Ana Silva",
  "birth_date": "2018-05-10",
  "cpf": "12345678901",
  "enrolled_at": "2026-02-01"
}
```

**Response 201**: Enrollment object with Student details.

### GET /students/:id

**Response 200**: Student object.

### PUT /students/:id

**Request**: Partial Student fields.

**Response 200**: Updated Student.

### DELETE /classes/:class_id/students/:student_id

Unenrolls a student from this class.

**Response 204**: No content.

---

## Grades & Attendance

> Write role: `teacher` (own classes only). Read: `school_admin`, `guardian` (own children).

### POST /classes/:class_id/grades

Batch upsert grades.

**Request**:
```json
{
  "subject": "Matemática",
  "period": "1B",
  "grades": [
    { "student_id": "uuid", "value": 8.5 },
    { "student_id": "uuid", "value": 7.0 }
  ]
}
```

**Response 200**: `{ "data": { "updated": 2 } }`

### GET /students/:id/grades

**Query**: `?subject=Matemática&period=1B&academic_year=2026`

**Response 200**: List of Grade objects.

### POST /classes/:class_id/attendance

Batch upsert attendance for a date.

**Request**:
```json
{
  "date": "2026-04-23",
  "records": [
    { "student_id": "uuid", "status": "present" },
    { "student_id": "uuid", "status": "absent", "note": "Medical certificate" }
  ]
}
```

**Response 200**: `{ "data": { "recorded": 30 } }`

### GET /students/:id/attendance

**Query**: `?from=2026-02-01&to=2026-04-30`

**Response 200**: List of AttendanceRecord objects plus summary (`{ "rate": 0.92, "total": 50, "present": 46 }`).

---

## Financial

> Role: `school_admin` (write, generate). `guardian` (read own children). `unit_staff` (read).

### POST /invoices/generate

Batch-generates invoices for all enrolled students in a unit for a reference period.

**Request**:
```json
{
  "unit_id": "uuid",
  "academic_year": 2026,
  "reference": "2026-05",
  "due_date": "2026-05-10",
  "amount": 1500.00
}
```

**Response 202**: `{ "data": { "queued": true, "estimated_count": 120 } }`

### GET /students/:id/invoices

**Query**: `?status=pending&year=2026`

**Response 200**: List of Invoice objects.

### GET /invoices/:id

**Response 200**: Invoice object.

### POST /invoices/:id/pay

Records a manual payment.

**Request**:
```json
{ "amount_paid": 1500.00, "method": "pix", "gateway_ref": "txn_123" }
```

**Response 200**: Updated Invoice + Payment object.

### GET /invoices/:id/receipt

**Response 200**: `application/pdf` receipt download.

### GET /schools/:school_id/financial/delinquency

**Query**: `?unit_id=uuid&days_overdue=5`

**Response 200**: List of students with overdue invoices.

---

## Re-Enrollment

> Role: `school_admin` (manage campaigns). `guardian` (submit response).

### POST /reenrollment/campaigns

**Request**:
```json
{
  "unit_id": "uuid",
  "academic_year": 2027,
  "deadline": "2026-11-30T23:59:59Z"
}
```

**Response 201**: ReenrollmentCampaign object.

### GET /reenrollment/campaigns/:id/dashboard

**Response 200**:
```json
{
  "data": {
    "total": 120,
    "confirmed": 85,
    "declined": 12,
    "not_started": 23
  }
}
```

### POST /reenrollment/campaigns/:campaign_id/respond

Guardian submits their response.

**Request**: `{ "student_id": "uuid", "status": "confirmed" }`

**Errors**: `409 OUTSTANDING_DEBT` (if student has overdue invoices), `400 DEADLINE_PASSED`

**Response 200**: Reenrollment object.

---

## School Queue (Waitlist)

> Role: anyone (register). `school_admin`, `unit_staff` (manage).

### POST /units/:unit_id/waitlist

**Auth**: Optional (anonymous registration allowed)

**Request**:
```json
{
  "grade_level": "1",
  "academic_year": 2027,
  "child_name": "Lucas Souza",
  "birth_date": "2019-03-15",
  "guardian_name": "Maria Souza",
  "guardian_email": "maria@email.com",
  "guardian_phone": "11988880000",
  "referral_code": "ABC12345"
}
```

**Response 201**: WaitlistEntry with `position` and `status: "waiting"`.

### GET /units/:unit_id/waitlist

**Query**: `?grade_level=1&academic_year=2027&status=waiting`

**Response 200**: Paginated list of WaitlistEntries.

### PUT /waitlist/:id/status

**Request**: `{ "status": "offer_made" | "accepted" | "declined" | "expired" }`

**Response 200**: Updated WaitlistEntry.

---

## Digital Agenda

> Write: `teacher` (own classes). Read: `guardian` (own children's classes), `school_admin`.

### GET /classes/:class_id/agenda

**Query**: `?type=homework&from=2026-04-01&to=2026-04-30`

**Response 200**: List of AgendaItems sorted by due_date.

### POST /classes/:class_id/agenda

**Request**:
```json
{
  "type": "homework",
  "title": "Exercícios pág. 45-50",
  "description": "Resolver os exercícios do capítulo 3",
  "due_date": "2026-04-25T18:00:00Z"
}
```

**Response 201**: Created AgendaItem.

### PUT /agenda/:id

**Request**: Partial AgendaItem fields.

**Response 200**: Updated AgendaItem.

### DELETE /agenda/:id

**Response 204**: No content.

---

## School Calendar

> Write: `school_admin`. Read: all authenticated users.

### GET /schools/:school_id/calendar

**Query**: `?from=2026-01-01&to=2026-12-31&unit_id=uuid`

**Response 200**: List of CalendarEvents (school-wide + unit-scoped if unit_id given).

### POST /schools/:school_id/calendar

**Request**:
```json
{
  "unit_id": null,
  "title": "Feriado Nacional",
  "type": "holiday",
  "start_date": "2026-09-07",
  "end_date": "2026-09-07"
}
```

**Response 201**: Created CalendarEvent.

### PUT /calendar/:id

**Request**: Partial CalendarEvent fields.

**Response 200**: Updated CalendarEvent.

### DELETE /calendar/:id

**Response 204**: No content.

---

## Support Desk (Atendimento)

> Write (tickets): `guardian`, `school_admin`. Reply: `unit_staff`, `school_admin`. Read: scoped.

### POST /tickets

**Request**:
```json
{
  "unit_id": "uuid",
  "subject": "Dúvida sobre boleto de maio",
  "body": "Gostaria de entender o cálculo...",
  "category": "financial"
}
```

**Response 201**: Ticket + first TicketMessage.

### GET /tickets

**Query**: `?status=open&unit_id=uuid&category=financial&page=1`

**Response 200**: Paginated list of Tickets.

### GET /tickets/:id

**Response 200**: Ticket object with all TicketMessages.

### POST /tickets/:id/reply

**Request**: `{ "body": "Segue a resposta..." }`

**Response 201**: TicketMessage. Ticket status auto-transitions to `in_progress` if staff replied.

### PUT /tickets/:id/status

**Request**: `{ "status": "resolved" | "closed" }` (staff/admin only)

**Response 200**: Updated Ticket.

### GET /schools/:school_id/tickets/report

**Query**: `?from=2026-04-01&to=2026-04-30`

**Response 200**:
```json
{
  "data": {
    "total": 45,
    "open": 8,
    "resolved": 35,
    "avg_resolution_hours": 6.2,
    "by_category": { "financial": 20, "academic": 15, "general": 10 }
  }
}
```

---

## Feed

> Write: `school_admin`, `unit_staff`. Read: all authenticated users.

### GET /schools/:school_id/feed

**Query**: `?unit_id=uuid&page=1&per_page=20`

Returns school-wide posts + unit posts if `unit_id` provided, sorted by `published_at` desc.

**Response 200**: Paginated list of FeedPosts.

### POST /feed

**Request**:
```json
{
  "school_id": "uuid",
  "unit_id": null,
  "body": "Comunicado: Semana Cultural de 20 a 24 de maio!",
  "image_url": "https://cdn.escola.com/img/banner.jpg"
}
```

**Response 201**: Created FeedPost.

### DELETE /feed/:id

**Response 204**: No content.

---

## Refer a Friend

### GET /me/referral-link

Returns the authenticated guardian's referral link.

**Response 200**:
```json
{
  "data": {
    "referral_code": "ABC12345",
    "link": "https://escola.com/cadastro?ref=ABC12345",
    "conversions": { "registered": 2, "enrolled": 1 }
  }
}
```

### GET /schools/:school_id/referrals

**Role**: `school_admin`. Returns referral tracking report.

**Query**: `?status=enrolled&page=1`

**Response 200**: Paginated list of Referrals with referrer name and status.

---

## Menu (Cardápio)

> Write: `school_admin`, `unit_staff`. Read: all authenticated users.

### GET /units/:unit_id/menu

**Query**: `?week_start=2026-04-21`

**Response 200**: Menu object with nested MenuItems grouped by day and meal_type.

### POST /units/:unit_id/menu

**Request**:
```json
{
  "week_start": "2026-04-21",
  "items": [
    { "day_of_week": "monday", "meal_type": "lunch", "description": "Arroz, feijão, frango grelhado, salada" },
    { "day_of_week": "monday", "meal_type": "snack", "description": "Fruta da época" }
  ]
}
```

**Response 201**: Created Menu with all items.

**Errors**: `409 MENU_ALREADY_EXISTS` for the same `(unit_id, week_start)`.

### PUT /menu/:id

Replace all items for an existing menu.

**Request**: `{ "items": [ ... ] }` (same format as POST items)

**Response 200**: Updated Menu.

---

## Curriculum Grid (Grade Curricular)

> Write: `school_admin`. Read: all authenticated users with access to the class.

### GET /classes/:class_id/curriculum

**Response 200**: List of CurriculumEntries for the class, sorted by day_of_week and start_time.

### POST /classes/:class_id/curriculum

Batch-creates or replaces the curriculum for a class.

**Request**:
```json
{
  "entries": [
    {
      "subject": "Matemática",
      "teacher_id": "uuid",
      "day_of_week": "monday",
      "start_time": "07:30",
      "end_time": "08:20"
    },
    {
      "subject": "Português",
      "teacher_id": "uuid",
      "day_of_week": "monday",
      "start_time": "08:20",
      "end_time": "09:10"
    }
  ]
}
```

**Response 201**: `{ "data": { "created": 10 } }`

**Errors**: `409 SCHEDULE_CONFLICT` if time slots overlap for the same class.

### PUT /curriculum/:id

**Request**: Partial CurriculumEntry fields.

**Response 200**: Updated CurriculumEntry.

### DELETE /curriculum/:id

**Response 204**: No content.

---

## Common Error Codes

| Code | HTTP Status | Description |
|------|------------|-------------|
| INVALID_CREDENTIALS | 401 | Wrong email/password |
| TOKEN_EXPIRED | 401 | JWT or refresh token expired |
| TOKEN_INVALID | 401 | Malformed or revoked token |
| FORBIDDEN | 403 | Role not allowed for this action |
| NOT_FOUND | 404 | Resource does not exist |
| CONFLICT | 409 | Duplicate resource or business rule violation |
| OUTSTANDING_DEBT | 409 | Re-enrollment blocked by overdue invoice |
| DEADLINE_PASSED | 400 | Campaign deadline has passed |
| SCHEDULE_CONFLICT | 409 | Curriculum time slot overlap |
| VALIDATION_ERROR | 422 | Request body fails validation |
| INTERNAL_ERROR | 500 | Unexpected server error |
