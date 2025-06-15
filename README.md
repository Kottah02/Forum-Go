# Forum-Go

Un forum moderne et interactif développé en Go, permettant aux utilisateurs de partager des idées, discuter de leur jeu préférer COUNTER-STRIKE.

## 🚀 Fonctionnalités

### Authentification et Gestion des Utilisateurs
- Inscription et connexion des utilisateurs 
    -Pseudo
    -Email
    -Mot de passe (12 caractères, Majuscule, Minuscule, Chiffre, caractères spécial)
    -Confirmation de Mot de passe
- Profil utilisateur personnalisé
- Gestion des sessions
- Protection des routes sensibles

### Gestion des Posts
- Création de posts avec titre et contenu
- Système de tags pour catégoriser les posts
- Suppression et modification de posts (uniquement par l'auteur)
- Affichage chronologique des posts
- Système de réactions (like/dislike)

### Système de Commentaires
- Ajout de commentaires sur les posts
- Affichage des commentaires par ordre chronologique
- Intégration avec le profil utilisateur

### Interface Utilisateur
- Design moderne et responsive
- Animations fluides pour une meilleure expérience utilisateur
- Navigation intuitive
- Affichage des statistiques (nombre de posts, commentaires, etc.)

## 🛠️ Prérequis

- Go 1.16 ou supérieur
- MySQL 5.7 ou supérieur
- Un navigateur web moderne

## 📦 Installation

1. Clonez le repository :
```bash
git clone https://github.com/Kottah02/Forum-Go.git
cd Forum-Go
```

2. Configurez la base de données MySQL :
- Créez une base de données nommée `website_db`
- Importez le schéma de la base de données (le fichier sera fourni)

3. Installez les dépendances :
```bash
go mod download
```

4. Lancez l'application :
```bash
go run cmd/main.go
```

L'application sera accessible à l'adresse : `http://localhost:8080`

Utilisateur créer : NOM D'UTILISATEUR : kottah , MDP : 98Y76em21.eddy

## 🏗️ Structure du Projet

```
Forum-Go/
├── cmd/
│   └── main.go           # Point d'entrée de l'application
├── internal/
│   ├── config/          # Configuration de l'application
│   ├── controllers/     # Gestionnaires de requêtes
│   ├── middleware/      # Middleware (auth, etc.)
│   ├── models/          # Modèles de données
│   └── routes/          # Définition des routes
├── static/
│   ├── css/            # Styles CSS
│   └── js/             # Scripts JavaScript
├── templates/          # Templates HTML
├── go.mod             # Dépendances Go
└── README.md          # Documentation
```

## 🔒 Sécurité

- Protection CSRF sur les formulaires
- Validation des entrées utilisateur
- Gestion sécurisée des sessions
- Protection des routes sensibles
- Transactions SQL pour l'intégrité des données

## 🎨 Personnalisation

Le projet utilise des fichiers CSS et JavaScript modulaires qui peuvent être facilement personnalisés :
- `static/css/style.css` : Styles principaux
- `static/js/animations.js` : Animations et interactions

## 🤝 Contribution

Les contributions sont les bienvenues ! N'hésitez pas à :
1. Fork le projet
2. Créer une branche pour votre fonctionnalité
3. Commiter vos changements
4. Pousser vers la branche
5. Ouvrir une Pull Request

## 📝 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 👥 Auteurs

- AMIR EDDY -LEBERRE ETHIENNE