🤖 Development Rule for AI Coding Agents

Before implementing any feature:

Read README.md.
Read the relevant documentation in docs/.
Follow existing architecture.
Do not introduce unnecessary dependencies.
Do not create features outside the current MVP without explicit instruction.
Do not expose secrets.
Validate AI-generated structured data before performing database operations.
Preserve backward compatibility when modifying existing APIs.
Prefer simple, maintainable implementations over unnecessary complexity.


---


# 2. `docs/AI_SPEC.md`


```markdown
# 🤖 AI Specification


## 1. Purpose


The AI is the intelligence layer of BATIQA AI Guest Assistant.


Its primary responsibilities are:


1. Understand guest messages.
2. Detect intent.
3. Extract entities.
4. Determine whether an operational action is required.
5. Route operational requests.
6. Determine priority.
7. Generate natural responses.


The AI is NOT allowed to directly modify the database.


All AI-generated actions must pass through backend validation.


---


# 2. AI Processing Pipeline


```text
Guest Message
     ↓
Language Detection
     ↓
Intent Classification
     ↓
Entity Extraction
     ↓
Context Resolution
     ↓
Action Decision
     ↓
Backend Validation
     ↓
Database Operation
     ↓
AI Response