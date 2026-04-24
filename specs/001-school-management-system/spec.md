# Feature Specification: School Management System

**Feature Branch**: `001-school-management-system`
**Created**: 2026-04-23
**Status**: Draft
**Input**: User description: "crie um projeto, sistema de gerenciamento de escola, com modulos: financeiro, agenda digital, atendimento, rematricula, indique um amigo, feed, fila de escola, notas e frequencia, calendario escolar, cardapio e grade curricular. Sistema precisa ter escola, unidade, sala de aula, turmas etc.."

---

## User Scenarios & Testing *(mandatory)*

### User Story 1 — Grades & Attendance Management (Priority: P1)

A teacher or school administrator logs in and records student grades and attendance for a given
class and period. Parents can view their child's grades and attendance history from their account.

**Why this priority**: Academic performance tracking is the core educational function. It blocks
parent and student engagement features and is required for re-enrollment decisions.

**Independent Test**: Create a school, a unit, a classroom, enroll a student, record grades and
attendance as a teacher, then verify the parent can view that data from a separate account.

**Acceptance Scenarios**:

1. **Given** a teacher authenticated to their class, **When** they submit attendance marks for
   the day, **Then** each student's attendance record is updated and available immediately.
2. **Given** a parent authenticated to their account, **When** they navigate to their child's
   academic summary, **Then** they see grades per subject and attendance rate per period.
3. **Given** a student with failing attendance (<75%), **When** the administrator views the
   attendance report, **Then** the student is flagged with a warning indicator.

---

### User Story 2 — Financial Module (Priority: P2)

The school financial team registers tuition fees, generates monthly invoices for students, and
records payments. Parents can view outstanding charges, payment history, and download receipts.

**Why this priority**: Revenue collection is operationally critical. It also gates re-enrollment
(students with outstanding debts may be blocked from re-enrolling).

**Independent Test**: Create a student, generate an invoice, simulate a payment, then verify the
invoice is marked as paid and appears in the parent's payment history.

**Acceptance Scenarios**:

1. **Given** a student enrolled in a paid plan, **When** the billing cycle runs, **Then** an
   invoice is created with correct amount and due date.
2. **Given** a parent with an outstanding invoice, **When** they view the financial module,
   **Then** they see the amount due, due date, and a way to initiate payment.
3. **Given** a payment is confirmed, **When** the parent accesses their account, **Then** the
   invoice is marked as paid and a downloadable receipt is available.
4. **Given** an invoice is overdue by 5+ days, **When** the financial team reviews the report,
   **Then** the student appears in the delinquency list with days overdue.

---

### User Story 3 — Re-Enrollment (Priority: P3)

At end-of-year, parents receive a re-enrollment notification. They confirm or decline renewal
for the next academic year. The school tracks re-enrollment status per student and unit.

**Why this priority**: Re-enrollment drives next-year capacity planning and revenue projection.
Financial clearance gates this flow.

**Independent Test**: Open a re-enrollment campaign, submit a re-enrollment request as a parent,
verify the student appears as "re-enrolled" in the school's dashboard.

**Acceptance Scenarios**:

1. **Given** an active re-enrollment campaign, **When** a parent logs in, **Then** they see a
   re-enrollment prompt for each eligible child.
2. **Given** a parent confirms re-enrollment, **When** their child has no outstanding invoices,
   **Then** the re-enrollment is recorded and a confirmation is sent.
3. **Given** a parent tries to re-enroll a child with an overdue invoice, **When** they submit,
   **Then** the system blocks re-enrollment and displays the outstanding debt.
4. **Given** the re-enrollment deadline passes, **When** a school administrator views the
   dashboard, **Then** they see confirmed, pending, and declined counts per class.

---

### User Story 4 — School Queue (Enrollment Waitlist) (Priority: P4)

Prospective parents can register their child on a school or unit waitlist. The school manages
the waitlist, contacts candidates when spots open, and converts accepted candidates to students.

