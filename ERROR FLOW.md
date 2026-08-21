Error Flow

If AI API fails:

Guest
 ↓
AI Request
 ↓
AI unavailable
 ↓
Fallback response

Example:

"Maaf, layanan AI sedang mengalami gangguan. Silakan coba kembali beberapa saat lagi."

The backend must not create a ticket unless the request has been successfully processed and validated.