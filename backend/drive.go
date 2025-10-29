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
	"strconv"
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
// === downloadFile ===
func downloadFile(c *gin.Context) {
	token := c.Query("token")
	fileIDStr := c.Query("file_id")
	username := c.Query("username")
	// if client explicitly asks download, force attachment
	downloadMode := c.Query("download") // "1" or "true" => attachment

	// session check: require valid token; if username not provided, use session.Username
	session, ok := sessions[token]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}
	if username == "" {
		username = session.Username
	}
	if session.Username != username || username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// parse file_id
	fid, err := strconv.Atoi(fileIDStr)
	if err != nil {
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

	// owner must match session user
	if ownerUsername != session.Username {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// build path on disk
	filePath := filepath.Join("uploads", ownerUsername)
	if file.FilePath != "" && file.FilePath != "/" {
		filePath = filepath.Join(filePath, file.FilePath)
	}
	filePath = filepath.Join(filePath, file.FileName)

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
	// read head bytes to detect if needed
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
	// Reset file offset and use ServeContent for range requests / proper headers
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