**Why this priority**: Manages demand for oversubscribed units and captures potential revenue.

**Independent Test**: Register a child on the waitlist, advance them to "offer made" status,
confirm acceptance, and verify the child is converted to a pending enrollment.

**Acceptance Scenarios**:

1. **Given** a prospective parent, **When** they submit a waitlist registration for a unit and
   grade, **Then** they receive a confirmation and appear in the school's waitlist view.
2. **Given** a spot opens, **When** the school contacts the next candidate on the waitlist,
   **Then** the candidate's status changes to "offer made" and a notification is sent.
3. **Given** an offer is accepted, **When** the parent confirms, **Then** the child moves to
   "enrollment pending" and is removed from the waitlist.

---

### User Story 5 — Digital Agenda & School Calendar (Priority: P5)

Teachers publish events, homework assignments, and reminders to a class agenda. Parents and
students can view the agenda filtered by their child's class. The school publishes the annual
calendar (holidays, exam periods, events) visible to all.

**Why this priority**: Improves communication between school and families; reduces missed
deadlines and confusion about school schedules.

**Independent Test**: A teacher creates an event for a class; a parent of a student in that
class views the agenda and sees the event. The school posts a holiday; all users see it on
the shared calendar.

**Acceptance Scenarios**:

1. **Given** a teacher in a class, **When** they create an agenda item (homework, event, or
   reminder), **Then** it is visible to all parents of students in that class.
2. **Given** a parent, **When** they open the agenda view, **Then** items are listed sorted by
   date and filterable by type (homework, event, reminder).
3. **Given** a school administrator, **When** they publish a school-wide event to the academic
   calendar, **Then** it appears for all users across all units.
4. **Given** an agenda item has a deadline approaching (within 48 hours), **When** the
   responsible parent has not acknowledged it, **Then** a push notification is sent.

---

### User Story 6 — Attendance & Support Desk (Priority: P6)

Parents and students open support tickets (atendimento) for general inquiries, complaints, or
requests. School staff responds within the platform. The school can track open, in-progress, and
resolved tickets per unit.

**Why this priority**: Centralizes parent communication; replaces untracked emails and phone
calls.

**Independent Test**: A parent opens a ticket; a staff member responds; the parent reads the
response and marks the ticket as resolved.

**Acceptance Scenarios**:

1. **Given** an authenticated parent, **When** they submit a support ticket with a subject and
   description, **Then** the ticket appears in the school's support queue.
2. **Given** a staff member views the ticket queue, **When** they reply to a ticket, **Then**
   the parent receives a notification and can read the response.
3. **Given** a ticket is resolved, **When** an administrator views the support report, **Then**
   it shows average resolution time and ticket volume by category.

---

### User Story 7 — School Feed & Refer a Friend (Priority: P7)

The school publishes posts (news, photos, announcements) to a feed visible to all authenticated
users of that school. Parents can refer new families using a unique referral link; referrals are
tracked and optionally rewarded.

**Why this priority**: Improves school-family engagement and drives organic growth.

**Independent Test**: The school posts an announcement; a parent sees it in their feed. A parent
generates a referral link; a new family uses it to register; the referral is attributed.

**Acceptance Scenarios**:

1. **Given** a school administrator, **When** they publish a post with text and optional image,
   **Then** it appears in the feed for all users of that school/unit.
2. **Given** a parent, **When** they share their referral link and a new family registers using
   it, **Then** the referral is recorded and attributed to the referring parent.
3. **Given** a referral converts to an enrolled student, **When** the school reviews referrals,
   **Then** the referring parent is shown as eligible for any configured reward.

---

### User Story 8 — Menu (Cardápio) & Curriculum Grid (Grade Curricular) (Priority: P8)

The school publishes the weekly/monthly canteen menu and the curriculum grid (subjects, schedule,
teacher assignments per class) for parents and students to view.

**Why this priority**: Informational modules with high parent satisfaction impact but no blocking
dependencies.

