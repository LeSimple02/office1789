package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// Proxy pour Roundcube - contourne les restrictions CORS et X-Frame-Options
func roundcubeProxy(c *gin.Context) {
	// URL de Roundcube (localhost:8081 car le backend Go tourne hors Docker)
	target, _ := url.Parse("http://localhost:8081")
	
	// Créer un reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Modifier la requête
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		// Enlever le préfixe /roundcube du path
		req.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/roundcube")
		if req.URL.Path == "" {
			req.URL.Path = "/"
		}
		req.URL.RawQuery = c.Request.URL.RawQuery
		req.Host = target.Host
		
		// Copier les headers
		for key, values := range c.Request.Header {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}
	
	// Modifier la réponse pour autoriser l'iframe
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Supprimer les headers qui bloquent l'iframe
		resp.Header.Del("X-Frame-Options")
		resp.Header.Del("Content-Security-Policy")
		
		// Ajouter les headers CORS
		resp.Header.Set("Access-Control-Allow-Origin", "*")
		resp.Header.Set("Access-Control-Allow-Credentials", "true")
		resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		return nil
	}
	
	// Exécuter le proxy
	proxy.ServeHTTP(c.Writer, c.Request)
}

// Gérer les OPTIONS (preflight CORS)
func roundcubeProxyOptions(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
	c.Status(http.StatusOK)
}

// Alternative : Proxy simple qui récupère le contenu et le sert
func roundcubeProxySimple(c *gin.Context) {
	// Construire l'URL Roundcube
	roundcubeURL := "http://localhost:8081" + c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		roundcubeURL += "?" + c.Request.URL.RawQuery
	}

	// Créer la requête
	req, err := http.NewRequest(c.Request.Method, roundcubeURL, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de proxy"})
		return
	}

	// Copier les headers
	for key, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur de connexion à Roundcube"})
		return
	}
	defer resp.Body.Close()

	// Copier les headers de réponse (sauf ceux qui bloquent l'iframe)
	for key, values := range resp.Header {
		if !strings.EqualFold(key, "X-Frame-Options") && 
		   !strings.EqualFold(key, "Content-Security-Policy") {
			for _, value := range values {
				c.Header(key, value)
			}
		}
	}

	// Ajouter les headers CORS
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")

	// Copier le corps de la réponse
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}
