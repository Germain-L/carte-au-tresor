# Carte au Trésor

Un simulateur de chasse au trésor en Go où des aventuriers explorent une carte pour collecter des trésors en évitant les montagnes.

## Description

Ce projet simule un jeu où des aventuriers se déplacent sur une carte quadrillée selon une séquence d'actions prédéfinies. Ils peuvent collecter des trésors, mais sont bloqués par les montagnes et les autres aventuriers.

## Structure du projet

```txt
cmd/carte-au-tresor/     # Point d'entrée principal
internal/
  ├── domain/
  │   ├── models/        # Entités (Aventurier, Trésor, Montagne)
  │   └── game/          # Moteur de simulation
  ├── parser/            # Analyseur de fichiers d'entrée
  ├── output/            # Générateur de sortie
  └── fileutil/          # Utilitaires I/O
inputs/                  # Fichiers de configuration
outputs/                 # Résultats de simulation
```

## Utilisation

### Exécution

```bash
go run cmd/carte-au-tresor/main.go
```

### Format d'entrée ([`inputs/example.txt`](inputs/example.txt))

```txt
C - 3 - 4                    # Carte (largeur-hauteur)
M - 1 - 0                    # Montagne (x-y)
T - 0 - 3 - 2                # Trésor (x-y-quantité)
A - Lara - 1 - 1 - S - AADADAGGA  # Aventurier (nom-x-y-orientation-actions)
```

### Actions disponibles

- `A` : Avancer
- `G` : Tourner à gauche  
- `D` : Tourner à droite

### Orientations

- `N` : Nord, `S` : Sud, `E` : Est, `O` : Ouest

## Tests

```bash
go test ./...
```

Le projet inclut des tests unitaires complets pour tous les composants principaux.

## Exemple

**Entrée** : Une carte 3x4 avec des montagnes, trésors et l'aventurière Lara  
**Sortie** : Position finale de Lara et trésors restants après simulation

Le fichier de sortie ([`outputs/example.txt`](outputs/example.txt)) contient l'état final de la carte.
