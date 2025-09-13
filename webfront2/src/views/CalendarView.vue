<script setup>
import { useI18n } from 'vue-i18n'
import { computed } from 'vue'
import FullCalendar from '@fullcalendar/vue3'
import dayGridPlugin from '@fullcalendar/daygrid'
import interactionPlugin from '@fullcalendar/interaction'

const { locale } = useI18n()

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, interactionPlugin],
  initialView: 'dayGridMonth',
  locale: locale.value,
  events: [
    { title: 'event 1', date: '2025-05-01' },
    { title: 'event 2', date: '2019-04-02' }
  ]
}))
</script>

<template>
  <FullCalendar :options="calendarOptions" />
</template>

<style scoped>
/* Style moderne pour FullCalendar */
:deep(.fc) {
  font-family: 'Roboto', Arial, sans-serif;
  background: rgba(245,245,247,0.85);
  border-radius: 24px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  padding: 24px;
}

:deep(.fc-toolbar) {
  background: linear-gradient(30deg, blue, red);
  font-family: roboto;
  border-radius: 16px;
  color: #fff;
  margin-bottom: 16px;
  padding: 12px 16px;
}

:deep(.fc-toolbar-title) {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 2px;
}

:deep(.fc-button) {
  background: #06c;
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  padding: 6px 16px;
  margin: 0 4px;
  transition: background 0.2s, box-shadow 0.2s;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  cursor: pointer;
}

:deep(.fc-button:hover),
:deep(.fc-button-active) {
  background: linear-gradient(30deg, blue, red);
  color: #fff;
  box-shadow: 0 4px 16px rgba(0,0,0,0.18);
}

:deep(.fc-daygrid-day) {
  border-radius: 8px;
  transition: background 0.2s;
}

:deep(.fc-daygrid-day:hover) {
  background: rgba(0,48,143,0.08);
}

:deep(.fc-event) {
  background: linear-gradient(30deg, blue, red);
  color: #fff;
  border: none;
  border-radius: 8px;
  font-size: 1rem;
  box-shadow: 0 2px 8px rgba(0,0,0,0.10);
  padding: 2px 8px;
}

:deep(.fc-daygrid-event-dot) {
  background: #06c;
}

@media (max-width: 700px) {
  :deep(.fc) {
    padding: 8px;
    border-radius: 12px;
  }
  :deep(.fc-toolbar-title) {
    font-size: 1.1rem;
  }
}
</style>