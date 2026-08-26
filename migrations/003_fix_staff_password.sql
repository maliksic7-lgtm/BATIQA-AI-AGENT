-- 003_fix_staff_password.sql
-- Replaces invalid bcrypt placeholder hashes from earlier seeds with the real
-- bcrypt hash of the demo password "batiqa123" (cost 10).
-- Idempotent: safe to run multiple times. DEMO ONLY - rotate before production.
INSERT INTO staff (name, email, password_hash, department) VALUES
('Admin BATIQA', 'admin@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'ADMIN'),
('Housekeeping Team', 'hk@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'HOUSEKEEPING'),
('Engineering Team', 'eng@batiqa.com', '$2a$10$WnP6TcAxzfFXrLBRKgoLWOgwNLEf7x9iGkuhmSO34kfBqyhsTRHCa', 'ENGINEERING')
ON DUPLICATE KEY UPDATE password_hash = VALUES(password_hash);