**Independent Test**: Upload a weekly menu; verify a parent can view it for their child's unit.
Publish a curriculum grid for a class; verify a parent of a student in that class can view it.

**Acceptance Scenarios**:

1. **Given** a school administrator, **When** they publish the weekly menu for a unit, **Then**
   parents of students in that unit can view it filtered by week.
2. **Given** a school administrator, **When** they configure the curriculum grid for a class
   (subjects, days, times, teacher), **Then** parents and students see an accurate timetable.

---

### Edge Cases

- What happens when a parent has children in different units of the same school? They MUST see
  aggregated data filtered per child.
- What happens when a teacher is assigned to multiple classes across units? They see all their
  classes in a unified view.
- What happens when a student is transferred between classes mid-year? Historical grades remain
  in the old class; new records start in the new class.
- What happens when a payment gateway is unavailable? Invoices remain viewable but payments are
  queued and retried without data loss.
- What happens when the re-enrollment deadline is extended? The school can change the deadline;
  parents already declined can be re-invited.

---

## Requirements *(mandatory)*

### Functional Requirements

**Organizational Structure**

- **FR-001**: The system MUST support a hierarchy: School → Unit → Classroom → Class (Turma)
  → Student enrollment.
- **FR-002**: A School MUST be able to have one or more Units (physical campuses).
- **FR-003**: A Unit MUST be able to have one or more Classrooms (Salas de Aula).
- **FR-004**: A Class (Turma) MUST be assigned to a Classroom, a grade level, and an academic
  year, and MUST have one or more enrolled Students.
- **FR-005**: A Teacher MUST be assignable to one or more Classes across one or more Units.

**Roles & Access**

- **FR-006**: The system MUST support at minimum four roles: School Admin, Unit Staff, Teacher,
  and Parent/Guardian.
- **FR-007**: A Parent/Guardian MUST be linked to one or more Students and MUST only access
  data for their linked children.
- **FR-008**: Role-based access control MUST prevent any user from reading or modifying data
  outside their assigned scope (school, unit, class).

**Grades & Attendance**

- **FR-009**: Teachers MUST be able to record attendance (present, absent, justified) per
  student per day.
- **FR-010**: Teachers MUST be able to record grades per student per subject per assessment
  period.
- **FR-011**: The system MUST calculate and display attendance rate and grade averages
  automatically.

**Financial**

- **FR-012**: The system MUST allow school admins to define fee structures (plans, amounts,
  due dates) per unit and grade.
- **FR-013**: The system MUST generate invoices automatically based on the fee structure and
  enrollment.
- **FR-014**: The system MUST record payments and update invoice status (pending, paid,
  overdue).
- **FR-015**: Parents MUST be able to download PDF receipts for paid invoices.

**Re-Enrollment**

- **FR-016**: School admins MUST be able to open a re-enrollment campaign with a deadline per
  unit or school-wide.
- **FR-017**: The system MUST block re-enrollment for students with overdue invoices.
- **FR-018**: Re-enrollment status MUST be trackable per student (not started, in progress,
  confirmed, declined).

**School Queue (Waitlist)**

- **FR-019**: Prospective parents MUST be able to register a child on a waitlist for a
  specific unit and grade.
- **FR-020**: The school MUST be able to manage waitlist order, make offers, and convert
  accepted candidates to pending students.

**Digital Agenda & Calendar**

- **FR-021**: Teachers MUST be able to create agenda items (homework, event, reminder) scoped
  to their class.
- **FR-022**: School admins MUST be able to publish school-wide calendar events visible across
  all units.
- **FR-023**: Agenda items with upcoming deadlines MUST trigger push notifications to
  relevant parents.

**Support Desk**

- **FR-024**: Parents MUST be able to open support tickets with a subject, description, and
  optional category.
- **FR-025**: School staff MUST be able to view, respond to, and close tickets.
- **FR-026**: The system MUST record and report ticket resolution time and volume.

