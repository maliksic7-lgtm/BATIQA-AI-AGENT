API Rules
API keys must never be returned.
AI provider credentials stay on backend.
Validate all request bodies.
Validate AI-generated structured output.
Do not trust frontend-provided department or priority without backend validation.
Staff endpoints require authentication.
Guest sessions must be isolated.


---


# 6. `docs/UI_SPEC.md`


```markdown
# 🎨 UI/UX Specification


## 1. Design Direction


The application should feel like a modern hotel digital service, not a generic AI chatbot.


Design principles:


- Premium
- Clean
- Minimal
- Trustworthy
- Hospitality-focused
- Mobile-first
- Fast
- Easy to understand


Avoid:


- Excessive gradients.
- Overly futuristic AI visuals.
- Excessive glassmorphism.
- Too many animations.
- Cluttered dashboards.
- Generic "AI generated" appearance.


---


# 2. Guest Application


Primary target:


```text
Mobile browser

The guest should be able to use the application comfortably with one hand.