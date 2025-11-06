# Création de fichiers Office dans DriveView

## Fonctionnalité

Permet de créer des nouveaux documents Office (Word, Excel, PowerPoint) directement depuis l'interface DriveView et de les ouvrir immédiatement dans OnlyOffice pour édition.

## Utilisation

### Interface utilisateur

1. Dans DriveView, cliquer sur le bouton **"📄 Nouveau fichier"**
2. Sélectionner le type de fichier:
   - 📄 Document Word (.docx)
   - 📊 Feuille Excel (.xlsx)
   - 📽️ Présentation PowerPoint (.pptx)
3. Entrer le nom du fichier (l'extension sera ajoutée automatiquement)
4. Cliquer sur **"Créer et ouvrir"**
5. Le fichier est créé et s'ouvre automatiquement dans OnlyOffice

### Backend API

**Endpoint:** `POST /api/drive/createFile`

**Requête:**
```json
{
  "username": "jean",
  "token": "session_token",
  "parent_path": "/",
  "file_name": "Mon Document",
  "file_type": "docx"
}
```

**Réponse (succès):**
```json
{
  "message": "file created",
  "file_id": 123,
  "file_name": "Mon Document.docx",
  "parent": "/"
}
```

**Types de fichiers supportés:**
- `docx` - Document Word
- `xlsx` - Feuille Excel
- `pptx` - Présentation PowerPoint

## Structure technique

### Frontend (Vue.js)

**Fichier:** `webfront2/src/views/DriveView.vue`

- **Modal:** `showNewFileModal` pour saisir le nom et type
- **État:** 
  - `newFileName` - nom du fichier
  - `newFileType` - type (docx/xlsx/pptx)
- **Fonction:** `createFile()` - appelle l'API et ouvre OnlyOffice

### Backend (Go)

**Fichier:** `backend/drive.go`

**Fonction principale:** `createFile(c *gin.Context)`

1. Valide la session utilisateur
2. Normalise le nom du fichier
3. Ajoute l'extension selon le type
4. Vérifie l'absence de conflit
5. Crée le fichier depuis un template
6. Insère l'entrée en base de données
7. Retourne l'ID du nouveau fichier

**Fonctions helpers:**
- `createEmptyDocx()` - charge `templates/empty.docx`
- `createEmptyXlsx()` - charge `templates/empty.xlsx`
- `createEmptyPptx()` - charge `templates/empty.pptx`

### Templates

**Emplacement:** `backend/templates/`

Fichiers vierges Office OpenXML créés avec:
- Microsoft Word pour `.docx`
- Microsoft Excel pour `.xlsx`
- Microsoft PowerPoint pour `.pptx`

**Génération:** Script PowerShell `create_templates.ps1`

## Types MIME

- DOCX: `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- XLSX: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- PPTX: `application/vnd.openxmlformats-officedocument.presentationml.presentation`

## Sécurité

- ✅ Validation de session obligatoire
- ✅ Sanitisation du nom de fichier (`filepath.Base`)
- ✅ Vérification de conflit (nom existant)
- ✅ Création dans le dossier utilisateur uniquement
- ✅ Interdiction de créer dans `.trash`

## Intégration OnlyOffice

Après création, le fichier est automatiquement:
1. Ajouté à la liste des fichiers
2. Sélectionné dans l'interface
3. Ouvert dans le modal OnlyOffice
4. Prêt pour édition collaborative

## Erreurs possibles

| Code | Message | Cause |
|------|---------|-------|
| 400 | invalid request | JSON malformé |
| 401 | invalid session | Token expiré/invalide |
| 404 | user not found | Utilisateur inexistant |
| 409 | file already exists | Nom déjà utilisé |
| 500 | db insert failed | Erreur base de données |
| 500 | cannot write file | Erreur système de fichiers |

## TODO / Améliorations futures

- [ ] Support de fichiers ODT/ODS/ODP (LibreOffice)
- [ ] Templates personnalisés par utilisateur
- [ ] Modèles de documents (facture, CV, etc.)
- [ ] Duplication de fichiers existants
- [ ] Import depuis Google Drive/OneDrive
