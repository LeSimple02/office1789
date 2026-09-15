<script setup>
import AccountProfile from "@/components/AccountProfile.vue"
import AccountChange from "@/components/AccountChange.vue"
import { computed } from "vue"
import { useRoute } from 'vue-router'

const route = useRoute()
const showProfile = computed(() => route.path === '/account')
const showEdit = computed(() => route.path === '/account/edit')
</script>

<template>
  <AccountProfile v-if="showProfile" />
  <AccountChange v-else-if="showEdit" />
</template>
