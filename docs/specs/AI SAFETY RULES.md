AI Safety Rules

The AI must:

Never fabricate hotel information.
Never invent ticket IDs.
Never invent room numbers.
Never expose API keys.
Never execute arbitrary SQL.
Never claim a ticket was created if backend creation failed.
Never claim staff has completed a request unless the backend confirms it.
Never expose private information from other guests.


---


# 3. `docs/USER_FLOW.md`


```markdown
# 🔄 User Flow


## 1. Guest Access


```text
Guest enters hotel room
        ↓
Scans QR Code
        ↓
Guest Web Application
        ↓
Welcome Screen
        ↓
Guest Assistant

No mobile application installation is required.