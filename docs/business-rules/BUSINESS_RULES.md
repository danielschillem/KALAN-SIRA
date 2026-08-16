# Regles metier MVP 1

## Etablissement

- Un etablissement est un tenant autonome.
- Une annee scolaire est unique par etablissement et par nom.
- Les classes appartiennent a une annee scolaire et un niveau.

## Eleve

- L'identifiant public eleve est stable dans son etablissement.
- Un eleve ne possede qu'une inscription par annee scolaire dans le MVP.
- Une reinscription cree une nouvelle `enrollment`; elle ne modifie pas l'historique precedent.

## Responsable

- Un eleve peut avoir plusieurs responsables.
- Un responsable peut avoir plusieurs enfants.
- Un responsable peut etre lie a des enfants de plusieurs etablissements.
- Une relation parent-enfant doit etre verifiee avant l'acces financier a distance.

## Tarification

- Une grille appartient a une annee scolaire et cible un niveau ou une classe.
- Les montants sont en XOF entier.
- Une modification de grille n'affecte pas les charges deja generees.
- Les remises, bourses et corrections sont des ajustements traces.

## Inscription

A la validation d'une inscription :

1. verifier l'annee et la classe;
2. selectionner la grille applicable;
3. creer le snapshot des charges;
4. calculer les echeances;
5. rendre l'inscription ACTIVE;
6. journaliser l'operation.

## Charges

Statuts : `UPCOMING`, `DUE`, `PARTIAL`, `PAID`, `OVERDUE`, `CANCELLED`.

`balance = net_amount - amount_paid` doit toujours etre vrai.

## Paiement

- Un PaymentIntent expire.
- Le montant est determine cote serveur.
- Un provider callback doit etre idempotent.
- `SUCCESS` est terminal sauf procedure explicite d'annulation/remboursement.
- Le surpaiement est refuse par defaut.
- Le paiement partiel depend de la politique de l'echeancier.

## Caisse

Le caissier peut enregistrer un paiement physique. Celui-ci utilise le meme moteur `Payment -> Allocation -> Receipt` que le paiement mobile afin d'eviter deux comptabilites paralleles.

## Recu

- Un paiement confirme genere au maximum un recu actif dans le MVP.
- Le recu porte une reference unique par etablissement.
- La verification publique ne doit pas exposer date de naissance, contacts ou autres donnees inutiles.

## Notifications

Evenements initiaux :

- rappel avant echeance;
- echeance du jour;
- retard;
- paiement confirme;
- paiement echoue.

Les relances ne doivent pas etre envoyees pour une charge deja soldee.

## Audit

Doivent etre audites au minimum :

- activation/annulation d'inscription;
- changement de classe;
- creation/modification de grille publiee;
- remise/bourse/correction;
- paiement manuel;
- annulation/remboursement;
- changement de role utilisateur.
