package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	downloadTokens     = map[string]downloadToken{}
	downloadTokensLock sync.RWMutex
)

// === Structures ===

type DriveFile struct {
	FileID       int       `json:"file_id"`
	UserID       int       `json:"user_id"`
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	FileSize     int64     `json:"file_size"`
	FileType     string    `json:"file_type"`
	DateUploaded time.Time `json:"date_uploaded"`
}

type FileRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	FileID   int    `json:"file_id,omitempty"`
	NewName  string `json:"new_name,omitempty"`
}

type FileListResponse struct {
	Files     []DriveFile `json:"files"`
	UserEmail string      `json:"user_email,omitempty"`
	Error     string      `json:"error,omitempty"`
}

// add token store for OnlyOffice temporary download tokens
type downloadToken struct {
	FileID   int
	Username string
	Expires  time.Time
}

// === getfiles ===
func getfiles(c *gin.Context) {
	var verif FileRequest
	if err := c.ShouldBindJSON(&verif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Vérification de session sécurisée
	if session, ok := sessions[verif.Token]; !ok || session.Username != verif.Username || verif.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", verif.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	rows, err := db.Query(`SELECT file_id, user_id, file_name, file_path, file_size, file_type, date_uploaded 
	                       FROM DriveFiles WHERE user_id=$1 ORDER BY date_uploaded DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var files []DriveFile
	for rows.Next() {
		var f DriveFile
		if err := rows.Scan(&f.FileID, &f.UserID, &f.FileName, &f.FilePath, &f.FileSize, &f.FileType, &f.DateUploaded); err == nil {
			files = append(files, f)
		}
	}

	c.JSON(http.StatusOK, FileListResponse{Files: files})
}

func uploadFile(c *gin.Context) {
	username := c.PostForm("username")
	token := c.PostForm("token")
	parentPath := strings.TrimSpace(c.PostForm("parent_path"))

	// normalize parent_path: store as "/" or "/path/.../" (server prefers leading+trailing slash)
	if parentPath == "" {
		parentPath = "/"
	} else {
		if !strings.HasPrefix(parentPath, "/") {
			parentPath = "/" + parentPath
		}
		if !strings.HasSuffix(parentPath, "/") {
			parentPath = parentPath + "/"
		}
	}
	// block uploads into trash
	if strings.HasPrefix(parentPath, "/.trash") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "upload into trash is not allowed"})
		return
	}

	// Vérif session
	if session, ok := sessions[token]; !ok || session.Username != username || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Récupération user_id
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Récupération des fichiers envoyés
	var fileHeaders []*multipart.FileHeader
	if form, err := c.MultipartForm(); err == nil && form != nil {
		for _, fhs := range form.File {
			for _, fh := range fhs {
				fileHeaders = append(fileHeaders, fh)
			}
		}
	}
	if len(fileHeaders) == 0 {
		if fh, err := c.FormFile("files"); err == nil {
			fileHeaders = append(fileHeaders, fh)
		}
	}
	if len(fileHeaders) == 0 {
		if fh, err := c.FormFile("file"); err == nil {
			fileHeaders = append(fileHeaders, fh)
		}
	}
	if len(fileHeaders) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no files uploaded"})
		return
	}

	// Créer le dossier utilisateur si besoin
	uploadDir := filepath.Join("uploads", username)
	if parentPath != "" && parentPath != "/" {
		// trim leading slash before Join so it's not treated as absolute
		p := strings.TrimPrefix(parentPath, "/")
		uploadDir = filepath.Join(uploadDir, p)
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}

	var uploaded []gin.H
	anyError := false

	for _, fh := range fileHeaders {
		// Garder le nom original
		originalName := filepath.Base(fh.Filename)
		dst := filepath.Join(uploadDir, originalName)

		// Si le fichier existe déjà, on ajoute un suffixe
		if _, err := os.Stat(dst); err == nil {
			base := strings.TrimSuffix(originalName, filepath.Ext(originalName))
			ext := filepath.Ext(originalName)
			for i := 1; ; i++ {
				candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
				candidatePath := filepath.Join(uploadDir, candidate)
				if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
					dst = candidatePath
					originalName = candidate
					break
				}
			}
		}

		// Ouverture du fichier uploadé
		in, err := fh.Open()
		if err != nil {
			uploaded = append(uploaded, gin.H{"file": fh.Filename, "error": err.Error()})
			continue
		}
		out, err := os.Create(dst)
		if err != nil {
			in.Close()
			uploaded = append(uploaded, gin.H{"file": fh.Filename, "error": err.Error()})
			continue
		}

		_, err = io.Copy(out, in)
		in.Close()
		out.Close()
		if err != nil {
			uploaded = append(uploaded, gin.H{"file": fh.Filename, "error": err.Error()})
			continue
		}

		// after writing, get actual size from disk
		fi, statErr := os.Stat(dst)
		var actualSize int64 = 0
		if statErr == nil {
			actualSize = fi.Size()
		} else {
			log.Println("stat after save failed:", statErr)
		}

		// Déterminer le type MIME
		contentType := fh.Header.Get("Content-Type")
		if contentType == "" {
			if ext := filepath.Ext(originalName); ext != "" {
				contentType = mime.TypeByExtension(ext)
			}
			if contentType == "" {
				contentType = "application/octet-stream"
			}
		}

		// Enregistrement en base de données — utilise la taille réelle
		dbPath := parentPath
		if dbPath == "" {
			dbPath = "/"
		}
		_, err = db.Exec(`INSERT INTO DriveFiles 
			(user_id, file_name, file_path, file_size, file_type, date_uploaded)
			VALUES ($1, $2, $3, $4, $5, NOW())`,
			userID, originalName, dbPath, actualSize, contentType)
		if err != nil {
			log.Println("db insert failed for", originalName, "err:", err)
			_ = os.Remove(dst)
			uploaded = append(uploaded, gin.H{"file": fh.Filename, "error": err.Error()})
			anyError = true
			continue
		}

		uploaded = append(uploaded, gin.H{"file": originalName, "status": "uploaded"})
	}

	// si au moins un fichier a échoué on renvoie 500 pour que le front marque erreur
	if anyError {
		c.JSON(http.StatusInternalServerError, gin.H{"uploaded": uploaded})
		return
	}
	c.JSON(http.StatusOK, gin.H{"uploaded": uploaded})
}

