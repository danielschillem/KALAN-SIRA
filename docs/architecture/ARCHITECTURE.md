# Architecture KALAN-SIRA

## 1. Frontieres du produit

KALAN-SIRA MVP 1 couvre l'administration de l'inscription et le cycle financier scolaire. Les notes, bulletins, emplois du temps, cours et gestion pedagogique sont hors scope.

## 2. Bounded contexts

- `school` : tenant, annee scolaire, niveaux, classes
- `student` : dossier eleve et responsables
- `enrollment` : inscription/reinscription annuelle
- `billing` : grille tarifaire, echeancier, charges, ajustements
- `payment` : intents, transactions, allocations, providers
- `notification` : SMS, rappels et confirmations
- `receipt` : recus et verification
- `audit` : tracabilite des operations sensibles
- `auth` : utilisateurs, roles et autorisations

## 3. Regles d'architecture

1. `school_id` est une frontiere de securite, pas seulement un filtre UI.
2. Le serveur derive le tenant autorise depuis l'identite authentifiee.
3. Les montants sont des entiers XOF (`BIGINT`).
4. Une grille tarifaire publiee est versionnee; une inscription genere des charges figees.
5. Aucun paiement n'est valide sur la seule base d'un retour navigateur.
6. Les callbacks providers doivent etre authentifies, idempotents et journalises.
7. Une transaction financiere validee n'est jamais supprimee; correction par operation compensatoire.
8. Les tokens de liens de paiement et de verification sensibles sont stockes sous forme de hash.
9. Toute modification financiere sensible produit un audit log.

## 4. Flux d'inscription

```text
Student -> Enrollment -> resolve FeeSchedule -> snapshot charges -> ACTIVE
```

La validation doit etre transactionnelle. Si la generation des charges echoue, l'inscription ne passe pas a `ACTIVE`.

## 5. Flux de paiement

```text
Guardian/Cashier
  -> PaymentIntent
  -> Provider adapter
  -> provider pending
  -> authenticated callback
  -> idempotency check
  -> Payment SUCCESS
  -> PaymentAllocations
  -> update StudentCharges
  -> Receipt
  -> Notification
```

## 6. Allocation

Un paiement peut couvrir plusieurs charges et une charge peut recevoir plusieurs paiements. `payment_allocations` est donc la source de verite du rapprochement paiement-creance.

L'ordre automatique par defaut sera :

1. charges obligatoires arrivees a echeance;
2. plus ancienne date d'echeance;
3. ordre de creation.

Une politique d'allocation specifique pourra etre ajoutee par etablissement.

## 7. Multi-parent et multi-ecole

Le dossier eleve appartient a son etablissement. Un compte parent peut etre relie a plusieurs dossiers eleves de plusieurs tenants. Une connaissance du `public_id` eleve ne suffit jamais pour creer la relation : validation OTP ou code d'association obligatoire.

## 8. Integration paiement

Le domaine utilise une interface provider stable :

```text
CreatePayment(intent)
GetStatus(reference)
VerifyCallback(headers, body)
ParseCallback(body)
```

Les implementations Orange Money et Moov Money restent derriere cette interface afin de ne pas contaminer le domaine scolaire avec des details operateurs.

## 9. Jobs asynchrones

Redis servira initialement a :

- relances d'echeances;
- SMS;
- generation de recus;
- reconciliation de transactions en attente;
- retries providers.

## 10. Securite

Les endpoints etablissement sont tenant-scoped et RBAC. Les endpoints parents ne retournent que les enfants dont la relation est verifiee. Les pages publiques de verification de recu exposent le minimum necessaire.
