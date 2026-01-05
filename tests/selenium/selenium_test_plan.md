# Selenium Test Plan for Pregnanzy

## Overview

This document outlines the comprehensive automated testing strategy. The goal is to verify all critical user journeys (CUJs), edge cases, and role-specific workflows.

## Prerequisite

- Python 3.x
- Selenium WebDriver
- Chrome/Edge Driver

## Test Structure

```
tests/
  selenium/
    config.py           # Configuration (Base URL, Test Credentials)
    driver_setup.py     # WebDriver Factory
    tests/
      test_auth_register.py
      test_doctor_patient_management.py
      test_doctor_medical_records.py
      test_doctor_appointments.py
      test_patient_dashboard.py
      test_patient_records.py
      test_profile_management.py
```

## Detailed Test Cases (30+ Scenarios)

### 1. Authentication & Registration (test_auth_register.py)

- **TC01: Register Success** -> Register new user, verify redirect to /login.
- **TC02: Register Duplicate Username** -> Try existing username, verify error alert.
- **TC03: Register Validation** -> Submit empty form, verify HTML5 validation/warnings.
- **TC04: Doctor Login Success** -> Login as Doctor, verify /doctor/dashboard.
- **TC05: Patient Login Success** -> Login as Patient, verify /.
- **TC06: Invalid Password** -> Login with wrong password, verify failure alert.
- **TC07: Invalid Username** -> Login with non-existent user, verify failure alert.
- **TC08: Logout** -> Click Logout, verify redirect to /login.

### 2. Doctor: Patient Management (test_doctor_patient_management.py)

- **TC09: Dashboard List Load** -> Login as Doctor, verify list is not empty.
- **TC10: Search Patient** -> Enter "Mommy", verify filtered list contains "Mommy".
- **TC11: View Patient Detail** -> Click patient, verify URL /doctor/patient/:id.
- **TC12: Create New Pregnancy** -> Fill LMP, verify EDC auto-calc, Save.
- **TC13: End Pregnancy** -> Open Modal, Fill Delivery Data, Submit. Verify status "Finished".
- **TC14: End Pregnancy Validation** -> Try submitting empty modal, verify validation.

### 3. Doctor: Medical Records (test_doctor_medical_records.py)

- **TC15: Edit Medical History (Preeclampsia)** -> Select "Yes", Fill text, Save. Verify persistence.
- **TC16: Edit Medical History (Toggle Off)** -> Select "No", Save. Verify persistence.
- **TC17: Vaccination - Add 1 Dose** -> Add Tetanus Dose 1, Save. Reload, verify 1 row.
- **TC18: Vaccination - Add Multiple Doses** -> Add 3 Doses, Save. Reload, verify 3 rows (Card/Table).
- **TC19: Vaccination - Delete Dose** -> Click Delete on row, Save. Verify removal.
- **TC20: Lab Result - Upload** -> Fill Form, Select PDF, Save. Verify entry in table working.
- **TC21: Antenatal Visit - Create** -> Add Weight/Pressure, Save. Verify entry in history list.
- **TC22: Previous Pregnancy - Add** -> Add new record, Save. Verify presence.

### 4. Doctor: Appointments (test_doctor_appointments.py)

- **TC23: Create Appointment** -> Fill Date/Time/Title. Verify Location defaults to "แผนกสูตินารีเวช อาคาร 2". Save.
- **TC24: Appointment Validation** -> Try save without Date, verify alert.

### 5. Patient: Dashboard Features (test_patient_dashboard.py)

- **TC25: Kick Count (Normal)** -> Add 10 kicks, Verify Blue Card + "Normal" text.
- **TC26: Kick Count (Danger)** -> Add 5 kicks, Verify Red Card + "Abnormal" text.
- **TC27: Kick Count (Daily Reset)** -> (Requires mock date if possible, otherwise verify Today's date display).
- **TC28: Baby Growth Card** -> Verify visibility and white card styling.
- **TC29: Check Appointments** -> Verify newly created appointment appears on dashboard.

### 6. Patient: Records & Profile (test_patient_records.py, test_profile_management.py)

- **TC30: Vaccination Read-Only** -> Go to /vaccine. Verify Card layout. Verify no input fields/buttons.
- **TC31: Doctor Profile View** -> Go to Profile. Verify inputs disabled (View Mode).
- **TC32: Doctor Profile Edit** -> Click Edit. Change Phone. Save. Verify Value.
- **TC33: Patient Husband Info** -> Edit Husband data, Save. Reload & Verify.

## Execution Plan

1. **Setup**: `driver_setup.py` handling Login helpers.
2. **Helpers**: `wait_for_element`, `fill_form` utilities.
3. **Run**: `pytest` or `python -m unittest discover tests/selenium`.
