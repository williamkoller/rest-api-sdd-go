# Data Model: School Management System

**Feature**: `001-school-management-system`
**Date**: 2026-04-23

---

## Entity Definitions

### School

Top-level tenant. Every other entity is scoped under a School.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK, generated |
| name | string | required, max 200 chars |
| cnpj | string | required, unique, 14 digits (Brazilian tax ID) |
| email | string | required, valid email |
| phone | string | optional |
| active | bool | default true |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### Unit

Physical campus belonging to a School.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School, required |
| name | string | required, max 200 chars |
| address | string | required |
| city | string | required |
| state | string | required, 2 chars (BR state code) |
| zip_code | string | required |
| phone | string | optional |
| active | bool | default true |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

**Constraints**: `(school_id, name)` unique.

---

### Classroom (Sala de Aula)

Physical room within a Unit.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| unit_id | UUID | FK → Unit, required |
| code | string | required, max 20 chars (e.g. "A101") |
| capacity | int | required, min 1 |
| active | bool | default true |

**Constraints**: `(unit_id, code)` unique.

---

### Class (Turma)

Academic grouping: a cohort of students for one academic year.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| unit_id | UUID | FK → Unit, required |
| classroom_id | UUID | FK → Classroom, optional (can change) |
| name | string | required, max 100 chars (e.g. "1º Ano A") |
| grade_level | string | required (e.g. "1", "2", "EF1", "EM3") |
| shift | enum | required: morning, afternoon, full |
| academic_year | int | required (e.g. 2026) |
| active | bool | default true |
| created_at | timestamp | auto |

**Constraints**: `(unit_id, name, academic_year)` unique.

---

### User

Single user table; role determines permissions.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School (null for super-admin) |
| name | string | required |
| email | string | required, unique |
| password_hash | string | required |
| role | enum | required: super_admin, school_admin, unit_staff, teacher, guardian |
| active | bool | default true |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

**Note**: A teacher may be scoped to one or more Units via the `TeacherUnit` join table.

---

### TeacherClass

Many-to-many: Teacher ↔ Class.

| Field | Type | Rules |
|-------|------|-------|
| teacher_id | UUID | FK → User (role=teacher) |
| class_id | UUID | FK → Class |
| subject | string | required (the subject this teacher teaches in this class) |

**PK**: `(teacher_id, class_id, subject)`.

---

### Student

Person enrolled in the school system.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| name | string | required |
| birth_date | date | required |
| cpf | string | optional, 11 digits |
| registration_number | string | required, unique per school |
| active | bool | default true |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### Enrollment

Student ↔ Class association per academic year.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| student_id | UUID | FK → Student |
| class_id | UUID | FK → Class |
| enrolled_at | timestamp | required |
| unenrolled_at | timestamp | null = currently active |
| status | enum | active, transferred, unenrolled |

**Constraints**: A student has at most one active enrollment per academic year.

---

### GuardianStudent

Guardian (parent) ↔ Student relationship.

| Field | Type | Rules |
|-------|------|-------|
| guardian_id | UUID | FK → User (role=guardian) |
| student_id | UUID | FK → Student |
| relationship | string | optional (e.g. "mother", "father", "guardian") |

**PK**: `(guardian_id, student_id)`.

---

### AttendanceRecord

Daily attendance mark per student.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| enrollment_id | UUID | FK → Enrollment |
| date | date | required |
| status | enum | required: present, absent, justified |
| note | string | optional |
| recorded_by | UUID | FK → User (teacher/staff) |
| created_at | timestamp | auto |

**Constraints**: `(enrollment_id, date)` unique.

---

### Grade (Nota)

Academic grade per student per subject per assessment period.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| enrollment_id | UUID | FK → Enrollment |
| subject | string | required |
| period | string | required (e.g. "1B", "2B", "FINAL") |
| value | decimal(5,2) | required, 0.00–10.00 |
| recorded_by | UUID | FK → User |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

**Constraints**: `(enrollment_id, subject, period)` unique.

---

### Invoice

Financial charge for a student's tuition cycle.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| student_id | UUID | FK → Student |
| school_id | UUID | FK → School |
| amount | decimal(10,2) | required, > 0 |
| due_date | date | required |
| reference | string | required (e.g. "2026-05") |
| status | enum | pending, paid, overdue, cancelled |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

**Constraints**: `(student_id, reference)` unique.

---

### Payment

Record of a confirmed payment against an Invoice.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| invoice_id | UUID | FK → Invoice |
| amount_paid | decimal(10,2) | required |
| paid_at | timestamp | required |
| method | enum | pix, boleto, credit_card, manual |
| gateway_ref | string | optional (external transaction ID) |
| created_at | timestamp | auto |

---

### Reenrollment

Re-enrollment response per student per campaign.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| student_id | UUID | FK → Student |
| campaign_id | UUID | FK → ReenrollmentCampaign |
| status | enum | not_started, confirmed, declined |
| responded_at | timestamp | null until response |
| created_at | timestamp | auto |

---

### ReenrollmentCampaign

School-wide or unit-scoped re-enrollment campaign.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| unit_id | UUID | FK → Unit, null = school-wide |
| academic_year | int | required |
| deadline | timestamp | required |
| status | enum | open, closed |
| created_at | timestamp | auto |

---

### WaitlistEntry

Prospective student registration for a Unit/grade.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| unit_id | UUID | FK → Unit |
| grade_level | string | required |
| academic_year | int | required |
| child_name | string | required |
| birth_date | date | required |
| guardian_name | string | required |
| guardian_email | string | required |
| guardian_phone | string | optional |
| status | enum | waiting, offer_made, accepted, declined, expired |
| position | int | required (queue position within unit+grade+year) |
| referral_id | UUID | FK → Referral, optional |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### Referral