**Feed**

- **FR-027**: School/unit admins MUST be able to publish text posts with optional images to
  the school feed.
- **FR-028**: The feed MUST be scoped: a post published to a unit is only visible to users
  of that unit; a school-wide post is visible to all.

**Refer a Friend**

- **FR-029**: Each authenticated parent MUST have a unique referral link tied to their school.
- **FR-030**: The system MUST track referral registration and conversion (waitlist → enrolled).

**Menu & Curriculum**

- **FR-031**: School admins MUST be able to publish weekly/monthly menus per unit.
- **FR-032**: School admins MUST be able to configure the curriculum grid (subjects, schedule,
  teacher) per class.

### Key Entities

- **School**: Top-level organization. Has a name, CNPJ, contact info.
- **Unit**: Physical campus belonging to a School. Has address, capacity.
- **Classroom** (Sala de Aula): Physical room within a Unit. Has a room code and capacity.
- **Class** (Turma): Academic grouping in a year, assigned to a Classroom. Has grade level,
  shift, and academic year.
- **Student**: Person enrolled in a Class. Linked to one or more Guardians.
- **Guardian** (Parent): User role linked to one or more Students.
- **Teacher**: User role assignable to one or more Classes.
- **Staff/Admin**: User role with school or unit-scoped administrative permissions.
- **Invoice**: Financial document for a Student's tuition cycle.
- **Payment**: Record of a confirmed payment against an Invoice.
- **Ticket**: Support request opened by a Guardian, handled by Staff.
- **AgendaItem**: Event, homework, or reminder attached to a Class.
- **CalendarEvent**: School-wide or unit-wide date event.
- **FeedPost**: Content post published by admin/staff to a school or unit feed.
- **WaitlistEntry**: Prospective student registration for a Unit/grade.
- **Referral**: Link between a referring Guardian and a newly registered family.
- **Menu**: Weekly/monthly canteen menu for a Unit.
- **CurriculumEntry**: Subject-time-teacher assignment within a Class timetable.

---

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A teacher can record attendance for an entire class (30 students) in under
  2 minutes.
- **SC-002**: A parent can view their child's grades and attendance within 3 taps/clicks
  of logging in.
- **SC-003**: Invoice generation for all students in a unit runs to completion within 5 minutes
  regardless of enrollment size.
- **SC-004**: Re-enrollment confirmation rate is trackable per campaign with 100% of responses
  captured (no offline gaps).
- **SC-005**: Support ticket first-response time is visible in reports; school can monitor
  average response time across all open tickets.
- **SC-006**: Waitlist operations (add, advance, convert) are completable by school staff
  without training beyond a 10-minute walkthrough.
- **SC-007**: All modules are accessible via a unified login; no module requires a separate
  user account.
- **SC-008**: The system correctly enforces data isolation: a user of School A cannot access
  any data from School B. [NEEDS CLARIFICATION: Is this a multi-school SaaS platform (multiple
  independent organizations) or a single organization's system?]

---

## Assumptions

- Each Student belongs to exactly one Class (Turma) at a time within an academic year;
  transfers are tracked historically.
- Academic years are annual; the system supports multiple concurrent academic years in
  read-only mode for historical data.
- The financial module records payment intent but delegates actual payment processing to an
  external payment gateway (specific gateway is an implementation detail).
- Push notifications are delivered via mobile app or web push; email is the fallback channel.
- The "Fila de escola" (School Queue) refers to an enrollment waitlist for prospective students,
  not a real-time physical queue. [NEEDS CLARIFICATION: Confirm this interpretation — or is it
  a real-time daily check-in/attendance queue for existing students?]
- The referral reward (if any) is configured by the school; this system only tracks eligibility,
  not payment of rewards.
- Menus are published as structured data (items per day) rather than uploaded PDF documents,
  enabling filtering and mobile display.
- The system supports Brazilian Portuguese as the primary language (LGPD compliance applies to
  personal data storage).
