# KALAN-SIRA

Plateforme SaaS multi-etablissements de gestion administrative, recouvrement et paiement de la scolarite, pensee pour les etablissements scolaires du Burkina Faso.

## Vision

KALAN-SIRA relie l'etablissement, l'eleve et ses responsables autour d'un compte de scolarite unique et auditable : inscription, tarification, echeances, relances, paiement a distance, caisse, recus et pilotage du recouvrement.

## MVP 1

1. Onboarding des etablissements
2. Annees scolaires, niveaux et classes
3. Inscription et reinscription des eleves
4. Responsables/parents et relations multi-enfants
5. Grilles tarifaires versionnees
6. Generation des dettes et echeances
7. Encaissement caisse + paiement mobile
8. Liens de paiement securises
9. Integration Orange Money / Moov Money via adaptateurs
10. Allocation des paiements aux creances
11. Recus electroniques verifiables
12. SMS de rappel et confirmation
13. Dashboard de recouvrement
14. RBAC et journal d'audit

## Stack

- Backend : Go
- API : REST `/api/v1`
- Base : PostgreSQL
- Acces DB : pgx + sqlc
- Cache / jobs : Redis
- Frontend : React + Vite
- Deploiement local : Docker Compose

## Architecture

```text
apps/
  api/                  API Go
  admin-web/            portail etablissement React/Vite
  parent-web/           portail parent/PWA React/Vite
internal/
  auth/
  school/
  student/
  enrollment/
  billing/
  payment/
  notification/
  receipt/
  audit/
database/
  migrations/
  queries/
  seeds/
docs/
  architecture/
  api/
  business-rules/
deployments/
  docker/
```

## Principe financier

Une inscription validee genere un snapshot des obligations financieres de l'eleve. Modifier une grille tarifaire ne modifie jamais retroactivement une dette existante.

```text
Enrollment
   -> Student Charges
      -> Installments
         -> Payment Intent
            -> Provider confirmation
               -> Payment
                  -> Payment Allocations
                     -> Receipt
```

Un paiement n'est considere comme acquis qu'apres confirmation serveur. Une redirection navigateur ne valide jamais une transaction.

Tous les montants XOF sont stockes en entiers (`BIGINT`), jamais en flottants.

## Multi-tenant

Chaque etablissement est un tenant. Les ressources metier portent un `school_id` et les controles d'acces sont appliques cote serveur. Le portail parent peut relier plusieurs enfants de plusieurs etablissements sans exposer les donnees d'un tenant a un autre.

## Documentation

- `docs/architecture/ARCHITECTURE.md`
- `docs/business-rules/BUSINESS_RULES.md`
- `docs/api/API_V1.md`
- `database/migrations/000001_init.up.sql`

## Etat

Phase actuelle : fondations MVP 1.