<script setup>
import { gls } from "@/stores/global"
import router from "@/router/index"
import { ref, onMounted, watch } from "vue"
import { useDark } from '@vueuse/core'

if (gls().log != 1) {
  router.push("login")
}

const mailUrl = ref("http://localhost:8081")
const loading = ref(true)
const mailIframe = ref(null)
const isDark = useDark()

// Surveiller les changements de mode dark
watch(isDark, (newValue) => {
  // Envoyer un message à l'iframe Roundcube pour changer le mode
  if (mailIframe.value && mailIframe.value.contentWindow) {
    mailIframe.value.contentWindow.postMessage({
      type: 'darkModeChange',
      isDark: newValue
    }, '*')
  }
})

onMounted(async () => {
  try {
    // Récupérer le token de session depuis le store global
    const sessionToken = gls().sessionT
    console.log("📧 [MailView] Session token:", sessionToken ? "✓ Présent" : "✗ Absent")
    console.log("📧 [MailView] Username:", gls().username)
    
    if (!sessionToken) {
      console.error("❌ [MailView] Aucun token de session, redirection vers login")
      router.push("/login")
      return
    }

    // Demander un token SSO au backend
    console.log("🔄 [MailView] Demande de token SSO au backend...")
    const response = await fetch("http://localhost:8080/api/mail/sso-token", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "Authorization": sessionToken
      }
    })

    console.log("📡 [MailView] Réponse backend:", response.status, response.statusText)

    if (response.ok) {
      const data = await response.json()
      console.log("✅ [MailView] Token SSO reçu:", data)
      
      // Ajouter le token SSO + mode dark à l'URL de Roundcube
      const darkParam = isDark.value ? '&dark=1' : ''
      mailUrl.value = `http://localhost:8081/?sso_token=${data.sso_token}${darkParam}`
      console.log("🌐 [MailView] URL Roundcube:", mailUrl.value)
    } else {
      const errorData = await response.json()
      console.error("❌ [MailView] Erreur lors de la récupération du token SSO:", errorData)
    }
  } catch (error) {
    console.error("💥 [MailView] Erreur SSO:", error)
  } finally {
    loading.value = false
    console.log("✓ [MailView] Chargement terminé")
  }
})
</script>

<template>
  <div id="mail-area">
    <div v-if="loading" class="loading">
      <p>Connexion à votre boîte mail...</p>
    </div>
    <iframe 
      v-else
      ref="mailIframe"
      :src="mailUrl" 
      class="mail-iframe"
      title="Roundcube Webmail"
      frameborder="0"
    ></iframe>
  </div>
</template>

<style scoped>
#mail-area {
  width: 100%;
  height: 100vh;
  display: flex;
  overflow: hidden;
  border-radius: 10px;
}

.mail-iframe {
  width: 100%;
  height: 100%;
  border: none;
}

.loading {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  font-family: roboto;
  font-size: 18px;
  background: -webkit-linear-gradient(30deg, blue, red);
  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}
</style>
