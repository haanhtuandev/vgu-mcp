# Tool Reference

`vgu-mcp` exposes 8 tools to your AI client. The AI selects and calls them automatically based on your prompt — you never need to call them by name.

---

## Tool Overview

| Tool | What it does |
|---|---|
| [`get_enrolled_courses`](#get_enrolled_courses) | Lists all your enrolled Moodle courses |
| [`get_course_contents`](#get_course_contents) | Fetches sections, modules, and files for a course |
| [`get_upcoming_deadlines`](#get_upcoming_deadlines) | Shows upcoming assignment deadlines |
| [`get_course_grades`](#get_course_grades) | Retrieves your grades for a course |
| [`read_course_announcements`](#read_course_announcements) | Reads the announcement forum for a course |
| [`download_course_material`](#download_course_material) | Downloads a lecture file to your local machine |
| [`extract_course_material_text`](#extract_course_material_text) | Extracts plain text from a PDF — no tools needed |
| [`stage_assignment_draft`](#stage_assignment_draft) | Uploads a file or text to an assignment as a Draft |

---

## `get_enrolled_courses`

Returns all Moodle courses you are currently enrolled in.

**Parameters:** none

**Example prompts:**
- *"What courses am I enrolled in this semester?"*
- *"List all my Moodle courses"*

---

## `get_course_contents`

Fetches the full content tree for a course — sections, modules (assignments, resources, forums), and downloadable files.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `courseid` | number | ✓ | Course ID from `get_enrolled_courses` |

**Example prompts:**
- *"Show me all the materials from my Distributed Systems course"*
- *"What files are in the Operating Systems course?"*
- *"Find the lab 2 specification for SWE 2"*

---

## `get_upcoming_deadlines`

Returns assignment deadlines from your Moodle calendar, sorted by due date.

**Parameters:** none

**Example prompts:**
- *"What assignments do I have due this week?"*
- *"Any deadlines coming up in the next 7 days?"*
- *"When is my next assignment due?"*

---

## `get_course_grades`

Returns grade items for a specific course.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `courseid` | number | ✓ | Course ID from `get_enrolled_courses` |

**Example prompts:**
- *"What are my grades in Operating Systems?"*
- *"Show me my current marks for SWE 2"*

---

## `read_course_announcements`

Reads the latest posts from the announcement/news forum of a course.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `courseid` | number | ✓ | Course ID from `get_enrolled_courses` |

**Example prompts:**
- *"Any new announcements in Distributed Systems?"*
- *"What's the latest news from my OS professor?"*

---

## `download_course_material`

Streams a file from Moodle to a local directory. Uses O(1) memory — the file is never fully buffered.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `file_url` | string | ✓ | The `fileurl` from a module's `contents` array in `get_course_contents` |
| `destination_dir` | string | | Local folder to save to (default: `./downloads`) |

**Example prompts:**
- *"Download the week 3 OS lecture slides"*
- *"Save the lab 2 specification PDF to my Downloads folder"*

---

## `extract_course_material_text`

Downloads a PDF from Moodle (or reads a local file) and extracts its text content using a pure-Go PDF reader — no `pdftotext`, no Python, no system packages required. Returns Markdown with a `## Page N` header before each page.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `file_url` | string | one of the two | `fileurl` from `get_course_contents` — fetches and extracts in one step |
| `local_filepath` | string | one of the two | Path to a locally downloaded PDF |

**Limitations:**
- PDF only. For DOCX, PPTX, or other formats, use `download_course_material` instead.
- Scanned/image-only PDFs have no text layer — extraction will return an error.

**Example prompts:**
- *"Read the week 3 OS lecture PDF and summarise it for me"*
- *"Explain the main concept from the Distributed Systems slides"*
- *"Quiz me on the Lab 2 specification"*
- *"What does page 4 of the lecture say?"*

---

## `stage_assignment_draft`

Uploads a local file (ZIP, PDF, code archive) and/or a text note to a Moodle assignment. The submission is saved as a **Draft** — it does **not** get submitted for grading. You must go to Moodle and click "Submit assignment" yourself to finalise.

**Parameters:**

| Name | Type | Required | Description |
|---|---|---|---|
| `assignment_id` | number | ✓ | The assignment's **instance ID** — use the `instance` field from a module in `get_course_contents` (NOT the `id` field) |
| `review_url` | string | | The assignment's Moodle page URL — use the `url` field from the same module |
| `file_path` | string | one of the two | Local path to the file to upload |
| `text_content` | string | one of the two | Online text or notes to attach |

> **`assignment_id` vs `id`:** Every Moodle assignment has two different numbers. The `id` field (e.g. `34852`) is the course module ID used in browser URLs. The `instance` field (e.g. `2432`) is the assignment record ID required by the submission API. Always use `instance` for `assignment_id`.

**What the AI returns:**

```json
{
  "status": "draft_staged",
  "assignment_id": 2432,
  "filename": "lab2.zip",
  "draft_item_id": 880413555,
  "moodle_review_url": "https://moodle.vgu.edu.vn/mod/assign/view.php?id=34852",
  "message": "File successfully uploaded and saved as a DRAFT. It has NOT been submitted for grading. Visit the review URL and click 'Submit assignment' to finalize."
}
```

**Example prompts:**
- *"Stage my lab2.zip to the component design assignment as a draft"*
- *"Upload solution.zip to the OS assignment for review"*
- *"Add a text note to my assignment draft saying 'Work in progress'"*