Link between a referring Guardian and a newly registered family.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| referrer_id | UUID | FK → User (guardian) |
| referral_code | string | unique, 8 chars |
| referred_email | string | optional (set when lead registers) |
| status | enum | pending, registered, enrolled |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### Ticket (Atendimento)

Support request opened by a Guardian or Student.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| unit_id | UUID | FK → Unit, optional |
| requester_id | UUID | FK → User |
| subject | string | required, max 200 chars |
| category | enum | general, financial, academic, administrative |
| status | enum | open, in_progress, resolved, closed |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |
| resolved_at | timestamp | null until resolved |

---

### TicketMessage

Messages within a Ticket conversation.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| ticket_id | UUID | FK → Ticket |
| sender_id | UUID | FK → User |
| body | text | required |
| created_at | timestamp | auto |

---

### AgendaItem

Event, homework, or reminder attached to a Class.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| class_id | UUID | FK → Class |
| created_by | UUID | FK → User (teacher) |
| type | enum | homework, event, reminder |
| title | string | required, max 200 chars |
| description | text | optional |
| due_date | timestamp | required |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### CalendarEvent

School-wide or unit-wide date event.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| unit_id | UUID | FK → Unit, null = school-wide |
| title | string | required |
| description | text | optional |
| type | enum | holiday, exam_period, event, recess |
| start_date | date | required |
| end_date | date | required, >= start_date |
| created_by | UUID | FK → User |
| created_at | timestamp | auto |

---

### FeedPost

Content post published by admin/staff.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| school_id | UUID | FK → School |
| unit_id | UUID | FK → Unit, null = school-wide |
| author_id | UUID | FK → User |
| body | text | required |
| image_url | string | optional |
| published_at | timestamp | auto |
| updated_at | timestamp | auto |

---

### Menu

Weekly/monthly canteen menu for a Unit.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| unit_id | UUID | FK → Unit |
| week_start | date | required (Monday of the week) |
| created_by | UUID | FK → User |
| created_at | timestamp | auto |
| updated_at | timestamp | auto |

**Constraints**: `(unit_id, week_start)` unique.

---

### MenuItem

Individual meal entry within a Menu (per day, per meal type).

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| menu_id | UUID | FK → Menu |
| day_of_week | enum | monday–friday |
| meal_type | enum | breakfast, lunch, snack, dinner |
| description | text | required |

---

### CurriculumEntry

Subject-time-teacher assignment within a Class timetable.

| Field | Type | Rules |
|-------|------|-------|
| id | UUID | PK |
| class_id | UUID | FK → Class |
| subject | string | required |
| teacher_id | UUID | FK → User (teacher) |
| day_of_week | enum | monday–friday |
| start_time | time | required |
| end_time | time | required, > start_time |

**Constraints**: `(class_id, day_of_week, start_time)` unique (no overlapping slots per class).

---

## Dev Tools Config Schemas (branch `002-makefile-fmt-lint-vulnerability-lefthook`)

### `.golangci.yml`

```yaml
linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - staticcheck
    - unused
  disable-all: false

linters-settings:
  errcheck:
    check-blank: true

issues:
  max-same-issues: 3
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

### `lefthook.yml`

```yaml
commit-msg:
  commands:
    validate-message:
      run: |
        MSG=$(cat "{1}")
        PATTERN='^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\([a-zA-Z0-9/_-]+\))?: .+'
        if ! echo "$MSG" | grep -qE "$PATTERN"; then
          echo "ERROR: Invalid commit message."
          echo "Expected: type(scope)?: description"
          echo "Types: feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert"
          exit 1
        fi
```

### Makefile additions

| Target | Command | Purpose |
|--------|---------|---------|
| `fmt` | `go fmt ./...` | Format all Go source files |
| `lint` | `golangci-lint run ./...` | Run configured linters |
| `vulnerability` | `govulncheck ./...` | Scan for known vulnerabilities |
| `tools` | install commands | Install golangci-lint, govulncheck, lefthook |
| `hooks` | `lefthook install` | Activate git hooks from lefthook.yml |

---

## State Transitions

### Invoice

```
pending → paid       (payment confirmed)
pending → overdue    (past due_date, scheduled job)
overdue → paid       (late payment confirmed)
pending → cancelled  (admin action)
overdue → cancelled  (admin action)
```

### Ticket

```
open → in_progress  (staff responds)
in_progress → open  (awaiting requester reply)
in_progress → resolved (staff resolves)
resolved → closed   (auto-close after 7 days or admin action)
```

### WaitlistEntry

```
waiting → offer_made   (school makes offer)
offer_made → accepted  (guardian accepts)
offer_made → declined  (guardian declines or deadline passes)
accepted → [converts to Enrollment]
```

### Reenrollment

```
not_started → confirmed (guardian confirms)
not_started → declined  (guardian declines)
```

---

## Entity Relationships Summary

```text
School
├── has many Units
│   ├── has many Classrooms
│   ├── has many Classes (Turmas)
│   │   ├── has many Enrollments → Students
│   │   ├── has many AttendanceRecords (via Enrollment)
│   │   ├── has many Grades (via Enrollment)
│   │   ├── has many AgendaItems
│   │   ├── has many CurriculumEntries
│   │   └── has many TeacherClass
│   ├── has many WaitlistEntries
│   ├── has many Menus
│   └── has many (scoped) CalendarEvents
├── has many Users (role-scoped)
├── has many Invoices → Payments
├── has many ReenrollmentCampaigns → Reenrollments
├── has many Tickets → TicketMessages
├── has many FeedPosts
├── has many CalendarEvents (school-wide)
└── has many Referrals
```