// === downloadFile ===
func downloadFile(c *gin.Context) {
	// debug
	log.Printf("download request: remote=%s query=%s headers=%v", c.ClientIP(), c.Request.URL.RawQuery, c.Request.Header)

	// Track if a valid temporary auth was used (download_token)
	var tempAuthUsed bool

	// Read incoming query values first (may be empty)
	fileIDStr := c.Query("file_id")
	username := c.Query("username")
	token := c.Query("token")
	downloadMode := c.Query("download") // "1" or "true" => attachment

	// If download_token present, prefer it and override fileIDStr/username directly
	if dt := c.Query("download_token"); dt != "" {
		downloadTokensLock.RLock()
		dentry, ok := downloadTokens[dt]
		downloadTokensLock.RUnlock()
		if ok {
			if time.Now().Before(dentry.Expires) {
				tempAuthUsed = true
				// override using the secure stored values — do NOT rely on rewriting RawQuery
				fileIDStr = strconv.Itoa(dentry.FileID)
				username = dentry.Username
				log.Printf("download_token accepted dt=%s file_id=%s username=%s expires=%s", dt, fileIDStr, dentry.Username, dentry.Expires)
			} else {
				downloadTokensLock.Lock()
				delete(downloadTokens, dt)
				downloadTokensLock.Unlock()
				log.Printf("download_token expired dt=%s", dt)
				c.JSON(http.StatusForbidden, gin.H{"error": "download token expired"})
				return
			}
		} else {
			log.Printf("download_token invalid dt=%s", dt)
			c.JSON(http.StatusForbidden, gin.H{"error": "invalid download token"})
			return
		}
	}

	// session check unless tempAuthUsed
	if token != "" && !tempAuthUsed {
		session, ok := sessions[token]
		if !ok {
			log.Printf("invalid session token while downloading file_id=%s username=%s", fileIDStr, username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}
		if username == "" {
			username = session.Username
		}
		if session.Username != username || username == "" {
			log.Printf("session username mismatch: session=%s requested=%s", session.Username, username)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
			return
		}
	} else if tempAuthUsed {
		// already populated username from token entry
	}

	// parse file_id (with clearer error log)
	fid, err := strconv.Atoi(fileIDStr)
	if err != nil {
		log.Printf("invalid or missing file_id: %q (tempAuthUsed=%v)", fileIDStr, tempAuthUsed)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file_id"})
		return
	}

	var file DriveFile
	var ownerUsername string
	err = db.QueryRow(`SELECT f.file_id, f.user_id, f.file_name, f.file_path, f.file_size, f.file_type, f.date_uploaded, u.username
                       FROM DriveFiles f 
                       JOIN Users u ON f.user_id=u.user_id 
                       WHERE f.file_id=$1`, fid).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath,
			&file.FileSize, &file.FileType, &file.DateUploaded, &ownerUsername)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	// authorize: owner must match requested username (username filled from session or token)
	if ownerUsername != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// build path on disk (trim leading slash)
	filePath := filepath.Join("uploads", ownerUsername)
	if file.FilePath != "" && file.FilePath != "/" {
		p := strings.TrimPrefix(file.FilePath, "/")
		if p != "" && p != "/" {
			filePath = filepath.Join(filePath, p)
		}
	}
	filePath = filepath.Join(filePath, file.FileName)
	log.Printf("Resolved filePath=%s (file_id=%d owner=%s)", filePath, fid, ownerUsername)

	fi, statErr := os.Stat(filePath)
	if os.IsNotExist(statErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk"})
		return
	} else if statErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot access file on disk"})
		return
	}

	f, err := os.Open(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file"})
		return
	}
	defer f.Close()

	// determine content type: prefer DB value, else detect from bytes, else by extension
	var contentType string
	if file.FileType != "" {
		contentType = file.FileType
	}
	head := make([]byte, 512)
	n, _ := f.Read(head)
	if contentType == "" || contentType == "application/octet-stream" {
		detected := http.DetectContentType(head[:n])
		if detected != "" && detected != "application/octet-stream" {
			contentType = detected
		} else {
			if ext := filepath.Ext(file.FileName); ext != "" {
				if t := mime.TypeByExtension(ext); t != "" {
					contentType = t
				}
			}
		}
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// headers: allow inline preview by default, allow forcing attachment via download param
	disposition := "inline"
	if downloadMode == "1" || strings.ToLower(downloadMode) == "true" {
		disposition = "attachment"
	}
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("%s; filename=%q", disposition, file.FileName))
	_, _ = f.Seek(0, 0)
	http.ServeContent(c.Writer, c.Request, file.FileName, fi.ModTime(), f)
}

