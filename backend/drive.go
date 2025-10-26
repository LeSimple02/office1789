package main

import (
	"database/sql"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

// === uploadFile ===
func uploadFile(c *gin.Context) {
	username := c.PostForm("username")
	token := c.PostForm("token")
	parentPath := c.PostForm("parent_path")

	if session, ok := sessions[token]; !ok || session.Username != username || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var userID int
	if err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", username).Scan(&userID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

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

	uploadDir := filepath.Join("uploads", username)
	if parentPath != "" && parentPath != "/" {
		uploadDir = filepath.Join(uploadDir, parentPath)
	}
	if err := os.MkdirAll(uploadDir, 0o755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot create upload dir"})
		return
	}

	var uploaded []gin.H
	for _, fh := range fileHeaders {
		safeName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fh.Filename))
		dst := filepath.Join(uploadDir, safeName)

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

		contentType := fh.Header.Get("Content-Type")
		if contentType == "" {
			if ext := filepath.Ext(fh.Filename); ext != "" {
				contentType = mime.TypeByExtension(ext)
			}
			if contentType == "" {
				contentType = "application/octet-stream"
			}
		}

		_, err = db.Exec(`INSERT INTO DriveFiles (user_id, file_name, file_path, file_size, file_type, date_uploaded)
		                  VALUES ($1, $2, $3, $4, $5, NOW())`,
			userID, safeName, parentPath, fh.Size, contentType)
		if err != nil {
			_ = os.Remove(dst)
			uploaded = append(uploaded, gin.H{"file": fh.Filename, "error": "db insert failed"})
			continue
		}

		uploaded = append(uploaded, gin.H{"file": fh.Filename, "status": "uploaded"})
	}

	c.JSON(http.StatusOK, gin.H{"uploaded": uploaded})
}

// === downloadFile ===
func downloadFile(c *gin.Context) {
	token := c.Query("token")
	fileID := c.Query("file_id")
	username := c.Query("username")

	if session, ok := sessions[token]; !ok || session.Username != username || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var file DriveFile
	var ownerUsername string
	err := db.QueryRow(`SELECT f.file_id, f.user_id, f.file_name, f.file_path, f.file_size, f.file_type, f.date_uploaded, u.username
	                    FROM DriveFiles f 
	                    JOIN Users u ON f.user_id=u.user_id 
	                    WHERE f.file_id=$1`, fileID).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath,
			&file.FileSize, &file.FileType, &file.DateUploaded, &ownerUsername)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db query failed"})
		return
	}

	if ownerUsername != username {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	filePath := filepath.Join("uploads", ownerUsername)
	if file.FilePath != "" && file.FilePath != "/" {
		filePath = filepath.Join(filePath, file.FilePath)
	}
	filePath = filepath.Join(filePath, file.FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk"})
		return
	}

	c.FileAttachment(filePath, file.FileName)
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

	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	var file DriveFile
	err := db.QueryRow(`SELECT file_id, user_id, file_name, file_path 
	                    FROM DriveFiles WHERE file_id=$1`, req.FileID).
		Scan(&file.FileID, &file.UserID, &file.FileName, &file.FilePath)
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

	userDir := filepath.Join("uploads", req.Username)
	trashDir := filepath.Join(userDir, ".trash")
	_ = os.MkdirAll(trashDir, 0o755)

	srcPath := filepath.Join(userDir)
	if file.FilePath != "" && file.FilePath != "/" && file.FilePath != "/.trash/" {
		srcPath = filepath.Join(srcPath, file.FilePath)
	}
	src := filepath.Join(srcPath, file.FileName)
	dst := filepath.Join(trashDir, file.FileName)

	if err := os.Rename(src, dst); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found on disk", "src": src})
		return
	}

	_, _ = db.Exec("UPDATE DriveFiles SET file_path=$1 WHERE file_id=$2", "/.trash/", req.FileID)

	c.JSON(http.StatusOK, gin.H{"message": "moved to trash", "file_id": req.FileID})
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
	_ = db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID)
	if file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	userDir := filepath.Join("uploads", req.Username)
	src := filepath.Join(userDir)
	if file.FilePath != "" && file.FilePath != "/" {
		src = filepath.Join(src, file.FilePath)
	}
	src = filepath.Join(src, file.FileName)

	_ = os.Remove(src)
	_, _ = db.Exec("DELETE FROM DriveFiles WHERE file_id=$1", req.FileID)

	c.JSON(http.StatusOK, gin.H{"message": "deleted permanently", "file_id": req.FileID})
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
