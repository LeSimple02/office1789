<script setup>
import { ref } from "vue"
import { gls } from "@/stores/global"
import router from "@/router/index"

let username = ref('')
let passf1 = ref('')
let passf2 = ref('')
let email = ref('')
let phonenumber = ref('')

let usernameR = ref('')
let emailR = ref('')
let phonenumberR = ref('')

function verif() {
  if (passf1.value == passf2.value && passf2.value !== "" && passf1.value !== "")
    connect()
}

function connect() {
  let d = new Date()
  let dc = `${d.getFullYear()}-${(d.getMonth() + 1).toString().padStart(2, '0')}-${d.getDate()} ${d.getHours()}:${d.getMinutes()}`

  fetch("http://127.0.0.1:8080/api/subscribe", {
    method: "POST",
    mode: "cors",
    credentials: "same-origin",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      "username": username.value,
      "password": passf1.value,
      "email": email.value,
      "phonenumber": phonenumber.value,
      "datejoined": dc,
      "lastlogin": dc
    })
  }).then((v) => v.json()).then((v) => {
    if (v["Username"] !== "no" && v["Phone"] !== "no" && v["Email"] !== "no") {
      localStorage.setItem("log", 1)
      gls().log = 1
      gls().username = v["Username"]
      gls().sessionT = v["Token"]

      document.cookie = `name=${v["Username"]}; expires=${v["Expiry"]}; Secure`
      document.cookie = `sessionToken=${v["Token"]}; expires=${v["Expiry"]}; Secure`

      router.push("/mail")
    } else {
      if (v["Username"]) usernameR.value = 1
      if (v["Email"]) emailR.value = 1
      if (v["Phone"]) phonenumberR.value = 1
    }
  })
}
</script>

<template>
  <div class="form-container">
    <h1 class="form-title">{{ $t("contact") }} :</h1>
    <div class="form-box">
      <ul class="form-grid">
        <!-- Ligne Email -->
        <li>
          <label class="form-label" for="email">Email :</label>
          <div>
            <input id="email" v-model="email" type="email" placeholder="Email" required />
            <span v-if="emailR" class="error-message">{{ $t('email_error') }}</span>
          </div>
        </li>

        <!-- Ligne Username -->
        <li>
          <label class="form-label" for="username">Username :</label>
          <div>
            <div class="username-input">
              <input id="username" v-model="username" type="text" required />
              <span>@office1789.com</span>
            </div>
            <span v-if="usernameR" class="error-message">{{ $t('username_error') }}</span>
          </div>
        </li>

        <!-- Ligne Téléphone -->
        <li>
          <label class="form-label" for="phone">Téléphone :</label>
          <div>
            <input id="phone" v-model="phonenumber" type="tel" placeholder="Téléphone" />
            <span v-if="phonenumberR" class="error-message">{{ $t('phone_error') }}</span>
          </div>
        </li>

        <!-- Ligne Message -->
        <li class="message-row">
          <label class="form-label" for="message">Message :</label>
          <div>
            <textarea id="message" placeholder="Message" style="height: 200px;"></textarea>
          </div>
        </li>
      </ul>
      <button class="submit-button" @click="verif()">{{ $t('submit') }}</button>
    </div>
  </div>
</template>

<style scoped>
.form-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 50px;
}

.form-box {
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  padding: 30px;
  width: 100%;
  max-width: 800px;
}
.dark .form-box{
	  background: #1C1C1E;
}
.form-title {
  font-family: 'Roboto', sans-serif;
  margin-bottom: 20px;
  text-align: center;
}

/* GRID principal */
.form-grid {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  grid-template-columns: 200px 1fr; /* colonne labels + inputs */
  row-gap: 25px;
  column-gap: 20px;
}

.form-grid li {
  display: contents; /* chaque li "disparaît" pour bien aligner */
}

.form-label {
  font-weight: 500;
  font-family: roboto;
  text-align: right;
  display: flex;
  justify-content: flex-end;
  align-items: center; /* centrage vertical */
  height: 100%;
  padding-right: 10px;
}

.form-grid input,
.form-grid textarea {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 16px;
  box-sizing: border-box;
}

.username-input {
  display: flex;
}

.username-input input {
  flex: 1;
  border-top-right-radius: 0;
  border-bottom-right-radius: 0;
}

.username-input span {
  padding: 12px;
  background: #f5f5f5;
  border: 1px solid #ddd;
  border-left: none;
  border-top-right-radius: 6px;
  border-bottom-right-radius: 6px;
}

.error-message {
  color: red;
  font-size: 12px;
  display: block;
  margin-top: 5px;
}

.submit-button {
  background: #000;
  color: white;
  border: none;
  padding: 12px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 16px;
  margin-top: 30px;
  transition: background 0.3s;
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.submit-button:hover {
  background: #333;
}

.dark .submit-button {
  background: white;
  color: black;
}
</style>
