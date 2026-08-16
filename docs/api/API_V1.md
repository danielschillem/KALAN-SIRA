# API REST v1

Base path : `/api/v1`

## Auth

```text
POST /auth/login
POST /auth/refresh
POST /auth/logout
POST /parent/auth/request-otp
POST /parent/auth/verify-otp
```

## Etablissements

```text
POST   /schools
GET    /schools/{schoolId}
PATCH  /schools/{schoolId}
GET    /schools/{schoolId}/dashboard
```

## Annees, niveaux, classes

```text
POST /schools/{schoolId}/years
GET  /schools/{schoolId}/years
POST /schools/{schoolId}/levels
GET  /schools/{schoolId}/levels
POST /schools/{schoolId}/classes
GET  /schools/{schoolId}/classes
```

## Eleves

```text
POST /schools/{schoolId}/students
GET  /schools/{schoolId}/students
GET  /schools/{schoolId}/students/{studentId}
PATCH /schools/{schoolId}/students/{studentId}
POST /schools/{schoolId}/students/{studentId}/guardians
```

## Inscriptions

```text
POST /schools/{schoolId}/enrollments
GET  /schools/{schoolId}/enrollments/{enrollmentId}
POST /schools/{schoolId}/enrollments/{enrollmentId}/activate
POST /schools/{schoolId}/enrollments/{enrollmentId}/cancel
```

`activate` declenche transactionnellement la generation des charges.

## Tarification

```text
POST /schools/{schoolId}/fee-schedules
GET  /schools/{schoolId}/fee-schedules
GET  /schools/{schoolId}/fee-schedules/{feeScheduleId}
POST /schools/{schoolId}/fee-schedules/{feeScheduleId}/items
POST /schools/{schoolId}/fee-schedules/{feeScheduleId}/installment-plans
POST /schools/{schoolId}/fee-schedules/{feeScheduleId}/publish
```

## Compte eleve

```text
GET  /schools/{schoolId}/students/{studentId}/account
GET  /schools/{schoolId}/students/{studentId}/charges
POST /schools/{schoolId}/charges/{chargeId}/adjustments
```

Reponse `account` cible :

```json
{
  "currency": "XOF",
  "total_charged": 213000,
  "total_adjustments": 0,
  "total_paid": 93000,
  "balance": 120000,
  "next_due_date": "2026-12-31"
}
```

## Paiements caisse

```text
POST /schools/{schoolId}/cash-payments
GET  /schools/{schoolId}/payments
GET  /schools/{schoolId}/payments/{paymentId}
```

## Paiement parent

```text
POST /payment-links/{token}/intent
GET  /payment-links/{token}
POST /payment-intents/{intentId}/start
GET  /payment-intents/{intentId}/status
```

Le client ne fournit jamais librement `school_id`, `student_id` et `amount` lorsqu'il utilise un lien signe : ces valeurs sont resolues par le serveur.

## Webhooks

```text
POST /webhooks/payments/orange
POST /webhooks/payments/moov
```

Contraintes : verification signature/secret, idempotence, persistance du payload, reponse rapide puis traitement robuste.

## Recus

```text
GET /schools/{schoolId}/receipts/{receiptId}
GET /receipts/verify/{verificationToken}
```

## Portail parent

```text
GET  /parent/children
POST /parent/children/link
POST /parent/children/link/verify
GET  /parent/children/{linkId}/account
GET  /parent/children/{linkId}/payments
GET  /parent/children/{linkId}/receipts
```

## Reporting

```text
GET /schools/{schoolId}/dashboard
GET /schools/{schoolId}/reports/collection
GET /schools/{schoolId}/reports/overdue
GET /schools/{schoolId}/reports/payments
```

## Conventions

- JSON UTF-8
- timestamps ISO-8601 UTC en API
- montants en entier XOF
- `Idempotency-Key` obligatoire pour creation d'intent/paiement sensible
- erreurs structurees avec `code`, `message`, `request_id`
- pagination cursor-based pour les grandes collections