// === restoreFile ===
func restoreFile(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// same session check style as other endpoints
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var file DriveFile
	err := db.QueryRow("SELECT file_id, user_id, file_name, file_path FROM DriveFiles WHERE file_id=$1",
		req.FileID).Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	if file.FilePath != "/.trash/" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is not in trash"})
		return
	}

	userDir := filepath.Join("uploads", req.Username)
	trashPath := filepath.Join(userDir, ".trash", file.FileName)
	destDir := userDir
	dest := filepath.Join(destDir, file.FileName)

	// avoid overwrite in destination: add suffix if necessary
	if _, err := os.Stat(dest); err == nil {
		base := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s_restored_%d%s", base, i, ext)
			candPath := filepath.Join(destDir, candidate)
			if _, err := os.Stat(candPath); os.IsNotExist(err) {
				dest = candPath
				file.FileName = candidate
				break
			}
		}
	}

	if err := os.Rename(trashPath, dest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot restore file", "err": err.Error()})
		return
	}

	// update DB: put back in root ("/") — we don't have the original path stored so we restore to "/" to be safe
	_, err = db.Exec("UPDATE DriveFiles SET file_path=$1, file_name=$2 WHERE file_id=$3", "/", file.FileName, req.FileID)
	if err != nil {
		// attempt rollback of file move
		_ = os.Rename(dest, trashPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restored", "file_id": req.FileID})
}

// === moveToTrash ===
func moveToTrash(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Token    string `json:"token"`
		FileID   int    `json:"file_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// session check (same style)
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// fetch file record
	var file DriveFile
	err := db.QueryRow(`SELECT file_id, user_id, file_name, file_path, file_type 
	                    FROM DriveFiles WHERE file_id=$1`, req.FileID).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	// owner check
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// build safe paths
	userDir := filepath.Clean(filepath.Join("uploads", req.Username))
	trashDir := filepath.Join(userDir, ".trash")
	if err := os.MkdirAll(trashDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create trash dir", "err": err.Error()})
		return
	}

	// construct src path
	srcBase := userDir
	if file.FilePath != "" && file.FilePath != "/" {
		// file.FilePath stored like "/a/b/" maybe - remove leading slash before Join
		p := strings.TrimPrefix(file.FilePath, "/")
		srcBase = filepath.Join(srcBase, p)
	}
	src := filepath.Join(srcBase, file.FileName)
	src = filepath.Clean(src)

	// ensure src exists
	info, statErr := os.Stat(src)
	if os.IsNotExist(statErr) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk", "src": src})
		return
	} else if statErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot stat source", "err": statErr.Error(), "src": src})
		return
	}

	// prepare dst (avoid collision)
	dst := filepath.Join(trashDir, file.FileName)
	dst = filepath.Clean(dst)
	if _, err := os.Stat(dst); err == nil {
		// collision: add suffix
		base := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
			candPath := filepath.Join(trashDir, candidate)
			if _, err := os.Stat(candPath); os.IsNotExist(err) {
				dst = candPath
				file.FileName = candidate // use new name in DB update later
				break
			}
		}
	}

	// attempt rename first (fast path)
	moveErr := os.Rename(src, dst)
	if moveErr != nil {
		// fallback: copy then remove (handles cross-device, insufficient rename perms in some cases)
		if info.IsDir() {
			if copyErr := copyDir(src, dst); copyErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to move directory to trash",
					"err":   copyErr.Error(),
					"src":   src,
					"dst":   dst,
				})
				return
			}
			// remove original dir
			if rmErr := os.RemoveAll(src); rmErr != nil {
				// log but continue with DB update; inform user
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "moved to trash but failed to remove original directory",
					"err":   rmErr.Error(),
					"src":   src,
					"dst":   dst,
				})
				return
			}
		} else {
			// file case
			if copyErr := copyFile(src, dst); copyErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "failed to move file to trash (copy fallback failed)",
					"err":   copyErr.Error(),
					"src":   src,
					"dst":   dst,
				})
				return
			}
			if rmErr := os.Remove(src); rmErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "moved to trash but failed to remove original file",
					"err":   rmErr.Error(),
					"src":   src,
					"dst":   dst,
				})
				return
			}
		}
	}

	// --- DB updates: only after successful move/copy ---
	// compute original prefix and new prefix (for folder children update)
	origPrefix := "/"
	if file.FilePath != "" && file.FilePath != "/" {
		origPrefix = file.FilePath
	}
	if !strings.HasSuffix(origPrefix, "/") {
		origPrefix = origPrefix + "/"
	}
	origPrefix = strings.TrimSuffix(origPrefix, "/") + "/" + file.FileName + "/"
	origPrefix = strings.ReplaceAll(origPrefix, "//", "/")

	newPrefix := "/.trash/"
	newPrefix = strings.TrimSuffix(newPrefix, "/") + "/" + file.FileName + "/"
	newPrefix = strings.ReplaceAll(newPrefix, "//", "/")

	// update the moved item itself (its file_path becomes "/.trash/" and file_name possibly changed above)
	if _, err := db.Exec("UPDATE DriveFiles SET file_path=$1, file_name=$2 WHERE file_id=$3", "/.trash/", file.FileName, file.FileID); err != nil {
		// attempt rollback: try to move back (best-effort)
		_ = os.Rename(dst, src)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed", "err": err.Error()})
		return
	}

	// if it was a folder, update children file_path prefixes
	isFolder := strings.Contains(strings.ToLower(file.FileType), "folder")
	if isFolder {
		likePattern := origPrefix + "%"
		if _, err := db.Exec(`UPDATE DriveFiles 
		                       SET file_path = replace(file_path, $1, $2) 
		                       WHERE user_id=$3 AND file_path LIKE $4`, origPrefix, newPrefix, userID, likePattern); err != nil {
			// best-effort rollback (try restore)
			_ = os.Rename(dst, src)
			_, _ = db.Exec("UPDATE DriveFiles SET file_path=$1, file_name=$2 WHERE file_id=$3", file.FilePath, file.FileName, file.FileID)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db update children failed", "err": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "moved to trash", "file_id": req.FileID})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// ensure parent dir exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(dst)
		return err
	}

	// copy file mode
	if fi, err := os.Stat(src); err == nil {
		_ = os.Chmod(dst, fi.Mode())
	}

	return nil
}

func copyDir(srcDir, dstDir string) error {
	// walk source
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dstDir, rel)
		if info.IsDir() {
			if err := os.MkdirAll(targetPath, info.Mode()); err != nil {
				return err
			}
			return nil
		}
		// file: copy
		if err := copyFile(path, targetPath); err != nil {
			return err
		}
		return nil
	})
}

func getTrashFiles(c *gin.Context) {
	var verif FileRequest
	if err := c.ShouldBindJSON(&verif); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// same session check as getfiles / upload
	if session, ok := sessions[verif.Token]; !ok || session.Username != verif.Username || verif.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", verif.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	rows, err := db.Query(`SELECT file_id, user_id, file_name, file_path, file_size, file_type, date_uploaded
                           FROM DriveFiles WHERE user_id=$1 AND file_path = $2 ORDER BY date_uploaded DESC`,
		userID, "/.trash/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	var files []DriveFile
	for rows.Next() {
		var f DriveFile
		if err := rows.Scan(&f.FileID, &f.UserID, &f.FileName, &f.FilePath, &f.FileSize, &f.FileType, &f.DateUploaded); err == nil {
			files = append(files, f)
		}
	}

	c.JSON(http.StatusOK, FileListResponse{Files: files})
}

// === deletePermanent ===
func deletePermanent(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var file DriveFile
	err := db.QueryRow("SELECT file_id, user_id, file_name, file_path, file_type FROM DriveFiles WHERE file_id=$1",
		req.FileID).Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	var userID int
	_ = db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID)
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	userDir := filepath.Join("uploads", req.Username)

	// si le fichier/dossier est dans la corbeille, on supprime récursivement
	// détecter si c'est un dossier
	isFolder := strings.Contains(strings.ToLower(file.FileType), "folder")

	// chemin sur disque du fichier/dossier
	srcBase := filepath.Join(userDir)
	if file.FilePath != "" && file.FilePath != "/" {
		srcBase = filepath.Join(srcBase, file.FilePath)
	}
	targetPath := filepath.Join(srcBase, file.FileName)

	if isFolder {
		// Si dossier : supprimer récursivement le dossier sur disque
		if err := os.RemoveAll(targetPath); err != nil {
			// log mais continuer à tenter suppression DB
			// c.JSON(http.StatusInternalServerError, gin.H{"error":"cannot remove folder on disk", "err": err.Error()})
			// return
		}

		// Construire le préfixe correspondant dans la corbeille : ex "/.trash/folder/"
		trashPrefix := file.FilePath
		if trashPrefix == "" || trashPrefix == "/" {
			trashPrefix = "/"
		}
		// ensure trailing slash
		if !strings.HasSuffix(trashPrefix, "/") {
			trashPrefix = trashPrefix + "/"
		}
		trashPrefix = strings.TrimSuffix(trashPrefix, "/") + "/" + file.FileName + "/"
		trashPrefix = strings.ReplaceAll(trashPrefix, "//", "/")

		likePattern := trashPrefix + "%"

		// Supprimer toutes les lignes DB correspondant au dossier et à ses enfants :
		// condition : file_path LIKE trashPrefix% OR (file_path = parent AND file_name = folderName)
		_, err = db.Exec(`DELETE FROM DriveFiles 
		                   WHERE user_id=$1 AND (file_path LIKE $2 OR (file_path=$3 AND file_name=$4))`,
			userID, likePattern, file.FilePath, file.FileName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db delete failed", "err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "folder and children deleted permanently", "file_id": req.FileID})
		return
	}

	// Si ce n'est pas un dossier : comportement existant (supprimer fichier sur disque + DB row)
	src := filepath.Join(userDir)
	if file.FilePath != "" && file.FilePath != "/" {
		src = filepath.Join(src, file.FilePath)
	}
	src = filepath.Join(src, file.FileName)

	_ = os.Remove(src)
	_, _ = db.Exec("DELETE FROM DriveFiles WHERE file_id=$1", req.FileID)

	c.JSON(http.StatusOK, gin.H{"message": "deleted permanently", "file_id": req.FileID})
}

// === renameFile ===
func renameFile(c *gin.Context) {
	var req FileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// auth comme ailleurs
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// récupérer infos du fichier (incluant file_path)
	var file DriveFile
	err := db.QueryRow("SELECT file_id, user_id, file_name, file_path FROM DriveFiles WHERE file_id=$1", req.FileID).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	// vérifier que l'utilisateur possède bien le fichier
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// sanitize new name
	newName := filepath.Base(req.NewName)
	if newName == ".trash" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_name"})
		return
	}
	if strings.TrimSpace(newName) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new_name"})
		return
	}

	// construire chemins sur disque en tenant compte de file_path
	userDir := filepath.Join("uploads", req.Username)
	srcDir := userDir
	if file.FilePath != "" && file.FilePath != "/" {
		p := strings.TrimPrefix(file.FilePath, "/")
		srcDir = filepath.Join(srcDir, p)
	}
	srcPath := filepath.Join(srcDir, file.FileName)
	dstPath := filepath.Join(srcDir, newName)

	// vérifier existence source
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk"})
		return
	}

	// éviter écrasement accidentel
	if _, err := os.Stat(dstPath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "target filename already exists"})
		return
	}

	// renommer sur le disque
	if err := os.Rename(srcPath, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot rename file on disk"})
		return
	}

	// mettre à jour la base ; revert si échec
	if _, err := db.Exec(`UPDATE DriveFiles SET file_name=$1 WHERE file_id=$2`, newName, req.FileID); err != nil {
		// tentative de rollback du renommage sur disque
		_ = os.Rename(dstPath, srcPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rename success", "file_id": req.FileID})
}

// === createFolder ===
func createFolder(c *gin.Context) {
	type reqT struct {
		Username   string `json:"username"`
		Token      string `json:"token"`
		ParentPath string `json:"parent_path"`
		FolderName string `json:"folder_name"`
	}
	var req reqT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// session check (same style)
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// sanitize folder name
	name := filepath.Base(strings.TrimSpace(req.FolderName))
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid folder_name"})
		return
	}
	// block creating the special .trash folder
	if name == ".trash" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "creation of folder named .trash is not allowed"})
		return
	}
	// normalize parent path
	parent := req.ParentPath
	if parent == "" {
		parent = "/"
	}
	// ensure trailing slash behavior consistent with your code: store parent paths without duplicate slashes
	if parent != "/" {
		parent = strings.TrimPrefix(parent, "/")
		if !strings.HasSuffix(parent, "/") {
			parent = parent + "/"
		}
		parent = "/" + parent // ensure leading slash
	}
	// do not allow creating under .trash
	if strings.HasPrefix(parent, "/.trash") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot create folders inside trash"})
		return
	}

	// find user id
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// check conflict: same name at same path
	var exists bool
	err := db.QueryRow(`SELECT EXISTS (
		SELECT 1 FROM DriveFiles WHERE user_id=$1 AND file_path=$2 AND file_name=$3
	)`, userID, parent, name).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "target already exists"})
		return
	}

	// create folder on disk (mirror current logic: uploads/<username>/<parent>/<name>/)
	userDir := filepath.Join("uploads", req.Username)
	destDir := userDir
	if parent != "" && parent != "/" {
		// parent begins with '/'
		p := strings.TrimPrefix(parent, "/")
		destDir = filepath.Join(destDir, p)
	}
	// create folder path: <destDir>/<name>
	folderPath := filepath.Join(destDir, name)
	if err := os.MkdirAll(folderPath, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create folder on disk"})
		return
	}

	// insert folder as DriveFiles entry with file_type='folder' and file_size=0
	_, err = db.Exec(`INSERT INTO DriveFiles (user_id, file_name, file_path, file_size, file_type, date_uploaded)
	                  VALUES ($1,$2,$3,$4,$5,NOW())`,
		userID, name, parent, 0, "folder")
	if err != nil {
		// try to cleanup disk folder (best-effort)
		_ = os.RemoveAll(folderPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "folder created", "folder": name, "parent": parent})
}

// === moveFile ===
func moveFile(c *gin.Context) {
	type reqT struct {
		Username        string `json:"username"`
		Token           string `json:"token"`
		FileID          int    `json:"file_id"`
		DestinationPath string `json:"destination_path"`
	}
	var req reqT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Vérif session
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Normaliser destination
	dest := strings.TrimSpace(req.DestinationPath)
	if dest == "" {
		dest = "/"
	}
	if !strings.HasPrefix(dest, "/") {
		dest = "/" + dest
	}
	if !strings.HasSuffix(dest, "/") {
		dest = dest + "/"
	}

	// Empêcher déplacement vers corbeille via ce endpoint
	if strings.HasPrefix(dest, "/.trash") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot move directly into trash"})
		return
	}

	// Récup info fichier
	var file DriveFile
	err := db.QueryRow(`SELECT file_id, user_id, file_name, file_path, file_type FROM DriveFiles WHERE file_id=$1`, req.FileID).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	// Vérif propriétaire
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Construire chemins disques
	userDir := filepath.Join("uploads", req.Username)
	srcBase := userDir
	if file.FilePath != "" && file.FilePath != "/" {
		p := strings.TrimPrefix(file.FilePath, "/")
		srcBase = filepath.Join(srcBase, p)
	}
	srcPath := filepath.Join(srcBase, file.FileName)

	destBase := filepath.Join(userDir)
	if dest != "/" {
		p := strings.TrimPrefix(dest, "/")
		destBase = filepath.Join(destBase, p)
	}
	if err := os.MkdirAll(destBase, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create destination"})
		return
	}
	destPath := filepath.Join(destBase, file.FileName)

	// Gérer collision
	if _, err := os.Stat(destPath); err == nil {
		base := strings.TrimSuffix(file.FileName, filepath.Ext(file.FileName))
		ext := filepath.Ext(file.FileName)
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s_%d%s", base, i, ext)
			candidatePath := filepath.Join(destBase, candidate)
			if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
				destPath = candidatePath
				file.FileName = candidate
				break
			}
		}
	}

	// Déplacement physique
	if err := os.Rename(srcPath, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot move file on disk", "err": err.Error()})
		return
	}

	// Mise à jour DB
	if _, err := db.Exec(`UPDATE DriveFiles SET file_path=$1, file_name=$2 WHERE file_id=$3`, dest, file.FileName, file.FileID); err != nil {
		// rollback du move sur disque si erreur DB
		_ = os.Rename(destPath, srcPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file moved successfully",
		"file_id":  req.FileID,
		"new_path": dest,
		"new_name": file.FileName,
	})
}

// === moveFolder ===
// Déplace un dossier (et tout son contenu) identifié par folder_path vers destination_path.
// JSON attendu:
//
//	{
//	  "username": "alice",
//	  "token": "abcd",
//	  "folder_path": "/src/folder/",
//	  "destination_path": "/dest/path/"
//	}
func moveFolder(c *gin.Context) {
	type reqT struct {
		Username        string `json:"username"`
		Token           string `json:"token"`
		FolderPath      string `json:"folder_path"`
		DestinationPath string `json:"destination_path"`
	}
	var req reqT
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Vérif session
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Normaliser source et destination
	src := strings.TrimSpace(req.FolderPath)
	if src == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder_path is required"})
		return
	}
	if !strings.HasPrefix(src, "/") {
		src = "/" + src
	}
	if !strings.HasSuffix(src, "/") {
		src = src + "/"
	}

	dest := strings.TrimSpace(req.DestinationPath)
	if dest == "" {
		dest = "/"
	}
	if !strings.HasPrefix(dest, "/") {
		dest = "/" + dest
	}
	if !strings.HasSuffix(dest, "/") {
		dest = dest + "/"
	}

	// Empêcher déplacement vers corbeille via ce endpoint
	if strings.HasPrefix(dest, "/.trash") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot move directly into trash"})
		return
	}

	// Récup user id et username
	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user lookup failed"})
		return
	}
	var username string
	if err := db.QueryRow("SELECT username FROM Users WHERE user_id=$1", userID).Scan(&username); err != nil {
		// fallback to provided username
		username = req.Username
	}

	userDir := filepath.Join("uploads", username)

	// chemin source physique: uploads/<username>/<src trimmed>
	srcBase := userDir
	if src != "/" {
		p := strings.TrimPrefix(src, "/")
		srcBase = filepath.Join(srcBase, p)
	}
	if _, err := os.Stat(srcBase); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "source folder not found on disk"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "stat source failed"})
		return
	}

	// destination base: uploads/<username>/<dest trimmed>
	destBase := userDir
	if dest != "/" {
		p := strings.TrimPrefix(dest, "/")
		destBase = filepath.Join(destBase, p)
	}
	if err := os.MkdirAll(destBase, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create destination"})
		return
	}

	// final destination folder path will be destBase/<folderName>
	folderName := filepath.Base(strings.TrimRight(src, "/"))
	finalDest := filepath.Join(destBase, folderName)

	// gérer collision de nom de dossier
	if _, err := os.Stat(finalDest); err == nil {
		for i := 1; ; i++ {
			candidate := fmt.Sprintf("%s_%d", folderName, i)
			candidatePath := filepath.Join(destBase, candidate)
			if _, err := os.Stat(candidatePath); os.IsNotExist(err) {
				finalDest = candidatePath
				folderName = candidate
				break
			}
		}
	}

	// déplacer dossier sur disque
	if err := os.Rename(srcBase, finalDest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot move folder on disk", "err": err.Error()})
		return
	}

	// Mise à jour DB: tous les enregistrements DriveFiles dont file_path commence par src
	// Ex: src = "/a/b/"  -> remplacer préfixe par dest + folderName + "/"
	newPrefix := dest
	if !strings.HasSuffix(newPrefix, "/") {
		newPrefix = newPrefix + "/"
	}
	if newPrefix == "/" {
		newPrefix = "/" + folderName + "/"
	} else {
		newPrefix = strings.TrimSuffix(newPrefix, "/") + "/" + folderName + "/"
	}

	oldPrefix := src

	// Transaction DB: mettre à jour les paths relatifs
	tx, err := db.Begin()
	if err != nil {
		// rollback physique
		_ = os.Rename(finalDest, srcBase)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db begin failed"})
		return
	}

	res, err := tx.Exec(`UPDATE DriveFiles SET file_path = $1 || substring(file_path from char_length($2)+1)
				WHERE user_id=$3 AND (file_path = $2 OR file_path LIKE $2 || '%')`, newPrefix, oldPrefix, userID)
	if err != nil {
		tx.Rollback()
		// tenter rollback disque
		_ = os.Rename(finalDest, srcBase)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db update failed", "err": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		// tenter rollback disque
		_ = os.Rename(finalDest, srcBase)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db commit failed"})
		return
	}

	affected, _ := res.RowsAffected()
	c.JSON(http.StatusOK, gin.H{
		"message":         "folder moved successfully",
		"moved_folder":    folderName,
		"new_parent":      dest,
		"db_rows_updated": affected,
	})
}

func saveDownloadTokens() {
	const tokenFile = "/tmp/download_tokens.json"

	downloadTokensLock.RLock()
	defer downloadTokensLock.RUnlock()

	// créer un fichier temporaire puis le remplacer (atomic write)
	tmpFile := tokenFile + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		log.Println("⚠️ erreur création fichier token:", err)
		return
	}
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(downloadTokens); err != nil {
		log.Println("⚠️ erreur écriture JSON tokens:", err)
		_ = f.Close()
		_ = os.Remove(tmpFile)
		return
	}
	f.Close()
	_ = os.Rename(tmpFile, tokenFile)
}

func onlyofficeConfig(c *gin.Context) {
	fileID := c.Query("file_id")
	token := c.Query("token")
	username := c.Query("username")

	// 1️⃣ Vérifier la session utilisateur
	session, ok := sessions[token]
	if !ok || session.Username != username || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// 2️⃣ Récupérer les infos du fichier
	var file DriveFile
	err := db.QueryRow(`
        SELECT file_id, user_id, file_name, file_path, file_type 
        FROM DriveFiles 
        WHERE file_id=$1
    `, fileID).Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	// 3️⃣ Déterminer le type de document
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.FileName), "."))
	fileType := ext
	documentType := "text"
	switch ext {
	case "xlsx", "xls", "csv":
		documentType = "spreadsheet"
	case "pptx", "ppt":
		documentType = "presentation"
	}

	// 4️⃣ Calculer le baseURL (priorité à l'env)
	baseURL := os.Getenv("ONLYOFFICE_BACKEND_URL")
	if baseURL == "" {
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		host := c.Request.Host
		if host == "" {
			host = "localhost:8080"
		}
		// corrige host pour être atteignable depuis le conteneur OnlyOffice
		if strings.HasPrefix(host, "localhost") || strings.HasPrefix(host, "127.") {
			host = "host.docker.internal:8080"
		}
		baseURL = fmt.Sprintf("%s://%s", scheme, host)
	}

	// 5️⃣ Créer le token temporaire et le persister
	dtok := uuid.NewString()
	newToken := downloadToken{
		FileID:   file.FileID,
		Username: username,
		Expires:  time.Now().Add(10 * time.Minute),
	}

	downloadTokensLock.Lock()
	downloadTokens[dtok] = newToken
	downloadTokensLock.Unlock()

	// Sauvegarde immédiate sur disque (sécurité anti-redémarrage)
	saveDownloadTokens()

	// 6️⃣ Construire les URLs OnlyOffice
	downloadURL := fmt.Sprintf("%s/api/drive/download?download_token=%s", baseURL, dtok)
	callbackURL := fmt.Sprintf("%s/api/onlyoffice/callback?file_id=%s&username=%s&download_token=%s", baseURL, fileID, username, dtok)

	// 7️⃣ Construire la config OnlyOffice
	config := gin.H{
		"document": gin.H{
			"fileType": fileType,
			"key":      fmt.Sprintf("%d-%d", file.UserID, file.FileID),
			"title":    file.FileName,
			"url":      downloadURL,
		},
		"documentType": documentType,
		"editorConfig": gin.H{
			"callbackUrl": callbackURL,
			"mode":        "edit",
			"user": gin.H{
				"id":   fmt.Sprintf("%d", file.UserID),
				"name": username,
			},
		},
		"type": "desktop",
	}

	log.Printf("onlyoffice config OK: user=%s file_id=%s token=%s", username, fileID, dtok)
	c.JSON(http.StatusOK, config)
}

// onlyofficeCallback gère les callbacks envoyés par DocumentServer (OnlyOffice).
// - Valide download_token (présent dans la query string)
// - Selon le status, télécharge le fichier depuis l'URL fournie par DocumentServer et remplace le fichier sur le disque
// - Répond { "error": 0 } à DocumentServer en cas de succès
func onlyofficeCallback(c *gin.Context) {
	// Query params
	downloadTokenStr := c.Query("download_token")
	fileIDQ := c.Query("file_id")
	usernameQ := c.Query("username")

	if downloadTokenStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": 1, "message": "missing download_token"})
		return
	}

	// Validate token exists
	dt, ok := downloadTokens[downloadTokenStr]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": 1, "message": "invalid download_token"})
		return
	}

	// Expiry check
	if time.Now().After(dt.Expires) {
		delete(downloadTokens, downloadTokenStr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": 1, "message": "download_token expired"})
		return
	}

	// Optionally verify file_id and username match the token (safer)
	if fileIDQ != "" && fmt.Sprintf("%d", dt.FileID) != fileIDQ {
		c.JSON(http.StatusBadRequest, gin.H{"error": 1, "message": "file_id mismatch with token"})
		return
	}
	if usernameQ != "" && dt.Username != usernameQ {
		c.JSON(http.StatusBadRequest, gin.H{"error": 1, "message": "username mismatch with token"})
		return
	}

	// Parse JSON body sent by OnlyOffice
	var payload struct {
		Status     int    `json:"status"`
		URL        string `json:"url"`        // URL where OnlyOffice stored the updated file
		FileURL    string `json:"fileUrl"`    // sometimes named fileUrl
		ChangesURL string `json:"changesurl"` // sometimes provided
		EndConvert bool   `json:"endConvert"`
	}
	if err := c.BindJSON(&payload); err != nil {
		// If body is empty or not JSON, still respond OK to avoid retries but log
		log.Printf("onlyoffice callback: invalid json body: %v", err)
		c.JSON(http.StatusOK, gin.H{"error": 0})
		return
	}

	log.Printf("onlyoffice callback: token=%s file_id=%d user=%s status=%d url=%s fileUrl=%s changesurl=%s",
		downloadTokenStr, dt.FileID, dt.Username, payload.Status, payload.URL, payload.FileURL, payload.ChangesURL)

	// Decide quelle URL utiliser pour récupérer le fichier (priorité : URL > FileURL > ChangesURL)
	sourceURL := payload.URL
	if sourceURL == "" {
		sourceURL = payload.FileURL
	}
	if sourceURL == "" {
		sourceURL = payload.ChangesURL
	}

	// OnlyOffice statuses: we handle the common ones that require saving:
	// 2 = MustSave (editing finished, save new version)
	// 3 = ForceSave (save anyway)
	// for other statuses, just acknowledge
	if payload.Status != 2 && payload.Status != 3 {
		// nothing to save — acknowledge
		c.JSON(http.StatusOK, gin.H{"error": 0})
		return
	}

	if sourceURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": 1, "message": "no source URL provided by OnlyOffice"})
		return
	}

	// Retrieve file meta from DB (path, name)
	var file DriveFile
	err := db.QueryRow(`
        SELECT file_id, user_id, file_name, file_path, file_type
        FROM DriveFiles
        WHERE file_id = $1
    `, dt.FileID).Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath, &file.FileType)
	if err != nil {
		log.Printf("onlyoffice callback: could not load file meta for file_id=%d: %v", dt.FileID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "file not found"})
		return
	}

	// Download the updated file from DocumentServer
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		log.Printf("onlyoffice callback: bad source URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "invalid source url"})
		return
	}

	// If your DocumentServer requires special headers, add here (normally not)
	// req.Header.Set("User-Agent", "MyApp/1.0")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("onlyoffice callback: error fetching file from OnlyOffice URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "failed to download updated file"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		log.Printf("onlyoffice callback: non-200 from source URL: %d body=%s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": fmt.Sprintf("source returned %d", resp.StatusCode)})
		return
	}

	// Write atomically: save to tmp then rename
	tmpPath := file.FilePath + ".oo_tmp"
	out, err := os.Create(tmpPath)
	if err != nil {
		log.Printf("onlyoffice callback: cannot create tmp file %s: %v", tmpPath, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "cannot create tmp file"})
		return
	}
	n, err := io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		log.Printf("onlyoffice callback: error writing tmp file: %v", err)
		_ = os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "failed to save file"})
		return
	}

	// atomically replace
	if err := os.Rename(tmpPath, file.FilePath); err != nil {
		log.Printf("onlyoffice callback: rename tmp -> final failed: %v", err)
		_ = os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": 1, "message": "failed to finalize file save"})
		return
	}

	// update DB (e.g. updated_at, optionally size)
	_, err = db.Exec(`UPDATE DriveFiles SET updated_at = NOW(), file_size = $1 WHERE file_id = $2`, n, file.FileID)
	if err != nil {
		log.Printf("onlyoffice callback: failed to update DB for file_id=%d: %v", file.FileID, err)
		// we succeeded saving the file; still acknowledge OnlyOffice but log DB error
		c.JSON(http.StatusOK, gin.H{"error": 0})
		// optionally return here
		delete(downloadTokens, downloadTokenStr)
		return
	}

	// success: invalidate the download token
	delete(downloadTokens, downloadTokenStr)

	log.Printf("onlyoffice callback: saved file_id=%d bytes=%d path=%s", file.FileID, n, file.FilePath)
	c.JSON(http.StatusOK, gin.H{"error": 0})
}
